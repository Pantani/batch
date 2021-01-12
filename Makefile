#! /usr/bin/make -f

# Project variables.
PACKAGE := github.com/Pantani/batch
VERSION := $(shell git describe --tags)
BUILD := $(shell git rev-parse --short HEAD)
DATETIME := $(shell date +"%Y.%m.%d-%H:%M:%S")
PROJECT_NAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOPKG := $(.)

# Go files
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# Use linker flags to provide version/build settings
LDFLAGS=-ldflags "-X=$(PACKAGE)/build.Version=$(VERSION) -X=$(PACKAGE)/build.Build=$(BUILD) -X=$(PACKAGE)/build.Date=$(DATETIME)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECT_NAME)-stderr.txt

# PID file will keep the process id of the server
PID_START := /tmp/.$(PROJECT_NAME).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

## install: Install missing dependencies. Runs `go get` internally. e.g; make install get=github.com/foo/bar
install: go-get

## start: Clean, compile and start simulation.
start:
	@bash -c "$(MAKE) clean compile start-simulation"

## start-simulation: Start alian simulation from binary.
start-simulation: stop
	@echo "  >  Starting $(PROJECT_NAME)"
	@-$(GOBIN)/$(PROJECT_NAME) 2>&1 & echo $$! > $(PID_START)
	@cat $(PID_START) | sed "/^/s/^/  \>  API PID: /"
	@echo "  >  Error log: $(STDERR)"

## stop: Stop the simulation.
stop:
	@-touch $(PID_START)
	@-kill `cat $(PID_START)` 2> /dev/null || true
	@-rm $(PID_START)

## compile: Compile the project.
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## exec: Run given command. e.g; make exec run="go test ./..."
exec:
	GOBIN=$(GOBIN) $(run)

## clean: Clean build files. Runs `go clean` internally.
clean:
	@-rm $(GOBIN)/$(PROJECT_NAME) 2> /dev/null
	@-$(MAKE) go-clean

## check: Run application check.
check: fmt govet golint

## test: Run all tests.
test: unit

## unit: Run all unit tests.
unit: go-unit

## fmt: Run `go fmt` for all go files.
fmt: go-fmt

## govet: Run go vet.
govet: go-vet

## golint: Run golint.
golint: go-lint

## install-swag: Install go-swagger.
install-swag:
ifeq (,$(wildcard test -f $(GOPATH)/bin/swag))
	@echo "  >  Installing swagger"
	@-bash -c "go get github.com/swaggo/swag/cmd/swag"
endif

## swag: Install and run go-swagger.
swag: install-swag
	@bash -c "$(GOPATH)/bin/swag init --parseDependency -g ./cmd/main.go -o ./docs"

go-compile: go-get go-build

go-build:
	@echo "  >  Building simulation binary..."
	GOBIN=$(GOBIN) go build $(LDFLAGS) -o $(GOBIN)/$(PROJECT_NAME) ./cmd

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	GOBIN=$(GOBIN) go get cmd/... $(get)

go-install:
	GOBIN=$(GOBIN) go install $(GOPKG)

go-clean:
	@echo "  >  Cleaning build cache"
	GOBIN=$(GOBIN) go clean

go-unit:
	@echo "  >  Running unit tests"
	GOBIN=$(GOBIN) go test -race -tags=functional -v ./...

go-fmt:
	@echo "  >  Format all go files"
	GOBIN=$(GOBIN) gofmt -w ${GOFMT_FILES}

go-vet:
	@echo "  >  Running go vet"
	GOBIN=$(GOBIN) go vet ./...

go-lint:
	@echo "  >  Running golint"
	GOBIN=$(GOBIN) golint ./...

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECT_NAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo