# Tech Spec — Rifa Online

## 1. Stack Tecnológica

| Camada | Tecnologia | Versão | Justificativa |
|--------|-----------|--------|---------------|
| Frontend | Vue 3 + TypeScript | ^3.5 | SPA reativa, ecossistema maduro |
| Build | Vite | ^6 | Build rápido, HMR nativo |
| UI | Tailwind CSS | ^4 | Utility-first, rápido protótipo |
| Estado | Pinia | ^2 | Store oficial Vue 3 |
| Backend | Go | 1.23+ | Performance, concorrência natural |
| API HTTP | Chi Router | v5 | Leve, padrão Go net/http |
| Database | MongoDB | 7+ | Schema flexível para rifas/números |
| ODM | MongoDB Go Driver | v1 | Driver oficial |
| Pagamento | InfinitePay Checkout API | — | Checkout hospedado por `handle`, PIX+cartão |
| Cache | Redis | 7+ | Locking de números no checkout |

## 2. Arquitetura

```
┌──────────────────────────────────────────────┐
│                 Frontend (Vue 3)              │
│  localhost:5173 (dev) / static build (prod)  │
└──────────────────┬───────────────────────────┘
                   │ HTTP REST (JSON)
                   ▼
┌──────────────────────────────────────────────┐
│             Backend (Go + Chi)                │
│  localhost:8080                               │
│                                               │
│  Handlers: Auth · Raffle · Payment ·          │
│            Webhook · Subscription · Admin     │
│       │           │              │            │
│  ┌────▼───────────▼──────────────▼────────┐   │
│  │          Service Layer                 │   │
│  │  auth · raffle · payment · subscription│   │
│  └────┬───────────┬──────────────┬────────┘   │
│       │           │              │            │
│  ┌────▼────┐ ┌────▼────┐  ┌──────▼──────┐     │
│  │ MongoDB │ │ Redis   │  │ InfinitePay │     │
│  └─────────┘ └─────────┘  │  Checkout   │     │
│                           └──────┬──────┘     │
└──────────────────────────────────┼───────────┘
                   ▲ Webhooks       │
                   │                │
      InfinitePay ─┘◄───────────────┘
```

Camadas: `handler` (HTTP) → `service` (regra de negócio) → `repository` (MongoDB) e `pkg/infinitepay` (cliente HTTP). Redis usado só para lock de reserva.

## 3. Modelagem de Dados (MongoDB)

### 3.1 Collection: `users`

```json
{
  "_id": "ObjectId",
  "name": "string",
  "email": "string (unique indexed)",
  "passwordHash": "string",
  "role": "string (ADMIN | USER)",
  "phone": "string (opcional)",
  "infinitePayHandle": "string (opcional, recebe pagamentos das rifas)",
  "subscriptionStatus": "string (ACTIVE | INACTIVE | PAST_DUE | CANCELLED)",
  "subscriptionExpiresAt": "ISODate | null",
  "subscriptionIsTrial": "boolean",
  "hasSubscriptionBefore": "boolean",
  "createdAt": "ISODate",
  "updatedAt": "ISODate"
}
```

Trial = 7 dias (`model.TrialDays`). `passwordHash` nunca é serializado em JSON.

### 3.2 Collection: `raffles`

```json
{
  "_id": "ObjectId",
  "organizerId": "ObjectId (ref: users)",
  "title": "string",
  "description": "string",
  "ticketPrice": "int (centavos)",
  "maxNumbers": "int",
  "drawDate": "ISODate",
  "imageUrl": "string (opcional)",
  "status": "string (ACTIVE | CANCELLED | DRAWN)",
  "externalId": "string (opcional)",
  "winnerNumber": "int | null",
  "createdAt": "ISODate",
  "updatedAt": "ISODate"
}
```

Indexes: `{ organizerId: 1, status: 1 }`, `{ status: 1, drawDate: 1 }`

### 3.3 Collection: `tickets`

