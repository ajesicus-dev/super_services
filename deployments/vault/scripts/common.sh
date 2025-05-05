#!/bin/sh

set -e

VAULT_ADDR="${VAULT_NODE1_ADDR:-http://vault-node1:8200}"
INIT_OUTPUT_DIR="/vault/files"
ROOT_TOKEN_FILE="${INIT_OUTPUT_DIR}/root.token"

wait_for_active_vault() {
  local addr=$1
  echo "Waiting for Vault at ${addr} to become active..."
  until [ "$(curl -s -o /dev/null -w "%{http_code}" ${addr}/v1/sys/health)" = "200" ]; do
    sleep 2
  done
  echo "Vault at ${addr} is active."
}

# Read the root token
read_root_token() {
  if [ ! -f "${ROOT_TOKEN_FILE}" ]; then
    echo "Root token file not found: ${ROOT_TOKEN_FILE}"
    exit 1
  fi
  cat "${ROOT_TOKEN_FILE}"
}

# Vault API GET
vault_get() {
  local path="$1"
  curl -s --header "X-Vault-Token: $(read_root_token)" "${VAULT_ADDR}/v1/${path}"
}

# Vault API POST
vault_post() {
  local path="$1"
  local data="$2"
  curl -s --header "X-Vault-Token: $(read_root_token)" --request POST --data "${data}" "${VAULT_ADDR}/v1/${path}"
}

vault_curl() {
  local method=$1
  local path=$2
  local data=$3
  curl --silent --fail --header "X-Vault-Token: ${VAULT_ROOT_TOKEN}" \
       --request "${method}" \
       --data "${data}" \
       "${VAULT_NODE1_ADDR}${path}"
}

# Helper: write secret
vault_write() {
  path=$1
  json=$2
  echo "Writing secret at path: $path"
  
  response=$(curl --silent --output /dev/stderr --write-out "%{http_code}" \
                  --header "X-Vault-Token: ${VAULT_TOKEN}" \
                  --request POST \
                  --data "{\"data\":${json}}" \
                  "${VAULT_ADDR}/v1/secret/data/${path}")
  
  if [ "$response" -ne 200 ]; then
    echo "Error: Failed to write secret at path: $path (HTTP $response)"
    exit 1
  fi
}