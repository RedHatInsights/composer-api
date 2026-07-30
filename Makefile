SHELL=/usr/bin/env bash

APPLICATION_NAME        = composer-api
APPLICATION_BINARY_NAME = composer-api

# Support both podman and docker.
DOCKER=$(shell which podman || which docker || echo 'docker')

.PHONY: build run test tidy lint image container clean

# Builds the project.
build:
	@echo "+$@"
	@mkdir -p bin
	@go build -o bin/$(APPLICATION_BINARY_NAME) ./cmd/$(APPLICATION_NAME)

# Runs the project after tidying and building it anew.
run: tidy build
	@echo "+$@"
	@echo "########### Running the application binary ############"
	@bin/$(APPLICATION_BINARY_NAME)

# Tests the whole project.
test:
	@echo "+$@"
	@go test -race -coverprofile=coverage.out -covermode=atomic ./...

# Runs "go mod tidy".
tidy:
	@echo "+$@"
	@go mod tidy

# Runs golangci-lint over the project.
lint:
	@echo "+$@"
	@golangci-lint run ./...

# Builds the container image.
image:
	@echo "+$@"
	@$(DOCKER) build --file Dockerfile --tag $(APPLICATION_NAME):latest .

# Runs the project container assuming the image is already built.
container:
	@echo "+$@"
	@echo "############### Removing old container ################"
	@$(DOCKER) rm -f $(APPLICATION_NAME)
	@echo "################ Running new container ################"
	@$(DOCKER) run --name $(APPLICATION_NAME) --detach --publish 8080:8080 \
		$(APPLICATION_NAME):latest

# Removes build artifacts.
clean:
	@echo "+$@"
	@rm -rf bin/ coverage.out
