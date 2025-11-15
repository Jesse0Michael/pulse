package main

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/jesse0michael/pkg/boot"
	"github.com/jesse0michael/pkg/config"
	"github.com/jesse0michael/pkg/http/server"
	"github.com/jesse0michael/pulse/internal/handler"
	"github.com/jesse0michael/pulse/internal/service"
)

type Config struct {
	config.AppConfig
	config.OpenTelemetryConfig
	Server server.Config
	Github service.GithubConfig
	AI     service.OpenAIConfig
}

type Server struct {
	server *server.Server
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run(ctx context.Context, cfg Config) error {
	github := service.NewGithub(cfg.Github)
	openAI := service.NewOpenAI(cfg.AI)
	pulser := service.NewPulser(github, openAI)
	routes := handler.New(pulser)
	s.server = server.New(cfg.Server, routes)
	go func() { log.Fatal(s.server.ListenAndServe()) }()
	slog.InfoContext(ctx, "started Pulse API", "port", cfg.Server.Port)
	return nil
}

func (s *Server) Close() error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(context.Background())
}

func main() {
	app := boot.NewApp[Config]()
	server := NewServer()
	if err := app.Run(server); err != nil {
		os.Exit(1)
	}
}
