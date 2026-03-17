package main

import (
	"context"
	"log"
	"os"

	"github.com/jesse0michael/pkg/config"
	"github.com/jesse0michael/pulse/internal/model/ollama"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model"
	adkgemini "google.golang.org/adk/model/gemini"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/genai"
)

type Config struct {
	// Gemini backend (cloud)
	GeminiModel  string `envconfig:"GEMINI_MODEL"   default:"gemini-2.5-flash"`
	GoogleAPIKey string `envconfig:"GOOGLE_API_KEY"`

	// Ollama backend (local)
	OllamaURL   string `envconfig:"OLLAMA_URL"   default:"http://localhost:11434"`
	OllamaModel string `envconfig:"OLLAMA_MODEL" default:"llama3.2"`

	// Set USE_OLLAMA=true to use local Ollama instead of Gemini
	UseOllama bool `envconfig:"USE_OLLAMA" default:"false"`
}

func main() {
	ctx := context.Background()

	cfg, err := config.New[Config]()
	if err != nil {
		log.Fatalf("failed to process config: %v", err)
	}

	var m model.LLM
	if cfg.UseOllama {
		log.Printf("using Ollama model %q at %s", cfg.OllamaModel, cfg.OllamaURL)
		m = ollama.NewModel(cfg.OllamaURL, cfg.OllamaModel)
	} else {
		log.Printf("using Gemini model %q", cfg.GeminiModel)
		m, err = adkgemini.NewModel(ctx, cfg.GeminiModel, &genai.ClientConfig{
			APIKey: cfg.GoogleAPIKey,
		})
		if err != nil {
			log.Fatalf("failed to create Gemini model: %v", err)
		}
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        "chat",
		Model:       m,
		Description: "A helpful AI assistant.",
		Instruction: "You are a helpful assistant. Answer questions clearly and concisely.",
	})
	if err != nil {
		log.Fatalf("failed to create agent: %v", err)
	}

	adkCfg := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	l := full.NewLauncher()
	if err = l.Execute(ctx, adkCfg, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
