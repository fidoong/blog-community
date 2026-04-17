.PHONY: all build dev test lint fmt ent wire docker-up docker-down

# Default target
all: build

# =======================
# Build
# =======================
build: build-backend build-frontend

build-backend:
	@echo "Building backend services..."
	cd backend && go build -o bin/api-service ./cmd/api-service

build-frontend:
	@echo "Building frontend..."
	cd frontend && pnpm build

# =======================
# Development
# =======================
dev-backend:
	cd backend && go run ./cmd/api-service

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
# Infrastructure (Orb VM: fidoo)
# =======================
infra-up:
	cd scripts && bash deploy-infra.sh

infra-down:
	ssh root@192.168.139.191 "cd /opt/blog-infra && docker compose down"

infra-logs:
	ssh root@192.168.139.191 "cd /opt/blog-infra && docker compose logs -f"

infra-status:
	ssh root@192.168.139.191 "cd /opt/blog-infra && docker compose ps"

infra-pull:
	ssh root@192.168.139.191 "cd /opt/blog-infra && docker compose pull"

# =======================
# Clean
# =======================
clean:
	rm -rf backend/bin
	rm -rf frontend/.next
	rm -rf frontend/out
