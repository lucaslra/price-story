# API Reference

The OpenAPI specification for the backend lives at:

- `cmd/price-story/openapi/openapi.yaml`

## Viewing the API

You can view or interact with the API spec using:

- Swagger Editor (online): upload the YAML file
- VS Code OpenAPI plugins
- `redocly` or `swagger-ui` tooling

## Key Endpoints

- `GET /api/health` — health check
- Users: `GET /api/users`, `POST /api/users`, `GET /api/users/{id}`, `PUT /api/users/{id}`
- Products: `GET /api/products`, `POST /api/products`, `GET /api/products/{id}`, `PUT /api/products/{id}`
- Price Stories: `GET /api/price-stories`, `POST /api/price-stories`, `GET /api/price-stories/{id}`, `PUT /api/price-stories/{id}`
- Price Points: `GET /api/price-points`, `POST /api/price-points`, `GET /api/price-points/{id}`, `PUT /api/price-points/{id}`

## Request Attribution

Create and update operations are attributed to a server-injected actor. Do not send `created_by_user_id` or `updated_by_user_id`; these are set by middleware and returned as `created_by_user` and `updated_by_user`.