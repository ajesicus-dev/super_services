#!/bin/sh

. /deployments/vault/scripts/common.sh

VAULT_NODE1_ADDR="${VAULT_NODE1_ADDR:-http://vault-node1:8200}"
VAULT_NODE2_ADDR="${VAULT_NODE2_ADDR:-http://vault-node2:8200}"
UNSEAL_KEY=$(cat ${INIT_OUTPUT_DIR}/unseal.key)

# Function to unseal a node
unseal_node() {
  NODE_ADDR=$1
  echo "Unsealing node at $NODE_ADDR..."
  
  if curl -s ${NODE_ADDR}/v1/sys/health | grep '"sealed":false'; then
    echo "Node already unsealed."
    return
  fi

  curl --request PUT --data "{\"key\":\"${UNSEAL_KEY}\"}" ${NODE_ADDR}/v1/sys/unseal
  echo "Node unsealed."
}

unseal_node "$VAULT_NODE1_ADDR"
unseal_node "$VAULT_NODE2_ADDR"

echo "All Vault nodes unsealed."