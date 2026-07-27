.PHONY: up down status logs check proto build tidy fmt test run-inventory run-order run-worker run-frontend check-loki check-kafka-ui

COMPOSE_FILE = infra/docker-compose.yml

# Go lives in two workspace modules (see go.work); ./... from the repo root
# matches neither, so Go targets address them explicitly.
GO_PKGS = ./backend/... ./proto/...

up:
	docker compose -f $(COMPOSE_FILE) up -d

down:
	docker compose -f $(COMPOSE_FILE) down

status:
	docker compose -f $(COMPOSE_FILE) ps

logs:
	docker compose -f $(COMPOSE_FILE) logs -f

proto:
	./scripts/gen-proto.sh

build:
	go build $(GO_PKGS)

tidy:
	cd proto && go mod tidy
	cd backend && go mod tidy

fmt:
	go fmt $(GO_PKGS)

test:
	go test $(GO_PKGS)

run-inventory:
	go run ./backend/cmd/inventory

run-order:
	go run ./backend/cmd/order

run-worker:
	go run ./backend/cmd/worker

run-frontend:
	cd frontend && npm run dev

check-loki:
	@echo "Checking Loki ready status:"
	@curl -s http://localhost:3100/ready
	@echo "\nChecking Loki indexed service_name labels:"
	@curl -s http://localhost:3100/loki/api/v1/label/service_name/values

check-kafka-ui:
	@echo "Checking Kafbat UI status..."
	@curl -s -I http://localhost:8080 | head -n 1

check:
	@echo "Checking Postgres..."
	@docker exec orderpulse-postgres pg_isready -U appuser -d orderpulse_db
	@echo "Checking Redis..."
	@docker exec orderpulse-redis redis-cli ping
	@echo "Checking RabbitMQ..."
	@docker exec orderpulse-rabbitmq rabbitmq-diagnostics ping
	@echo "Checking Loki..."
	@curl -s http://localhost:3100/ready
	@echo "Checking Kafbat UI..."
	@curl -s -I http://localhost:8080 | head -n 1
