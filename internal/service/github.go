package service

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/google/go-github/v54/github"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

type GithubConfig struct {
	URL   string `envconfig:"GITHUB_URL"`
	Token string `envconfig:"GITHUB_TOKEN"`
}

type Github struct {
	client *github.Client
}

func NewGithub(cfg GithubConfig) *Github {
	transport := cleanhttp.DefaultPooledTransport()
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = slog.Default()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Transport = transport

	var client *github.Client
	if cfg.Token != "" {
		client = github.NewTokenClient(context.Background(), cfg.Token)
	} else {
		client = github.NewClient(retryClient.StandardClient())
	}
	if cfg.URL != "" {
		client.BaseURL, _ = url.Parse(cfg.URL)
	}

	return &Github{
		client: client,
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
	if event == nil {
		return ""
	}

	payload, err := event.ParsePayload()
	if err != nil {
		return ""
	}

	switch event.GetType() {
	case "PushEvent":
		if e, ok := payload.(*github.PushEvent); ok && e != nil {
			return parsePushEvent(*e, event.Repo.GetName())
		}
	case "PullRequestEvent":
		if e, ok := payload.(*github.PullRequestEvent); ok && e != nil {
			return parsePullRequestEvent(*e, event.Repo.GetName())
		}
	}
	return ""
}

func parsePushEvent(e github.PushEvent, repository string) string {
	activity := fmt.Sprintf("pushed commits to repository %s", repository)
	if len(e.Commits) == 0 {
		return activity
	}

	commits := make([]string, len(e.Commits))
	for i, commit := range e.Commits {
		commits[i] = commit.GetMessage()
	}

	return fmt.Sprintf("%s\ncommit messages: %s", activity, strings.Join(commits, "\n"))
}

func parsePullRequestEvent(e github.PullRequestEvent, repository string) string {
	switch e.GetAction() {
	case "opened":
		activity := fmt.Sprintf("opened pull request in repository %s", repository)
		if e.PullRequest == nil || e.PullRequest.Body == nil {
			return activity
		}

		return fmt.Sprintf("%s\nbody: %s", activity, *e.PullRequest.Body)
	default:
		return ""
	}
}
