# OrderPulse

An event-driven order pipeline built as a Go microservices monorepo: gRPC between
services, Kafka for events, RabbitMQ for background work, and Grafana/Loki for
centralized structured logs. A Next.js app fronts it as a BFF.

## Layout

```
orderpulse/
├── go.work              Go workspace tying the two Go modules together
├── backend/             module github.com/orderpulse/backend
│   ├── cmd/             service entrypoints (inventory, order, worker)
│   ├── internal/        per-service domain / service / repository / transport
│   └── pkg/logger/      Zap logger with a Loki HTTP sink
├── proto/               module github.com/orderpulse/proto
│   ├── inventory/v1/    .proto contract + generated Go stubs
│   └── order/v1/
├── frontend/            Next.js 14 App Router BFF + UI
├── infra/               docker-compose + Grafana provisioning
├── scripts/             gen-proto.sh
└── docs/                step-by-step build instructions and checklists
```

`proto/` sits at the root rather than inside `backend/` because both sides consume
it: Go imports the generated stubs, and the Next.js route handlers load the raw
`.proto` files at runtime through `@grpc/proto-loader`.

## Architecture

```
Browser ──HTTP──▶ Next.js BFF (:3001) ──gRPC──▶ order-service (:50052)
                                                      │
                                       gRPC ──────────┼──▶ inventory-service (:50051) ──▶ Postgres + Redis
                                                      ├──▶ Kafka  "order-events"
                                                      └──▶ RabbitMQ "invoice_queue" ──▶ worker-service

all three services ──Zap JSON──▶ Loki (:3100) ──▶ Grafana (:3000)
```

## Prerequisites

- Go 1.26+
- Docker + Docker Compose
- Node.js 18+
- `protoc` (optional — `make proto` falls back to a Docker container)

## Quickstart

```bash
make up            # start infrastructure
make check         # verify every container answers

make run-inventory # terminal 1
make run-order     # terminal 2
make run-worker    # terminal 3
make run-frontend  # terminal 4  -> http://localhost:3001
```

Copy `.env.example` to `.env` first if you need to override any defaults; every
service falls back to the localhost values baked into its `config` package.

## Ports

| Service          | Port          |
| ---------------- | ------------- |
| Frontend         | 3001          |
| Grafana          | 3000          |
| Kafbat UI        | 8080          |
| Loki             | 3100          |
| inventory-service| 50051 (gRPC)  |
| order-service    | 50052 (gRPC)  |
| Postgres         | 5432          |
| Redis            | 6379          |
| RabbitMQ         | 5672 / 15672  |
| Kafka            | 9092          |

## Make targets

| Target | What it does |
| ------ | ------------ |
| `up` / `down` / `status` / `logs` | Docker Compose lifecycle |
| `check` | Ping every piece of infrastructure |
| `build` / `test` / `fmt` / `tidy` | Go workspace tasks across both modules |
| `proto` | Regenerate Go stubs from `.proto` |
| `run-*` | Run a single service or the frontend |

Because the Go code lives in two workspace modules, `./...` from the repo root
matches nothing — use `go build ./backend/... ./proto/...` (or just `make build`).

## Conventions

Go code follows the [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md);
project-specific rules live in [`.agents/AGENTS.md`](.agents/AGENTS.md).
