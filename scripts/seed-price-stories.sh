#!/usr/bin/env bash
set -euo pipefail

# Seed bogus price stories linked to existing products
# Usage:
#   bash scripts/seed-price-stories.sh [COUNT]
#   COUNT: optional number of price stories to create (default 10)

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

COUNT="${1:-10}"

banner "Seeding ${COUNT} bogus price stories"
log_env
trap 'err_exit $? $LINENO' ERR

random_suffix() {
  if command -v uuidgen >/dev/null 2>&1; then
    uuidgen | tr 'A-Z' 'a-z' | cut -c1-8
  else
    date +%s
  fi
}

ensure_products() {
  local existing
  existing=$(http GET "/products" | jq -r '.products[]?.id') || true
  if [ -z "${existing}" ]; then
    title "No products found; creating a seed product"
    local suffix name url img desc payload created pid
    suffix="$(random_suffix)"
    name="Seed-Product-${suffix}"
    url="https://example.com/products/${name}"
    img="https://picsum.photos/seed/${name}/400/400"
    desc="Seed product ${name} for price stories"
    payload=$(jq -nc --arg name "$name" --arg url "$url" --arg img "$img" --arg desc "$desc" '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
    created=$(http POST "/products" "$payload")
    pid=$(echo "$created" | jq -r '.id')
    echo "Created product id=${pid}"
  fi
}

get_product_ids() {
  http GET "/products" | jq -r '.products[]?.id'
}

create_price_story() {
  local pid="$1"; shift
  local payload created id
  payload=$(jq -nc --arg pid "$pid" '{product_id:$pid}')
  created=$(http POST "/price-stories" "$payload" | json)
  id=$(echo "$created" | jq -r '.id // empty')
  if [ -n "$id" ]; then
    echo "Created price story id=${id} for product=${pid}"
  fi
}

ensure_products

PRODUCT_IDS=()
while IFS= read -r pid; do
  if [ -n "$pid" ]; then
    PRODUCT_IDS+=("$pid")
  fi
done < <(get_product_ids)
if [ ${#PRODUCT_IDS[@]} -eq 0 ]; then
  echo "No product IDs available after ensure." >&2
  exit 1
fi

title "Creating price stories"
for i in $(seq 1 "$COUNT"); do
  # Pick a product id round-robin
  idx=$(( (i-1) % ${#PRODUCT_IDS[@]} ))
  pid="${PRODUCT_IDS[$idx]}"
  create_price_story "$pid" || echo "Failed to create price story $i" >&2
done

title "Listing price stories (latest)"
http GET "/price-stories" | jq '.priceStories | .[0:10]'

echo "\nSeed complete. API_BASE=${API_BASE}"