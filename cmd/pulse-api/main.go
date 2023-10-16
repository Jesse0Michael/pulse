package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jesse0michael/pulse/internal/server"
	"github.com/jesse0michael/pulse/internal/service"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	initLog()
	var cfg server.Config
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal("failed to process config")
	}

	// Setup context that will cancel on signalled termination
	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		slog.Info("termination signaled")
		cancel()
	}()

	github := service.NewGithub(cfg.Github)
	openAI := service.NewOpenAI(cfg.AI)
	pulser := service.NewPulser(github, openAI)
	srvr := server.New(cfg, pulser)
	go func() { log.Fatal(srvr.ListenAndServe()) }()
	slog.With("port", cfg.Port).Info("started Pulse API")

	// Exit safely
	<-ctx.Done()
	srvr.Close()
	slog.Info("exiting")
}

// initLog configures and sets the default logger.
func initLog() {
	handler := slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
