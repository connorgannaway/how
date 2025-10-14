package ai

import (
	"context"
	"fmt"

	"github.com/connorgannaway/how/internal/system"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAICompatibleProvider struct {
	client  *openai.Client
	model   string
	baseURL string
}

func NewOpenAICompatibleProvider(apiKey, model, baseURL string) *OpenAICompatibleProvider {
	opts := []option.RequestOption{option.WithBaseURL(baseURL)}

	// Only add API key if provided
	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}

	client := openai.NewClient(opts...)
	return &OpenAICompatibleProvider{
		client:  &client,
		model:   model,
		baseURL: baseURL,
	}
}

func (p *OpenAICompatibleProvider) Ask(ctx context.Context, question string, sysInfo *system.SystemInfo) (*Response, error) {
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
		return nil, fmt.Errorf("OpenAI-compatible API error: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI-compatible API")
	}

	responseText := chatCompletion.Choices[0].Message.Content
	return ParseResponse(responseText), nil
}

func (p *OpenAICompatibleProvider) GetName() string {
	return "OpenAI-Compatible"
}
