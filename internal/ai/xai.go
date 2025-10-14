package ai

import (
	"context"
	"fmt"

	"github.com/connorgannaway/how/internal/system"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type XAIProvider struct {
	client *openai.Client // No golang SDK for xAI, but is OpenAI-compatible
	model  string
}

func NewXAIProvider(apiKey, model string) *XAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey), option.WithBaseURL("https://api.x.ai/v1"))
	return &XAIProvider{
		client: &client,
		model:  model,
	}
}

func (p *XAIProvider) Ask(ctx context.Context, question string, sysInfo *system.SystemInfo) (*Response, error) {
	systemPrompt := BuildSystemPrompt(sysInfo)
	userPrompt := BuildUserPrompt(question)

	chatCompletion, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: p.model,
	})

	if err != nil {
		return nil, fmt.Errorf("xAI API error: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("no response from xAI")
	}

	responseText := chatCompletion.Choices[0].Message.Content
	return ParseResponse(responseText), nil
}

func (p *XAIProvider) GetName() string {
	return "xAI"
}