```json
{
  "_id": "ObjectId",
  "raffleId": "ObjectId (ref: raffles)",
  "number": "int (1..maxNumbers)",
  "status": "string (AVAILABLE | RESERVED | PAID)",
  "buyerName": "string | null",
  "buyerPhone": "string | null",
  "paymentId": "string | null",
  "reservedAt": "ISODate | null",
  "paidAt": "ISODate | null",
  "createdAt": "ISODate"
}
```

> A posse do ticket é identificada por **`buyerPhone`** (não email). `reservationExpiresIn` é calculado em runtime (não persistido).

Index: `{ raffleId: 1, number: 1 }` (unique compound), `{ raffleId: 1, status: 1 }`

### 3.4 Collection: `payments`

```json
{
  "_id": "ObjectId",
  "type": "string (RAFFLE | SUBSCRIPTION)",
  "raffleId": "ObjectId (opcional, RAFFLE)",
  "userId": "ObjectId (opcional, SUBSCRIPTION)",
  "ticketIds": ["ObjectId"],
  "buyerName": "string",
  "buyerEmail": "string (opcional)",
  "buyerPhone": "string",
  "checkoutUrl": "string",
  "invoiceSlug": "string (InfinitePay)",
  "transactionNsu": "string (InfinitePay)",
  "amount": "int (centavos)",
  "status": "string (PENDING | PAID | REFUNDED | EXPIRED)",
  "paymentMethod": "string (PIX | CARD | BOLETO)",
  "paidAt": "ISODate | null",
  "createdAt": "ISODate"
}
```

Um único modelo de pagamento atende rifas (`type=RAFFLE`) e assinaturas (`type=SUBSCRIPTION`, `amount=SubscriptionPrice=1000`).

### 3.5 Collection: `webhook_events`

```json
{
  "_id": "ObjectId",
  "eventId": "string (unique, idempotência)",
  "event": "string",
  "rawBody": "string",
  "processed": "boolean",
  "createdAt": "ISODate"
}
```

Index: `{ eventId: 1 }` (unique)

## 4. API REST (Endpoints)

Base: `/api/v1`. Rotas marcadas (auth) exigem JWT; (admin) exigem `role=ADMIN`.

### 4.1 Autenticação

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /auth/register | Cadastro (inicia trial 7d) |
| POST | /auth/login | Login (JWT) |
| POST | /auth/refresh | Renovar access token |

### 4.2 Perfil (auth)

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | /me | Perfil do usuário |
| PUT | /me | Atualizar perfil (nome, telefone) |
| PUT | /me/infinite-pay-handle | Definir handle InfinitePay |
| GET | /me/purchases | Compras do usuário |

### 4.3 Assinatura (auth)

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /subscription/checkout | Iniciar checkout de assinatura |
| POST | /subscription/dev-checkout | Ativar assinatura sem pagar (apenas dev) |
| GET | /subscription/status | Status da assinatura |

### 4.4 Rifas

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | /raffles | Listar rifas públicas |
| GET | /raffles/{id} | Detalhes + números |
| POST | /raffles/{id}/checkout | Iniciar compra de números |
| POST | /raffles | Criar rifa (auth, assinatura ativa) |
| PUT | /raffles/{id} | Editar rifa (auth, dono) |
| PATCH | /raffles/{id}/cancel | Cancelar rifa (auth) |
| DELETE | /raffles/{id} | Excluir rifa (auth) |
| POST | /raffles/{id}/draw | Sortear vencedor (auth) |
| GET | /raffles/my | Minhas rifas (auth) |
| GET | /raffles/{id}/stats | Estatísticas da rifa (auth) |

### 4.5 Pagamentos / Compra

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | /payments/my | Meus pagamentos (por telefone) |
| GET | /payments/{id} | Detalhe do pagamento |
| GET | /payments/my/tickets | Meus tickets (por telefone) |
| POST | /payments/{id}/confirm | Confirmar pagamento (polling InfinitePay) (auth) |

### 4.6 Dashboard / Admin

| Método | Rota | Descrição |
|--------|------|-----------|
| GET | /dashboard/stats | Estatísticas do organizador (auth) |
| GET | /admin/users | Listar/buscar usuários (admin) |
| GET | /admin/users/{id} | Detalhes do usuário (admin) |
| PUT | /admin/users/{id}/subscription | Alterar assinatura do usuário (admin) |
| GET | /admin/raffles | Listar rifas (admin) |
| GET | /admin/stats | Estatísticas globais (admin) |
| GET | /admin/contact-messages | Mensagens de contato recebidas (admin) |

