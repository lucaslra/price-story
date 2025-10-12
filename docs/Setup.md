# Setup Guide

This guide covers prerequisites, installation, configuration, and starting the system.

## Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Node.js 18+ and npm (for the frontend)

## Installation

```bash
git clone <repository-url>
cd price-story
go mod download
```

## Container-First Development

Start PostgreSQL, the Go API, and the frontend dev server via Docker Compose:

```bash
make docker-build
make docker-up
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:5173`

To view logs:

```bash
make docker-logs     # app logs
make db-logs         # postgres logs
```

To stop:

```bash
make docker-down
```

## Local (Non-Container) Development

Run the backend against a local PostgreSQL:

```bash
make db-up           # starts postgres via docker
make run             # runs the Go app locally
```

Start the frontend dev server:

```bash
cd frontend
npm install
npm run dev
```

## Configuration

Backend environment variables:

| Variable | Description | Default |
| --- | --- | --- |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `price_story_user` |
| `DB_PASSWORD` | PostgreSQL password | `price_story_password` |
| `DB_NAME` | Database name | `price_story` |
| `DB_SSLMODE` | SSL mode | `disable` |
| `PORT` | HTTP server port | `8080` |
| `READ_TIMEOUT` | HTTP read timeout | `15s` |
| `WRITE_TIMEOUT` | HTTP write timeout | `15s` |

Frontend dev proxy:

- The dev server proxies `/api/*` and `/health` to the backend.
- Proxy target can be set via `VITE_API_PROXY_TARGET` (Docker Compose uses `http://app:8080`).
- The frontend API client uses `API_BASE = '/api'` during development.