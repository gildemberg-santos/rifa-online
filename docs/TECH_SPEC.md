# Tech Spec — Rifa Online

## 1. Stack Tecnológica

| Camada | Tecnologia | Versão | Justificativa |
|--------|-----------|--------|---------------|
| Frontend | Vue 3 + TypeScript | ^3.5 | SPA reativa, ecossistema maduro |
| Build | Vite | ^6 | Build rápido, HMR nativo |
| UI | Tailwind CSS | ^4 | Utility-first, rápido protótipo |
| Backend | Go | 1.23+ | Performance, concorrência natural |
| API HTTP | Chi Router | v5 | Leve, padrão Go net/http |
| Database | MongoDB | 7+ | Schema flexível para rifas/números |
| ODM | MongoDB Go Driver | v1 | Driver oficial |
| Pagamento | AbacatePay API v2 | — | Gateway brasileiro, PIX+cartão |
| Cache | Redis | 7+ | Sessões de checkout e locking de números |

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
│  localhost:8080 (api.rifaonline.com.br)      │
│                                               │
│  ┌─────────┐ ┌──────────┐ ┌───────────────┐  │
│  │ Auth    │ │ Rifa     │ │ Pagamento     │  │
│  │ Handler │ │ Handler  │ │ Handler       │  │
│  └────┬────┘ └────┬─────┘ └──────┬────────┘  │
│       │           │              │            │
│  ┌────▼───────────▼──────────────▼────────┐   │
│  │          Service Layer                 │   │
│  └────┬───────────┬──────────────┬────────┘   │
│       │           │              │            │
│  ┌────▼────┐ ┌────▼────┐  ┌─────▼─────┐      │
│  │ MongoDB │ │ Redis   │  │AbacatePay │      │
│  └─────────┘ └─────────┘  │ API v2    │      │
│                           └───────────┘      │
└──────────────────────────────────────────────┘
                   ▲ Webhooks
                   │
      AbacatePay ──┘
