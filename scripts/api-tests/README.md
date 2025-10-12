# API Test Scripts

Reusable bash scripts to exercise the Price Story API. Scripts create sample data, verify responses, and clean up, so they can be safely re-run.

## Prerequisites

- `curl`
- `jq`
- API server running at `http://localhost:8080` (see container dev below)

## Configuration

- `API_BASE` (optional): Base API URL. Defaults to `http://localhost:8080/api`.

## Run All Endpoint Tests

Using the Makefile target:

```bash
make api-test
```

Or directly:

```bash
bash scripts/api-tests/run_all.sh
```

## Run Domain-Specific Tests

Each script sets up its own prerequisites and cleans up afterwards. Attribution for create/update operations is handled by middleware injecting a system user.

```bash
# Users CRUD
bash scripts/api-tests/users.sh

# Products CRUD
bash scripts/api-tests/products.sh

# Price Stories CRUD (creates a temporary product)
bash scripts/api-tests/price_stories.sh

# Price Points CRUD (creates a temporary product)
bash scripts/api-tests/price_points.sh
```

Or use Makefile shortcuts:

```bash
make api-test-users
make api-test-products
make api-test-price-stories
make api-test-price-points
```

## Flow in run_all.sh

- Health check
- Create and update a product
- Create and update a price story
- Create and update a price point
- List products, price stories, and price points
- Cleanup (delete created resources)

Each run generates unique emails and names, allowing repeated executions.