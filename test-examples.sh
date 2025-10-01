#!/bin/bash

# Test examples for the Managed Object Storage API
# Make sure the server is running: make run

BASE_URL="http://localhost:8080"

echo "=== Testing Health Check ==="
curl -s "$BASE_URL/healthz" | jq '.'

echo -e "\n=== Testing Presign Upload ==="
curl -s -X POST "$BASE_URL/v1/presign/upload" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "test-tenant",
    "object_key": "test-file.txt",
    "content_type": "text/plain",
    "max_size": 1048576
  }' | jq '.'

echo -e "\n=== Testing Presign Download ==="
curl -s -X POST "$BASE_URL/v1/presign/download" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "test-tenant",
    "object_key": "test-file.txt"
  }' | jq '.'

echo -e "\n=== Testing List Objects ==="
curl -s -X POST "$BASE_URL/v1/list" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "test-tenant",
    "prefix": "",
    "limit": 10
  }' | jq '.'

echo -e "\n=== Testing Delete Objects ==="
curl -s -X POST "$BASE_URL/v1/delete" \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "test-tenant",
    "object_keys": ["test-file.txt"]
  }' | jq '.'

echo -e "\n=== Test Complete ==="
