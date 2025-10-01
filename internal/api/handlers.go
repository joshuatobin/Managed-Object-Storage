package api

import (
	"net/http"
	"time"

	"managed-object-storage/internal/aws"
	"managed-object-storage/internal/models"

	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	s3Service *aws.S3Service
}

// NewHandler creates a new handler
func NewHandler(s3Service *aws.S3Service) *Handler {
	return &Handler{
		s3Service: s3Service,
	}
}

// PresignUpload handles POST /v1/presign/upload
func (h *Handler) PresignUpload(c *gin.Context) {
	var req models.PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate object key (prevent path traversal)
	if !isValidObjectKey(req.ObjectKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object key"})
		return
	}

	url, expiresAt, err := h.s3Service.GeneratePresignedUploadURL(
		c.Request.Context(),
		req.TenantID,
		req.ObjectKey,
		req.ContentType,
		req.MaxSize,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	headers := make(map[string]string)
	if req.ContentType != "" {
		headers["Content-Type"] = req.ContentType
	}

	response := models.PresignUploadResponse{
		URL:       url,
		Method:    "PUT",
		Headers:   headers,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

// PresignDownload handles POST /v1/presign/download
func (h *Handler) PresignDownload(c *gin.Context) {
	var req models.PresignDownloadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate object key (prevent path traversal)
	if !isValidObjectKey(req.ObjectKey) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object key"})
		return
	}

	url, expiresAt, err := h.s3Service.GeneratePresignedDownloadURL(
		c.Request.Context(),
		req.TenantID,
		req.ObjectKey,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := models.PresignDownloadResponse{
		URL:       url,
		ExpiresAt: expiresAt,
	}

	c.JSON(http.StatusOK, response)
}

// ListObjects handles POST /v1/list
func (h *Handler) ListObjects(c *gin.Context) {
	var req models.ListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set default limit if not provided
	limit := int64(1000)
	if req.Limit > 0 {
		limit = int64(req.Limit)
	}

	result, err := h.s3Service.ListObjects(
		c.Request.Context(),
		req.TenantID,
		req.Prefix,
		limit,
		req.Marker,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Convert mock objects to our response format
	var objects []models.ObjectInfo
	for _, obj := range result.Contents {
		// Remove tenant prefix from the key for response
		key := obj.Key
		tenantPrefix := "tenants/" + req.TenantID + "/"
		if len(key) > len(tenantPrefix) {
			key = key[len(tenantPrefix):]
		}

		objects = append(objects, models.ObjectInfo{
			Key:          key,
			Size:         obj.Size,
			LastModified: obj.LastModified,
			ETag:         obj.ETag,
		})
	}

	response := models.ListResponse{
		Objects:    objects,
		NextMarker: getStringValue(result.NextContinuationToken),
		Truncated:  result.IsTruncated,
	}

	c.JSON(http.StatusOK, response)
}

// DeleteObjects handles POST /v1/delete
func (h *Handler) DeleteObjects(c *gin.Context) {
	var req models.DeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate object keys
	for _, key := range req.ObjectKeys {
		if !isValidObjectKey(key) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid object key: " + key})
			return
		}
	}

	result, err := h.s3Service.DeleteObjects(
		c.Request.Context(),
		req.TenantID,
		req.ObjectKeys,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Extract deleted keys (remove tenant prefix)
	var deleted []string
	for _, obj := range result.Deleted {
		key := obj.Key
		tenantPrefix := "tenants/" + req.TenantID + "/"
		if len(key) > len(tenantPrefix) {
			key = key[len(tenantPrefix):]
		}
		deleted = append(deleted, key)
	}

	// Extract errors
	var errors []string
	for _, err := range result.Errors {
		errors = append(errors, err.Message)
	}

	response := models.DeleteResponse{
		Deleted: deleted,
		Errors:  errors,
	}

	c.JSON(http.StatusOK, response)
}

// Health handles GET /healthz
func (h *Handler) Health(c *gin.Context) {
	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

// Helper functions

func isValidObjectKey(key string) bool {
	// Basic validation to prevent path traversal
	if key == "" || key[0] == '/' || key[0] == '.' {
		return false
	}
	
	// Check for path traversal patterns
	for i := 0; i < len(key)-1; i++ {
		if key[i:i+2] == ".." {
			return false
		}
	}
	
	return true
}

func getStringValue(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
