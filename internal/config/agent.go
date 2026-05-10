package config

import "fmt"

// AgentConfig defines a named agent with its own instruction and optional overrides.
type AgentConfig struct {
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	Instruction string       `json:"instruction"`
	Model       string       `json:"model,omitempty"`
	Tools       *ToolsConfig `json:"tools,omitempty"`
}

// GetAgent returns the agent config matching the given name,
// filling in defaults from the top-level Settings.
func (s Settings) GetAgent(name string) (AgentConfig, error) {
	for _, ac := range s.Agents {
		if ac.Name == name {
			if ac.Model == "" {
				ac.Model = s.Model
			}
			if ac.Tools == nil {
				ac.Tools = &s.Tools
			}
			return ac, nil
		}
	}
	return AgentConfig{}, fmt.Errorf("agent %q not found in settings", name)
}

// AgentNames returns the list of configured agent names.
func (s Settings) AgentNames() []string {
	names := make([]string, 0, len(s.Agents))
	for _, ac := range s.Agents {
		names = append(names, ac.Name)
	}
	return names
}
