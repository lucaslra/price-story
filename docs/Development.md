# Development Guide

This guide outlines day-to-day workflows for backend and frontend development, testing, linting, and data seeding.

## Backend (Go)

- Build: `make build`
- Run (with Docker PostgreSQL): `make run`
- Tests: `make test`
- OpenAPI sync check: `make openapi-sync-check`
- Format: `make fmt`
- Lint (if `golangci-lint` installed): `make lint`

Migrations are embedded and run at startup. Manual control:

- Up: `make migrate-up`
- Down: `make migrate-down`
- Status: `make migrate-status`

## Frontend (React + Vite)

- Install: `make frontend-install`
- Dev server: `make frontend-dev`
- Tests (Vitest): `make frontend-test`
- Lint (ESLint): `make frontend-lint`
- Build: `make frontend-build`
- Preview build: `make frontend-preview`

Dev server runs on `http://localhost:5173` and proxies `/api/*` and `/health` to the backend.

## OpenAPI Synchronization

The router and handlers are validated against the OpenAPI spec to catch drift:

- Spec file: `cmd/price-story/openapi/openapi.yaml`
- Run validation: `make openapi-sync-check`

The validation ensures:
- Every route and method in the router is present in the spec, and vice versa.
- `POST`/`PUT` operations define a `requestBody`.
- `GET`/`PUT`/`DELETE` operations on `{id}` paths document the `id` path parameter.

## API Test Scripts

Reusable scripts exist under `scripts/api-tests/`. You can run all tests:

```bash
make api-test
```

Or domain-specific tests via:

```bash
make api-test-users
make api-test-products
make api-test-price-stories
make api-test-price-points
```

## Seeding Data

Seed helper scripts are available:

```bash
bash scripts/seed-products.sh [count]
bash scripts/seed-price-stories.sh
bash scripts/seed-price-points.sh
```

Set `API_BASE` to target a non-default server:

```bash
API_BASE=http://localhost:8080/api bash scripts/seed-products.sh 12
```

## Coding Style

- Go: follow `go fmt`; add comments for exported symbols where helpful.
- Frontend: TypeScript strict mode, React hooks dependency completeness, prefer functional components, avoid mutable shared state.
- Commit messages: imperative tone, concise summary, reference scopes where applicable.