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

type NewsletterMeta struct {
	Subject string   `json:"subject_line"`
	Preview string   `json:"preview_text"`
	Tags    []string `json:"tags"` // Go can now handle the array
}

func genNewsletter(notes string) (string, error) {
	// --- Step A: Fetch Metadata (using specific struct) ---
	// We can't use fetchJSON here because it returns map[string]string
	// We invoke callOllama directly with "json" format
	rawMeta, err := callOllama(PromptNewsMeta, notes, "json")
	if err != nil {
		return "", err
	}

	var meta NewsletterMeta
	// Clean the JSON (remove potential markdown blocks) and Unmarshal
	if err := json.Unmarshal([]byte(cleanJSON(rawMeta)), &meta); err != nil {
		return "", fmt.Errorf("meta parse failed: %w", err)
	}

	// --- Step B: Fetch Body (Raw Text) ---
	bodyText, err := fetchText(PromptNewsBody, notes)
	if err != nil {
		return "", err
	}

	// --- Step C: Combine into Final JSON ---
	// We convert the Tags array to a string for the final output if needed,
	// or keep it as an array if your frontend supports it.
	finalOutput := map[string]interface{}{
		"subject_line": meta.Subject,
		"preview_text": meta.Preview,
		"tags":         meta.Tags, // This stays an array now!
		"body":         bodyText,
	}

	jsonBytes, err := json.Marshal(finalOutput)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func fetchJSON(sysPrompt, userMsg string) (map[string]string, error) {
	raw, err := callLLM(sysPrompt, userMsg, "json")
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

// fetchText calls the selected LLM for raw text generation (Markdown, etc)
func fetchText(sysPrompt, userMsg string) (string, error) {
	return callLLM(sysPrompt, userMsg, "")
}

func cleanJSON(input string) string {
	start := strings.Index(input, "{")
	end := strings.LastIndex(input, "}")
	if start == -1 || end == -1 {
		return "{}"
	}
	return input[start : end+1]
}
