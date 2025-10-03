package main

import (
	"log"
	"os"

	"managed-object-storage/internal/api"
	"managed-object-storage/internal/aws"

	"github.com/gin-gonic/gin"
	"github.com/penglongli/gin-metrics/ginmetrics"
)

func main() {
	// Get configuration from environment variables
	bucketName := getEnv("S3_BUCKET_NAME", "aptible-objects-dev")
	region := getEnv("AWS_REGION", "us-east-1")
	port := getEnv("PORT", "8080")

	// Initialize S3 service
	s3Service, err := aws.NewS3Service(bucketName, region)
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	// Initialize handler
	handler := api.NewHandler(s3Service)

	// Setup Gin router
	router := gin.Default()

	// Setup gin-metrics
	m := ginmetrics.GetMonitor()
	m.SetMetricPath("/metrics")
	m.SetSlowTime(10)
	m.SetDuration([]float64{0.1, 0.3, 1.2, 5, 10})
	m.Use(router)

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API routes
	v1 := router.Group("/v1")
	{
		v1.POST("/presign/upload", handler.PresignUpload)
		v1.POST("/presign/download", handler.PresignDownload)
		v1.POST("/list", handler.ListObjects)
		v1.POST("/delete", handler.DeleteObjects)
	}

	// Health check
	router.GET("/healthz", handler.Health)

	// Start server
	log.Printf("Starting server on port %s", port)
	log.Printf("S3 Bucket: %s", bucketName)
	log.Printf("AWS Region: %s", region)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
