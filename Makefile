.DEFAULT_GOAL := help

run: format build
	@echo "Formatting, Building and Running"
	ENV=local ./bin/charades

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
