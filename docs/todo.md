# Step 1: Infrastructure Setup & Health Checks (`todo.md`)

## Objective
Establish the core infrastructure layer for **OrderPulse** using Docker Compose. Ensure all data stores, messaging brokers, and observability tools are running, healthy, accessible, and correctly networked before writing microservice code.

---

## 📋 Task Checklist

- [x] **1. Project Directory Initialization**
  - Create project root directory: `infra/` folder for infrastructure files.
  - Initialize Git repository and add a `.gitignore` ignoring volumes/temp files.

- [x] **2. Docker Compose Configuration (`infra/docker-compose.yml`)**
  - Create `infra/docker-compose.yml` with the provided services:
    - `postgres:16-alpine` (Port `5432`)
    - `redis:7-alpine` (Port `6379`)
    - `rabbitmq:3-management-alpine` (AMQP `5672`, Management UI `15672`)
    - `apache/kafka:3.7.0` in KRaft mode (Broker `9092`)
    - `grafana/loki:2.9.4` (Port `3100`)
    - `grafana/grafana:10.4.0` (Port `3000`)
  - Ensure volume persistence for Postgres and Redis.
  - Set healthchecks for Postgres, Redis, and RabbitMQ.

- [x] **3. Container Orchestration & Spin Up**
  - Execute `docker compose -f infra/docker-compose.yml up -d` to launch containers in detached mode.
  - Verify all container statuses with `docker compose -f infra/docker-compose.yml ps` and check for any crashing loops.

- [x] **4. Service Connection Verification**
  - **PostgreSQL**: Verify connection via `psql` or database GUI (DBeaver / DataGrip) using `appuser`/`apppassword` on `localhost:5432/orderpulse_db`.
  - **Redis**: Test ping/pong using `redis-cli -h localhost -p 6379 ping`.
  - **RabbitMQ**: Access `http://localhost:15672` in browser (Credentials: `guest` / `guest`).
  - **Kafka**: Test broker availability on `localhost:9092` using a local client tool or `kcat`.
  - **Grafana**: Access `http://localhost:3000` (Credentials: `admin` / `admin`). Configure Loki (`http://loki:3100`) as a default datasource in Grafana settings.

- [x] **5. Environment & Helper Setup**
  - Create a `.env.example` template with default connection strings for Go microservices.
  - Add a simple Makefile or shell script to wrap `docker compose up`, `down`, and `logs`.

---

## 📄 Docker Compose Reference

```yaml
version: '3.8'

services:
  # ---------------------------------------------------------------------------
  # Relational Database
  # ---------------------------------------------------------------------------
  postgres:
    image: postgres:16-alpine
    container_name: orderpulse-postgres
    environment:
      POSTGRES_USER: appuser
      POSTGRES_PASSWORD: apppassword
      POSTGRES_DB: orderpulse_db
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U appuser -d orderpulse_db"]
      interval: 5s
      timeout: 5s
      retries: 5

  # ---------------------------------------------------------------------------
  # Caching & Pub/Sub
  # ---------------------------------------------------------------------------
  redis:
    image: redis:7-alpine
    container_name: orderpulse-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  # ---------------------------------------------------------------------------
  # Task Queue (RabbitMQ)
  # ---------------------------------------------------------------------------
  rabbitmq:
    image: rabbitmq:3-management-alpine
    container_name: orderpulse-rabbitmq
    ports:
      - "5672:5672"   # AMQP protocol
      - "15672:15672" # Web Management Dashboard
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    healthcheck:
      test: ["CMD", "rabbitmq-diagnostics", "-q", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  # ---------------------------------------------------------------------------
  # Event Streaming & Management (Kafka KRaft + Kafbat UI)
  # ---------------------------------------------------------------------------
  kafka:
    image: apache/kafka:3.7.0
    container_name: orderpulse-kafka
    ports:
      - "9092:9092"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: 'broker,controller'
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: 'CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT,PLAINTEXT_HOST:PLAINTEXT'
      KAFKA_LISTENERS: 'PLAINTEXT://:29092,CONTROLLER://:9093,PLAINTEXT_HOST://localhost:9092'
      KAFKA_INTER_BROKER_LISTENER_NAME: 'PLAINTEXT'
      KAFKA_ADVERTISED_LISTENERS: 'PLAINTEXT://kafka:29092,PLAINTEXT_HOST://localhost:9092'
      KAFKA_CONTROLLER_LISTENER_NAMES: 'CONTROLLER'
      KAFKA_CONTROLLER_QUORUM_VOTERS: '1@kafka:9093'
      KAFKA_GROUP_INITIAL_REBALANCE_DELAY_MS: 0
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_LOG_DIRS: '/tmp/kraft-combined-logs'
      KAFKA_CLUSTER_ID: 'MkU3OEVBNTcwNTJENDM2Qk'

  kafka-ui:
    image: ghcr.io/kafbat/kafka-ui:latest
    container_name: orderpulse-kafka-ui
    ports:
      - "8080:8080"
    environment:
      KAFKA_CLUSTERS_0_NAME: orderpulse-cluster
      KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:29092
      DYNAMIC_CONFIG_ENABLED: 'true'
    depends_on:
      - kafka

  # ---------------------------------------------------------------------------
  # Log Aggregation & Visualization (Loki + Grafana)
  # ---------------------------------------------------------------------------
  loki:
    image: grafana/loki:2.9.4
    container_name: orderpulse-loki
    ports:
      - "3100:3100"
    command: -config.file=/etc/loki/local-config.yaml

  grafana:
    image: grafana/grafana:10.4.0
    container_name: orderpulse-grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
    depends_on:
      - loki

volumes:
  postgres_data:
  redis_data:
```

---

## ⚙️ Microservice Connection Matrix (`.env.example`)

```env
# Postgres Configuration
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=appuser
POSTGRES_PASSWORD=apppassword
POSTGRES_DB=orderpulse_db
POSTGRES_DSN="postgres://appuser:apppassword@localhost:5432/orderpulse_db?sslmode=disable"

# Redis Configuration
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=

# RabbitMQ Configuration
RABBITMQ_URL="amqp://guest:guest@localhost:5672/"

# Kafka Configuration
KAFKA_BROKERS=localhost:9092
KAFKA_UI_URL=http://localhost:8080

# Observability Configuration
LOKI_URL="http://localhost:3100/loki/api/v1/push"
```

---

## 🛠️ Developer Verification Commands

Execute these verification checks to confirm the infrastructure readiness:

```bash
# 1. Start all services
docker compose up -d

# 2. Check service status
docker compose ps

# 3. Test PostgreSQL
docker exec -it orderpulse-postgres pg_isready -U appuser -d orderpulse_db

# 4. Test Redis
docker exec -it orderpulse-redis redis-cli ping

# 5. Test RabbitMQ
docker exec -it orderpulse-rabbitmq rabbitmq-diagnostics ping

# 6. Check Loki health endpoint
curl -s http://localhost:3100/ready

# 7. Access Kafbat UI in browser
# http://localhost:8080
```

---

## Definition of Done (DoD)
- All 7 containers are running with `healthy` or active status.
- Port bindings on `5432`, `6379`, `5672`, `15672`, `9092`, `8080`, `3100`, and `3000` are active and responsive.
- Kafbat UI accessible on `http://localhost:8080` connected to `orderpulse-cluster`.
- Grafana is accessible via browser with Loki added as an active datasource.
