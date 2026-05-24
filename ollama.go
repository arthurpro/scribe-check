package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string         `json:"model"`
	Messages []chatMessage  `json:"messages"`
	Stream   bool           `json:"stream"`
	Format   any            `json:"format,omitempty"`
	Options  map[string]any `json:"options,omitempty"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
}

type Finding struct {
	Passage            string `json:"passage,omitempty"`
	Entity             string `json:"entity,omitempty"`
	Claim              string `json:"claim,omitempty"`
	LocationHint       string `json:"location_hint,omitempty"`
	Concern            string `json:"concern"`
	ClosestSourceMatch string `json:"closest_source_match,omitempty"`
}

type Report struct {
	QuotesFlagged    []Finding `json:"quotes_flagged"`
	NamesFlagged     []Finding `json:"names_flagged"`
	SpecificsFlagged []Finding `json:"specifics_flagged"`
	StyleDrift       []Finding `json:"style_drift"`
	TemporalDrift    []Finding `json:"temporal_drift"`
	Summary          string    `json:"summary"`
}

func hasFindings(r Report) bool {
	return len(r.QuotesFlagged)+len(r.NamesFlagged)+len(r.SpecificsFlagged)+
		len(r.StyleDrift)+len(r.TemporalDrift) > 0
}

// runReview calls Ollama once, retries once on malformed JSON with a corrective
// nudge, and returns the parsed Report plus the raw JSON body.
func runReview(host, model string, numCtx int, article Article, sources []Source) (Report, string, error) {
	user := buildUserPrompt(article, sources)
	messages := []chatMessage{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: user},
	}

	report, raw, err := callOllama(host, model, numCtx, messages)
	if err == nil {
		return report, raw, nil
	}

	// Retry once with a corrective system message.
	messages = append(messages,
		chatMessage{Role: "assistant", Content: raw},
		chatMessage{Role: "user", Content: "Your previous reply was not valid JSON matching the schema. Reply ONLY with the JSON object — no prose, no markdown fences, no commentary. Use empty arrays for categories with no findings."},
	)
	report, raw, err = callOllama(host, model, numCtx, messages)
	if err != nil {
		return Report{}, raw, fmt.Errorf("retry failed: %w", err)
	}
	return report, raw, nil
}

func callOllama(host, model string, numCtx int, messages []chatMessage) (Report, string, error) {
	reqBody := chatRequest{
		Model:    model,
		Messages: messages,
		Stream:   false,
		Format:   "json",
		Options: map[string]any{
			"temperature": 0.1,
			"num_ctx":     numCtx,
			"seed":        42,
		},
	}
	buf, err := json.Marshal(reqBody)
	if err != nil {
		return Report{}, "", err
	}

	client := &http.Client{Timeout: 15 * time.Minute}
	resp, err := client.Post(host+"/api/chat", "application/json", bytes.NewReader(buf))
	if err != nil {
		return Report{}, "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Report{}, "", err
	}
	if resp.StatusCode != 200 {
		return Report{}, string(body), fmt.Errorf("ollama HTTP %d: %s", resp.StatusCode, body)
	}

	var cr chatResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return Report{}, string(body), fmt.Errorf("decode ollama envelope: %w", err)
	}

	raw := strings.TrimSpace(cr.Message.Content)
	// Some models occasionally wrap JSON in fences despite format=json; strip.
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var report Report
	if err := json.Unmarshal([]byte(raw), &report); err != nil {
		return Report{}, raw, fmt.Errorf("decode report JSON: %w", err)
	}
	return report, raw, nil
}
