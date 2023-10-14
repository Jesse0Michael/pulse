package github

import (
	"context"
	"log/slog"

	"github.com/google/go-github/v54/github"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/go-retryablehttp"
)

type Client struct {
	client *github.Client
}

func NewClient() *Client {
	transport := cleanhttp.DefaultPooledTransport()
	retryClient := retryablehttp.NewClient()
	retryClient.Logger = slog.Default()
	retryClient.RetryMax = 10
	retryClient.HTTPClient.Transport = transport

	return &Client{
		client: github.NewClient(retryClient.StandardClient()),
	}
}

func (c *Client) UserActivity(ctx context.Context, username string) (string, error) {
	events, resp, err := c.client.Activity.ListEventsPerformedByUser(ctx, username, false, nil)
	if err != nil {
		return "", err
	}
	if resp != nil {
		slog.Debug(resp.Status)
	}

	var out string
	for _, event := range events {
		slog.Info("event", "id", event.GetID(), "type", event.GetType(), "repo", event.GetRepo(), "org", event.GetOrg())
		out += string(*event.RawPayload) + "\n"
	}

	return out, nil
}
