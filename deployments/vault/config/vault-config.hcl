# Storage: Raft (integrated storage for production-grade setups)
storage "raft" {
  path    = "/vault/data"
  
  retry_join {
    leader_api_addr = "http://vault-node1:8200"
  }

  retry_join {
    leader_api_addr = "http://vault-node2:8200"
  }
}

# Listener: TCP (HTTP for dev; in prod, use TLS)
listener "tcp" {
  address     = "0.0.0.0:8200"
  cluster_address  = "0.0.0.0:8201"
  tls_disable = true  # Only for local development!
}

# Default cluster configuration
cluster_name = "super_services_vault_cluster"

# disable locking memory (no swap) for secrets.
disable_mlock = true

# Enable Vault's built-in UI
ui = true

# Log Level
log_level = "info"

# Telemetry: Prometheus metrics settings
telemetry {
  prometheus_retention_time = "30s"
  disable_hostname = true
}
