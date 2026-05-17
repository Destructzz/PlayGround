-include .env
export

SWAG_MAIN=backend/cmd/api/main.go
MIGRATIONS_PATH=./backend/migrations
SQLC_PATH=./backend/

DB_HOST ?= localhost
DB_PORT ?= 5432

DATABASE_URL=postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(POSTGRES_DB)?sslmode=disable

.PHONY: all swag migrate-up migrate-down sqlc db-url

# 🔥 команда "всё сразу"
all: swag sqlc migrate-up
	@echo "✅ All done"

swag:
	@echo "Generating Swagger..."
	swag init -g $(SWAG_MAIN) -o backend/docs

sqlc:
	@echo "Generating SQLC code..."
	sqlc generate -f backend/sqlc.yaml

migrate-up:
	@echo "Running migrations up..."
	migrate -path $(MIGRATIONS_PATH) -database "$$DATABASE_URL" -verbose up

migrate-down:
	@echo "Running migrations down..."
	migrate -path $(MIGRATIONS_PATH) -database "$$DATABASE_URL" -verbose down

COMPOSE_FILE ?= docker-compose.yml

db-url:
	@echo $$DATABASE_URL

backup-export:
	@echo "Exporting database backup..."
	@mkdir -p backup
	docker compose -f $(COMPOSE_FILE) exec -T postgres pg_dump -c --if-exists -U $(POSTGRES_USER) -d $(POSTGRES_DB) > backup/backup_$$(date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup exported successfully to backup/ folder"

backup-import:
	@if [ -z "$(FILENAME)" ]; then \
		echo "❌ Error: FILENAME is required. Example: make backup-import FILENAME=backup_20260517_184000.sql"; \
		exit 1; \
	fi
	@if [ ! -f "backup/$(FILENAME)" ]; then \
		echo "❌ Error: File backup/$(FILENAME) does not exist."; \
		exit 1; \
	fi
	@echo "Importing database backup from backup/$(FILENAME)..."
	docker compose -f $(COMPOSE_FILE) exec -T postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB) < backup/$(FILENAME)
	@echo "✅ Backup imported successfully from backup/$(FILENAME)"