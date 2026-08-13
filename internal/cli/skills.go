package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Skill is a markdown-based instruction set that the user can trigger via
// the in-console `/skill <name>` command.
type Skill struct {
	Name string // filename without the .md extension
	Path string // location on disk
}

// LoadSkills discovers skills from ~/.pulse/skills/ and ./skills/.
// When the same name appears in both, the project copy (./skills/) wins.
func LoadSkills() []Skill {
	seen := map[string]Skill{}

	if home, err := os.UserHomeDir(); err == nil {
		readSkillDir(seen, filepath.Join(home, ".pulse", "skills"))
	}
	readSkillDir(seen, "skills")

	names := make([]string, 0, len(seen))
	for n := range seen {
		names = append(names, n)
	}
	sort.Strings(names)
	out := make([]Skill, 0, len(names))
	for _, n := range names {
		out = append(out, seen[n])
	}
	return out
}

func readSkillDir(seen map[string]Skill, dir string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		name := strings.TrimSuffix(e.Name(), ".md")
		seen[name] = Skill{
			Name: name,
			Path: filepath.Join(dir, e.Name()),
		}
	}
}

// Content reads the skill body from disk.
func (s Skill) Content() (string, error) {
	b, err := os.ReadFile(s.Path)
	if err != nil {
		return "", fmt.Errorf("reading skill %s: %w", s.Path, err)
	}
	return string(b), nil
}
