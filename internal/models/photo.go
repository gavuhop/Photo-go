package models

import (
	"time"

	"gorm.io/gorm"
)

type Photo struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	User        User           `json:"user" gorm:"foreignKey:UserID"`
	Title       string         `json:"title"`
	Description string         `json:"description"`
	FileName    string         `json:"file_name" gorm:"not null"`
	FilePath    string         `json:"file_path" gorm:"not null"`
	FileSize    int64          `json:"file_size"`
	MimeType    string         `json:"mime_type"`
	Width       int            `json:"width"`
	Height      int            `json:"height"`
	Thumbnail   string         `json:"thumbnail"`
	IsPublic    bool           `json:"is_public" gorm:"default:false"`
	Tags        []Tag          `json:"tags" gorm:"many2many:photo_tags;"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type Tag struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null"`
	Color     string         `json:"color"`
	Photos    []Photo        `json:"photos" gorm:"many2many:photo_tags;"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type PhotoUploadRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IsPublic    bool     `json:"is_public"`
	Tags        []string `json:"tags"`
}

type PhotoUpdateRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	IsPublic    bool     `json:"is_public"`
	Tags        []string `json:"tags"`
}

type PhotoResponse struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	FileSize    int64     `json:"file_size"`
	MimeType    string    `json:"mime_type"`
	Width       int       `json:"width"`
	Height      int       `json:"height"`
	Thumbnail   string    `json:"thumbnail"`
	IsPublic    bool      `json:"is_public"`
	Tags        []Tag     `json:"tags"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
