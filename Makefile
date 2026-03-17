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

up: 
	docker compose up

down:
	docker compose down
	
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
