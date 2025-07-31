# Photo-Go Rust Transcode Service

High-performance media processing microservice built with Rust and Actix-web.

## Features

### 🎬 Video Processing
- **FFmpeg Integration**: H.264/H.265 encoding, format conversion
- **Batch Processing**: Multiple videos với queue system
- **Metadata Extraction**: Codec, bitrate, frame rate, duration
- **Frame Extraction**: Extract specific frames từ video
- **Audio Separation**: Tách audio track từ video

### 🖼️ Advanced Image Processing
- **Filters & Effects**: 15+ built-in filters
  - Basic: Blur, Sharpen, Brightness, Contrast
  - Artistic: Sepia, Vintage, Grayscale
  - Advanced: Vignette, Emboss, Edge Detection
- **Transformations**: Resize, Rotate, Crop với smart algorithms
- **Watermarking**: 9 positions, opacity control, scaling
- **Format Support**: JPEG, PNG, WebP, GIF, BMP, TIFF
- **Quality Control**: Lossy/lossless compression options

### 🎵 Audio Processing
- **Format Transcoding**: MP3, AAC, OGG, FLAC, WAV
- **Audio Normalization**: Loudness standards compliance
- **Parameter Control**: Bitrate, sample rate, channels
- **Metadata Analysis**: Technical audio information
- **Quality Assessment**: Audio signal analysis

### 🤖 AI & Machine Learning
- **Object Detection**: Real-world object recognition
- **Face Detection**: Bounding boxes, demographics, emotions
- **Color Analysis**: Dominant colors extraction
- **Auto-tagging**: Smart tags based on content
- **Quality Assessment**: Image sharpness, noise, composition
- **Content Safety**: Adult content detection

### 📊 Metadata Extraction
- **EXIF Data**: Camera settings, GPS coordinates
- **Technical Info**: Dimensions, file size, creation date
- **Smart Analysis**: Automated content categorization

## API Endpoints

### Transcode Operations
```bash
# Video transcoding
POST /api/v1/transcode/video
{
  "input_path": "/path/to/input.mp4",
  "output_path": "/path/to/output.mp4",
  "format": "mp4",
  "codec": "h264",
  "bitrate": "1000k",
  "resolution": "1280x720",
  "fps": 30
}

# Image transcoding  
POST /api/v1/transcode/image
{
  "input_path": "/path/to/input.jpg",
  "output_path": "/path/to/output.png", 
  "format": "png",
  "quality": 85,
  "width": 800,
  "height": 600,
  "resize_mode": "fit"
}

# Audio transcoding
POST /api/v1/transcode/audio
{
  "input_path": "/path/to/input.mp3",
  "output_path": "/path/to/output.flac",
  "format": "flac",
  "bitrate": "320k",
  "sample_rate": 44100,
  "channels": 2
}
```

### Image Processing
```bash
# Apply filters
POST /api/v1/process/filter
{
  "input_path": "/path/to/input.jpg",
  "output_path": "/path/to/filtered.jpg",
  "filter_type": "sepia",
  "intensity": 0.8
}

# Add watermark
POST /api/v1/process/watermark
{
  "input_path": "/path/to/input.jpg",
  "output_path": "/path/to/watermarked.jpg",
  "watermark_path": "/path/to/watermark.png",
  "position": "bottom-right",
  "opacity": 0.5,
  "scale": 0.2
}

# Batch processing
POST /api/v1/process/batch
{
  "input_paths": ["/path/1.jpg", "/path/2.jpg"],
  "output_directory": "/output/dir",
  "operations": [
    {
      "operation_type": "resize",
      "parameters": {"width": 800, "height": 600}
    }
  ]
}
```

### AI Analysis
```bash
# Object detection
POST /api/v1/ai/detect-objects
{
  "image_path": "/path/to/image.jpg"
}

# Face detection
POST /api/v1/ai/detect-faces
{
  "image_path": "/path/to/image.jpg"
}

# Color analysis
POST /api/v1/ai/analyze-colors
{
  "image_path": "/path/to/image.jpg"
}
```

### Metadata Extraction
```bash
# Extract metadata
POST /api/v1/metadata/extract
{
  "file_path": "/path/to/media.jpg",
  "extract_exif": true,
  "extract_ai_tags": true
}

# Video analysis
POST /api/v1/metadata/analyze-video
{
  "file_path": "/path/to/video.mp4",
  "extract_frames": true,
  "frame_interval": 30,
  "extract_audio": true
}
```

## Dependencies

### Core Libraries
- **actix-web**: High-performance async web framework
- **tokio**: Async runtime
- **serde**: Serialization/deserialization
- **anyhow**: Error handling

### Media Processing
- **image**: Core image processing
- **imageproc**: Advanced image algorithms
- **photon-rs**: GPU-accelerated filters
- **fast_image_resize**: High-performance resizing
- **ffmpeg-next**: FFmpeg bindings for video/audio

### AI & Analysis
- **opencv**: Computer vision (optional)
- **candle-core**: ML inference (optional)
- **exif**: EXIF metadata extraction

### Audio
- **symphonia**: Audio decoding/encoding
- **hound**: WAV file support

## Installation

### Prerequisites
```bash
# Install FFmpeg
sudo apt install ffmpeg ffmpeg-dev

# Install OpenCV (optional, for advanced AI features)
sudo apt install libopencv-dev

# Rust toolchain
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

### Build
```bash
cd pkg/transcode
cargo build --release
```

### Run
```bash
# Development
cargo run

# Production
cargo run --release

# With custom port
PORT=8081 cargo run
```

## Configuration

Environment variables:
```bash
PORT=8081                    # Service port
RUST_LOG=info               # Log level
UPLOAD_DIR=./uploads        # File storage directory
FFMPEG_PATH=/usr/bin/ffmpeg # FFmpeg binary path
MAX_WORKERS=4               # Processing threads
```

## Performance

### Benchmarks
- **Image Resize**: 100+ images/second (1920x1080 → 800x600)
- **Filter Application**: 50+ images/second với complex filters
- **Video Transcode**: Real-time H.264 encoding (depends on hardware)
- **Metadata Extraction**: 1000+ files/second

### Memory Usage
- **Base Service**: ~10MB RAM
- **Image Processing**: ~50MB per concurrent job
- **Video Processing**: ~200MB per concurrent job

### Optimization Tips
- Use `fast_image_resize` cho bulk resizing
- Enable GPU acceleration với OpenCV
- Batch similar operations together
- Use async operations cho I/O bound tasks

## Error Handling

Service trả về structured errors:
```json
{
  "error": "processing_failed",
  "message": "Detailed error description",
  "details": {
    "file_path": "/path/to/file",
    "operation": "resize"
  }
}
```

Common error types:
- `file_not_found`: Input file doesn't exist
- `invalid_format`: Unsupported file format
- `processing_failed`: Operation failed
- `insufficient_resources`: Out of memory/disk space

## Testing

```bash
# Unit tests
cargo test

# Integration tests
cargo test --test integration

# Benchmark tests
cargo bench

# Test with sample files
./scripts/test_rust_service.sh
```

## Contributing

1. Follow Rust coding standards
2. Add tests for new features
3. Update documentation
4. Benchmark performance impact
5. Handle errors gracefully

## License

MIT License - see LICENSE file for details.