#!/usr/bin/env bash
set -euo pipefail

# Seed bogus products into the API for frontend testing
# Usage:
#   bash scripts/seed-products.sh [COUNT]
#   COUNT: optional number of products to create (default 8)

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

COUNT="${1:-8}"

banner "Seeding ${COUNT} bogus products"
log_env
trap 'err_exit $? $LINENO' ERR

random_suffix() {
  if command -v uuidgen >/dev/null 2>&1; then
    uuidgen | tr 'A-Z' 'a-z' | cut -c1-8
  else
    date +%s
  fi
}

create_product() {
  local name="$1"; shift
  local url="https://example.com/products/${name}"
  # Use Lorem Picsum for deterministic placeholder images per product name
  local img="https://picsum.photos/seed/${name}/400/400"
  local desc="Bogus product ${name} for UI testing"

  local payload
  payload=$(jq -nc \
    --arg name "$name" \
    --arg url "$url" \
    --arg img "$img" \
    --arg desc "$desc" \
    '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')

  local created
  created=$(http POST "/products" "$payload" | json)
  echo "$created" | jq -r '.id // empty' | awk '{ if ($0 != "") print "Created product id=" $0; }'
}

title "Creating products"
for i in $(seq 1 "$COUNT"); do
  suffix=$(random_suffix)
  name="Bogus-Product-${suffix}-${i}"
  create_product "$name" || echo "Failed to create product $i" >&2
done

title "Listing products (latest)"
http GET "/products" | jq '.products | .[0:10]'

echo "\nSeed complete. API_BASE=${API_BASE}"