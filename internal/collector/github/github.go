package github

import (
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

func (c *Client) UserActivity(username string) (*github.Repository, error) {
	// return c.client.Activity.Repositories.Get(context.Background(), owner, repo)
	return nil, nil
}
