# Rifa Online — Project Skill

## Stack
- **Backend:** Go 1.23+ / Chi Router v5 / MongoDB 7+ / Redis 7+
- **Frontend:** Vue 3 + TypeScript / Vite 6 / Tailwind CSS 4 / Pinia / Vue Router
- **Payments:** InfinitePay Checkout API (checkout hospedado, webhooks)

## Structure
```
backend/
├── cmd/server/main.go          # Entrypoint, DI, server startup
├── internal/
│   ├── auth/jwt.go             # JWT generation & validation
│   ├── config/config.go        # Env vars loader
│   ├── handler/                # HTTP handlers (auth, raffle, payment, webhook)
│   ├── middleware/auth.go      # JWT auth middleware
│   ├── model/                  # MongoDB structs (user, raffle, ticket, payment, webhook_event)
│   ├── repository/             # MongoDB CRUD + indexes + transactions
│   ├── service/                # Business logic layer
│   └── webhook/infinitepay.go  # Webhook payload parsing
└── pkg/infinitepay/             # InfinitePay HTTP client (checkout, types)
frontend/
└── src/
    ├── assets/                 # Images, CSS
    ├── components/             # NumberGrid, RaffleCard, Navbar, Footer, LoadingSpinner, Alert
    ├── layouts/                # Navbar, Footer
    ├── pages/                  # Home, RaffleDetail, Checkout, PaymentSuccess, PaymentPending
    │                           # Login, Register, Dashboard, CreateRaffle, EditRaffle, RaffleResult
    ├── router/index.ts         # Vue Router config
    ├── stores/                 # Pinia: auth, raffle, payment
    ├── utils/api.ts            # Axios/fetch wrapper with JWT interceptor
    └── views/                  # Additional views
```

## Commands
- `go run ./cmd/server` — start backend on :8080
- `npm run dev` — start frontend on :5173
- `go test ./...` — run all Go tests
- `npm run test` — run Vue tests (Vitest)
- `go vet ./...` — static analysis

## Conventions
- Go: stdlib `net/http` interfaces, Chi Router, BSON/JSON tags on all structs
- Vue: `<script setup lang="ts">`, Composition API, Pinia stores as `use*Store`
- API: `/api/v1/{resource}`, JSON request/response, snake_case fields
- MongoDB: `createdAt`/`updatedAt` timestamps on all documents
- JWT: access token 1h, refresh token 7d, bcrypt cost 12
- Errors: JSON `{ "error": "message" }`, appropriate HTTP status codes
- Redis: locks with 15min TTL for number reservation during checkout
- Concurrency: Redis lock -> create payment -> create InfinitePay checkout -> return URL

## Tasks (from docs/TASKS.md)
17 tasks in 5 phases. Each task has defined files, description, and acceptance criteria.
Phase 1 (base): Tasks 1.1–1.4 (6 days)
Phase 2 (raffles): Tasks 2.1–2.4 (8 days)
Phase 3 (draw/dashboard): Tasks 3.1–3.2 (2 days)
Phase 4 (frontend): Tasks 4.1–4.5 (9 days)
Phase 5 (finalization): Tasks 5.1–5.2 (4 days)

## PRD
docs/PRD.md — full product requirements
## Tech Spec
docs/TECH_SPEC.md — architecture, data models, API, security
