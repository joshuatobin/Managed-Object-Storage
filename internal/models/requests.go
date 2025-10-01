package models

import "time"

// PresignUploadRequest represents the request to get a presigned upload URL
type PresignUploadRequest struct {
	TenantID   string `json:"tenant_id" binding:"required"`
	ObjectKey  string `json:"object_key" binding:"required"`
	ContentType string `json:"content_type,omitempty"`
	MaxSize    int64  `json:"max_size,omitempty"`
}

// PresignUploadResponse represents the response with presigned upload URL
type PresignUploadResponse struct {
	URL        string    `json:"url"`
	Method     string    `json:"method"`
	Headers    map[string]string `json:"headers"`
	ExpiresAt  time.Time `json:"expires_at"`
}

// PresignDownloadRequest represents the request to get a presigned download URL
type PresignDownloadRequest struct {
	TenantID  string `json:"tenant_id" binding:"required"`
	ObjectKey string `json:"object_key" binding:"required"`
}

// PresignDownloadResponse represents the response with presigned download URL
type PresignDownloadResponse struct {
	URL       string    `json:"url"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ListRequest represents the request to list objects
type ListRequest struct {
	TenantID string `json:"tenant_id" binding:"required"`
	Prefix   string `json:"prefix,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Marker   string `json:"marker,omitempty"`
}

// ListResponse represents the response with object list
type ListResponse struct {
	Objects   []ObjectInfo `json:"objects"`
	NextMarker string      `json:"next_marker,omitempty"`
	Truncated bool         `json:"truncated"`
}

// ObjectInfo represents basic object information
type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ETag         string    `json:"etag"`
}

// DeleteRequest represents the request to delete objects
type DeleteRequest struct {
	TenantID   string   `json:"tenant_id" binding:"required"`
	ObjectKeys []string `json:"object_keys" binding:"required"`
}

// DeleteResponse represents the response for delete operation
type DeleteResponse struct {
	Deleted []string `json:"deleted"`
	Errors  []string `json:"errors,omitempty"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
