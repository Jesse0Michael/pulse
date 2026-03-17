.PHONY: test build

# Source a local .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

#################################################################################
# BUILD COMMANDS
#################################################################################
generate:
	go generate ./...

build-cli: 
	go build -o ./bin/pulse ./cmd/pulse

#################################################################################
# RUN COMMANDS
#################################################################################
run-agent:
	go run cmd/agent/main.go web api webui

# Run agent with local Ollama (USE_OLLAMA=true + Ollama must be running)
run-chat:
	USE_OLLAMA=true go run cmd/agent/main.go web api webui

up:
	docker compose up -d

down:
	docker compose down

# Pull a model into Ollama (run once after `make up`). Override: make pull-model MODEL=qwen2.5
pull-model:
	docker exec pulse-ollama-1 ollama pull $(or $(MODEL),llama3.2)

# Register a local GGUF from ./models/ as an Ollama model.
# Usage: make register-model NAME=my-model FILE=my-model.gguf
register-model:
	@echo 'FROM /models/$(FILE)' | docker exec -i pulse-ollama-1 ollama create $(NAME) -f -
	
#################################################################################
# TEST COMMANDS
#################################################################################
test:
	go test -cover ./... 

lint:
	go tool golangci-lint run ./...

cover:
	go test -coverpkg ./internal/... -coverprofile coverage.out ./... && go tool cover -html=coverage.out

vuln: 
	go tool govulncheck -test ./...

behavior: build-cli
	go test ./test/...  -tags=behavior
