package main

import (
	"context"
	"log"
	"os"

	"github.com/jesse0michael/pkg/config"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"golang.org/x/oauth2"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/mcptoolset"
	"google.golang.org/genai"
)

type GeminiConfig struct {
	Model        string `envconfig:"GEMINI_MODEL" default:"gemini-2.5-flash"`
	GoogleAPIKey string `envconfig:"GOOGLE_API_KEY"`
}

type Config struct {
	GeminiConfig
}

func main() {
	ctx := context.Background()

	cfg, err := config.New[Config]()
	if err != nil {
		log.Fatalf("failed to process config: %v", err)
	}

	m, err := gemini.NewModel(ctx, cfg.Model, &genai.ClientConfig{
		APIKey: cfg.GoogleAPIKey,
	})
	if err != nil {
		log.Fatalf("failed to create Gemini model: %v", err)
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_PAT")},
	)
	mcpToolSet, err := mcptoolset.New(mcptoolset.Config{
		Transport: &mcp.StreamableClientTransport{
			Endpoint:   "https://api.githubcopilot.com/mcp/",
			HTTPClient: oauth2.NewClient(ctx, ts),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create MCP tool set: %v", err)
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "status-agent",
		Model:       m,
		Description: "Get the status update for the team update.",
		Instruction: "You are a helpful assistant that will lookup what a team member has done in a given time frame and summarize an update.",
		Toolsets: []tool.Toolset{
			mcpToolSet,
		},
	})

	config := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
