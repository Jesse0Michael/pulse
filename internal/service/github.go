package service

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/go-github/v54/github"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

type GithubConfig struct {
}

type Github struct {
	client *github.Client
}

func NewGithub(_ GithubConfig) *Github {
	transport := cleanhttp.DefaultPooledTransport()
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = slog.Default()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Transport = transport

	return &Github{
		client: github.NewClient(retryClient.StandardClient()),
	}
}

func (g *Github) UserActivity(ctx context.Context, username string) (string, error) {
	opts := github.ListOptions{}
	events, resp, err := g.client.Activity.ListEventsPerformedByUser(ctx, username, false, &opts)
	if err != nil {
		return "", err
	}
	if resp != nil {
		slog.Debug(resp.Status)
	}

	var out string
	for _, event := range events {
		slog.Info("event", "id", event.GetID(), "type", event.GetType(), "repo", event.GetRepo(), "org", event.GetOrg())
		if activity := eventActivity(event); activity != "" {
			out += activity + "\n"
		}
	}

	return out, nil
}

func eventActivity(event *github.Event) string {
	switch event.GetType() {
	case "PushEvent":
		activity := fmt.Sprintf("pushed commits to repository %s", *event.Repo.Name)
		payload, err := event.ParsePayload()
		e, ok := payload.(*github.PushEvent)
		if err != nil || !ok {
			slog.With("error", err).Error("error parsing github event")
			return activity
		}

		commits := make([]string, len(e.Commits))
		for i, commit := range e.Commits {
			commits[i] = commit.GetMessage()
		}

		return fmt.Sprintf("%s\ncommit messages: %s", activity, strings.Join(commits, "\n"))
	case "PullRequestEvent":
		activity := fmt.Sprintf("pushed commits to repository %s", *event.Repo.Name)
		payload, err := event.ParsePayload()
		e, ok := payload.(*github.PullRequestEvent)
		if err != nil || !ok {
			slog.With("error", err).Error("error parsing github event")
			return activity
		}

		switch e.GetAction() {
		case "opened":
			activity := fmt.Sprintf("opened pull request in repository %s", *event.Repo.Name)
			if e.PullRequest == nil || e.PullRequest.Body == nil {
				return activity
			}

			return fmt.Sprintf("%s\nbody: %s", activity, *e.PullRequest.Body)
		default:
			return ""
		}
	default:
		return ""
	}
}
