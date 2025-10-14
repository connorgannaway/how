package ai

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/connorgannaway/how/internal/system"
)

type AnthropicProvider struct {
	client anthropic.Client
	model  string
}

func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &AnthropicProvider{
		client: client,
		model:  model,
	}
}

func (p *AnthropicProvider) Ask(ctx context.Context, question string, sysInfo *system.SystemInfo) (*Response, error) {
	systemPrompt := BuildSystemPrompt(sysInfo)
	userPrompt := BuildUserPrompt(question)

	message, err := p.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.Model(p.model),
		MaxTokens: 1024,
		System: []anthropic.TextBlockParam{
			{
				Type: "text",
				Text: systemPrompt,
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(userPrompt)),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("anthropic API error: %w", err)
	}

	if len(message.Content) == 0 {
		return nil, fmt.Errorf("no response from Anthropic")
	}

	var responseText string
	for _, block := range message.Content {
		textBlock := block.AsText()
		if textBlock.Type == "text" {
			responseText = textBlock.Text
			break
		}
	}

	if responseText == "" {
		return nil, fmt.Errorf("no text content in Anthropic response")
	}

	return ParseResponse(responseText), nil
}

func (p *AnthropicProvider) GetName() string {
	return "Anthropic"
}
