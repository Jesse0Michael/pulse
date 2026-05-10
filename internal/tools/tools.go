package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

// ShellArgs is the input for the shell tool.
type ShellArgs struct {
	Command string `json:"command"`
}

// ShellResult is the output of the shell tool.
type ShellResult struct {
	Output   string `json:"output"`
	ExitCode int    `json:"exit_code"`
}

// NewShellTool creates a tool that executes shell commands.
func NewShellTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "shell",
		Description: "Execute a shell command in the current working directory and return its output.",
	}, func(_ tool.Context, args ShellArgs) ShellResult {
		cmd := exec.Command("sh", "-c", args.Command)
		out, err := cmd.CombinedOutput()
		exitCode := 0
		if err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				exitCode = exitErr.ExitCode()
			} else {
				return ShellResult{Output: fmt.Sprintf("error: %v", err), ExitCode: 1}
			}
		}
		return ShellResult{Output: string(out), ExitCode: exitCode}
	})
}

// ReadFileArgs is the input for the read_file tool.
type ReadFileArgs struct {
	Path string `json:"path"`
}

// ReadFileResult is the output of the read_file tool.
type ReadFileResult struct {
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

// NewReadFileTool creates a tool that reads file contents.
func NewReadFileTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "read_file",
		Description: "Read the contents of a file at the given path.",
	}, func(_ tool.Context, args ReadFileArgs) ReadFileResult {
		path := args.Path
		if !filepath.IsAbs(path) {
			cwd, _ := os.Getwd()
			path = filepath.Join(cwd, path)
		}
		path = filepath.Clean(path)
		data, err := os.ReadFile(path)
		if err != nil {
			return ReadFileResult{Error: err.Error()}
		}
		return ReadFileResult{Content: string(data)}
	})
}

// ListDirArgs is the input for the list_dir tool.
type ListDirArgs struct {
	Path string `json:"path"`
}

// ListDirResult is the output of the list_dir tool.
type ListDirResult struct {
	Entries []string `json:"entries"`
	Error   string   `json:"error,omitempty"`
}

// NewListDirTool creates a tool that lists directory contents.
func NewListDirTool() (tool.Tool, error) {
	return functiontool.New(functiontool.Config{
		Name:        "list_dir",
		Description: "List the contents of a directory. Returns file and directory names (directories end with /).",
	}, func(_ tool.Context, args ListDirArgs) ListDirResult {
		path := args.Path
		if path == "" {
			path = "."
		}
		if !filepath.IsAbs(path) {
			cwd, _ := os.Getwd()
			path = filepath.Join(cwd, path)
		}
		path = filepath.Clean(path)

		entries, err := os.ReadDir(path)
		if err != nil {
			return ListDirResult{Error: err.Error()}
		}
		var names []string
		for _, e := range entries {
			name := e.Name()
			if strings.HasPrefix(name, ".") {
				continue
			}
			if e.IsDir() {
				name += "/"
			}
			names = append(names, name)
		}
		return ListDirResult{Entries: names}
	})
}
