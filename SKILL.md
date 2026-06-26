# Rifa Online — Project Skill

## Stack
- **Backend:** Go 1.23+ / Chi Router v5 / MongoDB 7+ / Redis 7+
- **Frontend:** Vue 3 + TypeScript / Vite 6 / Tailwind CSS 4 / Pinia / Vue Router
- **Payments:** InfinitePay Checkout API (checkout hospedado por `handle`, webhook + polling)
- **Model:** SaaS — organizadores assinam (R$10/mês, trial 7d) p/ criar rifas; admin bypass

## Structure
```
backend/
├── cmd/server/main.go          # Entrypoint, DI, server startup
├── internal/
│   ├── auth/jwt.go             # JWT generation & validation
│   ├── config/config.go        # Env vars loader
│   ├── handler/                # HTTP handlers (auth, raffle, payment, webhook, subscription, admin)
│   ├── middleware/             # JWT auth, subscription gate, admin-only
│   ├── model/                  # MongoDB structs (user, raffle, ticket, payment, webhook_event)
│   ├── repository/             # MongoDB CRUD + indexes + transactions
│   ├── service/                # Business logic (auth, raffle, payment, subscription)
│   └── webhook/infinitepay.go  # Webhook payload parsing
└── pkg/infinitepay/             # InfinitePay HTTP client (checkout, payment-check, types)
frontend/
└── src/
    ├── assets/                 # Images, CSS
    ├── components/             # NumberGrid, RaffleCard, LoadingSpinner, Alert
    ├── layouts/                # Navbar, Footer
    ├── pages/                  # Home, RaffleDetail, Checkout, PaymentSuccess, PaymentPending,
    │                           # Login, Register, Dashboard, CreateRaffle, EditRaffle, RaffleResult,
    │                           # Profile, Subscription, MyPurchases, Admin
    ├── router/index.ts         # Vue Router config
    ├── stores/                 # Pinia: auth, raffle, payment
    └── utils/api.ts            # fetch wrapper with JWT interceptor
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
- JWT: access token 1h, refresh token 7d, bcrypt cost 12 (password max 72 chars)
- Errors: JSON `{ "error": "message" }`, appropriate HTTP status codes
- Redis: `SetNX` lock 5min TTL for number reservation; pending cleanup at 10min (`reservationTTL`)
- Concurrency: Redis lock -> create payment -> create InfinitePay checkout -> return URL
- Ticket ownership keyed by `buyerPhone` (not email); participants buy without an account
- Payments unified: `type` RAFFLE | SUBSCRIPTION; confirm via webhook OR `/payments/{id}/confirm` polling

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
