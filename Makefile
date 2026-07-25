.PHONY: up down status logs check proto

COMPOSE_FILE = infra/docker-compose.yml

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

check:
	@echo "Checking Postgres..."
	@docker exec orderpulse-postgres pg_isready -U appuser -d orderpulse_db
	@echo "Checking Redis..."
	@docker exec orderpulse-redis redis-cli ping
	@echo "Checking RabbitMQ..."
	@docker exec orderpulse-rabbitmq rabbitmq-diagnostics ping
	@echo "Checking Loki..."
	@curl -s http://localhost:3100/ready
