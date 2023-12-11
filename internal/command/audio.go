package command

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jesse0michael/pulse/internal/service"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

type AudioConfig struct {
	AI service.OpenAIConfig
}

type Audio struct {
	openAI *service.OpenAI
	output string
	Files  []string
}

// NewAudio creates a new Audio service.
func NewAudio() *Audio {
	return &Audio{}
}

// Command will return the cobra command structure that can be executed.
func (c *Audio) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "audio [file(s)]",
		Short:   "generate a summary of the audio files",
		PreRunE: c.Init,
		PostRun: c.Output,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.Run(cmd.Context())
		},
	}

	cmd.SetUsageTemplate(cmd.UsageTemplate() + `
Environment:
  OPENAI_URL         the url for accessing the OpenAI API
  OPENAI_TOKEN       the authentication token to use with the OpenAI API
`)

	return cmd
}

// Init will initialize audio dependencies.
func (c *Audio) Init(_ *cobra.Command, args []string) error {
	var cfg AudioConfig
	if err := envconfig.Process("", &cfg); err != nil {
		return err
	}

	c.openAI = service.NewOpenAI(cfg.AI)

	c.Files = args

	return nil
}

// Output will run after the execution of the command to write the results to StdOut.
func (c *Audio) Output(_ *cobra.Command, _ []string) {
	log.SetOutput(os.Stdout)
	log.Println(c.output)
}

// Run will execute the audio transcribe pulse summary generation process.
func (c *Audio) Run(ctx context.Context) error {
	for _, file := range c.Files {
		output, err := c.openAI.Transcribe(ctx, file)
		if err != nil {
			return err
		}

		c.output += fmt.Sprintf("%s\n", output)
	}
	return nil
}
