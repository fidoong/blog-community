.PHONY: all build dev test lint fmt ent wire docker-up docker-down

# Default target
all: build

# =======================
# Build
# =======================
build: build-backend build-frontend

build-backend:
	@echo "Building backend services..."
	cd backend && go build -o bin/user-service ./cmd/user-service

build-frontend:
	@echo "Building frontend..."
	cd frontend && pnpm build

# =======================
# Development
# =======================
dev-backend:
	cd backend && go run ./cmd/user-service

dev-frontend:
	cd frontend && pnpm dev

# =======================
# Testing
# =======================
test: test-backend test-frontend

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && pnpm test

# =======================
# Lint / Format
# =======================
lint: lint-backend lint-frontend

lint-backend:
	cd backend && go vet ./...

lint-frontend:
	cd frontend && pnpm lint

fmt:
	cd backend && gofmt -w .
	cd frontend && pnpm format

# =======================
# Code Generation
# =======================
ent:
	cd backend/internal/user && go generate ./ent/...

wire:
	cd backend/internal/user && go run github.com/google/wire/cmd/wire@latest

# =======================
# Infrastructure
# =======================
docker-up:
	docker compose -f deployments/docker-compose.yml up -d

docker-down:
	docker compose -f deployments/docker-compose.yml down

# =======================
# Clean
# =======================
clean:
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/out
