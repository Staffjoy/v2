#!/usr/bin/env bash
set -euo pipefail

ENV_FILE="${1:-.env}"
if [[ ! -f "$ENV_FILE" ]]; then
  echo "error: $ENV_FILE not found. Create it from .env.example"
  exit 1
fi

required=(
  ENV
  DATABASE_URL
  SESSION_SECRET
  LOG_LEVEL
)

missing=0
for key in "${required[@]}"; do
  if ! grep -Eq "^${key}=.+" "$ENV_FILE"; then
    echo "missing required variable: ${key}"
    missing=1
  fi
done

if [[ "$missing" -ne 0 ]]; then
  echo "env validation failed"
  exit 1
fi

echo "env validation passed for $ENV_FILE"
