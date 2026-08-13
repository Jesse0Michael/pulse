package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// AgentConfig defines a named agent with its own instruction and optional overrides.
//
// Instruction may be set inline or loaded from PromptFile (relative to cwd or
// absolute; a leading `~/` expands to the user's home directory). If both are
// set, Instruction wins.
type AgentConfig struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Instruction string       `json:"instruction,omitempty"`
	PromptFile  string       `json:"prompt_file,omitempty"`
	Model       string       `json:"model,omitempty"`
	Tools       *ToolsConfig `json:"tools,omitempty"`
}

// GetAgent returns the agent config matching the given name,
// filling in defaults from the top-level Settings.
func (s Settings) GetAgent(name string) (AgentConfig, error) {
	for _, ac := range s.Agents {
		if ac.Name != name {
			continue
		}
		if ac.Model == "" {
			ac.Model = s.Model
		}
		if ac.Tools == nil {
			ac.Tools = &s.Tools
		}
		if ac.Instruction == "" && ac.PromptFile != "" {
			body, err := loadPromptFile(ac.PromptFile)
			if err != nil {
				return AgentConfig{}, fmt.Errorf("agent %q: %w", name, err)
			}
			ac.Instruction = body
		}
		return ac, nil
	}
	return AgentConfig{}, fmt.Errorf("agent %q not found in settings", name)
}

func loadPromptFile(p string) (string, error) {
	if strings.HasPrefix(p, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			p = filepath.Join(home, p[2:])
		}
	}
	data, err := os.ReadFile(p)
	if err != nil {
		return "", fmt.Errorf("reading prompt file %s: %w", p, err)
	}
	return string(data), nil
}

// AgentNames returns the list of configured agent names.
func (s Settings) AgentNames() []string {
	names := make([]string, 0, len(s.Agents))
	for _, ac := range s.Agents {
		names = append(names, ac.Name)
	}
	return names
}
