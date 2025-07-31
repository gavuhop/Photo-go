#!/bin/bash

# Development script for Photo-Go backend
# This script starts both Go API server and Rust microservices

set -e

echo "🚀 Starting Photo-Go Development Environment"

# Check if required tools are installed
check_requirements() {
    echo "📋 Checking requirements..."
    
    if ! command -v go &> /dev/null; then
        echo "❌ Go is not installed"
        exit 1
    fi
    
    if ! command -v cargo &> /dev/null; then
        echo "❌ Rust/Cargo is not installed"
        exit 1
    fi
    
    if ! command -v ffmpeg &> /dev/null; then
        echo "⚠️  FFmpeg is not installed (video transcode will not work)"
    fi
    
    echo "✅ Requirements check passed"
}

# Build Rust services
build_rust_services() {
    echo "🔨 Building Rust services..."
    
    cd pkg/transcode
    cargo build --release
    cd ../..
    
    echo "✅ Rust services built successfully"
}

# Start services
start_services() {
    echo "🌐 Starting services..."
    
    # Start Go API server in background
    echo "📡 Starting Go API server on port 8080..."
    cd cmd/api-server
    go run main.go &
    GO_PID=$!
    cd ../..
    
    # Start Rust transcode service in background
    echo "🎬 Starting Rust transcode service on port 8081..."
    cd pkg/transcode
    cargo run &
    RUST_PID=$!
    cd ../..
    
    echo "✅ Services started successfully"
    echo "📊 Service status:"
    echo "   - Go API Server: http://localhost:8080"
    echo "   - Rust Transcode Service: http://localhost:8081"
    echo ""
    echo "Press Ctrl+C to stop all services"
    
    # Wait for interrupt
    trap 'cleanup' INT
    wait
}

# Cleanup function
cleanup() {
    echo ""
    echo "🛑 Stopping services..."
    
    if [ ! -z "$GO_PID" ]; then
        kill $GO_PID 2>/dev/null || true
        echo "✅ Go API server stopped"
    fi
    
    if [ ! -z "$RUST_PID" ]; then
        kill $RUST_PID 2>/dev/null || true
        echo "✅ Rust transcode service stopped"
    fi
    
    echo "👋 Development environment stopped"
    exit 0
}

# Main execution
main() {
    check_requirements
    build_rust_services
    start_services
}

main 