SHELL := /bin/bash

PROJECT := "wb-l0"

BIN := "./bin/api"
SRC := "./cmd/api"

.PHONY: setup
setup: ## Install all the build and lint dependencies
	bash scripts/setup.sh

.PHONY: build
build: ## Build project
	bash scripts/build.sh $(BIN) $(SRC)

.PHONY: run
run: build ## Run project in local environment
	$(BIN)

.PHONY: up
up: ## Run project in docker environment
	bash scripts/up.sh $(PROJECT)

.PHONY: down
down: ## Stop project in docker environment
	docker-compose -f deployments/docker-compose.yml -p $(PROJECT) --env-file .env down

.PHONY: logs
logs: ## View project logs from the docker container
	docker-compose -f deployments/docker-compose.yml -p $(PROJECT) logs

.PHONY: version
version: build ## Build and view project version
	$(BIN) version

.PHONY: fmt
fmt: ## Run format tools on all go files
	bash scripts/fmt.sh

.PHONY: lint
lint: ## Run all the linters
	golangci-lint run ./...

.PHONY: test.unit
test.unit: ## Run all unit tests
	@echo 'mode: atomic' > coverage.txt
	go test -v -race ./cmd/...
	go test -covermode=atomic -coverprofile=coverage.txt -v -race ./internal/...

.PHONY: cover
cover: test.unit ## Run all the tests and opens the coverage report
	go tool cover -html=coverage.txt

.PHONY: generate
gen: ## Generate files
	go generate ./...

.PHONY: ci
ci: lint test.unit ## Run all the tests and code checks

.PHONY: clean
clean: ## Remove temporary files
	@go clean
	@rm -rf bin/
	@rm -rf coverage.txt
	@echo "SUCCESS!"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL:= help