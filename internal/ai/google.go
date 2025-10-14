package ai

import (
	"context"
	"fmt"

	"github.com/connorgannaway/how/internal/system"
	"google.golang.org/genai"
)

type GoogleProvider struct {
	client *genai.Client
	model  string
}

func NewGoogleProvider(apiKey, model string) *GoogleProvider {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		// If client creation fails, return a provider with nil client
		// The error will be caught and displayed by UI when Ask is called
		return &GoogleProvider{
			client: nil,
			model:  model,
		}
	}
	return &GoogleProvider{
		client: client,
		model:  model,
	}
}

func (p *GoogleProvider) Ask(ctx context.Context, question string, sysInfo *system.SystemInfo) (*Response, error) {
	if p.client == nil {
		return nil, fmt.Errorf("google client not initialized")
	}

	systemPrompt := BuildSystemPrompt(sysInfo)
	userPrompt := BuildUserPrompt(question)

	// Combine system and user prompts for Gemini
	fullPrompt := fmt.Sprintf("%s\n\nUser question: %s", systemPrompt, userPrompt)

	// Create contents using Text helper
	contents := genai.Text(fullPrompt)

	resp, err := p.client.Models.GenerateContent(ctx, p.model, contents, nil)
	if err != nil {
		return nil, fmt.Errorf("google API error: %w", err)
	}

	if len(resp.Candidates) == 0 {
		return nil, fmt.Errorf("no response from Google")
	}

	candidate := resp.Candidates[0]
	if candidate.Content == nil || len(candidate.Content.Parts) == 0 {
		return nil, fmt.Errorf("no content in Google response")
	}

	var responseText string
	for _, part := range candidate.Content.Parts {
		if part.Text != "" {
			responseText += part.Text
		}
	}

	if responseText == "" {
		return nil, fmt.Errorf("no text content in Google response")
	}

	return ParseResponse(responseText), nil
}

func (p *GoogleProvider) GetName() string {
	return "Google"
}
