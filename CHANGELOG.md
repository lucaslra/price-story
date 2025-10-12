# Changelog

All notable changes to this project are documented here.

## Unreleased

- Documentation overhaul: added `docs/Overview.md`, `docs/Setup.md`, `docs/Development.md`, `docs/Troubleshooting.md`, and `docs/API.md`
- Clarify frontend dev proxy behavior and Docker Compose proxy target
- Note on test strategy for frontend: mock API client to avoid relative `fetch` pitfalls

## 0.1.0 – Initial public version

- Backend: Users, Products, Price Stories, Price Points CRUD; health endpoint
- Embedded migrations and OpenAPI spec
- Frontend: React + Vite app with products listing
- Dev environment: Docker Compose with PostgreSQL, Go API, and Vite dev server
- Makefile tasks for build, run, tests, linting, and migrations