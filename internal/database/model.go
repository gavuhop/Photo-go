package database

type Media struct {
	ID        uint   `gorm:"primaryKey"`
	Type      string // video, image
	Path      string
	CreatedAt int64
	UpdatedAt int64
	// TODO: Thêm các trường metadata (duration, quality, ...)
}
