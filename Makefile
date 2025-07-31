.PHONY: help build run dev test clean docker-build docker-run

# Default target
help:
	@echo "Photo-Go Backend Management"
	@echo ""
	@echo "Available commands:"
	@echo "  make build        - Build Go API server"
	@echo "  make build-rust   - Build Rust transcode service"
	@echo "  make run          - Run Go API server"
	@echo "  make dev          - Run development environment"
	@echo "  make test         - Run tests"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make docker-build - Build Docker images"
	@echo "  make docker-run   - Run with Docker Compose"
	@echo "  make docker-stop  - Stop Docker services"

# Build Go API server
build:
	@echo "🔨 Building Go API server..."
	@go build -o bin/api-server ./cmd/api-server
	@echo "✅ Go API server built successfully"

# Build Rust transcode service
build-rust:
	@echo "🔨 Building Rust transcode service..."
	@cd pkg/transcode && cargo build --release
	@echo "✅ Rust transcode service built successfully"

# Run Go API server
run:
	@echo "🚀 Running Go API server..."
	@go run ./cmd/api-server/main.go

# Run development environment
dev:
	@echo "🚀 Starting development environment..."
	@./scripts/dev.sh

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test ./...
	@cd pkg/transcode && cargo test

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf bin/
	@cd pkg/transcode && cargo clean
	@echo "✅ Clean completed"

# Build Docker images
docker-build:
	@echo "🐳 Building Docker images..."
	@docker-compose build
	@echo "✅ Docker images built successfully"

# Run with Docker Compose
docker-run:
	@echo "🐳 Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "✅ Services started successfully"
	@echo "📊 Service URLs:"
	@echo "   - Go API Server: http://localhost:8080"
	@echo "   - Rust Transcode Service: http://localhost:8081"
	@echo "   - PostgreSQL: localhost:5432"
	@echo "   - Redis: localhost:6379"

# Stop Docker services
docker-stop:
	@echo "🛑 Stopping Docker services..."
	@docker-compose down
	@echo "✅ Services stopped successfully"

# Install dependencies
install:
	@echo "📦 Installing Go dependencies..."
	@go mod tidy
	@echo "📦 Installing Rust dependencies..."
	@cd pkg/transcode && cargo build
	@echo "✅ Dependencies installed successfully"

# Database migrations
migrate:
	@echo "🗄️ Running database migrations..."
	@go run ./cmd/migrate/main.go

# Generate API documentation
docs:
	@echo "📚 Generating API documentation..."
	@swag init -g cmd/api-server/main.go -o docs
	@echo "✅ API documentation generated"

# Security check
security:
	@echo "🔒 Running security checks..."
	@go vet ./...
	@cd pkg/transcode && cargo audit
	@echo "✅ Security checks completed" 