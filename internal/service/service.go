package service

import (
	"context"
	"fmt"
	"log/slog"
)

type SummaryRequest struct {
	Username string
}

type Pulser struct {
	github *Github
	openAI *OpenAI
}

func NewPulser(github *Github, openAI *OpenAI) *Pulser {
	return &Pulser{
		github: github,
		openAI: openAI,
	}
}

func (p *Pulser) Summary(ctx context.Context, req SummaryRequest) (string, error) {
	content, err := p.github.UserActivity(ctx, req.Username)
	if err != nil {
		return "", err
	}
	slog.With("content", content).InfoContext(ctx, "github user activity")

	content = fmt.Sprintf("Summarize the following github activity for the username %s:\n%s", req.Username, content)

	return p.openAI.Summarize(ctx, content)
}
