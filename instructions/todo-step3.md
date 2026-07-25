# Step 3: Uber-Style Go Microservice Core (Inventory Service) (`todo-step3.md`)

## Objective
Build the **Inventory Service** as the first microservice using Uber's Go style guidelines (`uber-go/fx` for dependency injection and `uber-go/zap` for structured logging). Implement the gRPC server based on the generated protobuf definitions from Step 2, connected to PostgreSQL and Redis.

---

## 📋 Task Checklist

- [x] **1. Define Project Directory Structure (Uber / Clean Architecture Style)**
  - Create `cmd/inventory/main.go`
  - Create `internal/inventory/config/config.go`
  - Create `internal/inventory/domain/item.go`
  - Create `internal/inventory/repository/postgres.go` & `redis.go`
  - Create `internal/inventory/service/inventory.go`
  - Create `internal/inventory/grpc/handler.go`
  - Create `internal/inventory/module.go` (Fx dependency injection container)

- [x] **2. Setup Uber Fx & Zap Logging Infrastructure**
  - Initialize `zap.Logger` configured for JSON structured logs.
  - Implement environment configuration loading (`config.go`).
  - Wire application modules cleanly using `go.uber.org/fx`.

- [x] **3. Database & Caching Repositories**
  - Implement PostgreSQL repository for reading/updating item stock counts (`inventory_items` table).
  - Implement Redis caching layer for fast stock availability lookups (`CheckStock`).

- [x] **4. Implement gRPC Handler (`grpc/handler.go`)**
  - Implement `inventoryv1.UnimplementedInventoryServiceServer`.
  - Wire `CheckStock` to query Redis cache first, falling back to PostgreSQL.
  - Wire `ReserveStock` to decrement stock inside a PostgreSQL database transaction and invalidate/update the Redis cache.

- [x] **5. gRPC Server Lifecycle Management**
  - Register the gRPC server with Uber Fx lifecycle hooks (`fx.Hook` on `OnStart` and `OnStop`).
  - Expose gRPC service on port `:50051`.

---

## 📁 Directory & File Layout Template

```
orderpulse/
├── cmd/
│   └── inventory/
│       └── main.go
├── internal/
│   └── inventory/
│       ├── config/
│       │   └── config.go
│       ├── domain/
│       │   └── item.go
│       ├── repository/
│       │   ├── postgres.go
│       │   └── redis.go
│       ├── service/
│       │   └── inventory.go
│       ├── grpc/
│       │   └── handler.go
│       └── module.go
├── proto/
│   └── inventory/
│       └── v1/
```

---

## 📄 Code Reference Snippets

### `cmd/inventory/main.go`
```go
package main

import (
	"go.uber.org/fx"
	"orderpulse/internal/inventory"
)

func main() {
	fx.New(
		inventory.Module,
	).Run()
}
```

### `internal/inventory/module.go`
```go
package inventory

import (
	"context"
	"net"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"orderpulse/internal/inventory/config"
	inventorygrpc "orderpulse/internal/inventory/grpc"
	"orderpulse/internal/inventory/repository"
	inventorysvc "orderpulse/internal/inventory/service"
	inventoryv1 "orderpulse/proto/inventory/v1"
)

var Module = fx.Options(
	fx.Provide(
		config.Load,
		zap.NewProduction,
		repository.NewPostgresRepository,
		repository.NewRedisRepository,
		inventorysvc.NewInventoryService,
		inventorygrpc.NewHandler,
	),
	fx.Invoke(RegisterGRPCServer),
)

func RegisterGRPCServer(
	lc fx.Lifecycle,
	logger *zap.Logger,
	cfg *config.Config,
	handler *inventorygrpc.Handler,
) {
	server := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(server, handler)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			lis, err := net.Listen("tcp", cfg.GRPCPort)
			if err != nil {
				return err
			}
			logger.Info("Starting Inventory gRPC Server", zap.String("port", cfg.GRPCPort))
			go server.Serve(lis)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Inventory gRPC Server")
			server.GracefulStop()
			return nil
		},
	})
}
```

---

## 🛠️ Verification Commands

```bash
# 1. Install Uber Fx and Zap dependencies
go get go.uber.org/fx go.uber.org/zap

# 2. Run the Inventory Service locally
go run cmd/inventory/main.go

# 3. Test gRPC endpoint using grpcurl
grpcurl -plaintext -d '{"product_id": "prod-101", "quantity": 2}'   localhost:50051 inventory.v1.InventoryService/CheckStock
```

---

## Definition of Done (DoD)
- Clean package structure adhering to Uber Go style guide and clean architecture.
- Service bootstraps using `uber-go/fx` with graceful shutdown hooks.
- Structured logging provided via `uber-go/zap`.
- gRPC server starts on port `:50051` and successfully responds to `CheckStock` and `ReserveStock` calls.
