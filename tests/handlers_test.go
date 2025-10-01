package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"managed-object-storage/internal/api"
	"managed-object-storage/internal/aws"
	"managed-object-storage/internal/models"

	"github.com/gin-gonic/gin"
)

func TestPresignUpload(t *testing.T) {
	// Mock S3 service
	s3Service, _ := aws.NewS3Service("test-bucket", "us-east-1")
	handler := api.NewHandler(s3Service)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/presign/upload", handler.PresignUpload)

	// Test request
	reqBody := models.PresignUploadRequest{
		TenantID:   "test-tenant",
		ObjectKey:  "test-file.txt",
		ContentType: "text/plain",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/v1/presign/upload", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Basic assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response models.PresignUploadResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response.URL == "" {
		t.Error("Expected URL to be present")
	}
	
	if response.Method != "PUT" {
		t.Errorf("Expected method PUT, got %s", response.Method)
	}
}

func TestPresignDownload(t *testing.T) {
	// Mock S3 service
	s3Service, _ := aws.NewS3Service("test-bucket", "us-east-1")
	handler := api.NewHandler(s3Service)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/v1/presign/download", handler.PresignDownload)

	// Test request
	reqBody := models.PresignDownloadRequest{
		TenantID:  "test-tenant",
		ObjectKey: "test-file.txt",
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/v1/presign/download", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Basic assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response models.PresignDownloadResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response.URL == "" {
		t.Error("Expected URL to be present")
	}
}

func TestHealth(t *testing.T) {
	// Mock S3 service
	s3Service, _ := aws.NewS3Service("test-bucket", "us-east-1")
	handler := api.NewHandler(s3Service)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/healthz", handler.Health)

	req, _ := http.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Basic assertions
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	var response models.HealthResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}
	
	if response.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got %s", response.Status)
	}
}
