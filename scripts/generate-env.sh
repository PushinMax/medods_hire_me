#!/bin/bash

ENV_FILE=".env"

generate_jwt_secret() {
  openssl rand -base64 64 | tr -d '\n'
}

if [ ! -f "$ENV_FILE" ]; then
  cat > "$ENV_FILE" <<EOF
# PostgreSQL
DB_USER=postgres
DB_PASSWORD=postgres

# JWT
JWT_ACCESS_SECRET=$(generate_jwt_secret)
EOF
  echo "Created .env file with generated secrets"
else
  echo ".env file already exists. Skipping generation."
fi