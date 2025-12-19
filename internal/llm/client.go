package llm

import (
	"fmt"
)

const (
	// CONFIG: Local Llama (Ollama)
	OllamaURL = "http://localhost:11434/v1/chat/completions"
	Model     = "gemma3"
)

type ollamaReq struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
	Format   string    `json:"format,omitempty"` // "json" or empty
	Stream   bool      `json:"stream"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaResp struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

const (
	TypeTwitter    = "twitter"
	TypeLinkedIn   = "linkedin"
	TypeInstagram  = "instagram"
	TypeNewsletter = "newsletter"
)

// Unified Entry Point
func GenerateContent(feedType, userNotes string) (string, error) {
	switch feedType {
	case TypeTwitter:
		return genTwitter(userNotes)
	case TypeLinkedIn:
		return genLinkedIn(userNotes)
	case TypeInstagram:
		return genInstagram(userNotes)
	case TypeNewsletter:
		return genNewsletter(userNotes)
	default:
		return "", fmt.Errorf("unsupported feed type: %s", feedType)
	}
}
