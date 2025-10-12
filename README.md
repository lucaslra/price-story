# Price Story API

A Go API for tracking and analyzing price histories of products.

## Project Structure

```
price-story/
├── cmd/
│   └── price-story/  # Application entry point
├── pkg/
│   ├── config/       # Application configuration
│   ├── database/     # Database connection and queries
│   ├── handlers/     # HTTP request handlers
│   ├── middleware/   # HTTP middleware components
│   ├── models/       # Data models
│   ├── repository/   # Data access layer
│   └── router/       # HTTP router setup
└── Makefile          # Build automation
```

## Getting Started

### Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose (for PostgreSQL)

### Installation

1. Clone the repository:
   ```
   git clone <repository-url>
   cd price-story
   ```

2. Install dependencies:
   ```
   go mod download
   ```

### Running the Application (Container-First)

The recommended way to develop is to run both the API and PostgreSQL in containers using Docker Compose.

```bash
# Build the app image
make docker-build

# Start app + PostgreSQL
make docker-up

# View app logs
make docker-logs

# Stop and remove containers
make docker-down
```

The server listens on `http://localhost:8080`. Migrations run automatically at startup.

### Request Attribution

Create and update operations are attributed to a server-injected actor via middleware. A persistent system user (default: `system@example.com`) is ensured at startup and its ID is injected into each request context. Clients should not send `created_by_user_id` or `updated_by_user_id` in payloads; these fields are set automatically and returned in responses as `created_by_user` and `updated_by_user`.

### Running Locally (Optional)

You can still run the app locally if preferred:

```bash
# Start PostgreSQL database
make db-up

# Run the application locally
make run
```

### Available Commands

#### Database Commands
- `make db-up` - Start PostgreSQL database
- `make db-down` - Stop PostgreSQL database  
- `make db-logs` - View database logs
- `make db-reset` - Reset database (removes all data)

#### Migration Commands
- `make migrate-up` - Apply all pending migrations
- `make migrate-down` - Rollback the last applied migration
- `make migrate-status` - Show current migration status

#### Application Commands
- `make docker-build` - Build the application image
- `make docker-up` - Start the application and PostgreSQL in containers
- `make docker-logs` - Tail application logs
- `make docker-down` - Stop and remove containers
- `make run` - Start PostgreSQL and run the application locally
- `make build` - Build the application binary
- `make test` - Run tests
- `make clean` - Clean build artifacts

## Frontend (React + Vite)

The project includes a small React + TypeScript frontend in `frontend/` that consumes the API and renders products.

### Quick Start

```bash
# Install dependencies
make frontend-install

# Start the Vite dev server (http://localhost:5173)
make frontend-dev
```

Alternatively, from the `frontend/` directory:

```bash
npm install
npm start
```

### Dev Proxy

- The frontend expects the API at `http://localhost:8080`.
- The dev server proxies requests from `/api/*` to the backend, configured in `frontend/vite.config.ts`.
- The API client uses `API_BASE = '/api'`.

### Helpful Commands

```bash
# In the repo root
make frontend-test     # Run Vitest
make frontend-lint     # Run ESLint
make frontend-build    # Type-check and build
make frontend-preview  # Preview the production build
```

### Seed Bogus Data

You can quickly seed some products to see real data in the UI:

```bash
# Backend should be running (via make run or make docker-up)
bash scripts/seed-products.sh         # default 8 products
bash scripts/seed-products.sh 12      # custom count

# If your API is on a different host/port
API_BASE=http://localhost:8080/api bash scripts/seed-products.sh 10
```

### Notes

- TypeScript is pinned to `5.5.x` to align with `@typescript-eslint` support.
- If `npm install` warns about vulnerabilities, they don’t affect dev usage. Run `npm audit fix` if desired.

## Environment Variables

Configure the application using these environment variables (PostgreSQL-only):

| Variable | Description | Default |
| --- | --- | --- |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `price_story_user` |
| `DB_PASSWORD` | PostgreSQL password | `price_story_password` |
| `DB_NAME` | Database name | `price_story` |
| `DB_SSLMODE` | SSL mode (`disable`, `require`, etc.) | `disable` |
| `PORT` | HTTP server port | `8080` |
| `READ_TIMEOUT` | HTTP read timeout | `15s` |
| `WRITE_TIMEOUT` | HTTP write timeout | `15s` |

Migrations are embedded and run automatically on startup. Use `make migrate-up`, `make migrate-down`, and `make migrate-status` for manual control.