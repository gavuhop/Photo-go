package api

import (
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, db *gorm.DB, redisClient *redis.Client) {
	// API v1 group
	v1 := r.Group("/api/v1")

	// Health check
	v1.GET("/health", healthCheck)

	// Auth routes
	auth := v1.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.POST("/refresh", refreshToken)
	}

	// Photo routes (protected)
	photos := v1.Group("/photos")
	photos.Use(authMiddleware())
	{
		photos.GET("", getPhotos)
		photos.POST("", uploadPhoto)
		photos.GET("/:id", getPhoto)
		photos.PUT("/:id", updatePhoto)
		photos.DELETE("/:id", deletePhoto)
	}

	// Transcode routes (protected)
	transcode := v1.Group("/transcode")
	transcode.Use(authMiddleware())
	{
		transcode.POST("/video", transcodeVideo)
		transcode.POST("/image", transcodeImage)
		transcode.GET("/status/:job_id", getTranscodeStatus)
	}
}

// Health check endpoint
func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"status":  "ok",
		"message": "Photo-Go API is running",
	})
}

// Placeholder functions - will be implemented
func register(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Register endpoint"})
}

func login(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Login endpoint"})
}

func refreshToken(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Refresh token endpoint"})
}

func getPhotos(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get photos endpoint"})
}

func uploadPhoto(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Upload photo endpoint"})
}

func getPhoto(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get photo endpoint"})
}

func updatePhoto(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Update photo endpoint"})
}

func deletePhoto(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Delete photo endpoint"})
}

func transcodeVideo(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Transcode video endpoint"})
}

func transcodeImage(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Transcode image endpoint"})
}

func getTranscodeStatus(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Get transcode status endpoint"})
}

// Auth middleware placeholder
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// TODO: Implement JWT authentication
		c.Next()
	}
}