```

## 3. Modelagem de Dados (MongoDB)

### 3.1 Collection: `users`

```json
{
  "_id": "ObjectId",
  "name": "string",
  "email": "string (unique indexed)",
  "passwordHash": "string",
  "createdAt": "ISODate",
  "updatedAt": "ISODate"
}
```

### 3.2 Collection: `raffles`

```json
{
  "_id": "ObjectId",
  "organizerId": "ObjectId (ref: users)",
  "title": "string",
  "description": "string",
  "ticketPrice": "int (centavos)",
  "maxNumbers": "int (ex: 500)",
  "drawDate": "ISODate",
  "imageUrl": "string (opcional)",
  "status": "string (ACTIVE | CANCELLED | DRAWN)",
  "externalId": "string (único para produto AbacatePay)",
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
  "buyerEmail": "string | null",
  "paymentId": "string | null (bill_xxx do AbacatePay)",
  "reservedAt": "ISODate | null",
  "paidAt": "ISODate | null",
  "createdAt": "ISODate"
}
```

Index: `{ raffleId: 1, number: 1 }` (unique compound), `{ raffleId: 1, status: 1 }`

### 3.4 Collection: `payments`

```json
{
  "_id": "ObjectId",
  "raffleId": "ObjectId (ref: raffles)",
  "ticketIds": ["ObjectId"],
  "buyerName": "string",
  "buyerEmail": "string",
  "abacateCheckoutId": "string (bill_xxx)",
  "abacateCheckoutUrl": "string",
  "amount": "int (centavos)",
  "status": "string (PENDING | PAID | REFUNDED)",
  "paymentMethod": "string (PIX | CARD | BOLETO)",
  "paidAt": "ISODate | null",
  "createdAt": "ISODate"
}
```

Index: `{ abacateCheckoutId: 1 }` (unique), `{ buyerEmail: 1 }`

### 3.5 Collection: `webhook_events`

```json
{
  "_id": "ObjectId",
  "eventId": "string (unique, para idempotência)",
  "event": "string (checkout.completed, etc)",
  "rawBody": "string",
  "processed": "boolean",
  "createdAt": "ISODate"
}
```

Index: `{ eventId: 1 }` (unique)

## 4. API REST (Endpoints)

### 4.1 Autenticação

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /api/v1/auth/register | Cadastro |
| POST | /api/v1/auth/login | Login (JWT) |

### 4.2 Rifas

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /api/v1/raffles | Criar rifa (auth) |
| GET | /api/v1/raffles | Listar rifas públicas |
| GET | /api/v1/raffles/:id | Detalhes da rifa + números |
| PUT | /api/v1/raffles/:id | Editar rifa (auth, dono) |
| PATCH | /api/v1/raffles/:id/cancel | Cancelar rifa (auth) |
| POST | /api/v1/raffles/:id/draw | Sortear vencedor (auth) |
| GET | /api/v1/raffles/my | Minhas rifas (auth, organizador) |

### 4.3 Tickets / Compra

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /api/v1/raffles/:id/checkout | Iniciar compra de números |
| GET | /api/v1/raffles/:id/tickets | Listar números (público) |
| GET | /api/v1/payments/my | Meus pagamentos (email) |

### 4.4 Webhooks

| Método | Rota | Descrição |
|--------|------|-----------|
| POST | /api/v1/webhooks/abacatepay | Receber eventos AbacatePay |

## 5. Integração com AbacatePay

### 5.1 Fluxo de Compra

```
1. POST /api/v1/raffles/:id/checkout
   Request: { numbers: [5, 12, 33], name, email }
   - Backend reserva números no Redis (lock 15min)
   - Cria payment no MongoDB (status PENDING)
   - POST /checkouts/create (AbacatePay) com os números como itens
   - Salva abacateCheckoutId no payment
   Response: { checkoutUrl: "https://app.abacatepay.com/pay/bill_xxx" }

2. Frontend redireciona para checkoutUrl

3. Cliente paga (PIX/Cartão)

4. AbacatePay chama webhook: POST /api/v1/webhooks/abacatepay
   Evento: checkout.completed
   - Verifica HMAC signature
   - Verifica idempotência (eventId)
   - Atualiza payment.status = PAID
   - Atualiza tickets.status = PAID
   - Libera locks do Redis

5. Frontend redireciona para completionUrl
   - Tela de confirmação com números comprados
