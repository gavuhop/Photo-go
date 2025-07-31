# Photo-Go API Documentation

## Overview

Photo-Go là một backend API server sử dụng Go và Rust cho ứng dụng quản lý ảnh. API được thiết kế để hỗ trợ cả web và mobile app.

## Base URL

- **Development**: `http://localhost:8080/api/v1`
- **Production**: `https://api.photo-go.com/api/v1`

## Authentication

API sử dụng JWT (JSON Web Token) cho authentication. Tất cả các protected endpoints yêu cầu header:

```
Authorization: Bearer <access_token>
```

## Endpoints

### Health Check

#### GET /health

Kiểm tra trạng thái của API server.

**Response:**
```json
{
  "status": "ok",
  "message": "Photo-Go API is running"
}
```

### Authentication

#### POST /auth/register

Đăng ký tài khoản mới.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "username": "username",
  "full_name": "Full Name"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "username": "username",
    "full_name": "Full Name",
    "avatar": "",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_in": 86400
}
```

#### POST /auth/login

Đăng nhập.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** Tương tự như register.

#### POST /auth/refresh

Làm mới access token.

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Photos

#### GET /photos

Lấy danh sách ảnh của user.

**Query Parameters:**
- `page` (optional): Số trang (default: 1)
- `limit` (optional): Số item mỗi trang (default: 20)
- `search` (optional): Tìm kiếm theo title
- `tags` (optional): Lọc theo tags (comma-separated)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "title": "My Photo",
      "description": "A beautiful photo",
      "file_name": "photo.jpg",
      "file_path": "/uploads/photo.jpg",
      "file_size": 1024000,
      "mime_type": "image/jpeg",
      "width": 1920,
      "height": 1080,
      "thumbnail": "/uploads/thumbnails/photo.jpg",
      "is_public": false,
      "tags": [
        {
          "id": 1,
          "name": "nature",
          "color": "#4CAF50"
        }
      ],
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

#### POST /photos

Upload ảnh mới.

**Request:** Multipart form data
- `file`: File ảnh
- `title` (optional): Tiêu đề ảnh
- `description` (optional): Mô tả
- `is_public` (optional): Công khai (default: false)
- `tags` (optional): Tags (comma-separated)

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "My Photo",
  "description": "A beautiful photo",
  "file_name": "photo.jpg",
  "file_path": "/uploads/photo.jpg",
  "file_size": 1024000,
  "mime_type": "image/jpeg",
  "width": 1920,
  "height": 1080,
  "thumbnail": "/uploads/thumbnails/photo.jpg",
  "is_public": false,
  "tags": [],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

#### GET /photos/{id}

Lấy chi tiết ảnh.

**Response:** Tương tự như POST /photos.

#### PUT /photos/{id}

Cập nhật thông tin ảnh.

**Request Body:**
```json
{
  "title": "Updated Title",
  "description": "Updated description",
  "is_public": true,
  "tags": ["nature", "landscape"]
}
```

#### DELETE /photos/{id}

Xóa ảnh.

**Response:**
```json
{
  "message": "Photo deleted successfully"
}
```

### Transcode

#### POST /transcode/video

Transcode video.

**Request Body:**
```json
{
  "input_path": "/uploads/video.mp4",
  "output_path": "/uploads/video_compressed.mp4",
  "format": "mp4",
  "codec": "h264",
  "bitrate": "1000k",
  "resolution": "1280x720",
  "fps": 30
}
```

**Response:**
```json
{
  "job_id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "message": "Video transcode job created"
}
```

#### POST /transcode/image

Transcode ảnh.

**Request Body:**
```json
{
  "input_path": "/uploads/image.jpg",
  "output_path": "/uploads/image_compressed.jpg",
  "format": "jpeg",
  "quality": 85,
  "width": 800,
  "height": 600,
  "resize_mode": "fit"
}
```

#### GET /transcode/status/{job_id}

Kiểm tra trạng thái transcode job.

**Response:**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "completed",
  "input_path": "/uploads/video.mp4",
  "output_path": "/uploads/video_compressed.mp4",
  "format": "mp4",
  "progress": 100.0,
  "error_message": null,
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:05:00Z"
}
```

