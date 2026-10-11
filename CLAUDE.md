# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go microservices monorepo for a ride-hailing (代驾) platform, built on [Kratos v2](https://go-kratos.dev/) (v2.8.3). Each service is an independent Go module with its own `go.mod`, dual gRPC + HTTP transports, Protobuf API definitions, and Wire dependency injection.

## Build & Development Commands

All Makefile targets run from **within a service directory** (e.g., `cd backend/customer`):

```bash
# Install protoc plugins and Wire (one-time setup)
make init

# Generate proto stubs (API + config + errors) and Wire code
make all          # api + config + errors + generate (Wire)
make api          # proto → Go + gRPC + HTTP + OpenAPI stubs
make config       # internal conf.proto → Go struct
make errors       # proto → Kratos error codes
make generate     # go generate ./... (runs Wire)

# Build
make build        # outputs to ./bin/

# Run (from service dir)
kratos run

# Tests
go test ./...
go test -run TestName ./path/to/package
```

Prereqs: `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`, `protoc-gen-go-http` (Kratos), `protoc-gen-go-errors` (Kratos), `protoc-gen-openapiv2`, `wire`.

## Architecture

### Repo Layout

```
backend/
  customer/          # Customer service (HTTP :8100, gRPC :9100)
  verifyCode/        # Verification code service (HTTP :8000, gRPC :9000)
driver/              # Driver service (placeholder)
manager/             # Manager service (placeholder)
consumer/            # Consumer service (placeholder)
```

Each service under `backend/` follows the standard Kratos layout:

```
cmd/<service>/       # main.go + wire.go (entry point, DI injector)
configs/             # config.yaml (loaded via -conf flag)
api/                 # Proto definitions + generated Go stubs (service-local)
internal/
  conf/              # Config proto → generated Bootstrap/Server/Data structs
  server/            # HTTP + gRPC server creation, middleware, route registration
  service/           # Adapts proto RPCs to biz layer calls
  biz/               # Business logic, defines repo interfaces
  data/              # Implements repo interfaces (DB, cache)
third_party/         # Shared proto imports (google/api, validate, openapi, errors)
```

### Clean Architecture Flow

`proto RPC → service → biz → data`

- **biz** defines repository interfaces; **data** implements them. Biz never imports data.
- **service** translates between proto types and biz domain types.
- Each layer exports a `wire.ProviderSet`; `wire.go` combines them all.

### Inter-Service Communication

The customer service holds a `verifyCode` proto client definition (`backend/customer/api/verifyCode/`) for calling the verifyCode service via gRPC. This is the pattern for cross-service calls: import the target service's proto as a gRPC client.

### Config

Defined as Protobuf (`internal/conf/conf.proto`), deserialized from YAML at runtime. The `Bootstrap` message contains `Server` (HTTP/gRPC addr + timeout) and `Data` (database + Redis connection info).

### Adding a New Service

1. Copy an existing service directory (e.g., `backend/customer`) as template
2. Update `go.mod` module name
3. Define API protos in `api/`
4. Implement layers: service → biz → data
5. Register handlers in `server/http.go` and `server/grpc.go`
6. Run `make all` to generate stubs and Wire code
