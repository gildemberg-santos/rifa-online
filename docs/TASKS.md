# Tasks — Rifa Online

## Fase 1: Projeto Base (Semanas 1-2)

### Task 1.1 — Inicializar Backend Go

**Tempo estimado:** 1 dia

**Arquivos a criar:**
- `backend/go.mod`
- `backend/go.sum`
- `backend/cmd/server/main.go`
- `backend/internal/config/config.go`

**Descrição:**
- Inicializar módulo Go `github.com/user/rifa-online`
- Configurar Chi Router, conexão MongoDB e Redis
- Estruturar `main.go` com injeção de dependências
- Ler variáveis de ambiente do `.env`

**Critérios de aceitação:**
- `go run cmd/server/main.go` sobe servidor em `:8080`
- Health check `GET /health` retorna `200 OK`

**Dependências:**
- go.mod: `github.com/go-chi/chi/v5`, `go.mongodb.org/mongo-driver`, `github.com/redis/go-redis/v9`, `github.com/joho/godotenv`, `golang.org/x/crypto`

---

### Task 1.2 — Modelos MongoDB

**Tempo estimado:** 1 dia

**Arquivos a criar:**
- `backend/internal/model/user.go`
- `backend/internal/model/raffle.go`
- `backend/internal/model/ticket.go`
- `backend/internal/model/payment.go`
- `backend/internal/model/webhook_event.go`

**Descrição:**
- Definir structs Go para cada collection
- Implementar timestamps automáticos (CreateAt, UpdatedAt)
- Tags BSON e JSON em todos os campos

**Critérios de aceitação:**
- Structs representam fielmente o schema definido no TECH_SPEC
- Compilam sem erros

---

### Task 1.3 — Repositórios MongoDB

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/internal/repository/user_repo.go`
- `backend/internal/repository/raffle_repo.go`
- `backend/internal/repository/ticket_repo.go`
- `backend/internal/repository/payment_repo.go`
- `backend/internal/repository/webhook_repo.go`

**Descrição:**
- Implementar CRUD básico para cada collection
- Garantir índices únicos e compostos
- Transações atômicas para criação de tickets + pagamento

**Critérios de aceitação:**
- Testes de unidade para cada método
- Inserção, busca, atualização funcionando contra MongoDB

---

### Task 1.4 — Autenticação (Go)

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/internal/auth/jwt.go`
- `backend/internal/middleware/auth.go`
- `backend/internal/handler/auth_handler.go`
- `backend/internal/service/auth_service.go`

**Descrição:**
- Hash de senha com bcrypt (cost 12)
- Geração e validação de JWT (access + refresh)
- Middleware de autenticação (extrai user_id do token)
- Handlers: POST /register, POST /login

**Critérios de aceitação:**
- Cadastro com email duplicado retorna 409
- Login com credenciais inválidas retorna 401
- Token válido permite acesso a rotas protegidas

---

## Fase 2: Módulo de Rifas (Semanas 2-4)

### Task 2.1 — CRUD de Rifas

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/internal/handler/raffle_handler.go`
- `backend/internal/service/raffle_service.go`

**Descrição:**
- Criar rifa (organizador autenticado)
- Listar rifas públicas (com paginação)
- Detalhes de rifa + grid de números
- Editar rifa (apenas dono, antes do sorteio)
- Cancelar rifa (apenas dono)

**Critérios de aceitação:**
- Organizador cria rifa com dados válidos
- Usuário não autenticado vê lista pública
- Grid mostra números disponíveis/reservados/pagos
- Edição só permite antes do sorteio
- Cancelar muda status para CANCELLED

---

### Task 2.2 — AbacatePay Client (Go)

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/pkg/abacatepay/client.go`
- `backend/pkg/abacatepay/types.go`
- `backend/pkg/abacatepay/checkout.go`
- `backend/pkg/abacatepay/product.go`

**Descrição:**
- Cliente HTTP para API v2 da AbacatePay
- Criar produto via `POST /products/create`
- Criar checkout via `POST /checkouts/create`
- Buscar checkout por ID via `GET /checkouts/{id}`
- Tratamento de erros da API (401, 403, 429)
- Response format `{ data, error, success }`