### 4.7 Contato

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /contact | Enviar mensagem pela página pública de contato |

Persistido na collection `contact_messages` (`{ name, contact, message, ip, handled, createdAt }`). Validação: nome obrigatório (≤150), mensagem entre 10 e 2000 caracteres. Sem exposição de e-mail — leitura pelo painel admin.

### 4.8 Webhooks

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /webhooks/infinitepay | Receber eventos InfinitePay |

## 5. Integração com InfinitePay

InfinitePay usa **checkout hospedado por `handle`** (a `$tag` do recebedor). Não há API key — o pagamento é roteado pelo handle. Base URL: `https://api.checkout.infinitepay.io` (`INFINITEPAY_BASE_URL`).

### 5.1 Fluxo de Compra (rifa)

```
1. POST /api/v1/raffles/{id}/checkout
   Request: { numbers: [5,12,33], name, phone }
   - Backend reserva números no Redis (lock 5min, SetNX)
   - Cria payment (type=RAFFLE, status=PENDING)
   - POST /checkout (InfinitePay) com handle do organizador,
     items (quantity, price, description), order_nsu, redirect_url, webhook_url
   - Salva checkoutUrl no payment
   Response: { checkoutUrl }

2. Frontend redireciona para checkoutUrl

3. Cliente paga (PIX/Cartão)

4. Confirmação por dois caminhos (idempotentes):
   a) Webhook: POST /api/v1/webhooks/infinitepay
      Payload: invoice_slug, order_nsu, transaction_nsu, amount, paid_amount, ...
   b) Polling: POST /api/v1/payments/{id}/confirm  (CheckPayment via order_nsu)
   - Verifica idempotência (webhook_events.eventId)
   - payment.status = PAID, tickets.status = PAID (buyerPhone)
   - Libera locks do Redis

5. Frontend exibe confirmação (PaymentSuccess / PaymentPending)
```

### 5.2 Tipos principais (`pkg/infinitepay`)

- `CreateCheckoutRequest{ handle, items[], order_nsu, redirect_url, webhook_url, customer }`
- `CreateCheckoutResponse{ url }`
- `PaymentCheckRequest{ handle, order_nsu, transaction_nsu, slug }` → `PaymentCheckResponse{ success, paid, amount, paid_amount, ... }`
- `WebhookPayload{ invoice_slug, order_nsu, transaction_nsu, amount, paid_amount, receipt_url, items[] }`

### 5.3 Assinatura

`POST /subscription/checkout` cria um payment `type=SUBSCRIPTION` (`amount=1000`) e um checkout InfinitePay para o **handle da plataforma**. Confirmação ativa a assinatura (`subscriptionStatus=ACTIVE`, define `subscriptionExpiresAt`). Em dev, `POST /subscription/dev-checkout` ativa sem pagar.

## 6. Tratamento de Concorrência

Problema: dois participantes selecionarem o mesmo número simultaneamente.

Solução:
1. **Redis lock** (`SetNX`, TTL **5min**) ao iniciar checkout
2. Ao criar o checkout, valida que os números seguem disponíveis; senão 409
3. TTL expira o lock automaticamente
4. Reservas pendentes são limpas (varredura de **10min** — `reservationTTL` / `ExpirePendingOlderThan`); MongoDB é a source of truth

## 7. Segurança

- JWT (access 1h, refresh 7d)
- Senhas com bcrypt (cost 12; entrada limitada a 72 chars — limite do bcrypt)
- Middleware de auth + middleware de assinatura (bloqueia criação de rifa sem plano ativo; admin faz bypass)
- Role-based: rotas `/admin` exigem `role=ADMIN`
- CORS restrito ao frontend
- Idempotência de webhook via `webhook_events.eventId`
- Validação de limites de tamanho de campos (backend + `maxlength` no frontend)
- HTTPS obrigatório em produção

## 8. Frontend (Vue 3)

