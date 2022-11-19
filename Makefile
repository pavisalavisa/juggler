SHELL := $(shell which bash)
OSTYPE := $(shell uname)
DOCKER := $(shell command -v docker)
VERSION ?= $(shell git describe --tags --always)
args = $(foreach a,$($(subst -,_,$1)_args),$(if $(value $a),$a="$($a)"))

help: ## Show this help
	@echo "Help"
	@grep -E '^[0-9a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[93m %s\n", $$1, $$2}'

.PHONY: default
default: help

.PHONY: docker-dev
docker-dev: ## Builds the development docker image.
	IMAGE=juggler-dev VERSION=${VERSION} DOCKER_FILE=Dockerfile ./scripts/build-image.sh

.PHONY: build
build: ## Builds the production binary for linux amd64
	./scripts/build.sh juggler

.PHONY: m1
m1: ## Builds the binary for local development in macos m1
	./scripts/build.sh juggler macos

.PHONY: test
test: ## Launch tests with core emulator. Used in juggler CI
	./scripts/test.sh

.PHONY: clean
clean: ## Cleans binary output dir
	@rm -rf ./bin

.PHONY: lint
lint: ## Invoke go linter
	golangci-lint run

.PHONY: run-m1
run-m1:  ## Builds the binary for local development in macos m1 and runs it
	make m1
	./bin/juggler-darwin-arm64

.PHONY: configure-git-hooks
configure-git-hooks: ## Configure git hooks
	./scripts/configure-git-hooks.sh
