# Include environment variables from .env file
ifneq (,$(wildcard ./.env.development))
	include .env.development
	export
endif

# push docker image to aws ecr
push-docker-image:
	@docker tag zolaris-backend-app-stage 864981729345.dkr.ecr.ap-south-1.amazonaws.com/zolaris-go-app:latest
	@docker push 864981729345.dkr.ecr.ap-south-1.amazonaws.com/zolaris-go-app:latest

# Database migration commands
.PHONY: migrate-up migrate-down migrate-create

# Run migrations up
migrate-up:
	@migrate -path ./internal/db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB_NAME}?sslmode=${POSTGRES_SSL_MODE}" up

# Run migrations down
migrate-down:
	@migrate -path ./internal/db/migrations -database "postgresql://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DB_NAME}?sslmode=${POSTGRES_SSL_MODE}" down

# Create a new migration file
migrate-create:
	@migrate create -ext sql -dir ./migrations -seq $(name)

start-dev:
	@docker compose -f docker-compose.dev.yml up -d

stop-dev:
	@docker compose -f docker-compose.dev.yml down
