.PHONY: test build

# Source a local .env
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

#################################################################################
# BUILD COMMANDS
#################################################################################
dependencies:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	go install github.com/golang/mock/mockgen@v1.6.0

gen: dependencies
	go generate ./...

build-cli: 
	go build -o ./bin/pulse ./cmd/pulse
	
#################################################################################
# TEST COMMANDS
#################################################################################
test:
	go test -cover ./... 

lint:
	golangci-lint run ./...

cover:
	go test -coverpkg ./internal/... -coverprofile coverage.out ./... && go tool cover -html=coverage.out

vuln: dependencies
	govulncheck -test ./...
