include .env
.PHONY: create-migration up down down-to-zero down-to up-from

# Set the directory containing migration files
MIGRATIONS_DIR = internal/migrations

# Create a new migration file using goose command
create-migration:
	atlas migrate diff --env gorm

# Set the command to run goose
GOOSE_CMD = goose -dir $(MIGRATIONS_DIR) sqlite3 $(SQLITE_DB_PATH)

# Display the migration status
status:
	$(GOOSE_CMD) status

# Apply all pending migrations
up:
	atlas migrate apply --env gorm

# Rollback the last migration
down:
	atlas migrate down --env gorm

# Rollback all migrations
down-to-zero:
	$(GOOSE_CMD) down-to 0

# Rollback to a specific migration version
down-to:
	$(GOOSE_CMD) down-to $(filter-out $@,$(MAKECMDGOALS))

# Run migrations from a specific migration version
up-from:
	$(GOOSE_CMD) up-by-one $(filter-out $@,$(MAKECMDGOALS))
