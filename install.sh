#!/bin/bash

set -e  # Exit on any error

echo "Installing Go dependencies..."

# Download Go modules
echo "Downloading Go modules..."
go mod download

# Tidy dependencies
echo "Tidying Go modules..."
go mod tidy

# Build the application
echo "Building application..."
mkdir -p bin
go build -o bin/api cmd/main.go

echo "✅ Installation complete!"
echo "Run the application with: ./bin/api"
