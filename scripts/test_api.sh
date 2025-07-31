#!/bin/bash

# Script to test Photo-Go API endpoints
# Make sure the services are running before executing this script

set -e

echo "🧪 Testing Photo-Go API Endpoints"

API_BASE="http://localhost:8080/api/v1"
RUST_BASE="http://localhost:8081/api/v1"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to test endpoint
test_endpoint() {
    local method=$1
    local url=$2
    local data=$3
    local description=$4
    
    echo -e "${YELLOW}Testing: $description${NC}"
    echo "  $method $url"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$url")
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$url")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [[ $http_code -ge 200 && $http_code -lt 300 ]]; then
        echo -e "  ${GREEN}✓ Success ($http_code)${NC}"
        if [ ! -z "$body" ]; then
            echo "  Response: $(echo "$body" | jq -c . 2>/dev/null || echo "$body")"
        fi
    else
        echo -e "  ${RED}✗ Failed ($http_code)${NC}"
        if [ ! -z "$body" ]; then
            echo "  Error: $body"
        fi
    fi
    echo ""
}

echo ""
echo "=== Go API Server Tests ==="

# Health check
test_endpoint "GET" "$API_BASE/health" "" "Health Check"

# Auth endpoints (these will fail without proper setup, but we test the endpoints)
test_endpoint "POST" "$API_BASE/auth/register" '{
    "email": "test@example.com",
    "password": "password123",
    "username": "testuser",
    "full_name": "Test User"
}' "User Registration"

test_endpoint "POST" "$API_BASE/auth/login" '{
    "email": "test@example.com",
    "password": "password123"
}' "User Login"

echo ""
echo "=== Rust Transcode Service Tests ==="

# Health check
test_endpoint "GET" "$RUST_BASE/health" "" "Rust Service Health Check"

# Video transcode
test_endpoint "POST" "$RUST_BASE/transcode/video" '{
    "input_path": "/tmp/input.mp4",
    "output_path": "/tmp/output.mp4",
    "format": "mp4",
    "codec": "h264",
    "bitrate": "1000k",
    "resolution": "1280x720",
    "fps": 30
}' "Video Transcode"

# Image transcode
test_endpoint "POST" "$RUST_BASE/transcode/image" '{
    "input_path": "/tmp/input.jpg",
    "output_path": "/tmp/output.png",
    "format": "png",
    "quality": 85,
    "width": 800,
    "height": 600,
    "resize_mode": "fit"
}' "Image Transcode"

# Audio transcode
test_endpoint "POST" "$RUST_BASE/transcode/audio" '{
    "input_path": "/tmp/input.mp3",
    "output_path": "/tmp/output.flac",
    "format": "flac",
    "bitrate": "320k",
    "sample_rate": 44100,
    "channels": 2
}' "Audio Transcode"

# Metadata extraction
test_endpoint "POST" "$RUST_BASE/metadata/extract" '{
    "file_path": "/tmp/test.jpg",
    "extract_exif": true,
    "extract_ai_tags": false
}' "Metadata Extraction"

# Video analysis
test_endpoint "POST" "$RUST_BASE/metadata/analyze-video" '{
    "file_path": "/tmp/test.mp4",
    "extract_frames": true,
    "frame_interval": 30,
    "extract_audio": true
}' "Video Analysis"

echo ""
echo "=== Image Processing Tests ==="

# Apply filter
test_endpoint "POST" "$RUST_BASE/process/filter" '{
    "input_path": "/tmp/input.jpg",
    "output_path": "/tmp/filtered.jpg",
    "filter_type": "sepia",
    "intensity": 0.8
}' "Apply Sepia Filter"

# Watermark
test_endpoint "POST" "$RUST_BASE/process/watermark" '{
    "input_path": "/tmp/input.jpg",
    "output_path": "/tmp/watermarked.jpg",
    "watermark_path": "/tmp/watermark.png",
    "position": "bottom-right",
    "opacity": 0.5,
    "scale": 0.2
}' "Add Watermark"

# Batch processing
test_endpoint "POST" "$RUST_BASE/process/batch" '{
    "input_paths": ["/tmp/image1.jpg", "/tmp/image2.jpg"],
    "output_directory": "/tmp/processed",
    "operations": [
        {
            "operation_type": "resize",
            "parameters": {"width": 800, "height": 600}
        }
    ]
}' "Batch Processing"

echo ""
echo "=== AI Analysis Tests ==="

# Object detection
test_endpoint "POST" "$RUST_BASE/ai/detect-objects" '{
    "image_path": "/tmp/test.jpg"
}' "Object Detection"

# Face detection
test_endpoint "POST" "$RUST_BASE/ai/detect-faces" '{
    "image_path": "/tmp/test.jpg"
}' "Face Detection"

# Color analysis
test_endpoint "POST" "$RUST_BASE/ai/analyze-colors" '{
    "image_path": "/tmp/test.jpg"
}' "Color Analysis"

echo ""
echo "=== Job Status Check ==="

# Check transcode status
test_endpoint "GET" "$RUST_BASE/transcode/status/550e8400-e29b-41d4-a716-446655440000" "" "Transcode Job Status"

echo ""
echo -e "${GREEN}🎉 API Testing Completed!${NC}"
echo ""
echo "Notes:"
echo "- Some endpoints may fail if services are not running"
echo "- File-based operations need actual files to exist"
echo "- Database operations require PostgreSQL connection"
echo "- AI features need additional model files"
echo ""
echo "To run the services:"
echo "  make dev          # Development mode"
echo "  make docker-run   # Docker mode"