#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
source "${ROOT_DIR}/scripts/api-tests/common.sh"

banner "API Tests: Users"
log_env
trap 'err_exit $? $LINENO' ERR

EMAIL="user-$(date +%s)@example.com"
PASS="password-$(date +%s)"

title "Users: create"
USER_JSON=$(http POST "/users" "$(jq -nc --arg email "$EMAIL" --arg pass "$PASS" '{email:$email, password_hash:$pass}')")
echo "$USER_JSON" | jq .
USER_ID=$(echo "$USER_JSON" | jq -r .id)

title "Users: get by id"
http GET "/users/${USER_ID}" | jq .

title "Users: list"
http GET "/users" | jq .

title "Users: update"
NEW_EMAIL="updated-${EMAIL}"
NEW_PASS="updated-${PASS}"
UPDATED_JSON=$(http PUT "/users/${USER_ID}" "$(jq -nc --arg email "$NEW_EMAIL" --arg pass "$NEW_PASS" '{email:$email, password_hash:$pass}')")
echo "$UPDATED_JSON" | jq .

title "Users: delete"
http DELETE "/users/${USER_ID}" && echo "Deleted user ${USER_ID}"

title "Users domain tests completed"