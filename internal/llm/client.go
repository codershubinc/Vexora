package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	// CONFIG: Local Llama (Ollama)
	OllamaURL = "http://localhost:11434/v1/chat/completions"
	Model     = "llama3.2"
)

// --- DATA STRUCTURES ---

type OpenAIRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Format   string    `json:"format,omitempty"` // Force JSON
	Stream   bool      `json:"stream"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type JournalResponse struct {
	Summary      string   `json:"summary"`
	PolishedNote string   `json:"polished_note"`
	TwitterDraft string   `json:"twitter_draft"`
	Tags         []string `json:"tags"`
}

// --- PUBLIC API ---

func MessageLlama(projectName, diff string) (*JournalResponse, error) {
	// NEW PROMPT: The "Editor" Persona
	systemPrompt := `
# Role
You are Vexora, an expert Tech Editor.

# Task
The user will provide rough notes or a Copilot-generated summary from a coding session.
Your job is to rewrite this into a polished, engaging Developer Log or Social Media update.

# Tone
Professional but authentic. Fix grammar, improve flow, and make it sound like a passionate engineer.

# Output Format (JSON Only)
{
  "summary": "Clean 1-sentence summary",
  "polished_note": "The corrected, better version of the user's notes.",
  "twitter_draft": "A short, engaging tweet based on the note. ALWAYS end with: \n\nCreated with Vexora Studio @codershubinc",
  "tags": ["#tag1", "#tag2"]
}
`
	reqBody := OpenAIRequest{
		Model:  Model,
		Format: "json", // Strict JSON
		Stream: false,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: diff},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

	maxRetries := 3
	var lastErr error

	for i := 0; i < maxRetries; i++ {
		if i > 0 {
			fmt.Printf("⚠️ LLM returned empty response. Retrying (Attempt %d/%d)...\n", i+1, maxRetries)
			time.Sleep(2 * time.Second)
		}

		// Re-create request for each attempt because Body is consumed
		req, _ := http.NewRequest("POST", OllamaURL, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		// Wait up to 60s for analysis
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		defer resp.Body.Close()

		var result OpenAIResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			lastErr = err
			continue
		}

		if len(result.Choices) == 0 || result.Choices[0].Message.Content == "" {
			lastErr = fmt.Errorf("empty response from LLM")
			continue
		}

		content := result.Choices[0].Message.Content
		cleanContent := cleanJSON(content)

		var journal JournalResponse
		if err := json.Unmarshal([]byte(cleanContent), &journal); err != nil {
			lastErr = fmt.Errorf("failed to parse LLM JSON: %v", err)
			continue
		}

		return &journal, nil
	}

	return nil, fmt.Errorf("LLM failed after %d attempts: %v", maxRetries, lastErr)
}

// --- HELPERS ---

func cleanJSON(input string) string {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 {
		return "{}"
	}
	return input[start : end+1]
}
