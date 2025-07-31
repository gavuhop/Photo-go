package main

import (
	"log"

	"photo-go/internal/database"
	"photo-go/internal/models"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Running database migrations...")

	// Auto migrate all models
	err = db.AutoMigrate(
		&models.User{},
		&models.Photo{},
		&models.Tag{},
		&models.MediaFile{},
		&models.TranscodeJob{},
		&models.Album{},
		&models.MediaMetadata{},
	)

	if err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed successfully!")

	// Create indexes for better performance
	createIndexes(db)

	log.Println("Database setup completed!")
}

func createIndexes(db *gorm.DB) {
	log.Println("Creating database indexes...")

	// Photos indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_photos_user_id ON photos(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_photos_created_at ON photos(created_at)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_photos_is_public ON photos(is_public)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_photos_mime_type ON photos(mime_type)")

	// Media files indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_media_files_user_id ON media_files(user_id)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_media_files_type ON media_files(file_type)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_media_files_status ON media_files(status)")

	// Transcode jobs indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_transcode_jobs_status ON transcode_jobs(status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_transcode_jobs_created_at ON transcode_jobs(created_at)")

	// Tags indexes
	db.Exec("CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name)")

	log.Println("Database indexes created successfully!")
}
