.PHONY: dev-backend dev-frontend dev build-backend build-frontend build \
        test-backend test-frontend test lint-backend lint-frontend lint \
        docker-up docker-down docker-build db-up db-down setup clean

# ─── Development ────────────────────────────────────────────────

dev-backend:
	cd backend && go run ./cmd/server/main.go

dev-backend-hot:
	cd backend && air

dev-frontend:
	cd frontend && npm run dev

dev: dev-backend dev-frontend

dev-hot: dev-backend-hot dev-frontend

# ─── Build ──────────────────────────────────────────────────────

build-backend:
	cd backend && go build -o bin/server ./cmd/server/main.go

build-frontend:
	cd frontend && npm run build

build: build-backend build-frontend

# ─── Test ───────────────────────────────────────────────────────

test-backend:
	cd backend && go test ./...

test-frontend:
	cd frontend && npx vitest run

test: test-backend test-frontend

# ─── Lint ───────────────────────────────────────────────────────

lint-backend:
	cd backend && test -z "$$(go vet ./... 2>&1)" || go vet ./...

lint-frontend:
	cd frontend && npx vue-tsc --noEmit 2>/dev/null; npx eslint . --ext .vue,.ts 2>/dev/null || true

lint: lint-backend lint-frontend

# ─── Docker ─────────────────────────────────────────────────────

docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-logs:
	docker compose logs -f

docker: docker-up

# ─── Infrastructure ─────────────────────────────────────────────

db-up:
	docker compose up -d mongodb redis

db-down:
	docker compose down mongodb redis

# ─── Setup ──────────────────────────────────────────────────────

setup:
	cd backend && cp -n .env.example .env 2>/dev/null || true
	cd frontend && npm install
	cd backend && go mod tidy

# ─── Clean ──────────────────────────────────────────────────────

clean:
	rm -rf backend/bin
	rm -rf frontend/dist
