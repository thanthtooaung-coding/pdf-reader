# PDF Reader Backend

Monolithic Go backend for the PDF Reader app. Follows the same layered layout as `cmm-backend` (`handler → service → repository → models`, with `mapper / request / response` for I/O).

Schema and flows are derived from [documentation/ERD.jpg](../documentation/ERD.jpg), [System Flow.jpg](../documentation/System%20Flow.jpg), and [User Flow.jpg](../documentation/User%20Flow.jpg).

## Stack

- Go 1.24, [Fiber v2](https://gofiber.io)
- [GORM v2](https://gorm.io) on PostgreSQL (NeonDB)
- Redis for OTP (register + login)
- Local filesystem storage for uploaded PDFs

## Project layout

```
backend/
├── cmd/main.go                 # Fiber bootstrap + graceful shutdown
├── db/init.sql                 # PostgreSQL schema (manual init option)
├── docker-compose.yml          # Redis + backend (DB via NeonDB in .env)
├── .env.example
└── internal/
    ├── config/
    ├── database/
    ├── models/                 # users, roles, workspaces, files, comments, ai_jobs
    ├── repository/
    ├── service/
    ├── handler/
    ├── routes/
    ├── mapper/
    ├── request/
    ├── response/
    ├── middleware/             # JWT auth
    └── otp/                    # Redis OTP store
└── pkg/
    ├── email/                  # Resend OTP mailer
    ├── storage/                # Local PDF storage
    └── utils/
```

## Prerequisites

- Go **1.24+**
- Docker + Docker Compose (optional, runs Redis + backend)
- [NeonDB](https://neon.tech) (or any PostgreSQL) — set `DATABASE_URL` in `.env`
- Redis 7 (local or via docker compose)

## Quick start (Docker)

```bash
cd backend
cp .env.example .env
# Set DATABASE_URL to your NeonDB connection string in .env
docker compose up --build
```

API: `http://localhost:8080`  
Health: `GET /healthz`

## Quick start (local)

```bash
cd backend
cp .env.example .env
# Set DATABASE_URL in .env (NeonDB). Start Redis locally or: docker compose up redis

go mod tidy
go run ./cmd
```

Set `DB_AUTO_MIGRATE=true` to let GORM create tables, or apply `db/init.sql` manually and keep `DB_AUTO_MIGRATE=false`.

With `OTP_DEBUG_RETURN=true` and no Resend API key, OTP codes are returned in the API response for local development.

## API overview

| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/api/v1/auth/register/request` | No | Start registration, send OTP |
| POST | `/api/v1/auth/register/verify` | No | Verify OTP, create account + JWT |
| POST | `/api/v1/auth/register/resend` | No | Resend registration OTP |
| POST | `/api/v1/auth/login` | No | Validate credentials, send OTP |
| POST | `/api/v1/auth/login/verify` | No | Verify login OTP, return JWT |
| GET | `/api/v1/users/me` | JWT | Current user profile |
| POST | `/api/v1/workspaces` | JWT | Create workspace (history entry) |
| GET | `/api/v1/workspaces` | JWT | List workspace history |
| GET | `/api/v1/workspaces/:id` | JWT | Workspace detail |
| POST | `/api/v1/workspaces/:workspaceId/files` | JWT | Upload PDF (`multipart/form-data`, field `file`) |
| GET | `/api/v1/workspaces/:workspaceId/files` | JWT | List files in workspace |
| GET | `/api/v1/files/:id` | JWT | File metadata |
| GET | `/api/v1/files/:id/download` | JWT | Download PDF |
| POST | `/api/v1/files/:fileId/comments` | JWT | Add comment on file |
| GET | `/api/v1/files/:fileId/comments` | JWT | List comments |
| POST | `/api/v1/ai-jobs` | JWT | Create AI job (`TRANSLATE`, `SUMMARIZE`, `COMMENT`) |
| GET | `/api/v1/ai-jobs/:id` | JWT | Poll AI job status/result |
| GET | `/api/v1/workspaces/:workspaceId/ai-jobs` | JWT | List AI jobs for workspace |

All JSON responses use the envelope:

```json
{ "success": true, "message": "...", "data": { ... } }
```

## Default admin

When `SEED_DEFAULT_ADMIN=true`, boot seeds roles (`ADMIN`, `USER`) and an admin account from `DEFAULT_ADMIN_EMAIL` / `DEFAULT_ADMIN_PASSWORD`.

## AI jobs

`TRANSLATE` and `SUMMARIZE` jobs extract text from the uploaded PDF and send it to the OpenAI Chat Completions API.

Configure in `.env`:

```env
OPENAI_API_KEY=sk-...
OPENAI_MODEL=gpt-4o-mini
# optional, for Azure OpenAI or compatible gateways:
OPENAI_BASE_URL=https://api.openai.com/v1
```

- **TRANSLATE** — set `input` to the target language (e.g. `my`, `English`, `Japanese`).
- **SUMMARIZE** — `input` is optional.

If `OPENAI_API_KEY` is not set, translate/summarize fall back to stub output (used in CI). `COMMENT` jobs still use a local stub.

## Tests

API integration tests live in `tests/integration/` and cover every route (auth, users, workspaces, files, comments, AI jobs).

They require PostgreSQL and Redis. CI uses a local Postgres service container; for local dev set `DATABASE_URL` to your NeonDB string in `.env`.

Local run:

```bash
cd backend
# DATABASE_URL from .env (NeonDB)
export REDIS_HOST=localhost
export REDIS_PORT=6379
export DB_AUTO_MIGRATE=true
export OTP_DEBUG_RETURN=true
export JWT_SECRET=test-secret

go test ./tests/integration/... -v -count=1
```

If `DATABASE_URL` is unset, the integration package exits successfully without running tests.
