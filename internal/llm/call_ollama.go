package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func callOllama(sysPrompt, userMsg, format string) (string, error) {
	reqBody := ollamaReq{
		Model:  Model,
		Format: format,
		Stream: false,
		Messages: []message{
			{Role: "system", Content: sysPrompt},
			{Role: "user", Content: userMsg},
		},
	}

	jsonData, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", OllamaURL, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// 120s timeout for longer generations
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("❌ Ollama Connection Error: %v", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("❌ Ollama returned non-200 status: %s", resp.Status)
		return "", fmt.Errorf("ollama status: %s", resp.Status)
	}

	var oResp ollamaResp
	if err := json.NewDecoder(resp.Body).Decode(&oResp); err != nil {
		return "", err
	}

	if len(oResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from ollama")
	}

	return oResp.Choices[0].Message.Content, nil
}
