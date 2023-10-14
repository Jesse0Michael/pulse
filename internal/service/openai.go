package service

import (
	"context"
	"log/slog"

	"github.com/sashabaranov/go-openai"
)

type OpenAIConfig struct {
	Token string `envconfig:"OPENAI_TOKEN"`
}

type OpenAI struct {
	client *openai.Client
}

func NewOpenAI(cfg OpenAIConfig) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(cfg.Token),
	}
}

func (o *OpenAI) Summarize(ctx context.Context, activity string) (string, error) {
	resp, err := o.client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	var out string
	for _, choice := range resp.Choices {
		slog.Info("choice message", "content", choice.Message.Content)
		out += choice.Message.Content + "\n"
	}

	return out, nil
}
