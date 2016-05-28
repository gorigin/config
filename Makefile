# Makefile configuration
.DEFAULT_GOAL := help
.PHONY: help fmt vet test deps cyclo

ok: fmt vet cyclo test ## Prepares codebase (fmt+vet+test)

fmt: ## Golang code formatting tool
	@echo "Running formatting tool"
	@gofmt -s -w .

vet: ## Check code against common errors
	@echo "Running code inspection tools"
	@go vet ./...

cyclo: ## Check cyclomatic complexity
	@echo "Running cyclomatic complexity test"
	@${GOPATH}/bin/gocyclo -over 15 .

test: ## Run tests
	@echo "Running unit tests"
	@go test ./...

deps: ## Download required dependencies
	go get gopkg.in/yaml.v2
	go get github.com/stretchr/testify/assert

help:
	@grep --extended-regexp '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-10s\033[0m %s\n", $$1, $$2}'
