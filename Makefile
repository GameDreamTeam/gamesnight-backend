.DEFAULT_GOAL := help

# Load .env file in the application, not in the Makefile
# include .env
# export

# Define 'run' as depending on 'format' and 'build'
run: format build
	@echo "Formatting, Building and Running"
	./bin/charades

format:
	@echo "Formatting"
	go fmt ./...

build:
	@echo "Building"
	go build -o ./bin/charades ./cmd/charades

help:
	@echo "Available commands:"
	@echo "  make format   - Formatting Repository"
	@echo "  make build    - Building Repository"
	@echo "  make run      - Starting Server"
