package v1

import (
	"photo-go/internal/database"

	"gorm.io/gorm"
)

type MediaRepository interface {
	Create(media *database.Media) error
	FindByID(id uint) (*database.Media, error)
	List(offset, limit int) ([]database.Media, error)
}

type GormMediaRepository struct {
	DB *gorm.DB
}

func NewGormMediaRepository(db *gorm.DB) *GormMediaRepository {
	return &GormMediaRepository{DB: db}
}

func (r *GormMediaRepository) Create(media *database.Media) error {
	return r.DB.Create(media).Error
}

func (r *GormMediaRepository) FindByID(id uint) (*database.Media, error) {
	var m database.Media
	err := r.DB.First(&m, id).Error
	return &m, err
}

func (r *GormMediaRepository) List(offset, limit int) ([]database.Media, error) {
	var ms []database.Media
	err := r.DB.Offset(offset).Limit(limit).Find(&ms).Error
	return ms, err
}
