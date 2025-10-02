# Managed Object Storage API - MVP Mock

A simple Go-based API server that mocks a multi-tenant object storage control-plane. It returns realistic, stubbed S3-style presigned URLs and object listings without contacting AWS. Intended for design validation and a live demo.

## Quick Start (Runnable)
```bash
make build
make run            # starts on :8080
# In another terminal
make health         # quick ping
make curl           # exercise all endpoints
```

## API Endpoints
- GET `/healthz` – Health check
- POST `/v1/presign/upload` – Mock upload URL
- POST `/v1/presign/download` – Mock download URL
- POST `/v1/list` – Mock object listing
- POST `/v1/delete` – Mock delete

## Demo Scope
- This is a mock control-plane: no AWS calls, no auth. Responses are deterministic for demo/testing.
- Basic input validation and tenant-scoped key formatting mirror the ADR.

## Why this implementation
- Validate the ADR’s control-plane API quickly without cloud setup/creds.
- Demonstrate presigned-URL flow and tenant prefixing while keeping the control-plane out of the data path.
- Keep dependencies minimal for a reliable live demo.

## Considerations, decisions, assumptions
- Decision: Single bucket + tenant prefixes and presigned URLs (mocked) per ADR MVP.
- Assumption: AuthN/Z (tenant/role) is out-of-scope for the mock; would be enforced before presign in production.
- Consideration: Responses mimic S3 presigned query params and TTLs to ease client integration later.
- Non-goals (here): No AWS SDK/KMS/ACLs/RBAC; focused on API contract only.

## V2 improvements 
- Replace mocks with real S3 presign/list/delete
  - Challenge: IAM least-privilege, SSE-KMS, VPC endpoints, retries, error surfaces.
- Add AuthN/Z (tokens or OIDC → `reader|writer|admin`)
  - Challenge: tenant membership mapping, per-prefix enforcement, rate limits, audit trails.
- Multipart uploads + constraints (size, content-type, checksum)
  - Challenge: client compat, enforcing constraints at S3, validating from control-plane.
- Quotas and cost attribution per tenant
  - Challenge: metering (bytes/objects), inventory reconciliation, backfill/drift.
- Observability & operations
  - Challenge: metrics, alerting, CloudTrail data event costs, DR and incident response.

## Make Targets
- `make build` – Build the server
- `make run` – Start the server
- `make test` – Run tests
- `make health` – Quick health check
- `make curl` – Exercise all endpoints
- `make clean` – Remove build artifacts

## Note
This is a mock implementation for MVP demo purposes. It returns fake URLs/data and does not contact AWS.

## Workflows

- Presign Upload (source: docs/mermaid/presign_upload.mmd)
```mermaid
sequenceDiagram
  autonumber
  participant Client
  participant API as Control-Plane API
  participant Auth as AuthNZ
  participant S3 as S3 (bucket)

  Client->>API: POST /v1/presign/upload {tenant_id, object_key, content_type, max_size}
  API->>Auth: Validate token -> (tenant_id, role)
  Auth-->>API: OK (role permits PUT)
  API->>API: Normalize and validate key: tenants/TENANT_ID/OBJECT_KEY
  API->>S3: Generate presigned PUT URL (TTL <= 10m)
  S3-->>API: {url, headers, expires_at}
  API-->>Client: 200 {url, method: PUT, headers, expires_at}
```

- Presign Download (source: docs/mermaid/presign_download.mmd)
```mermaid
sequenceDiagram
  autonumber
  participant Client
  participant API as Control-Plane API
  participant Auth as AuthNZ
  participant S3 as S3 (bucket)

  Client->>API: POST /v1/presign/download {tenant_id, object_key}
  API->>Auth: Validate token -> (tenant_id, role)
  Auth-->>API: OK (role permits GET/HEAD)
  API->>API: Normalize and validate key: tenants/TENANT_ID/OBJECT_KEY
  API->>S3: Generate presigned GET URL (TTL <= 10m)
  S3-->>API: {url, expires_at}
  API-->>Client: 200 {url, expires_at}
```

- List Objects (source: docs/mermaid/list_objects.mmd)
```mermaid
sequenceDiagram
  autonumber
  participant Client
  participant API as Control-Plane API
  participant Auth as AuthNZ
  participant S3 as S3 (bucket)

  Client->>API: POST /v1/list {tenant_id, prefix, limit, marker}
  API->>Auth: Validate token -> (tenant_id, role)
  Auth-->>API: OK (role permits list)
  API->>API: Build prefix: tenants/TENANT_ID/PREFIX
  API->>S3: ListObjectsV2(Prefix, MaxKeys, ContinuationToken)
  S3-->>API: {Contents[], NextContinuationToken, IsTruncated}
  API-->>Client: 200 {objects[], next_marker, truncated}
```

- Delete Objects (source: docs/mermaid/delete_objects.mmd)
```mermaid
sequenceDiagram
  autonumber
  participant Client
  participant API as Control-Plane API
  participant Auth as AuthNZ
  participant S3 as S3 (bucket)

  Client->>API: POST /v1/delete {tenant_id, object_keys[]}
  API->>Auth: Validate token -> (tenant_id, role)
  Auth-->>API: OK (role permits delete)
  API->>API: Validate keys and map to tenants/TENANT_ID/KEY
  API->>S3: DeleteObjects(Identifiers[])
  S3-->>API: {Deleted[], Errors[]}
  API-->>Client: 200 {deleted[], errors[]}
```
