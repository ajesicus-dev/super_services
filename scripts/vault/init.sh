#!/bin/sh

. /scripts/vault/common.sh

# Check if Vault is already initialized
if vault_get "sys/init" | grep '"initialized":true'; then
  echo "Vault already initialized."
  exit 0
fi

echo "Initializing Vault..."

response=$(curl --silent --request PUT --data '{"secret_shares":1,"secret_threshold":1}' ${VAULT_ADDR}/v1/sys/init)

unseal_key=$(echo $response | jq -r '.keys_base64[0]')
root_token=$(echo $response | jq -r '.root_token')

echo "Saving unseal key and root token..."

echo "$unseal_key" > ${INIT_OUTPUT_DIR}/unseal.key
echo "$root_token" > ${INIT_OUTPUT_DIR}/root.token

echo "Vault initialized successfully."
