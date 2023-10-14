package command

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/jesse0michael/pulse/internal/service"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type Config struct {
	Pool int `env:"POOL" default:"1"`
	AI   service.OpenAIConfig
}

type Github struct {
	cfg      Config
	client   *service.Github
	ai       *service.OpenAI
	output   string
	Username string
}

// NewGithub creates a new Github service
func NewGithub() *Github {
	return &Github{
		cfg:    Config{},
		client: service.NewGithub(service.GithubConfig{}),
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

	fmt.Println(c.cfg)
	c.ai = service.NewOpenAI(c.cfg.AI)

	if len(args) > 0 {
		c.Username = args[0]
	}

	return nil
}

// Output will run after the execution of the command to write the results to StdOut
func (c *Github) Output(cmd *cobra.Command, _ []string) {
	l := slog.New(slog.NewTextHandler(os.Stdout, nil))
	l.InfoContext(cmd.Context(), c.output)
}

// Run will execute the github pulse summary generation process
func (c *Github) Run(ctx context.Context) error {
	content, err := c.client.UserActivity(ctx, c.Username)
	if err != nil {
		return err
	}

	c.output, err = c.ai.Summarize(ctx, content)

	return err
}
