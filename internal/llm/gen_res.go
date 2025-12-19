package llm

import (
	"encoding/json"
	"fmt"
	"strings"
)

// --- Platform Specific Logic ---

func genTwitter(notes string) (string, error) {
	return fetchText(PromptTwitter, notes)
}

func genLinkedIn(notes string) (string, error) {
	return fetchText(PromptLinkedIn, notes)
}

func genInstagram(notes string) (string, error) {
	return fetchText(PromptInstagram, notes)
}

// The Complex One: Newsletter (Two-Pass Strategy)
func genNewsletter(notes string) (string, error) {
	result := make(map[string]string)

	// Pass 1: Metadata (JSON)
	metaPrompt := `
    # Task
    Generate metadata for a technical newsletter.
    # JSON Output: {"subject": "Catchy Title", "preview": "Short hook", "tags": "#tag1"}
    `
	metaData, err := fetchJSON(metaPrompt, notes)
	if err != nil {
		return "", err
	}
	for k, v := range metaData {
		result[k] = v
	}

	// Pass 2: Body (Raw Text) - Reusing the "Senior Engineer" Prompt
	bodyPrompt := `
    # Role
    Senior Principal Engineer.
    # Task
    Write a Masterclass-level newsletter deep dive.
    - Include Code Snippets.
    - Focus on architectural trade-offs.
    - Output: Raw Markdown ONLY.
    `
	bodyText, err := fetchText(bodyPrompt, notes)
	if err != nil {
		return "", err
	}
	result["body"] = bodyText

	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func fetchJSON(sysPrompt, userMsg string) (map[string]string, error) {
	raw, err := callOllama(sysPrompt, userMsg, "json")
	if err != nil {
		return nil, err
	}

	clean := cleanJSON(raw)
	var result map[string]string
	if err := json.Unmarshal([]byte(clean), &result); err != nil {
		return nil, fmt.Errorf("json parse failed: %w", err)
	}
	return result, nil
}

// fetchText calls Ollama for raw text generation (Markdown, etc)
func fetchText(sysPrompt, userMsg string) (string, error) {
	return callOllama(sysPrompt, userMsg, "")
}

func cleanJSON(input string) string {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 {
		return "{}"
	}
	return input[start : end+1]
}
