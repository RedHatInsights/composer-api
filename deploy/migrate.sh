#!/bin/bash
set -euo pipefail

ACG_CONFIG="${ACG_CONFIG:-/cdapp/cdappconfig.json}"

if [ ! -f "$ACG_CONFIG" ]; then
  echo "ERROR: Clowder config not found at $ACG_CONFIG"
  exit 1
fi

DB_HOST=$(jq -r '.database.hostname' "$ACG_CONFIG")
DB_PORT=$(jq -r '.database.port' "$ACG_CONFIG")
DB_NAME=$(jq -r '.database.name' "$ACG_CONFIG")
DB_USER=$(jq -r '.database.username' "$ACG_CONFIG")
DB_PASS=$(jq -r '.database.password' "$ACG_CONFIG")
DB_SSLMODE=$(jq -r '.database.sslMode // "disable"' "$ACG_CONFIG")

DATABASE_URL="postgres://${DB_USER}:${DB_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"

echo "Running database migrations..."
migrate -path /migrations -database "$DATABASE_URL" up
echo "Migrations completed successfully."
