# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Vision

This is a **Go-based AI playground** built around [Google ADK (Agent Development Kit)](https://google.github.io/adk-docs/). The goal is a personal AI system with:
- **Agents** triggered to perform specific tasks (GitHub activity, summarization, etc.)
- **Local LLM chat** via Ollama with a web UI or terminal UI
- **TTS** in the web UI (future)

The existing `internal/` code (GitHub/OpenAI services, REST API, CLI) is legacy and will likely be replaced by ADK-based agents.

## Commands

```bash
make run-agent      # Run the ADK agent (web UI + API + web UI combined)
make up / make down # Start/stop Ollama via Docker Compose
make test           # Run unit tests
make lint           # Run golangci-lint
make vuln           # Run govulncheck
make generate       # Run go generate (for mocks)
```

Run individual tests:
```bash
go test ./internal/service/... -v
go test ./test/... -tags=behavior   # BDD/Gherkin behavior tests
```

## ADK Agent Architecture

The primary entry point is [cmd/agent/main.go](cmd/agent/main.go). It uses the Google ADK Go SDK (`google.golang.org/adk`).

**How the launcher works:**

```go
l := full.NewLauncher()
l.Execute(ctx, config, os.Args[1:])
// args: "web" = web UI, "api" = REST API, "webui" = browser UI
```

The `make run-agent` command passes `web api webui` to run all three modes simultaneously.

**Current agent:** `status-agent` — uses Gemini + GitHub Copilot MCP tools to summarize team activity.

**ADK core concepts:**
- **`llmagent.New`** — creates an LLM-backed agent with a model, instructions, and toolsets
- **`tool.Toolset`** — groups of tools the agent can call (e.g., MCP, custom functions)
- **`mcptoolset`** — connects to any MCP server (Model Context Protocol) as a toolset
- **`services.NewSingleAgentLoader`** — registers one agent with the ADK runtime
- **`full.NewLauncher`** — launches web UI + REST API together

## Configuration

Env vars are loaded via `envconfig` from a `.env` file (not committed). Key vars:

```
GOOGLE_API_KEY=     # For Gemini model
GEMINI_MODEL=       # Default: gemini-2.5-flash
GITHUB_PAT=         # For GitHub Copilot MCP tools
GITHUB_TOKEN=       # For GitHub API (legacy)
OPENAI_URL=         # Can point to Ollama for local LLMs
OPENAI_TOKEN=       # API token
```

## Local LLM (Ollama)

`docker-compose.yml` runs an Ollama service. ADK supports connecting to local models — the pattern is to configure the model backend to point at the local Ollama endpoint instead of Gemini/OpenAI.

## Adding New Agents / Tools

1. Create a new `llmagent.New(...)` with a name, model, instructions, and toolsets
2. Register custom Go functions as tools using `tool.NewFunctionTool` or similar ADK tool constructors
3. Use `mcptoolset.New(...)` to connect to any MCP server as a toolset
4. Register agents via `services.NewSingleAgentLoader` (single) or a multi-agent loader

## Testing

- Unit tests: standard `*_test.go` alongside source files
- Behavior tests: Gherkin feature files + godog in `test/godog/`, use `go-rest-assured` for HTTP mocking
- Mocks are generated with `mockgen` via `go generate`
