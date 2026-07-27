# Step 4: Order Service & Asynchronous Messaging (Kafka + RabbitMQ) (`todo-step4.md`)

## Objective
Build the **Order Service** and **Worker Service**. The Order Service will accept order creation requests, verify/reserve stock with the Inventory Service via **gRPC**, publish an `OrderPlaced` event to **Kafka** for analytics, and enqueue an async `ProcessInvoice` task into **RabbitMQ** for processing by the Worker Service.

---

## 📋 Task Checklist

- [x] **1. Define Directory Layout for Order & Worker Services**
  - Create `cmd/order/main.go` & `cmd/worker/main.go`
  - Create `internal/order/` (config, domain, service, grpc, kafka, rabbitmq, module.go)
  - Create `internal/worker/` (config, consumer, module.go)

- [x] **2. Kafka Producer & RabbitMQ Publisher Setup in Order Service**
  - Setup Kafka producer (`segmentio/kafka-go`) to publish `order-events` to Kafka.
  - Setup RabbitMQ client (`amqp091-go`) to declare exchange `order_exchange` and publish tasks to queue `invoice_queue`.

- [x] **3. Implement Order Creation Workflow (`OrderService`)**
  - Receive `CreateOrder` request via gRPC.
  - Call **Inventory Service** via **gRPC** (`ReserveStock`).
  - Write initial Order record into **PostgreSQL** (`orders` table).
  - Emits `OrderPlaced` JSON payload event to **Kafka** topic `order-events`.
  - Enqueues `GenerateInvoiceTask` payload into **RabbitMQ** queue `invoice_queue`.

- [x] **4. Build Background Worker Service (`cmd/worker/main.go`)**
  - Connect to **RabbitMQ** and subscribe as a consumer on `invoice_queue`.
  - Parse incoming task messages, simulate PDF invoice generation/email sending, and acknowledge (`Ack`) messages upon completion.
  - Wire worker lifecycle into `uber-go/fx`.

- [x] **5. Logging & Traceability**
  - Attach structured `uber-go/zap` logs across Order Service and Worker Service, passing a unified `trace_id` / `order_id`.

---

## 📁 Directory Layout Template

```
orderpulse/
├── cmd/
│   ├── order/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/
│   ├── order/
│   │   ├── config/
│   │   ├── grpc/
│   │   ├── kafka/
│   │   ├── rabbitmq/
│   │   ├── service/
│   │   └── module.go
│   └── worker/
│       ├── config/
│       ├── consumer/
│       └── module.go
```

---

## 📄 Service Interaction Flow

```
[ Next.js / API Client ]
        │
        ▼
[ Order Service ]
   ├── 1. gRPC Call ─────────► [ Inventory Service ]
   ├── 2. Write Order ────────► [ PostgreSQL ]
   ├── 3. Publish Event ──────► [ Kafka: order-events ]
   └── 4. Enqueue Task ───────► [ RabbitMQ: invoice_queue ]
                                       │
                                       ▼
                                [ Worker Service ]
```

---

## 🛠️ Verification Commands

```bash
# 1. Start Inventory Service (from Step 3)
go run cmd/inventory/main.go

# 2. Start Order Service and Worker Service
go run cmd/order/main.go
go run cmd/worker/main.go

# 3. Trigger CreateOrder via gRPC / grpcurl
grpcurl -plaintext -d '{"user_id": "usr-1", "product_id": "prod-101", "quantity": 1}'   localhost:50052 order.v1.OrderService/CreateOrder

# 4. Verify Kafka event emission
docker exec -it orderpulse-kafka kafka-console-consumer.sh   --bootstrap-server localhost:9092 --topic order-events --from-beginning

# 5. Check RabbitMQ Management UI at http://localhost:15672 to verify queue processing
```

---

## Definition of Done (DoD)
- Order Service successfully calls Inventory Service over gRPC.
- Order Service publishes events to Kafka (`order-events`) and enqueues jobs into RabbitMQ (`invoice_queue`).
- Worker Service consumes messages from RabbitMQ, logs structured output via `zap`, and acknowledges task completion.
- All services bootstrap via `uber-go/fx`.
