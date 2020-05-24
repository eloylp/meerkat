PROJECT_NAME := $(shell basename "$(PWD)")
BINARY_NAME := meerkat
VERSION := $(shell git describe --tags)
GO_VERSION := 1.14.2
GO_LINT_CI_VERSION := v1.26.0
GO_LINT_CI_PATH := $(shell go env GOPATH)/bin
TIME := $(shell date +%Y-%m-%dT%T%z)
BUILD := $(shell git rev-parse --short HEAD)
DIST_FOLDER := ./dist
BINARY_OUTPUT := $(DIST_FOLDER)/$(BINARY_NAME)
LDFLAGS=-ldflags "-s -w \
		-X=main.Name=$(PROJECTNAME) \
		-X=main.Version=$(VERSION) \
		-X=main.Build=$(BUILD) \
		-X=main.BuildTime=$(TIME)"
FLAGS=-trimpath

.DEFAULT_GOAL := build

lint:
	golangci-lint run -v
lint-fix:
	golangci-lint run -v --fix
linter-install:
	wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GO_LINT_CI_PATH) $(GO_LINT_CI_VERSION)
all: lint test build

test: test-unit test-integration test-race test-bench

test-unit:
	go test -race -v --tags="unit" ./...
test-integration:
	go test -race -v --tags="integration" ./...
test-race:
	go test -race -v --tags="race" ./...
test-bench:
	go test -v -bench=. ./...
build:
	mkdir -p $(DIST_FOLDER)
	CGO_ENABLED=0 go build $(FLAGS) $(LDFLAGS) -o $(BINARY_OUTPUT)
	@echo "Binary output at $(BINARY_OUTPUT)"
build-docker:
	docker run -e "CGO_ENABLED=0" --rm -v ${CURDIR}:/usr/src/code -w /usr/src/code golang:$(GO_VERSION) make build
install:
	go install
clean:
	rm -rf $(DIST_FOLDER)