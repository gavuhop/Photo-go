package models

import (
	"time"

	"gorm.io/gorm"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
)

type MediaStatus string

const (
	MediaStatusUploading  MediaStatus = "uploading"
	MediaStatusProcessing MediaStatus = "processing"
	MediaStatusReady      MediaStatus = "ready"
	MediaStatusFailed     MediaStatus = "failed"
)

// MediaFile represents any media file (image, video, audio)
type MediaFile struct {
	ID           uint        `json:"id" gorm:"primaryKey"`
	UserID       uint        `json:"user_id" gorm:"not null"`
	User         User        `json:"user" gorm:"foreignKey:UserID"`
	Title        string      `json:"title"`
	Description  string      `json:"description"`
	FileName     string      `json:"file_name" gorm:"not null"`
	OriginalName string      `json:"original_name"`
	FilePath     string      `json:"file_path" gorm:"not null"`
	FileSize     int64       `json:"file_size"`
	MimeType     string      `json:"mime_type"`
	FileType     MediaType   `json:"file_type" gorm:"not null"`
	Status       MediaStatus `json:"status" gorm:"default:uploading"`

	// Dimensions for images and videos
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	// Duration for videos and audio (in seconds)
	Duration float64 `json:"duration,omitempty"`

	// Thumbnails and previews
	ThumbnailPath string `json:"thumbnail_path,omitempty"`
	PreviewPath   string `json:"preview_path,omitempty"`

	// Processing results
	ProcessedFiles []ProcessedFile `json:"processed_files" gorm:"foreignKey:MediaFileID"`

	// Metadata
	Metadata *MediaMetadata `json:"metadata,omitempty" gorm:"foreignKey:MediaFileID"`

	// Relationships
	AlbumID *uint  `json:"album_id,omitempty"`
	Album   *Album `json:"album,omitempty" gorm:"foreignKey:AlbumID"`
	Tags    []Tag  `json:"tags" gorm:"many2many:media_file_tags;"`

	// Permissions
	IsPublic bool `json:"is_public" gorm:"default:false"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// ProcessedFile represents different versions/formats of a media file
type ProcessedFile struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	MediaFileID uint       `json:"media_file_id" gorm:"not null"`
	MediaFile   *MediaFile `json:"media_file,omitempty" gorm:"foreignKey:MediaFileID"`

	// Processing info
	ProcessType string `json:"process_type"` // "thumbnail", "compressed", "watermarked", etc.
	Format      string `json:"format"`       // "jpeg", "mp4", "webp", etc.
	Quality     int    `json:"quality,omitempty"`

	// File info
	FilePath string `json:"file_path" gorm:"not null"`
	FileSize int64  `json:"file_size"`

	// Dimensions
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`

	// Video/Audio specific
	Duration float64 `json:"duration,omitempty"`
	Bitrate  int     `json:"bitrate,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Album for organizing media files
type Album struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	UserID      uint   `json:"user_id" gorm:"not null"`
	User        User   `json:"user" gorm:"foreignKey:UserID"`
	Name        string `json:"name" gorm:"not null"`
	Description string `json:"description"`
	CoverPath   string `json:"cover_path"`
	IsPublic    bool   `json:"is_public" gorm:"default:false"`

	MediaFiles []MediaFile `json:"media_files" gorm:"foreignKey:AlbumID"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// MediaMetadata stores extracted metadata from media files
type MediaMetadata struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	MediaFileID uint      `json:"media_file_id" gorm:"not null;uniqueIndex"`
	MediaFile   MediaFile `json:"media_file" gorm:"foreignKey:MediaFileID"`

	// Common metadata
	Creator     string `json:"creator,omitempty"`
	Copyright   string `json:"copyright,omitempty"`
	Keywords    string `json:"keywords,omitempty"`
	Description string `json:"description,omitempty"`

	// Image specific metadata (EXIF)
	CameraMake   string  `json:"camera_make,omitempty"`
	CameraModel  string  `json:"camera_model,omitempty"`
	Lens         string  `json:"lens,omitempty"`
	FocalLength  float64 `json:"focal_length,omitempty"`
	Aperture     float64 `json:"aperture,omitempty"`
	ShutterSpeed string  `json:"shutter_speed,omitempty"`
	ISO          int     `json:"iso,omitempty"`
	Flash        bool    `json:"flash,omitempty"`

	// Location data
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
	Altitude  *float64 `json:"altitude,omitempty"`
	Location  string   `json:"location,omitempty"`

	// Date taken (different from created_at which is upload time)
	DateTaken *time.Time `json:"date_taken,omitempty"`

	// Video/Audio specific
	Codec      string  `json:"codec,omitempty"`
	Bitrate    int     `json:"bitrate,omitempty"`
	FrameRate  float64 `json:"frame_rate,omitempty"`
	SampleRate int     `json:"sample_rate,omitempty"`
	Channels   int     `json:"channels,omitempty"`

	// AI analysis results
	Objects string `json:"objects,omitempty"`   // JSON array of detected objects
	Faces   string `json:"faces,omitempty"`     // JSON array of face data
	Colors  string `json:"colors,omitempty"`    // JSON array of dominant colors
	Tags    string `json:"auto_tags,omitempty"` // JSON array of auto-generated tags

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TranscodeJob tracks processing jobs
type TranscodeJob struct {
	ID          string     `json:"id" gorm:"primaryKey"` // UUID
	MediaFileID uint       `json:"media_file_id" gorm:"not null"`
	MediaFile   *MediaFile `json:"media_file,omitempty" gorm:"foreignKey:MediaFileID"`

	JobType  string  `json:"job_type"` // "transcode", "thumbnail", "compress", etc.
	Status   string  `json:"status"`   // "pending", "processing", "completed", "failed"
	Progress float32 `json:"progress"` // 0-100

	// Input/Output paths
	InputPath  string `json:"input_path"`
	OutputPath string `json:"output_path"`

	// Processing parameters
	Parameters string `json:"parameters"` // JSON string of processing parameters

	// Results
	ErrorMessage *string `json:"error_message,omitempty"`
	ResultData   string  `json:"result_data,omitempty"` // JSON string of results

	// Timing
	StartedAt   *time.Time `json:"started_at,omitempty"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// Request/Response DTOs
type MediaUploadRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	AlbumID     *uint    `json:"album_id"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

type MediaProcessRequest struct {
	ProcessType string                 `json:"process_type"` // "resize", "compress", "watermark", etc.
	Parameters  map[string]interface{} `json:"parameters"`
}

type MediaResponse struct {
	ID             uint            `json:"id"`
	Title          string          `json:"title"`
	Description    string          `json:"description"`
	FileName       string          `json:"file_name"`
	FilePath       string          `json:"file_path"`
	FileSize       int64           `json:"file_size"`
	MimeType       string          `json:"mime_type"`
	FileType       MediaType       `json:"file_type"`
	Status         MediaStatus     `json:"status"`
	Width          int             `json:"width,omitempty"`
	Height         int             `json:"height,omitempty"`
	Duration       float64         `json:"duration,omitempty"`
	ThumbnailPath  string          `json:"thumbnail_path,omitempty"`
	ProcessedFiles []ProcessedFile `json:"processed_files,omitempty"`
	Album          *Album          `json:"album,omitempty"`
	Tags           []Tag           `json:"tags"`
	IsPublic       bool            `json:"is_public"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type AlbumCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

type AlbumResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CoverPath   string    `json:"cover_path"`
	IsPublic    bool      `json:"is_public"`
	MediaCount  int       `json:"media_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