```

### 5.2 Produtos no AbacatePay

Cada rifa será criada como um **produto** no AbacatePay via `POST /products/create`. Os números vendidos serão itens individuais no checkout com `quantity` igual à quantidade de números.

Alternativa mais simples: não criar produtos no AbacatePay, usando `externalId` no `POST /checkouts/create` com `items` contendo os dados da rifa diretamente. A AbacatePay permite passar produtos avulsos no checkout sem pré-cadastro usando o `items[].id` com suporte a produtos existentes ou criados dinamicamente.

**Decisão:** Usar produtos pré-cadastrados para rifas na AbacatePay. Ao criar uma rifa, criamos também um produto no AbacatePay e armazenamos o `prod_id` retornado. No checkout, usamos esse `prod_id` com `quantity` igual ao número de tickets selecionados.

### 5.3 Webhook Security

A AbacatePay envia:
- `?webhookSecret=SEU_SECRET` na query string
- `X-Webhook-Signature` header com HMAC-SHA256

Validação HMAC em Go usando a chave pública da AbacatePay (hardcoded) e `crypto/hmac` com `timingSafeEqual`.

## 6. Tratamento de Concorrência

Problema: dois participantes selecionarem o mesmo número simultaneamente.

Solução:
1. **Redis lock** ao selecionar números (TTL 15min)
2. Ao criar checkout, verificar se números ainda estão disponíveis
3. Se já reservados, retornar erro 409 Conflict
4. Se expirou, liberar automaticamente

## 7. Segurança

- JWT com refresh token (access: 1h, refresh: 7d)
- Senhas com bcrypt (cost 12)
- CORS restrito ao domínio do frontend
- Rate limiting (100 req/min por IP)
- Webhook validation (HMAC + secret)
- Validação de entrada em todas as rotas
- HTTPS obrigatório em produção

## 8. Frontend (Vue 3)

### 8.1 Páginas

| Rota | Componente | Descrição |
|------|-----------|-----------|
| / | Home.vue | Landing + lista de rifas ativas |
| /raffles/:id | RaffleDetail.vue | Detalhes + grid de números |
| /raffles/:id/checkout | Checkout.vue | Formulário de compra |
| /payment/success | PaymentSuccess.vue | Confirmação de pagamento |
| /payment/pending | PaymentPending.vue | Aguardando pagamento |
| /login | Login.vue | Login do organizador |
| /register | Register.vue | Cadastro do organizador |
| /dashboard | Dashboard.vue | Painel do organizador |
| /dashboard/raffles/new | CreateRaffle.vue | Criar nova rifa |
| /dashboard/raffles/:id/edit | EditRaffle.vue | Editar rifa |
| /raffles/:id/result | RaffleResult.vue | Resultado do sorteio |

### 8.2 Componentes Compartilhados

- `NumberGrid.vue` — Grid visual dos números (disponível/reservado/pago)
- `RaffleCard.vue` — Card de rifa para listagem
- `Navbar.vue` — Barra de navegação
- `Footer.vue` — Rodapé
- `LoadingSpinner.vue` — Spinner de carregamento
- `Alert.vue` — Alertas e mensagens

### 8.3 Estado Global (Pinia)

- `useAuthStore` — Estado de autenticação
- `useRaffleStore` — Dados de rifas
- `usePaymentStore` — Estado de pagamentos

## 9. Variáveis de Ambiente

### Backend (.env)

```env
PORT=8080
MONGODB_URI=mongodb://localhost:27017/rifaonline
REDIS_URI=redis://localhost:6379
JWT_SECRET=seu-secret-aqui
ABACATEPAY_API_KEY=sk_...
ABACATEPAY_WEBHOOK_SECRET=whsec_...
FRONTEND_URL=http://localhost:5173
```

### Frontend (.env)

```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 10. Decisões Técnicas (ADRs)

### ADR-001: MongoDB sobre PostgreSQL

**Contexto:** Rifas têm dados semi-estruturados (números como sub-documentos, histórico de pagamentos). MongoDB oferece schema flexível que simplifica a modelagem de tickets como documentos individuais com status mutável.

**Decisão:** MongoDB.

**Consequência:** Necessário gerenciar índices e consistência via aplicação.

### ADR-002: Chi Router sobre Gin

**Contexto:** Gin é mais popular mas Chi segue o padrão `net/http` nativo do Go, facilitando manutenção e reduzindo dependências.

**Decisão:** Chi Router.

### ADR-003: Checkout hospedado (redirect) vs Transparente

**Contexto:** AbacatePay oferece checkout hospedado (redirect) e transparente (embed). Checkout hospedado é mais seguro e rápido de implementar no MVP.

**Decisão:** Checkout hospedado no MVP. Transparente como melhoria futura.

### ADR-004: Redis para locking em vez de apenas MongoDB

**Contexto:** Dois usuários podem selecionar o mesmo número simultaneamente. Usar apenas MongoDB com `findAndModify` funciona mas Redis oferece TTL automático e menor latência.

**Decisão:** Redis para locking temporário + MongoDB como source of truth.

## 11. Infraestrutura (MVP)

| Componente | Serviço |
|-----------|---------|
| Frontend | Vercel / static hosting |
| Backend | Railway / Fly.io |
| Database | MongoDB Atlas (free tier) |
| Redis | Upstash (serverless) |
| Domínio | opcional |
