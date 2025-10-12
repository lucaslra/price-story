# Price Story Overview

Price Story is a Go backend with a React + Vite frontend for tracking products and their price histories.

## Features

- Users: list, create, update
- Products: list, create, update
- Price Stories: list, create, update
- Price Points: list, create, update
- Health endpoint for basic system checks
- Dev-friendly setup using Docker Compose and Makefile targets
- Frontend dev proxy that routes `/api/*` to the backend during development

## Architecture

- Backend: Go 1.25, Gorilla Mux, PostgreSQL
- Frontend: React, TypeScript, Vite
- Infrastructure: Docker Compose for container-first development
- API Spec: OpenAPI (`cmd/price-story/openapi/openapi.yaml`)

## Project Layout

```
price-story/
├── cmd/price-story/          # Application entry point and embedded assets
├── pkg/                      # Backend packages (config, db, router, repository, models)
├── frontend/                 # Web UI
├── scripts/                  # CLI helpers and API test scripts
├── docs/                     # Documentation
└── Makefile                  # Build and dev workflows
```