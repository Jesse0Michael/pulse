package config

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/jesse0michael/pkg/logger"
)

const (
	// HomeDir is the pulse config directory under the user's home.
	HomeDir = ".pulse"
	// HomeSettingsFile is the settings file name inside HomeDir.
	HomeSettingsFile = "settings.json"
	// LocalSettingsFile is the project-local settings file.
	LocalSettingsFile = "pulse.json"
)

//go:embed default_settings.json
var defaultSettingsTmpl string

// Settings holds the pulse configuration.
//
// Sources are layered in this order (later wins):
//  1. struct-tag defaults
//  2. environment variables
//  3. settings files (~/.pulse/settings.json then ./pulse.json)
//  4. CLI flags (via pkg/config + go-arg)
type Settings struct {
	Logger logger.Config `json:"logger,omitzero" arg:"-"`

	Model     string      `envconfig:"PULSE_MODEL"      default:"qwen3.6"               json:"model"      arg:"--model"`
	OllamaURL string      `envconfig:"PULSE_OLLAMA_URL" default:"http://localhost:11434" json:"ollama_url" arg:"--ollama-url"`
	Agent     string      `envconfig:"PULSE_AGENT"      default:"pulse"                 json:"agent"      arg:"--agent,-a"`
	Tools     ToolsConfig `json:"tools" arg:"-"`

	Agents []AgentConfig `json:"agents,omitempty" arg:"-"`
}

// ToolsConfig controls which agent tools are available.
type ToolsConfig struct {
	Shell    bool `envconfig:"PULSE_TOOL_SHELL"     default:"true" json:"shell"`
	ReadFile bool `envconfig:"PULSE_TOOL_READ_FILE" default:"true" json:"read_file"`
	ListDir  bool `envconfig:"PULSE_TOOL_LIST_DIR"  default:"true" json:"list_dir"`
}

// SettingsFile returns the home settings file path: ~/.pulse/settings.json.
func SettingsFile() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, HomeDir, HomeSettingsFile)
}

// EnsureHomeDir creates ~/.pulse/ and ~/.pulse/logs/ with 0o700 if missing.
// On first run it also writes ~/.pulse/settings.json from the embedded default template.
func EnsureHomeDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	for _, sub := range []string{"", "logs"} {
		dir := filepath.Join(home, HomeDir, sub)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return fmt.Errorf("creating %s: %w", dir, err)
		}
	}

	settingsPath := filepath.Join(home, HomeDir, HomeSettingsFile)
	if _, err := os.Stat(settingsPath); err == nil {
		return nil
	}

	tmpl, err := template.New("settings").Parse(defaultSettingsTmpl)
	if err != nil {
		return fmt.Errorf("parsing default settings template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{"HomeDir": home}); err != nil {
		return fmt.Errorf("executing default settings template: %w", err)
	}

	if err := os.WriteFile(settingsPath, buf.Bytes(), 0o600); err != nil {
		return fmt.Errorf("writing default settings to %s: %w", settingsPath, err)
	}
	return nil
}
