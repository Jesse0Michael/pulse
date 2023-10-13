package command

import (
	"context"
	"log/slog"
	"os"

	"github.com/jesse0michael/pulse/internal/collector/github"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type Config struct {
	Pool int `env:"POOL" default:"1"`
}

type Github struct {
	cfg      Config
	client   *github.Client
	Username string
}

// NewGithub creates a new Github service
func NewGithub() *Github {
	return &Github{
		client: github.NewClient(),
	}
}

// Command will return the cobra command structure that can be executed
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

	return cmd
}

// Init will initialize export dependencies
func (c *Github) Init(cmd *cobra.Command, args []string) error {
	if err := envconfig.Process("", &c.cfg); err != nil {
		return err
	}

	if len(args) > 0 {
		c.Username = args[0]
	}

	return nil
}

// Output will run after the execution of the command to write the results to StdOut
func (c *Github) Output(cmd *cobra.Command, _ []string) {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	l.InfoContext(cmd.Context(), "complete")
}

// Run will execute the github pulse summary generation process
func (c *Github) Run(ctx context.Context) error {
	summary, err := c.client.UserActivity(c.Username)
	if err != nil {
		return err
	}

	return nil
}