### 8.1 Páginas (`src/pages`)

| Rota | Componente | Descrição |
|------|-----------|-----------|
| / | Home.vue | Landing + rifas ativas |
| /raffles/:id | RaffleDetail.vue | Detalhes + grid de números |
| /raffles/:id/checkout | Checkout.vue | Formulário de compra (nome + telefone) |
| /payment/success | PaymentSuccess.vue | Confirmação |
| /payment/pending | PaymentPending.vue | Aguardando pagamento |
| /login | Login.vue | Login |
| /register | Register.vue | Cadastro |
| /dashboard | Dashboard.vue | Painel do organizador (gráficos) |
| /dashboard/raffles/new | CreateRaffle.vue | Criar rifa |
| /dashboard/raffles/:id/edit | EditRaffle.vue | Editar rifa |
| /raffles/:id/result | RaffleResult.vue | Resultado do sorteio |
| /profile | Profile.vue | Perfil (telefone, handle InfinitePay) |
| /subscription | Subscription.vue | Assinatura |
| /purchases | MyPurchases.vue | Minhas compras (por telefone) |
| /admin | Admin.vue | Painel administrativo |

### 8.2 Componentes (`src/components`)

`NumberGrid.vue`, `RaffleCard.vue`, `Alert.vue`, `LoadingSpinner.vue`. Layout em `src/layouts` (Navbar, Footer).

### 8.3 Estado Global (Pinia, `src/stores`)

- `useAuthStore` — autenticação (restaura usuário do localStorage no refresh)
- `useRaffleStore` — dados de rifas
- `usePaymentStore` — pagamentos

## 9. Variáveis de Ambiente

### Backend

```env
PORT=8080
APP_ENV=development            # development | production
MONGODB_URI=mongodb://localhost:27017
MONGODB_DB_NAME=rifaonline
REDIS_URI=redis://localhost:6379
JWT_SECRET=change-me
INFINITEPAY_HANDLE=$sua_tag    # handle da plataforma (assinaturas)
INFINITEPAY_BASE_URL=https://api.checkout.infinitepay.io
FRONTEND_URL=http://localhost:5173
LOG_LEVEL=info                 # debug | info | warn | error
LOG_FORMAT=text                # text | json
```

### Frontend

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 10. Decisões Técnicas (ADRs)

### ADR-001: MongoDB sobre PostgreSQL
Rifas têm dados semi-estruturados; tickets como documentos com status mutável. **Decisão:** MongoDB. Índices e consistência geridos pela aplicação.

### ADR-002: Chi Router sobre Gin
Chi segue o `net/http` nativo, reduz dependências. **Decisão:** Chi Router.

### ADR-003: Checkout hospedado (redirect)
Mais seguro e rápido para o MVP. **Decisão:** checkout hospedado InfinitePay.

### ADR-004: Redis para locking
TTL automático e menor latência que `findAndModify`. **Decisão:** Redis para lock temporário (5min) + MongoDB como source of truth.

### ADR-005: InfinitePay sobre AbacatePay
Migrado de AbacatePay para **InfinitePay** (checkout por `handle`, sem API key). Cada organizador define seu próprio handle para receber os pagamentos das rifas; a plataforma usa seu handle para as assinaturas. Confirmação por webhook **e** polling (`/confirm`) para robustez.

### ADR-006: Posse de ticket por telefone
Participantes compram sem conta. O **telefone** (`buyerPhone`) identifica o comprador e é usado em "Minhas compras", em vez de email.

### ADR-007: Modelo SaaS com trial
Organizadores assinam (R$10/mês via InfinitePay) com trial de 7 dias. Middleware bloqueia criação de rifa sem plano ativo; admin faz bypass.

## 11. Infraestrutura

| Componente | Serviço |
|-----------|---------|
| Frontend | Vercel / static hosting |
| Backend | Railway / Fly.io |
| Database | MongoDB Atlas |
| Redis | Upstash (serverless) |

Deploy local/unificado via `docker-compose.yml` + `nginx.unified.conf` (frontend estático + proxy para o backend). Ver `Dockerfile`, `Makefile`, `start.sh`.
