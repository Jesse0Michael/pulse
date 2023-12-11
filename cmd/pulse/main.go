package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jesse0michael/pulse/internal/command"
	"github.com/spf13/cobra"
)

func main() {
	initLog()

	ctx, cancel := context.WithCancelCause(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-sig
		slog.Info("termination signaled")
		cancel(nil)
	}()

	root := &cobra.Command{
		SilenceErrors: true,
		SilenceUsage:  true,
		Short:         "CLI for pulse: AI Empowered Insights",
		Args:          cobra.MinimumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	root.SetHelpCommand(&cobra.Command{Hidden: true})
	root.AddCommand(
		command.NewGithub().Command(),
		command.NewAudio().Command(),
	)

	err := root.ExecuteContext(ctx)
	if err != nil {
		slog.With("error", err).ErrorContext(ctx, "failed to execute command")
		os.Exit(1)
	}
}

// initLog configures and sets the default logger.
func initLog() {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}
