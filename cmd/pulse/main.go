package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/jesse0michael/pkg/boot"
	"github.com/jesse0michael/pulse/internal/cli"
	"github.com/jesse0michael/pulse/internal/config"
	"github.com/jesse0michael/pulse/internal/model/ollama"
	"github.com/jesse0michael/pulse/internal/tools"
	"github.com/spf13/cobra"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher/adk"
	"google.golang.org/adk/cmd/launcher/universal"
	"google.golang.org/adk/cmd/launcher/web"
	"google.golang.org/adk/cmd/launcher/web/a2a"
	"google.golang.org/adk/cmd/launcher/web/api"
	"google.golang.org/adk/cmd/launcher/web/webui"
	"google.golang.org/adk/server/restapi/services"
	"google.golang.org/adk/tool"
)

func main() {
	// Pre-boot: ensure ~/.pulse/, logs dir, and default settings exist.
	if err := config.EnsureHomeDir(); err != nil {
		slog.Error("pulse", "err", err)
		os.Exit(1)
	}

	// boot.NewApp handles: context + signal handling, config loading
	// (defaults < env < settings files), and logger setup.
	app := boot.NewApp[config.Settings](
		boot.WithConfigFile(config.SettingsFile()),
		boot.WithConfigFile(config.LocalSettingsFile),
	)

	if err := app.Run(&pulse{}); err != nil {
		fmt.Fprintf(os.Stderr, "pulse: %v\n", err)
		os.Exit(1)
	}
}

// pulse implements boot.Runner to bridge boot.App with cobra.
type pulse struct{}

func (r *pulse) Run(ctx context.Context, settings config.Settings) error {
	root := &cobra.Command{
		Use:                "pulse [mode...]",
		Short:              "AI Empowered Insights",
		Args:               cobra.ArbitraryArgs,
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAgent(cmd.Context(), settings, args)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	return root.ExecuteContext(ctx)
}

func (r *pulse) Close() error { return nil }

func runAgent(ctx context.Context, settings config.Settings, args []string) error {
	// Filter out flag arguments (already parsed by boot/config); keep only
	// positional args used as launcher modes (e.g. "console", "web").
	var modes []string
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			// Skip the flag and its value (e.g. "-a claude-code" or "--agent=foo").
			if !strings.Contains(args[i], "=") && i+1 < len(args) {
				i++ // skip the next arg (flag value)
			}
			continue
		}
		modes = append(modes, args[i])
	}

	agentCfg, err := settings.GetAgent(settings.Agent)
	if err != nil {
		return fmt.Errorf("%w\navailable agents: %v", err, settings.AgentNames())
	}

	m := ollama.NewModel(settings.OllamaURL, agentCfg.Model)
	slog.InfoContext(ctx, "pulse", "agent", agentCfg.Name, "model", agentCfg.Model, "ollama", settings.OllamaURL)

	agentTools, err := buildTools(agentCfg.Tools)
	if err != nil {
		return err
	}

	agent, err := llmagent.New(llmagent.Config{
		Name:        agentCfg.Name,
		Model:       m,
		Description: agentCfg.Description,
		Instruction: agentCfg.Instruction,
		Tools:       agentTools,
	})
	if err != nil {
		return fmt.Errorf("creating agent: %w", err)
	}

	launcherArgs := modes
	if len(launcherArgs) == 0 {
		launcherArgs = []string{"console"}
	}

	adkCfg := &adk.Config{
		AgentLoader: services.NewSingleAgentLoader(agent),
	}
	console := cli.New()
	console.AgentName = agentCfg.Name
	console.ModelName = agentCfg.Model
	l := universal.NewLauncher(
		console,
		web.NewLauncher(api.NewLauncher(), a2a.NewLauncher(), webui.NewLauncher()),
	)
	return l.Execute(ctx, adkCfg, launcherArgs)
}

func buildTools(tc *config.ToolsConfig) ([]tool.Tool, error) {
	if tc == nil {
		return nil, nil
	}
	var ts []tool.Tool
	if tc.Shell {
		t, err := tools.NewShellTool()
		if err != nil {
			return nil, fmt.Errorf("creating shell tool: %w", err)
		}
		ts = append(ts, t)
	}
	if tc.ReadFile {
		t, err := tools.NewReadFileTool()
		if err != nil {
			return nil, fmt.Errorf("creating read_file tool: %w", err)
		}
		ts = append(ts, t)
	}
	if tc.ListDir {
		t, err := tools.NewListDirTool()
		if err != nil {
			return nil, fmt.Errorf("creating list_dir tool: %w", err)
		}
		ts = append(ts, t)
	}
	return ts, nil
}
