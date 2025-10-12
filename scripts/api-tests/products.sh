#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

banner "API Tests: Products"
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

title "Products: actor injected via middleware (no user prerequisites)"
echo "Requests are attributed to the injected system user."

title "Products: create"
CREATE=$(jq -nc --arg name "$PRODUCT_NAME" --arg url "$PRODUCT_URL" --arg img "$PRODUCT_IMAGE_URL" --arg desc "$PRODUCT_DESCRIPTION" '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
PJSON=$(http POST "/products" "$CREATE")
echo "$PJSON" | jq .
PID=$(echo "$PJSON" | jq -r .id)

title "Products: get by id"
http GET "/products/${PID}" | jq .

title "Products: list"
http GET "/products" | jq .

title "Products: update"
UPDATE=$(jq -nc --arg name "$PRODUCT_NAME" --arg url "$PRODUCT_URL" --arg img "$PRODUCT_IMAGE_URL" --arg desc "$PRODUCT_DESCRIPTION" '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
UPDATED=$(http PUT "/products/${PID}" "$UPDATE")
echo "$UPDATED" | jq .

title "Products: cleanup"
http DELETE "/products/${PID}" && echo "Deleted product ${PID}"

title "Products domain tests completed"