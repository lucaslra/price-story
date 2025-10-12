#!/usr/bin/env bash
set -euo pipefail

# Base API URL, override with API_BASE env var, e.g.
# API_BASE=http://localhost:8080/api
API_BASE="${API_BASE:-http://localhost:8080/api}"

# Verbose logging toggles (override via env)
VERBOSE="${VERBOSE:-1}"
TRACE_HTTP="${TRACE_HTTP:-0}"

timestamp() {
  date -u +"%Y-%m-%dT%H:%M:%SZ"
}

banner() {
  echo ""
  echo "=== $(timestamp) | $1 ==="
  echo "API_BASE=${API_BASE}"
  echo ""
}

log_env() {
  echo "curl: $(curl --version | head -n1)"
  echo "jq: $(jq --version)"
}

err_exit() {
  local code="$1"; shift || true
  local line="$1"; shift || true
  echo "ERROR: exit code ${code} at line ${line}" >&2
}

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "Missing required command: $1" >&2
    exit 1
  fi
}

need_cmd curl
need_cmd jq

json() {
  # Print compact JSON from jq input or passthrough
  jq -c . 2>/dev/null || cat
}

http() {
  # http METHOD PATH [JSON]
  local method="$1"; shift
  local path="$1"; shift
  local url="${API_BASE}${path}"
  if [ $# -gt 0 ]; then
    # Sanitize payload: drop any *_id fields that are empty strings
    # This reflects server expectations: omit or null, never "" for UUIDs
    local raw="$1"
    local body
    body=$(echo "$raw" | jq -c 'if type=="object" then with_entries(select(((.key|endswith("_id")) and (.value=="")) | not)) else . end')
    if [ "${TRACE_HTTP}" = "1" ]; then
      >&2 echo "HTTP ${method} ${url}"
      >&2 echo "Payload: ${body}"
    fi
    curl -fsS -H "Content-Type: application/json" -X "$method" "$url" -d "$body"
  else
    if [ "${TRACE_HTTP}" = "1" ]; then
      >&2 echo "HTTP ${method} ${url}"
    fi
    curl -fsS -H "Content-Type: application/json" -X "$method" "$url"
  fi
}

title() {
  echo ""; echo "==> $(timestamp) $1"; echo ""
}