**Critérios de aceitação:**
- Criação de produto na AbacatePay funcional
- Criação de checkout retorna URL para redirecionamento
- Rate limit handle (backoff)
- Testes com Dev Mode da AbacatePay

---

### Task 2.3 — Processo de Compra

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/internal/handler/payment_handler.go`
- `backend/internal/service/payment_service.go`

**Descrição:**
- Endpoint POST /raffles/:id/checkout
- Reservar números no Redis (lock 15min TTL)
- Criar payment no MongoDB
- Criar checkout na AbacatePay
- Retornar URL de redirecionamento
- Endpoint GET /payments/my (por email)

**Critérios de aceitação:**
- Números válidos e disponíveis → checkout criado
- Número já reservado → erro 409
- Lock expirado libera automaticamente
- URL do AbacatePay retornada ao frontend

---

### Task 2.4 — Webhook Handler

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/internal/webhook/abacatepay.go`
- `backend/internal/handler/webhook_handler.go`

**Descrição:**
- Endpoint POST /webhooks/abacatepay
- Verificação HMAC (X-Webhook-Signature)
- Validação do webhookSecret na query string
- Idempotência (eventId)
- Processamento de eventos:
  - `checkout.completed` → atualizar payment + tickets
  - `checkout.refunded` → atualizar status payment + tickets
- Resposta 200 após processamento completo

**Critérios de aceitação:**
- Webhook inválido → 401
- Evento duplicado → 200 (não reprocessa)
- Pagamento confirmado → tickets marcados como PAID
- Simulação com Dev Mode da AbacatePay

---

## Fase 3: Sorteio e Dashboard (Semanas 4-5)

### Task 3.1 — Sorteio

**Tempo estimado:** 1 dia

**Descrição:**
- Endpoint POST /raffles/:id/draw (organizador autenticado)
- Validar que rifa está ACTIVE
- Validar que pelo menos 1 número foi vendido
- Sortear número aleatório entre os PAID
- Atualizar raffle.winnerNumber e raffle.status = DRAWN

**Critérios de aceitação:**
- Apenas organizador da rifa pode sortear
- Rifa já sorteada retorna 400
- Vencedor escolhido aleatoriamente

---

### Task 3.2 — Dashboard

**Tempo estimado:** 1 dia

**Descrição:**
- GET /raffles/my → minhas rifas com métricas
- GET /raffles/:id/stats → estatísticas detalhadas
  - Total de números vendidos
  - Total arrecadado
  - Percentual vendido
  - Lista de participantes

**Critérios de aceitação:**
- Métricas corretas baseadas nos payments PAID
- Apenas organizador vê seus dados

---

## Fase 4: Frontend Vue (Semanas 5-7)

### Task 4.1 — Setup do Frontend

**Tempo estimado:** 1 dia

**Arquivos/alterações:**
- `frontend/package.json`
- `frontend/vite.config.ts`
- `frontend/src/main.ts`
- `frontend/src/App.vue`
- `frontend/src/router/index.ts`
- `frontend/src/stores/auth.ts`
- `frontend/tailwind.config.js`
- `frontend/index.html`

**Descrição:**
- Inicializar Vue 3 + TypeScript + Vite
- Configurar Tailwind CSS, Pinia, Vue Router
- Configurar proxy para API em dev
- Layout base (Navbar + Footer)

**Critérios de aceitação:**
- `npm run dev` sobe em :5173
- Navbar com links e estado de login
- Proxy para backend em /api

---

### Task 4.2 — Páginas Públicas

**Tempo estimado:** 3 dias

**Arquivos a criar:**
- `frontend/src/pages/Home.vue`
- `frontend/src/pages/RaffleDetail.vue`
- `frontend/src/pages/RaffleResult.vue`
- `frontend/src/components/NumberGrid.vue`
- `frontend/src/components/RaffleCard.vue`

