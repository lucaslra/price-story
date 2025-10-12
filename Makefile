.PHONY: build run test clean db-up db-down db-logs db-reset migrate-up migrate-down migrate-status \
        docker-build docker-up docker-down docker-logs docker-restart \
        api-test api-test-users api-test-products api-test-price-stories api-test-price-points \
        fmt lint env openapi-sync-check \
        frontend-install frontend-dev frontend-build frontend-test frontend-lint frontend-preview

# Build the application
build:
	go build -o ./bin/price-story ./cmd/price-story

# Run the application with PostgreSQL
run: db-up
	@echo "Waiting for PostgreSQL to be ready..."
	@until docker-compose exec -T postgres pg_isready -U price_story_user -d price_story >/dev/null 2>&1; do \
		sleep 2; \
	done
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=price_story_user \
	DB_PASSWORD=price_story_password \
	DB_NAME=price_story \
	DB_SSLMODE=disable \
	go run ./cmd/price-story/main.go


# Container-first development
docker-build:
	docker-compose build app

docker-up:
	docker-compose up -d postgres app frontend

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f app

docker-restart:
	docker-compose down && docker-compose up -d postgres app frontend


# Run tests
test:
	go test -v ./...

# Validate OpenAPI spec synchronization with router
openapi-sync-check:
	go test -v ./pkg/openapi -run TestOpenAPI

# Clean build artifacts
clean:
	rm -rf ./bin

# Start PostgreSQL database
db-up:
	docker-compose up -d postgres

# Stop PostgreSQL database
db-down:
	docker-compose down

# View database logs
db-logs:
	docker-compose logs -f postgres

# Reset database (stop, remove volumes, and start fresh)
db-reset:
	docker-compose down -v
	docker-compose up -d postgres

# Run database migrations up
migrate-up:
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=price_story_user \
	DB_PASSWORD=price_story_password \
	DB_NAME=price_story \
	DB_SSLMODE=disable \
	go run ./cmd/price-story/main.go -migrate up

# Rollback last migration
migrate-down:
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=price_story_user \
	DB_PASSWORD=price_story_password \
	DB_NAME=price_story \
	DB_SSLMODE=disable \
	go run ./cmd/price-story/main.go -migrate down

# Show migration status
migrate-status:
	DB_HOST=localhost \
	DB_PORT=5432 \
	DB_USER=price_story_user \
	DB_PASSWORD=price_story_password \
	DB_NAME=price_story \
	DB_SSLMODE=disable \
	go run ./cmd/price-story/main.go -migrate status

# Default target
all: clean build

# Formatting
fmt:
	go fmt ./...

# Linting (requires golangci-lint installed)
lint:
	golangci-lint run

# Show effective environment (if set)
env:
	@printenv | grep -E '^(DB_HOST|DB_PORT|DB_USER|DB_PASSWORD|DB_NAME|DB_SSLMODE|PORT|READ_TIMEOUT|WRITE_TIMEOUT)$$' | sort || true

# Run all API endpoint tests
api-test:
	@bash ./scripts/api-tests/run_all.sh

# Run domain-specific API test scripts
api-test-users:
	@bash ./scripts/api-tests/users.sh

api-test-products:
	@bash ./scripts/api-tests/products.sh

api-test-price-stories:
	@bash ./scripts/api-tests/price_stories.sh

api-test-price-points:
	@bash ./scripts/api-tests/price_points.sh

# --- Frontend helpers ---
frontend-install:
	cd frontend && npm install

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-test:
	cd frontend && npm test

frontend-lint:
	cd frontend && npm run lint

frontend-preview:
	cd frontend && npm run preview