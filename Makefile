include .env

.PHONY: create-migration up down reset setup-db


# Create a new migration file using atlas migrate diff command
create-migration:
	atlas migrate diff --env gorm

# Apply all pending migrations
up:
	atlas schema apply --env gorm

# Rollback the last migration
down:
	atlas schema migrate down --env gorm

# Reset the database by dropping all objects and re-applying all migrations
reset:
	atlas schema clean --env gorm
	atlas schema apply --env gorm

# Setup database (create it if not exists, useful for local development)
setup-db:
	atlas db create --env gorm
	atlas schema apply --env gorm

inspect:
	atlas schema inspect --env gorm