**Descrição:**
- Home: lista de rifas ativas com cards
- RaffleDetail: grid de números (coloridos por status), botão "Comprar"
- RaffleResult: exibe vencedor e número sorteado
- NumberGrid: componente reutilizável com cores:
  - Verde: disponível
  - Amarelo: reservado
  - Azul: pago
  - Cinza: indisponível

**Critérios de aceitação:**
- Navegação entre páginas
- Grid responsivo
- Seleção de números com feedback visual

---

### Task 4.3 — Fluxo de Pagamento

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `frontend/src/pages/Checkout.vue`
- `frontend/src/pages/PaymentSuccess.vue`
- `frontend/src/pages/PaymentPending.vue`
- `frontend/src/stores/payment.ts`

**Descrição:**
- Checkout: formulário (nome, email) + números selecionados + valor total
- Ao submeter: chamar POST /raffles/:id/checkout
- Redirecionar para URL do AbacatePay
- PaymentSuccess: tela de confirmação (after redirect)
- PaymentPending: polling do status via GET /payments/my

**Critérios de aceitação:**
- Formulário valida dados antes de enviar
- Redireciona para AbacatePay
- Tela de sucesso exibe números comprados

---

### Task 4.4 — Dashboard (Frontend)

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `frontend/src/pages/Dashboard.vue`
- `frontend/src/pages/CreateRaffle.vue`
- `frontend/src/pages/EditRaffle.vue`
- `frontend/src/stores/raffle.ts`

**Descrição:**
- Dashboard: listagem das rifas do organizador com métricas
- CreateRaffle: formulário de criação
- EditRaffle: formulário de edição
- Botão "Realizar Sorteio" no dashboard

**Critérios de aceitação:**
- Criar rifa redireciona para dashboard
- Editar rifa pré-preenche formulário
- Sorteio mostra confirmação antes de executar

---

### Task 4.5 — Autenticação (Frontend)

**Tempo estimado:** 1 dia

**Arquivos a criar:**
- `frontend/src/pages/Login.vue`
- `frontend/src/pages/Register.vue`
- `frontend/src/utils/api.ts` (axios/fetch wrapper)

**Descrição:**
- Tela de login e cadastro
- Armazenar JWT no localStorage
- Interceptor para incluir token nas requisições
- Rotas protegidas (guard no router)

**Critérios de aceitação:**
- Login salva token e redireciona
- Logout limpa token e redireciona
- Rotas protegidas redirecionam para /login

---

## Fase 5: Finalização (Semana 8)

### Task 5.1 — Testes

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `backend/internal/handler/*_test.go`
- `backend/internal/service/*_test.go`
- `backend/internal/repository/*_test.go`
- `frontend/src/**/*.spec.ts` (Vitest)

**Descrição:**
- Testes unitários para serviços e repositórios (Go)
- Testes de handler HTTP (httptest)
- Testes de componente Vue (Vitest + Vue Test Utils)
- Mínimo 60% de cobertura

**Critérios de aceitação:**
- `go test ./...` passa
- `npm run test` passa

---

### Task 5.2 — Docker e Deploy

**Tempo estimado:** 2 dias

**Arquivos a criar:**
- `Dockerfile` (backend multi-stage build)
- `Dockerfile` (frontend nginx)
- `docker-compose.yml` (backend + mongo + redis)
- `backend/.env.example`
- `frontend/.env.example`

**Descrição:**
- Dockerfile multi-stage para backend Go (build + distroless)
- Dockerfile para frontend (Vite build + nginx)
- docker-compose para ambiente local completo
- Scripts de deploy

**Critérios de aceitação:**
- `docker compose up` sobe todo o stack
- Frontend → Backend → MongoDB → Redis funcionam

---

## Resumo de Estimativas

| Fase | Tasks | Dias |
|------|-------|------|
| Fase 1: Projeto Base | 4 | 6 |
| Fase 2: Módulo de Rifas | 4 | 8 |
| Fase 3: Sorteio e Dashboard | 2 | 2 |
| Fase 4: Frontend Vue | 5 | 9 |
| Fase 5: Finalização | 2 | 4 |
| **Total** | **17** | **29** |
