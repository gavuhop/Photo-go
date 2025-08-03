# Photo-Go Backend

Backend API server sử dụng Go và Rust cho ứng dụng web và mobile.

## Kiến trúc

```
Photo-go/
├── cmd/
│   ├── api-server/          # Go API server chính
│   └── rust-services/       # Rust microservices
├── internal/
│   ├── api/                 # API handlers
│   ├── models/              # Data models
│   ├── services/            # Business logic
│   └── database/            # Database layer
├── pkg/
│   ├── transcode/           # Rust transcode service
│   └── image-processing/    # Rust image processing
├── configs/                 # Configuration files
├── scripts/                 # Build & deployment scripts
└── docs/                    # API documentation
```

## Công nghệ sử dụng

### Go (API Server chính)
- **Gin** - Web framework
- **GORM** - ORM
- **JWT** - Authentication
- **Redis** - Caching
- **PostgreSQL** - Database chính

### Rust (Microservices cho tác vụ nặng)
- **Actix-web** - Web framework cho microservices
- **FFmpeg-next** - Video/audio transcode
- **Image** - Image processing
- **OpenCV** - Computer vision
- **SQLx** - Async database operations

## Cài đặt và chạy

### Yêu cầu
- Go 1.21+
- Rust 1.70+
- PostgreSQL
- Redis
- FFmpeg

### Chạy development
```bash
# Chạy Go API server
cd cmd/api-server
go run main.go

# Chạy Rust transcode service
cd pkg/transcode
cargo run

# Chạy Rust image processing service
cd pkg/image-processing
cargo run
```

## API Endpoints

### Authentication
- `POST /api/auth/login`
- `POST /api/auth/register`
- `POST /api/auth/refresh`

### Photos
- `GET /api/photos`
- `POST /api/photos`
- `GET /api/photos/{id}`
- `PUT /api/photos/{id}`
- `DELETE /api/photos/{id}`

### Transcode
- `POST /api/transcode/video`
- `POST /api/transcode/image`
- `GET /api/transcode/status/{job_id}`

## Microservices

### Transcode Service (Rust)
- Video compression
- Format conversion
- Thumbnail generation
- Batch processing

### Image Processing Service (Rust)
- Image resizing
- Filters và effects
- Face detection
- OCR processing
