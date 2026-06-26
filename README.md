# Rifa Online

Plataforma SaaS para **criação e gestão de rifas online**, com pagamento via PIX/cartão pela
[InfinitePay](https://www.infinitepay.io/). Organizadores assinam o serviço para criar rifas e
realizar sorteios; participantes compram números sem precisar de conta.

- **Frontend:** Vue 3 + TypeScript + Vite + Tailwind CSS
- **Backend:** Go (Chi) + MongoDB + Redis
- **Pagamentos:** InfinitePay Checkout (PIX e cartão)

> Documentação de produto e técnica detalhada em [`docs/`](#-documentação). Tutoriais de uso para
> o usuário final ficam dentro do app, na **Central de Ajuda** (`/ajuda`).

---

## 📁 Estrutura do repositório

```
.
├── backend/                 # API em Go
│   ├── cmd/server/main.go   # ponto de entrada (rotas, bootstrap)
│   ├── internal/
│   │   ├── handler/         # camada HTTP
│   │   ├── service/         # regras de negócio
│   │   ├── repository/      # acesso ao MongoDB
│   │   ├── model/           # modelos de dados
│   │   ├── middleware/      # auth, admin, assinatura, logging
│   │   ├── migrations/      # criação de índices
│   │   └── config/          # carregamento de env
│   └── pkg/infinitepay/     # cliente HTTP da InfinitePay
├── frontend/                # SPA em Vue 3
│   └── src/
│       ├── pages/           # telas (inclui legal/ e Central de Ajuda)
│       ├── components/      # componentes reutilizáveis
│       ├── layouts/         # Navbar, Footer
│       ├── data/help.ts     # conteúdo da Central de Ajuda
│       ├── stores/          # Pinia
│       └── router/          # rotas
├── docs/                    # PRD, Tech Spec, Tasks
├── docker-compose.yml
└── Makefile
```

---

## ✅ Pré-requisitos

- **Go** 1.23+
- **Node.js** 20+ e npm
- **Docker** + Docker Compose (para MongoDB/Redis ou para subir tudo)

---

## 🚀 Como rodar

### Opção A — Tudo via Docker (mais simples)

```bash
cp .env.example .env          # ajuste INFINITEPAY_HANDLE e JWT_SECRET
make docker-up                # sobe mongodb, redis, backend e frontend
```

- Frontend: http://localhost:5173
- Backend:  http://localhost:8080  (health: `/health`)

### Opção B — Desenvolvimento local (backend e frontend nativos)

```bash
make setup                    # copia backend/.env, instala deps do front, go mod tidy
make db-up                    # sobe apenas mongodb e redis via Docker
```

> ⚠️ **Atenção às portas:** o `docker-compose` expõe o MongoDB em **27018** e o Redis em **6380**
> no host (para não conflitar com instâncias locais). Ao rodar o backend nativo, ajuste o
> `backend/.env`:
>
> ```env
> MONGODB_URI=mongodb://localhost:27018
> REDIS_URI=redis://localhost:6380
> ```

```bash
make dev-backend              # API em :8080 (ou `make dev-backend-hot` com air)
make dev-frontend             # Vite em :5173
```

Usuário admin padrão é criado automaticamente: **admin@email.com** / **123456** (apenas dev).

---

## 🔐 Variáveis de ambiente (backend)

| Variável | Padrão | Descrição |
|----------|--------|-----------|
| `PORT` | `8080` | Porta da API |
| `MONGODB_URI` | `mongodb://localhost:27017` | Conexão MongoDB |
| `MONGODB_DB_NAME` | `rifaonline` | Nome do banco |
| `REDIS_URI` | `redis://localhost:6379` | Conexão Redis (lock de reservas) |
| `JWT_SECRET` | `change-me` | Segredo dos tokens JWT — **troque em produção** |
| `APP_ENV` | `development` | `development` habilita rotas de dev (ex.: `dev-checkout`) |
| `INFINITEPAY_HANDLE` | _(vazio)_ | Handle da plataforma (recebe as assinaturas) |
| `INFINITEPAY_BASE_URL` | `https://api.checkout.infinitepay.io` | Base da API InfinitePay |
| `FRONTEND_URL` | `http://localhost:5173` | URL do front (CORS + base do `webhook_url`) |
| `LOG_LEVEL` | `info` | `debug` \| `info` \| `warn` \| `error` |
| `LOG_FORMAT` | `text` | `text` \| `json` |

O frontend usa `VITE_API_URL` (opcional) para apontar a base da API; vazio = mesma origem (`/api/v1`).

---

## 🛠️ Comandos úteis (Makefile)

| Comando | O que faz |
|---------|-----------|
| `make dev-backend` / `make dev-frontend` | Sobe backend / frontend em modo dev |
| `make build` | Compila backend (`bin/server`) e frontend (`dist/`) |
| `make test` | Roda testes de backend (`go test`) e frontend (`vitest`) |
| `make lint` | `go vet` + `vue-tsc`/eslint |
| `make docker-up` / `make docker-down` | Sobe / derruba toda a stack |
| `make db-up` / `make db-down` | Sobe / derruba apenas MongoDB e Redis |
| `make clean` | Remove artefatos de build |

---

## 🔌 API

Base: `/api/v1`. Autenticação por **JWT** (access + refresh); rotas admin exigem `role=ADMIN`.

A **referência completa de endpoints** (auth, perfil, assinatura, rifas, pagamentos, dashboard,
admin, contato e webhooks) está em [`docs/TECH_SPEC.md` §4](docs/TECH_SPEC.md). A **integração de
pagamento** (criação de checkout, webhook, `payment_check`, ciclo de vida) está em
[`docs/TECH_SPEC.md` §5](docs/TECH_SPEC.md) e, em linguagem acessível, na Central de Ajuda
(`/ajuda/integracao-pagamento-infinitepay`).

---

## 📚 Documentação

### Técnica / produto (`docs/`)
- [PRD.md](docs/PRD.md) — visão de produto, personas, funcionalidades e fluxos
- [TECH_SPEC.md](docs/TECH_SPEC.md) — stack, arquitetura, modelagem de dados, API e integração
- [TASKS.md](docs/TASKS.md) — acompanhamento de tarefas

### No aplicativo (usuário final)
- **Central de Ajuda** — `/ajuda` (tutoriais por perfil: organizador, participante, admin, pagamentos)
- **Contato** — `/contato` (canal de atendimento; mensagens lidas no painel admin)
- **Documentos legais** — `/termos-de-uso`, `/politica-de-privacidade`, `/politica-de-cookies`, `/termo-do-organizador`

---

## 🧱 Decisões de arquitetura (resumo)

- **Camadas:** `handler` → `service` → `repository`; `pkg/infinitepay` isola o cliente de pagamento.
- **Concorrência:** reserva de números via **lock no Redis** (`SetNX`, TTL) para evitar venda dupla.
- **Pagamentos:** a plataforma **não custodia valores** nem **armazena dados de cartão** — tudo é
  processado pela InfinitePay e roteado pelo `handle` do organizador.
- **Confirmação idempotente:** webhook + polling (`payment_check`), com deduplicação por `eventId`.

Para detalhes, veja [`docs/TECH_SPEC.md`](docs/TECH_SPEC.md).
