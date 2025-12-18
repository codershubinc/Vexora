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
	Model     = "gemma3"
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
	Summary          string   `json:"summary"`
	PolishedNote     string   `json:"polished_note"`
	TwitterDraft     string   `json:"twitter_draft"`
	LinkedinDraft    string   `json:"linkedin_draft"`
	InstagramCaption string   `json:"instagram_caption"`
	Tags             []string `json:"tags"`
}

// --- PUBLIC API ---

func MessageLlama(projectName, diff string) (*JournalResponse, error) {
	// NEW PROMPT: The "Editor" Persona

	systemPrompt := `
# Role
You are Vexora, an expert Tech Editor and Developer Brand Manager.

# Task
Rewrite the user's rough coding notes into polished content for multiple platforms.

# Content Guidelines
1. **Polished Note:** A clean, grammar-perfect version of the log for the internal archive.
2. **Twitter/X:** Short, punchy, "building in public" energy. Under 280 chars. 
   - *Requirement:* MUST end with: "\n\nCreated with Vexora Studio @codershubinc"
3. **LinkedIn:** Professional, "Engineering Thought Leadership" tone. 
   - Structure: Problem -> Solution -> Impact.
   - Use bullet points if listing features.
4. **Instagram:** Casual "Behind the Scenes" vibe.
   - Assume the image will be a screenshot of code or the terminal.
   - Use a hook like "Late night coding session..." or "Finally cracked this bug...".

# Output Format (JSON Only)
{
  "summary": "Clean 1-sentence summary",
  "polished_note": "The corrected internal log version.",
  "twitter_draft": "The tweet content.",
  "linkedin_draft": "The professional post content.",
  "instagram_caption": "The IG caption content.",
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
