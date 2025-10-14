package ai

import (
	"context"
	"fmt"

	"github.com/connorgannaway/how/internal/system"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type OpenAIProvider struct {
	client *openai.Client
	model  string
}

func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	client := openai.NewClient(option.WithAPIKey(apiKey))
	return &OpenAIProvider{
		client: &client,
		model:  model,
	}
}

func (p *OpenAIProvider) Ask(ctx context.Context, question string, sysInfo *system.SystemInfo) (*Response, error) {
	systemPrompt := BuildSystemPrompt(sysInfo)
	userPrompt := BuildUserPrompt(question)

	chatCompletion, err := p.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: openai.ChatModel(p.model),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(chatCompletion.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	responseText := chatCompletion.Choices[0].Message.Content
	return ParseResponse(responseText), nil
}

func (p *OpenAIProvider) GetName() string {
	return "OpenAI"
}
