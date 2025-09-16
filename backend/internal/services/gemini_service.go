package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GeminiServiceSimple is a minimal interface used by some components/tests.
type GeminiServiceSimple interface {
	GenerateItinerary(ctx context.Context, prompt string, contextBlocks []string) (string, error)
}

// geminiClient is a thin adapter that implements GeminiServiceSimple and
// delegates to the existing GeminiService where possible. It also provides
// a standalone REST call fallback when an API key is available.
type geminiClient struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// NewGeminiServiceSimple creates a new minimal Gemini client.
// If apiKey is empty, the client will return deterministic mock responses.
func NewGeminiServiceSimple(apiKey, baseURL string) GeminiServiceSimple {
	if baseURL == "" {
		baseURL = "https://generativelanguage.googleapis.com/v1beta"
	}
	return &geminiClient{
		apiKey:  apiKey,
		baseURL: baseURL,
		client:  &http.Client{Timeout: 20 * time.Second},
	}
}

// GenerateItinerary sends prompt + context to Gemini and returns the raw model output (string).
func (g *geminiClient) GenerateItinerary(ctx context.Context, prompt string, contextBlocks []string) (string, error) {
	// If no API key, return deterministic mock for local development
	if strings.TrimSpace(g.apiKey) == "" {
		return mockItineraryResponse(), nil
	}

	// Build request payload following the Generative API simple pattern
	payload := map[string]interface{}{
		"prompt":      prompt,
		"context":     contextBlocks,
		"maxTokens":   1500,
		"temperature": 0.2,
	}
	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/models/gemini-pro:generateContent?key=%s", strings.TrimRight(g.baseURL, "/"), g.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, io.NopCloser(strings.NewReader(string(body))))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("gemini api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// Try to unmarshal a common response shape used elsewhere in the codebase
	var parsed struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(respBody, &parsed); err != nil {
		// If parsing fails, return raw text as a fallback
		return string(respBody), nil
	}

	if len(parsed.Candidates) == 0 || len(parsed.Candidates[0].Content.Parts) == 0 {
		return "", errors.New("no content found in gemini response")
	}

	return parsed.Candidates[0].Content.Parts[0].Text, nil
}

func mockItineraryResponse() string {
	// deterministic, small JSON for tests
	return `{"variant":"mock","total_cost":1000,"currency":"INR","confidence":0.5,"days":[]}`
}
