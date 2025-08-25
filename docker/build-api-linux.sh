#!/bin/bash

# Build script for Promptly API (Linux)
set -e

echo "Building Promptly API for Linux..."

# Get the script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Change to project root
cd "$PROJECT_ROOT"

# Build the Go binary for Linux (linux/amd64) - for Docker container
GOOS=linux GOARCH=amd64 go build -o docker/promptly-api-linux cmd/promptly/main.go

# Copy config file to docker directory if it exists
if [ -f "config.yaml" ]; then
    cp config.yaml docker/
    echo "📋 Config file copied to docker/"
fi

# Create data directory structure in docker folder based on config
if [ -f "config.yaml" ]; then
    DATA_PATH=$(grep "^data:" config.yaml | cut -d' ' -f2)
    if [ ! -z "$DATA_PATH" ]; then
        mkdir -p "docker/$(dirname "$DATA_PATH")"
        chmod 755 "docker/$(dirname "$DATA_PATH")"
        echo "📁 Created data directory: docker/$(dirname "$DATA_PATH")"
    fi
fi

echo "✅ Binary built successfully: docker/promptly-api-linux"
# echo "To run from docker directory: cd docker && ./promptly-api-linux serve"