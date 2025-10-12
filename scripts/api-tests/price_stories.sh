#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

banner "API Tests: Price Stories"
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

title "Price Stories: actor injected via middleware (no user prerequisites)"
echo "Requests are attributed to the injected system user."

title "Price Stories: create prerequisite product"
P_CREATE=$(jq -nc --arg name "$PRODUCT_NAME" --arg url "$PRODUCT_URL" --arg img "$PRODUCT_IMAGE_URL" --arg desc "$PRODUCT_DESCRIPTION" '{product_name:$name, product_url:$url, product_image_url:$img, product_description:$desc}')
PJSON=$(http POST "/products" "$P_CREATE")
PID=$(echo "$PJSON" | jq -r .id)

title "Price Stories: create"
PS_CREATE=$(jq -nc --arg pid "$PID" '{product_id:$pid}')
PS_JSON=$(http POST "/price-stories" "$PS_CREATE")
echo "$PS_JSON" | jq .
PS_ID=$(echo "$PS_JSON" | jq -r .id)

title "Price Stories: get by id"
http GET "/price-stories/${PS_ID}" | jq .

title "Price Stories: list"
http GET "/price-stories" | jq .

title "Price Stories: update"
PS_UPDATE=$(jq -nc --arg pid "$PID" '{product_id:$pid}')
PS_UPDATED=$(http PUT "/price-stories/${PS_ID}" "$PS_UPDATE")
echo "$PS_UPDATED" | jq .

title "Price Stories: cleanup"
http DELETE "/price-stories/${PS_ID}" && echo "Deleted price story ${PS_ID}"
http DELETE "/products/${PID}" && echo "Deleted product ${PID}"

title "Price Stories domain tests completed"