## Error Responses

### 400 Bad Request
```json
{
  "error": "validation_error",
  "message": "Invalid input data",
  "details": {
    "email": ["Email is required"],
    "password": ["Password must be at least 6 characters"]
  }
}
```

### 401 Unauthorized
```json
{
  "error": "unauthorized",
  "message": "Invalid or expired token"
}
```

### 403 Forbidden
```json
{
  "error": "forbidden",
  "message": "Access denied"
}
```

### 404 Not Found
```json
{
  "error": "not_found",
  "message": "Resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal_error",
  "message": "Something went wrong"
}
```

## Rate Limiting

API có rate limiting để tránh spam:
- **Authentication endpoints**: 5 requests per minute
- **Photo upload**: 10 requests per minute
- **Other endpoints**: 100 requests per minute

## File Upload Limits

- **Maximum file size**: 50MB
- **Supported formats**: JPEG, PNG, GIF, WebP, MP4, MOV, AVI
- **Image dimensions**: Tự động resize nếu quá lớn

## Rust Microservices

### Transcode Service

**Base URL**: `http://localhost:8081/api/v1`

#### Transcode Endpoints
- `POST /transcode/video` - Transcode video với FFmpeg
- `POST /transcode/image` - Xử lý ảnh với image crate  
- `POST /transcode/audio` - Transcode audio với FFmpeg
- `GET /transcode/status/{job_id}` - Kiểm tra trạng thái

#### Metadata Extraction
- `POST /metadata/extract` - Trích xuất metadata từ media files
- `POST /metadata/analyze-video` - Phân tích video chi tiết

#### Image Processing
- `POST /process/filter` - Áp dụng filters (blur, sharpen, sepia, etc.)
- `POST /process/watermark` - Thêm watermark
- `POST /process/batch` - Xử lý hàng loạt

#### AI Analysis
- `POST /ai/detect-objects` - Phát hiện đối tượng trong ảnh
- `POST /ai/detect-faces` - Phát hiện khuôn mặt
- `POST /ai/analyze-colors` - Phân tích màu sắc chủ đạo

### Rust Service Features

#### 🎬 Video Processing
- H.264/H.265 encoding với FFmpeg
- Batch video compression
- Video analysis và metadata extraction
- Frame extraction
- Audio track separation

#### 🖼️ Image Processing  
- **Filters**: Blur, Sharpen, Sepia, Grayscale, Vintage, Brightness, Contrast, Saturation
- **Effects**: Vignette, Emboss, Edge Detection
- **Transformations**: Resize, Rotate, Crop
- **Watermarking**: Multi-position với opacity control
- **Format conversion**: JPEG, PNG, WebP, GIF với quality settings

#### 🎵 Audio Processing
- MP3, AAC, OGG, FLAC, WAV transcoding
- Audio normalization
- Bitrate và sample rate conversion  
- Audio extraction từ video
- Metadata analysis

#### 🤖 AI Features
- **Object Detection**: YOLO-style object recognition
- **Face Detection**: Bounding boxes, age/gender estimation, emotion analysis
- **Color Analysis**: Dominant colors extraction
- **Image Quality Assessment**: Sharpness, brightness, contrast, noise analysis
- **Auto-tagging**: Tự động tạo tags dựa trên AI analysis

#### 📊 Metadata Extraction
- **EXIF Data**: Camera info, GPS, technical settings
- **Video Metadata**: Codec, bitrate, frame rate, duration  
- **Audio Metadata**: Sample rate, channels, codec
- **File Information**: Dimensions, size, creation date

## Development

### Chạy locally

```bash
# Install dependencies
make install

# Run development environment
make dev

# Hoặc chạy riêng lẻ
make run          # Go API server
cd pkg/transcode && cargo run  # Rust service
```

### Docker

```bash
# Build và chạy với Docker
make docker-build
make docker-run
``` 