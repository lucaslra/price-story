#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

banner "API Tests: Price Points"
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

title "Price Points: actor injected via middleware (no user prerequisites)"
echo "Requests are attributed to the injected system user."

title "Price Points: create prerequisite product"
P_CREATE=$(jq -nc --arg name "$PRODUCT_NAME" --arg url "$PRODUCT_URL" --arg img "$PRODUCT_IMAGE_URL" --arg desc "$PRODUCT_DESCRIPTION" '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
PJSON=$(http POST "/products" "$P_CREATE")
PID=$(echo "$PJSON" | jq -r .id)

title "Price Points: create"
NOW=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
PP_CREATE=$(jq -nc --arg price "19.99" --arg ts "$NOW" --arg pid "$PID" '{price:($price|tonumber), timestamp:$ts, product_id:$pid}')
PP_JSON=$(http POST "/price-points" "$PP_CREATE")
echo "$PP_JSON" | jq .
PP_ID=$(echo "$PP_JSON" | jq -r .id)

title "Price Points: get by id"
http GET "/price-points/${PP_ID}" | jq .

title "Price Points: list"
http GET "/price-points" | jq .

title "Price Points: update"
PP_UPDATE=$(jq -nc --arg price "21.49" --arg ts "$NOW" --arg pid "$PID" '{price:($price|tonumber), timestamp:$ts, product_id:$pid}')
PP_UPDATED=$(http PUT "/price-points/${PP_ID}" "$PP_UPDATE")
echo "$PP_UPDATED" | jq .

title "Price Points: cleanup"
http DELETE "/price-points/${PP_ID}" && echo "Deleted price point ${PP_ID}"
http DELETE "/products/${PID}" && echo "Deleted product ${PID}"

title "Price Points domain tests completed"