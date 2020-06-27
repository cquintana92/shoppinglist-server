ARTIFACT_NAME:=shopping-list

SHELL:=/bin/sh
MAKEFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
PROJECT_ROOT := $(dir $(MAKEFILE_PATH))
BIN_DIRECTORY := ${PROJECT_ROOT}/bin

GO:=go

DISABLE_CACHE:=-count=1
VERBOSE:=-v
TEST_COMMAND:=${GO} test ./...
TEST_NO_CACHE_COMMAND:=${TEST_COMMAND} ${DISABLE_CACHE}
BUILD_COMMAND:=${GO} build
SRC_DIR:=${PROJECT_ROOT}

VERSION=`git describe --tags --always --dirty`
VERSION_VARIABLE_NAME=shoppinglistserver/constants.VERSION
VERSION_VARIABLE_BUILD_FLAG=-ldflags "-X ${VERSION_VARIABLE_NAME}=${VERSION}"
BUILD_WITH_VERSION_COMMAND=${BUILD_COMMAND} ${VERSION_VARIABLE_BUILD_FLAG}

.PHONY: default
default: help

.PHONY: deps
deps: ## Setup dependencies
	@ ${GO} get  ./...

.PHONY: build
build: deps ## Build
	@ ${BUILD_WITH_VERSION_COMMAND} -o ${BIN_DIRECTORY}/${ARTIFACT_NAME} cmd/server/main.go

.PHONY: build_all
build_all: build client ## Build all the possible binaries

.PHONY: fmt
fmt: ## Apply linting and formatting
	@ ${GO} fmt ./...

.PHONY: test
test: ## Run tests
	@ ${GO} test ./...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

