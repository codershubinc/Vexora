package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type geminiReq struct {
	SystemInstruction *geminiInstruction `json:"system_instruction,omitempty"`
	Contents          []geminiContent    `json:"contents"`
	GenerationConfig  *geminiConfig      `json:"generationConfig,omitempty"`
}

type geminiInstruction struct {
	Parts geminiPart `json:"parts"`
}

type geminiContent struct {
	Role  string       `json:"role,omitempty"`
	Parts []geminiPart `json:"parts"`
}

type geminiPart struct {
	Text string `json:"text"`
}

type geminiConfig struct {
	ResponseMimeType string `json:"response_mime_type,omitempty"`
}

type geminiResp struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func callGemini(sysPrompt, userMsg, format string) (string, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY environment variable not set")
	}

	model := os.Getenv("GEMINI_MODEL")
	model = "gemini-2.5-flash-lite"
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", model, apiKey)

	reqBody := geminiReq{
		Contents: []geminiContent{
			{
				Role: "user",
				Parts: []geminiPart{
					{Text: userMsg},
				},
			},
		},
	}

	if sysPrompt != "" {
		reqBody.SystemInstruction = &geminiInstruction{
			Parts: geminiPart{Text: sysPrompt},
		}
	}

	if format == "json" {
		reqBody.GenerationConfig = &geminiConfig{
			ResponseMimeType: "application/json",
		}
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Gemini Connection Error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	var gResp geminiResp
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		log.Printf("❌ Gemini JSON Decode Error: %v", err)
		return "", err
	}

	if gResp.Error != nil {
		log.Printf("❌ Gemini API Error: %s", gResp.Error.Message)
		return "", fmt.Errorf("gemini error: %s", gResp.Error.Message)
	}

	if len(gResp.Candidates) == 0 || len(gResp.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("empty response from gemini")
	}

	return gResp.Candidates[0].Content.Parts[0].Text, nil
}
