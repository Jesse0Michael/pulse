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
