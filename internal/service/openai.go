package service

import (
	"context"
	"strings"

	"github.com/sashabaranov/go-openai"
)

type OpenAIConfig struct {
	URL   string `envconfig:"OPENAI_URL"`
	Token string `envconfig:"OPENAI_TOKEN"`
}

type OpenAI struct {
	client *openai.Client
}

func NewOpenAI(cfg OpenAIConfig) *OpenAI {
	config := openai.DefaultConfig(cfg.Token)
	if cfg.URL != "" {
		config.BaseURL = cfg.URL
	}

	return &OpenAI{
		client: openai.NewClientWithConfig(config),
	}
}

func (o *OpenAI) Summarize(ctx context.Context, content string) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	out := make([]string, len(resp.Choices))
	for i, choice := range resp.Choices {
		out[i] = choice.Message.Content
	}

	return strings.Join(out, "\n"), nil
}
