package command

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jesse0michael/pulse/internal/service"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type GithubConfig struct {
	Github service.GithubConfig
	AI     service.OpenAIConfig
}

type Github struct {
	pulser       *service.Pulser
	output       string
	Username     string
	Organization string
	Repository   string
	StartDate    *time.Time
	EndDate      *time.Time
}

// NewGithub creates a new Github service.
func NewGithub() *Github {
	return &Github{}
}

// Command will return the cobra command structure that can be executed.
func (c *Github) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "github [username]",
		Short:   "generate a summary of a github user's activity",
		PreRunE: c.Init,
		PostRun: c.Output,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(cmd.Context())
		},
	}

	cmd.Flags().StringVar(&c.Organization, "organization", "", "filter results on an organization")
	cmd.Flags().StringVar(&c.Repository, "repository", "", "filter results on a repository")
	cmd.Flags().String("startDate", "", "filter results after a start date")
	cmd.Flags().String("endDate", "", "filter results before an end date")

	cmd.SetUsageTemplate(cmd.UsageTemplate() + `
Environment:
  GITHUB_URL         the url for accessing the GitHub API
  GITHUB_TOKEN       the authentication token to use with the GitHub API
  OPENAI_URL         the url for accessing the OpenAI API
  OPENAI_TOKEN       the authentication token to use with the OpenAI API
`)

	return cmd
}

// Init will initialize GitHub dependencies.
func (c *Github) Init(cmd *cobra.Command, args []string) error {
	var cfg GithubConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}

	github := service.NewGithub(cfg.Github)
	openAI := service.NewOpenAI(cfg.AI)
	c.pulser = service.NewPulser(github, openAI)

	if len(args) > 0 {
		c.Username = args[0]
	}

	if b, err := cmd.Flags().GetString("startDate"); err == nil {
		if d, err := time.Parse(time.RFC3339, b); err == nil {
			c.StartDate = &d
		}
	}
	if b, err := cmd.Flags().GetString("endDate"); err == nil {
		if d, err := time.Parse(time.RFC3339, b); err == nil {
			c.EndDate = &d
		}
	}

	return nil
}

// Output will run after the execution of the command to write the results to StdOut.
func (c *Github) Output(_ *cobra.Command, _ []string) {
	log.SetOutput(os.Stdout)
	log.Println(c.output)
}

// Run will execute the github pulse summary generation process.
func (c *Github) Run(ctx context.Context) error {
	var err error
	c.output, err = c.pulser.Summary(ctx, service.SummaryRequest{
		Username:     c.Username,
		Organization: c.Organization,
		Repository:   c.Repository,
		StartDate:    c.StartDate,
		EndDate:      c.EndDate,
	})
	return err
}
