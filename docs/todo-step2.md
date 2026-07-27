# Step 2: Protobufs & gRPC Code Generation (`todo-step2.md`)

## Objective
Define the gRPC service contracts using Protocol Buffers (`.proto`) for internal communication between the **Order Service** and **Inventory Service**. Set up the code generation pipeline to compile `.proto` files into usable Go code.

---

## 📋 Task Checklist

- [x] **1. Create Protocol Buffer Directory Layout**
  - Create directory: `proto/inventory/v1/`
  - Create directory: `proto/order/v1/`

- [x] **2. Define Inventory Service Spec (`proto/inventory/v1/inventory.proto`)**
  - Create protobuf file with `syntax = "proto3";`.
  - Define package `inventory.v1`.
  - Set Go package option: `option go_package = "github.com/orderpulse/proto/inventory/v1;inventoryv1";`.
  - Add `CheckStockRequest` (`product_id`, `quantity`) & `CheckStockResponse` (`is_available`, `available_quantity`, `price`).
  - Add `ReserveStockRequest` (`order_id`, `product_id`, `quantity`) & `ReserveStockResponse` (`success`, `message`).
  - Declare service `InventoryService` with `CheckStock` and `ReserveStock` RPCs.

- [x] **3. Define Order Service Spec (`proto/order/v1/order.proto`)**
  - Create protobuf file with `syntax = "proto3";`.
  - Define package `order.v1`.
  - Set Go package option: `option go_package = "github.com/orderpulse/proto/order/v1;orderv1";`.
  - Add `CreateOrderRequest` (`user_id`, `product_id`, `quantity`) & `CreateOrderResponse` (`order_id`, `status`, `total_amount`).
  - Add `GetOrderStatusRequest` (`order_id`) & `GetOrderStatusResponse` (`order_id`, `status`, `created_at`).
  - Declare service `OrderService` with `CreateOrder` and `GetOrderStatus` RPCs.

- [x] **4. Code Generation Pipeline**
  - Verify local installation of `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc`.
  - Add a Makefile target or shell generation script (`scripts/gen-proto.sh`) to compile proto files:
    ```bash
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/inventory/v1/inventory.proto
    
    protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/order/v1/order.proto
    ```

- [x] **5. Verification & Tidying**
  - Execute generation script and verify created `*.pb.go` and `*_grpc.pb.go` files inside `proto/inventory/v1/` and `proto/order/v1/`.
  - Initialize Go module if not present (`go mod init orderpulse`) and run `go mod tidy` to resolve gRPC dependencies.

---

## 📄 Protocol Buffer Specification Templates

### `proto/inventory/v1/inventory.proto`
```protobuf
syntax = "proto3";

package inventory.v1;

option go_package = "github.com/orderpulse/proto/inventory/v1;inventoryv1";

service InventoryService {
  rpc CheckStock (CheckStockRequest) returns (CheckStockResponse);
  rpc ReserveStock (ReserveStockRequest) returns (ReserveStockResponse);
}

message CheckStockRequest {
  string product_id = 1;
  int32 quantity = 2;
}

message CheckStockResponse {
  bool is_available = 1;
  int32 available_quantity = 2;
  double price = 3;
}

message ReserveStockRequest {
  string order_id = 1;
  string product_id = 2;
  int32 quantity = 3;
}

message ReserveStockResponse {
  bool success = 1;
  string message = 2;
}
```

### `proto/order/v1/order.proto`
```protobuf
syntax = "proto3";

package order.v1;

option go_package = "github.com/orderpulse/proto/order/v1;orderv1";

service OrderService {
  rpc CreateOrder (CreateOrderRequest) returns (CreateOrderResponse);
  rpc GetOrderStatus (GetOrderStatusRequest) returns (GetOrderStatusResponse);
}

message CreateOrderRequest {
  string user_id = 1;
  string product_id = 2;
  int32 quantity = 3;
}

message CreateOrderResponse {
  string order_id = 1;
  string status = 2;
  double total_amount = 3;
}

message GetOrderStatusRequest {
  string order_id = 1;
}

message GetOrderStatusResponse {
  string order_id = 1;
  string status = 2;
  string created_at = 3;
}
```

---

## 🛠️ Verification Commands

```bash
# 1. Run protobuf generation
chmod +x scripts/gen-proto.sh && ./scripts/gen-proto.sh

# 2. Check that generated files exist
ls -la proto/inventory/v1/*.pb.go
ls -la proto/order/v1/*.pb.go

# 3. Tidy Go dependencies
go mod tidy
```

---

## Definition of Done (DoD)
- Clean `v1` versioned directory structures containing `.proto` files.
- Successful code generation producing `*.pb.go` and `*_grpc.pb.go` files without syntax/import errors.
- `go mod tidy` resolves all Go gRPC & Protobuf dependencies cleanly.
