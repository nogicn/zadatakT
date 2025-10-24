# Simple Makefile for a Go project

# Build the application
all: build test

build:
	@echo "Building..."
	
	
	@CGO_ENABLED=1 GOOS=linux go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go
# Create DB container
docker-run:
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "Using Docker Compose V1"; \
		docker-compose up --build; \
	elif docker compose version >/dev/null 2>&1; then \
		echo "Using Docker Compose V2"; \
		docker compose up --build; \
	else \
		echo "Error: Neither Docker Compose V1 nor V2 found. Please install Docker Compose."; \
		exit 1; \
	fi

# Shutdown DB container
docker-down:
	@if command -v docker-compose >/dev/null 2>&1; then \
		echo "Using Docker Compose V1"; \
		docker-compose down; \
	elif docker compose version >/dev/null 2>&1; then \
		echo "Using Docker Compose V2"; \
		docker compose down; \
	else \
		echo "Error: Neither Docker Compose V1 nor V2 found. Please install Docker Compose."; \
		exit 1; \
	fi

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi

.PHONY: all build run test clean watch
