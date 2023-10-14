package openai

import (
	"context"
	"log/slog"

	"github.com/sashabaranov/go-openai"
)

type Config struct {
	Token string `envconfig:"OPENAI_TOKEN"`
}

type Client struct {
	client *openai.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		client: openai.NewClient(cfg.Token),
	}
}

func (c *Client) Summarize(ctx context.Context, activity string) (string, error) {
	resp, err := c.client.CreateChatCompletion(ctx,
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
