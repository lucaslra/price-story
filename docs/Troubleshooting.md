# Troubleshooting Guide

Common issues and resolutions for the Price Story system.

## Frontend fetch in tests fails with relative URLs

Symptom: Vitest errors like "Failed to parse URL from /api/..." or "Cannot read properties of undefined (reading 'ok')".

Cause: Node/JSDOM `fetch` implementations expect absolute URLs. The app uses relative paths (e.g., `/api/products`).

Fixes:
- Mock the API client (`@/api/client`) in tests to avoid reliance on `fetch`.
- Alternatively, polyfill `fetch` that supports relative URLs, or construct absolute URLs in tests.

## Dev server not proxying API requests

Symptom: Frontend fails to load data; network requests go to `http://localhost:5173/api/...` and return 404.

Cause: Vite proxy misconfiguration or missing proxy target.

Checks:
- Verify `frontend/vite.config.ts` includes proxies for `/api` and `/health`.
- In Docker Compose, ensure `VITE_API_PROXY_TARGET` points to `http://app:8080`.

## Port conflicts (8080 / 5173 / 5432)

Symptom: Docker services fail to start due to ports already in use.

Fixes:
- Stop other local services on those ports or edit `docker-compose.yml` to use alternate host ports.
- For PostgreSQL, ensure no local instance is already running on `5432`.

## Database connection failures

Symptom: API logs show failures to connect to PostgreSQL.

Checks:
- Confirm postgres container is healthy (`make db-logs`) and credentials match `docker-compose.yml`.
- Ensure `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, and `DB_SSLMODE` are correctly set.

## Migrations not applied

Symptom: API starts but queries fail due to missing tables.

Fixes:
- Restart the app to re-run embedded migration logic.
- Use `make migrate-up` to apply migrations manually.

## CORS errors (when bypassing dev proxy)

Symptom: Browser blocks requests to `http://localhost:8080/api` from `http://localhost:5173`.

Fixes:
- Use the dev proxy (`/api/*`) during development so requests are same-origin to the dev server.

## Seed scripts fail

Symptom: Scripts under `scripts/` error out or hit wrong server.

Fixes:
- Ensure backend is reachable; set `API_BASE` explicitly (e.g., `API_BASE=http://localhost:8080/api`).
- Check `jq` and `curl` are installed and available in PATH.