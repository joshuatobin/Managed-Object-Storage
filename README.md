# Managed Object Storage API - MVP Mock

A simple Go-based API server for multi-tenant object storage using mock S3 responses. This is an MVP implementation for testing and demonstration purposes.

## Quick Start

```bash
# Build and run
make build
make run

# Test endpoints
make health
make curl
```

## API Endpoints

- `GET /healthz` - Health check
- `POST /v1/presign/upload` - Get mock upload URL
- `POST /v1/presign/download` - Get mock download URL  
- `POST /v1/list` - List mock objects
- `POST /v1/delete` - Delete mock objects

## Example Usage

```bash
# Health check
curl http://localhost:8080/healthz

# Get upload URL
curl -X POST http://localhost:8080/v1/presign/upload \
  -H "Content-Type: application/json" \
  -d '{"tenant_id": "test-tenant", "object_key": "file.txt", "content_type": "text/plain"}'
```

## Make Targets

- `make build` - Build the server
- `make run` - Start the server
- `make test` - Run tests
- `make health` - Quick health check
- `make curl` - Test all endpoints
- `make clean` - Clean build artifacts

## Note

This is a **mock implementation** that returns fake URLs and data. No actual S3 calls are made.
