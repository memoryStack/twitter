# Project prompts

## Prompt 1 — Go backend scaffold

Create a new Go project in a folder named `backend`. The backend is a Go server with these dependencies:

1. **go fiber** — routing
2. **gorm** — database
3. **godotenv** — environment variables
4. **compiledaemon** — file watching and server restarts

**Folders**

1. `controllers` — route handlers and business logic
2. `helpers` — reusable helpers
3. `initializers` — env loading, DB connection, etc., invoked from `main`
4. `models` — data models
5. `middlewares` — HTTP middlewares (with recommendations for a robust API)

**Files**

1. `main.go` — routes, DB init, server start

**Environments:** `development` and `production`, each with its own env file loaded when that mode runs.

**Database:** PostgreSQL.

---

## Prompt 2 — Documentation and simpler configuration

1. Record all prompts in a file in the project (`PROMPTS.md`).
2. Avoid fallback-heavy configuration (`BACKEND_ROOT`, `APP_ENV`, many optional env knobs). Prefer a **simple, conservative** setup.

**What was implemented**

- **Removed:** `BACKEND_ROOT`, reading env from anywhere other than process CWD, `APP_ENV` and implicit default environment, dual DB config (`DATABASE_URL` *or* many `DB_*` fields), pool tuning via env, HTTP timeout via env, dev-time `CORS_ALLOW_ORIGINS` env.
- **Kept:** Required CLI flag `-env`, single `DATABASE_URL` for Postgres, fixed connection pool in code, fixed 30s handler timeout, dev CORS `*` in code, production CORS from env only (required at startup).

---

## Prompt 3 — Committed env files

Create real `.env.development` and `.env.production` (not only `.example`) with everything `LoadEnv` needs. Stop gitignoring those files so **other developers get the same baseline**; document that production files should hold **placeholders** until deploy, not real secrets in public repos.

**What was implemented**

- Added committed `backend/.env.development` and `backend/.env.production`.
- Removed those paths from `backend/.gitignore`.
- Removed separate `.env.*.example` files (single source of truth).
