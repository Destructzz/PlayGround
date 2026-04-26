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
	swag init -g $(SWAG_MAIN)

sqlc:
	@echo "Generating SQLC code..."
	sqlc generate -f backend/sqlc.yaml

migrate-up:
	@echo "Running migrations up..."
	migrate -path $(MIGRATIONS_PATH) -database "$$DATABASE_URL" -verbose up

migrate-down:
	@echo "Running migrations down..."
	migrate -path $(MIGRATIONS_PATH) -database "$$DATABASE_URL" -verbose down

db-url:
	@echo $$DATABASE_URL