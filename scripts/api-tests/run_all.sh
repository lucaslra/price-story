#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

banner "API Tests: Full End-to-End Suite"
log_env
trap 'err_exit $? $LINENO' ERR

random_suffix() {
  if command -v uuidgen >/dev/null 2>&1; then
    uuidgen | tr 'A-Z' 'a-z' | cut -c1-8
  else
    date +%s
  fi
}

PRODUCT_NAME="Widget-$(random_suffix)"
PRODUCT_URL="https://example.com/products/${PRODUCT_NAME}"
PRODUCT_IMAGE_URL="https://example.com/images/${PRODUCT_NAME}.png"
PRODUCT_DESCRIPTION="Demo product ${PRODUCT_NAME}"

title "Health check"
http GET "/health" | jq .

title "Actor injection via middleware (no user creation)"
echo "Requests are attributed to the injected system user."

title "Create product"
PRODUCT_CREATE=$(jq -nc \
  --arg name "$PRODUCT_NAME" \
  --arg url "$PRODUCT_URL" \
  --arg img "$PRODUCT_IMAGE_URL" \
  --arg desc "$PRODUCT_DESCRIPTION" \
  '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
PRODUCT_JSON=$(http POST "/products" "$PRODUCT_CREATE")
echo "$PRODUCT_JSON" | jq .
PRODUCT_ID=$(echo "$PRODUCT_JSON" | jq -r .id)

title "Update product"
PRODUCT_UPDATE=$(jq -nc \
  --arg name "$PRODUCT_NAME" \
  --arg url "$PRODUCT_URL" \
  --arg img "$PRODUCT_IMAGE_URL" \
  --arg desc "$PRODUCT_DESCRIPTION" \
  '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
PRODUCT_UPDATED=$(http PUT "/products/${PRODUCT_ID}" "$PRODUCT_UPDATE")
echo "$PRODUCT_UPDATED" | jq .

title "Create price story"
PS_CREATE=$(jq -nc --arg pid "$PRODUCT_ID" '{product_id:$pid}')
PS_JSON=$(http POST "/price-stories" "$PS_CREATE")
echo "$PS_JSON" | jq .
PS_ID=$(echo "$PS_JSON" | jq -r .id)

title "Update price story"
PS_UPDATE=$(jq -nc --arg pid "$PRODUCT_ID" '{product_id:$pid}')
PS_UPDATED=$(http PUT "/price-stories/${PS_ID}" "$PS_UPDATE")
echo "$PS_UPDATED" | jq .

title "Create price point"
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
PP_CREATE=$(jq -nc --arg price "19.99" --arg ts "$NOW" --arg pid "$PRODUCT_ID" '{price:($price|tonumber), timestamp:$ts, product_id:$pid}')
PP_JSON=$(http POST "/price-points" "$PP_CREATE")
echo "$PP_JSON" | jq .
PP_ID=$(echo "$PP_JSON" | jq -r .id)

title "Update price point"
PP_UPDATE=$(jq -nc --arg price "21.49" --arg ts "$NOW" --arg pid "$PRODUCT_ID" '{price:($price|tonumber), timestamp:$ts, product_id:$pid}')
PP_UPDATED=$(http PUT "/price-points/${PP_ID}" "$PP_UPDATE")
echo "$PP_UPDATED" | jq .

title "List endpoints"
echo "Products:" && http GET "/products" | jq .
echo "Price Stories:" && http GET "/price-stories" | jq .
echo "Price Points:" && http GET "/price-points" | jq .

title "Cleanup (delete in reverse dependency order)"
http DELETE "/price-points/${PP_ID}" && echo "Deleted price point ${PP_ID}"
http DELETE "/price-stories/${PS_ID}" && echo "Deleted price story ${PS_ID}"
http DELETE "/products/${PRODUCT_ID}" && echo "Deleted product ${PRODUCT_ID}"

title "All endpoint tests completed successfully"