package llm

import (
	"fmt"
	"log"
	"os"
)

const (
	// CONFIG: Local Llama (Ollama)
	OllamaURL = "http://localhost:11434/v1/chat/completions"
	Model     = "gemma3"

	ProviderOllama = "ollama"
	ProviderGemini = "gemini"
)

func getProvider() string {
	p := os.Getenv("LLM_PROVIDER")
	if p == "" {
		return "gemini"
	}
	return "gemini"
}

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

func callLLM(sysPrompt, userMsg, format string) (string, error) {
	provider := getProvider()
	log.Printf("🤖 Using LLM Provider: %s", provider)

	var res string
	var err error

	switch provider {
	case ProviderGemini:
		res, err = callGemini(sysPrompt, userMsg, format)
	case ProviderOllama:
		res, err = callOllama(sysPrompt, userMsg, format)
	default:
		log.Printf("⚠️ Unknown provider '%s', falling back to Ollama", provider)
		res, err = callOllama(sysPrompt, userMsg, format)
	}

	if err != nil {
		log.Printf("❌ LLM Error (%s): %v", provider, err)
	}
	return res, err
}
