include .env
.PHONY: create-migration up down down-to-zero down-to up-from

# Set the directory containing migration files
MIGRATIONS_DIR = internal/migrations

# Create a new migration file using goose command
create-migration:
	goose -dir $(MIGRATIONS_DIR) create $(filter-out $@,$(MAKECMDGOALS)) sql

# Set the command to run goose
GOOSE_CMD = goose -dir $(MIGRATIONS_DIR) sqlite3 $(SQLITE_DB_PATH)

# Apply all pending migrations
up:
	$(GOOSE_CMD) up

# Rollback the last migration
down:
	$(GOOSE_CMD) down

# Rollback all migrations
down-to-zero:
	$(GOOSE_CMD) down-to 0

# Rollback to a specific migration version
down-to:
	$(GOOSE_CMD) down-to $(filter-out $@,$(MAKECMDGOALS))

# Run migrations from a specific migration version
up-from:
	$(GOOSE_CMD) up-by-one $(filter-out $@,$(MAKECMDGOALS))
