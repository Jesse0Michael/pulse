package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

type SummaryRequest struct {
	Username     string
	Organization string
	Repository   string
	StartDate    *time.Time
	EndDate      *time.Time
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
	content, err := p.github.UserActivity(ctx, req.Username, req.Organization, req.Repository, req.StartDate, req.EndDate)
	if err != nil {
		return "", err
	}
	if content == "" {
		return "", fmt.Errorf("no activity found for the username %s", req.Username)
	}

	slog.With("content", content).InfoContext(ctx, "github user activity")

	content = fmt.Sprintf("Summarize the following github activity for the username %s:\n%s", req.Username, content)

	return p.openAI.Summarize(ctx, content)
}
