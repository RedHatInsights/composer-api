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
DB_RDSCERT=$(jq -r '.database.rdsCa // empty' "$ACG_CONFIG")

for var in DB_HOST DB_PORT DB_NAME DB_USER DB_PASS; do
  val="${!var}"
  if [[ -z "$val" || "$val" == "null" ]]; then
    echo "ERROR: Required config value $var is not set in $ACG_CONFIG"
    exit 1
  fi
done

ENCODED_USER=$(printf '%s' "$DB_USER" | jq -Rr @uri)
ENCODED_PASS=$(printf '%s' "$DB_PASS" | jq -Rr @uri)

SSL_PARAMS="sslmode=${DB_SSLMODE}"
if [[ -n "$DB_RDSCERT" ]]; then
  CA_PATH="/tmp/db-ca.crt"
  printf '%s\n' "$DB_RDSCERT" > "$CA_PATH"
  chmod 0600 "$CA_PATH"
  SSL_PARAMS="${SSL_PARAMS}&sslrootcert=${CA_PATH}"
fi

DATABASE_URL="postgres://${ENCODED_USER}:${ENCODED_PASS}@${DB_HOST}:${DB_PORT}/${DB_NAME}?${SSL_PARAMS}"

echo "Running database migrations..."
migrate -path /migrations -database "$DATABASE_URL" up
echo "Migrations completed successfully."
