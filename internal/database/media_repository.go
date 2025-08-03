package database

import (
	"gorm.io/gorm"
)

type Media struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string // video, image
	Path      string
	CreatedAt int64
	UpdatedAt int64
	// TODO: Thêm các trường metadata (duration, quality, ...)
}

type MediaRepository interface {
	Create(media *Media) error
	FindByID(id uint) (*Media, error)
	List(offset, limit int) ([]Media, error)
}

type GormMediaRepository struct {
	DB *gorm.DB
}

func NewGormMediaRepository(db *gorm.DB) *GormMediaRepository {
	return &GormMediaRepository{DB: db}
}

func (r *GormMediaRepository) Create(media *Media) error {
	return r.DB.Create(media).Error
}

func (r *GormMediaRepository) FindByID(id uint) (*Media, error) {
	var m Media
	err := r.DB.First(&m, id).Error
	return &m, err
}

func (r *GormMediaRepository) List(offset, limit int) ([]Media, error) {
	var ms []Media
	err := r.DB.Offset(offset).Limit(limit).Find(&ms).Error
	return ms, err
}
