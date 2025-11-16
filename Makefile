.PHONY: help build up start down stop restart logs ps test-integration

help:
	@echo "Available commands:"
	@echo "  make build    - Build containers"
	@echo "  make up       - Start containers in background"
	@echo "  make start    - Start containers"
	@echo "  make down     - Stop and remove containers"
	@echo "  make stop     - Stop containers"
	@echo "  make restart  - Restart containers"
	@echo "  make logs     - Show logs (follow mode)"
	@echo "  make ps       - Show container status"
	@echo ""
	@echo "Add service name: make up c=service_name"

build:
	docker compose -f docker-compose.yml build $(c)

up:
	docker compose -f docker-compose.yml up -d $(c)

start:
	docker compose -f docker-compose.yml start $(c)

down:
	docker compose -f docker-compose.yml down $(c)

stop:
	docker compose -f docker-compose.yml stop $(c)

restart:
	docker compose -f docker-compose.yml stop $(c)
	docker compose -f docker-compose.yml up -d $(c)

logs:
	docker compose -f docker-compose.yml logs --tail=100 -f $(c)

ps:
	docker compose -f docker-compose.yml ps