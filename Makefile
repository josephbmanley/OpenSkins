#include .env

PROJECTNAME="OpenSkins"
STANDALONE_LOCATION="../OpenSkins-Standalone-Plugin/bin/OpenSkins-Standalone-Plugin"

# Go related variables.
GOBASE=$(shell pwd)
GOPATH=$(GOBASE)/vendor:$(GOBASE)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Redirect error output to a file, so we can show it in development mode.
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will store the server process id when it's running on development mode
PID=/tmp/.$(PROJECTNAME)-api-server.pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

go-compile: go-clean go-get go-build

go-build:
	@echo "  >  Building binary..."
	@go build -o $(GOBIN)/$(PROJECTNAME) $(GOFILES)

go-generate:
	@echo "  >  Generating dependency files..."
	@go generate $(generate)

go-get:
	@echo "  >  Checking if there is any missing dependencies..."
	@go get $(get)

go-install:
	@echo "  >  Running go install..."
	@go install $(GOFILES)

go-clean:
	@echo "  >  Cleaning build cache"
	@go clean

go-test:
	@echo "  >  Running tests..."
	@go test

go-run:
	@echo "  >  Running ${PROJECTNAME}"
	@-(cd $(GOBIN); ./$(PROJECTNAME))

openskins-common:
	@echo "  >  Updating common library..."
	@go get -u github.com/josephbmanley/OpenSkins-Common


## install: downloads and installs dependencies
install: openskins-common go-get

## clean: Runs go clean
clean:
	@(MAKEFILE) go-clean

## compile: cleans project, installs dependencies, and builds project
compile:
	@-touch $(STDERR)
	@-rm $(STDERR)
	@-$(MAKE) -s go-compile 2> $(STDERR)
	@cat $(STDERR) | sed -e '1s/.*/\nError:\n/'  | sed 's/make\[.*/ /' | sed "/^/s/^/     /" 1>&2

## watch: Runs command on code update
watch:
	@yolo -i . -e vendor -e bin -c $(run)

## build: Runs go build
build: go-build

## run: Compiles and executes project binary
run: go-compile go-run

## test: Run unit tests
test: go-test

## standalone: Compile project & install standalone plugin *REQUIRES BUILD IN `../OpenSkins-Standalone-Plugin/bin`*
standalone: compile
	@mkdir -p ./bin/plugins
	@rm -f ./bin/plugins/standalone.so
	@cp $(STANDALONE_LOCATION) ./bin/plugins/standalone.so

## help: Displays help text for make commands
.DEFAULT_GOAL := help
all: help
help: Makefile
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'