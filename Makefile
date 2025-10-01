.PHONY: build run test clean curl

# Build the application
build:
	go build -o bin/server cmd/server/main.go

# Run the application
run: build
	./bin/server

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Install dependencies
deps:
	go mod tidy

# Test API endpoints with curl
curl:
	@echo "=== Testing Health Check ==="
	@curl -s http://localhost:8080/healthz
	@echo -e "\n"
	@echo "=== Testing Presign Upload ==="
	@curl -s -X POST http://localhost:8080/v1/presign/upload \
		-H "Content-Type: application/json" \
		-d '{"tenant_id": "test-tenant", "object_key": "test-file.txt", "content_type": "text/plain", "max_size": 1048576}'
	@printf "\n"
	@echo "=== Testing Presign Download ==="
	@curl -s -X POST http://localhost:8080/v1/presign/download \
		-H "Content-Type: application/json" \
		-d '{"tenant_id": "test-tenant", "object_key": "test-file.txt"}'
	@printf "\n"
	@echo "=== Testing List Objects ==="
	@curl -s -X POST http://localhost:8080/v1/list \
		-H "Content-Type: application/json" \
		-d '{"tenant_id": "test-tenant", "prefix": "", "limit": 10}'
	@printf "\n"
	@echo "=== Testing Delete Objects ==="
	@curl -s -X POST http://localhost:8080/v1/delete \
		-H "Content-Type: application/json" \
		-d '{"tenant_id": "test-tenant", "object_keys": ["test-file.txt"]}'
	@printf "\n"
	@echo "=== Test Complete ==="

# Health check endpoint
health:
	@echo "=== Health Check ==="
	@curl -s http://localhost:8080/healthz
	@echo -e "\n"
