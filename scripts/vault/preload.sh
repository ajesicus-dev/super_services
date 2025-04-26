#!/bin/sh

. /scripts/vault/common.sh

if [ -z "$VAULT_ROOT_TOKEN" ]; then
  echo "Error: VAULT_ROOT_TOKEN not found!"
  exit 1
fi

# Wait for Vault to be unsealed and active (status 200)
wait_for_active_vault "$VAULT_ADDR"

# List all enabled secret engines
if vault_get "sys/mounts" | grep '"secret/":'; then
  echo "KV engine already enabled at 'secret/'."
else
  echo "Enabling KV secrets engine at path 'secret/'..."

  vault_post "sys/mounts/secret" '{"type": "kv", "options": {"version": "2"}}'

  echo "KV engine enabled successfully."
fi

echo "Checking if secrets already exist..."
if vault_get "secret/data/super_services/config" | grep '"data"'; then
  echo "Secrets already exist. Skipping preload."
  exit 0
fi

echo "Preloading secrets..."
vault_post "secret/data/super_services/config" '{"data":{"username":"admin","password":"supersecret"}}'

# ==================
# Load your secrets
# ==================

vault_write "auth/config" '{
  "jwt_secret": "super-secret-jwt-key",
  "oauth_client_id": "client-id",
  "oauth_client_secret": "client-secret"
}'

vault_write "billing/config" '{
  "stripe_api_key": "sk_test_xxx",
  "billing_secret": "billing-secret-key"
}'

vault_write "database/creds" '{
  "username": "postgres",
  "password": "postgres"
}'

vault_write "gateway/config" '{
  "internal_api_key": "gateway-internal-secret"
}'

vault_write "search/config" '{
  "elastic_api_key": "elastic-key"
}'

vault_write "notification/config" '{
  "sse_secret": "sse-super-secret"
}'

echo "Vault preload completed."
