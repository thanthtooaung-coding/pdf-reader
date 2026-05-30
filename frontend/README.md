# DevDoc Companion (Frontend)

Developer-focused PDF reading assistant with translation, summarization, and file comments — wired to the Go backend.

## Prerequisites

- Node.js 18+
- Backend running at `http://localhost:8080` (see [`../backend/README.md`](../backend/README.md))

For local dev, start the backend with Docker:

```bash
cd backend
cp .env.example .env
docker compose up --build
```

Ensure `OTP_DEBUG_RETURN=true` in backend `.env` so OTP codes appear in API responses during development.

Default seeded admin (when `SEED_DEFAULT_ADMIN=true`):

- Email: `admin@pdfreader.local`
- Password: `PdfAdmin123@`
- Username: check backend seed or use the username you registered with

## Tech Stack

- React + TypeScript (Vite)
- Tailwind CSS, Framer Motion
- react-pdf, Axios, TanStack React Query, React Hook Form, Zod
- react-router-dom

## Getting Started

```bash
npm install
cp .env.example .env   # optional; Vite proxy works without it
npm run dev
```

Open `http://localhost:5173`. The Vite dev server proxies `/api` to `http://localhost:8080`.

## Environment

| Variable | Description |
|----------|-------------|
| `VITE_API_BASE_URL` | Backend base URL. Leave empty in dev to use the Vite proxy. Set to `http://localhost:8080` if not using the proxy. |

## App Flow

1. **Login / Register** — OTP verification, JWT stored locally
2. **Workspaces** — create and select a workspace
3. **Files** — upload PDFs to a workspace
4. **Reader** — view PDF, run document-level Translate/Summarize (async AI jobs), save selection as file comments

## Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start development server |
| `npm run build` | Type-check and build for production |
| `npm run preview` | Preview production build |
