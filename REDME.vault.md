# Vault Integration Guide

## Overview
Super Services uses HashiCorp Vault for managing secrets across all microservices in a secure, scalable, and production-ready way.
Vault handles:

- Service API keys
- Database credentials
- External service tokens (e.g., NATS, Redis, Postgres, Matrix Synapse)
- Configuration secrets for multi-tenant environments

### Cluster Setup

```less
 +------------------------+
  |   vault-node1 (server)  |
  +-----------+------------+
              |
              v
       +------+------+
       |   Raft Data |
       |  (/vault/data) |
       +------+------+
              ^
              |
  +-----------+------------+
  |   vault-node2 (server)  |
  +------------------------+


       (Initialization and Preloading)

           +------------------+
           | vault-bootstrap   |
           +------------------+
                    |
           +------------------+
           | vault-preload     |
           +------------------+

   vault-bootstrap --> vault-node1 & vault-node2
   vault-preload   --> vault-node1 & vault-node2
```

- vault-node1 and vault-node2 share the same Raft cluster for storage.
- vault-bootstrap connects to vault-node1 (or vault-node2) to:
  - Initialize Vault (only if it's not already initialized).
  - Unseal Vault.
- vault-preload loads some example secrets after initialization.

## Running Vault (Development)

Vault runs inside Docker in dev mode for local development.

### Start Vault with all services:

  ```bash
  make up
  ```

- Vault URL: http://localhost:8200
- Default Root Token: root
- Dev mode storage (file-based)

### Check Vault health:

```bash
curl http://localhost:8200/v1/sys/health
```

### Folder structure

super_services/
├── pkg/
│   └── vault/
│       └── vault.go  # Vault client wrapper
We use a custom Vault Go client wrapper for easy secret management inside services.

## Using Vault in Go Services

1. Initialize Vault Client 
    ```go
    import "super_services/pkg/vault"

    client, err := vault.NewVaultClient("http://localhost:8200", "root")
    if err != nil {
      log.Fatalf("failed to connect Vault: %v", err)
    }
    ```
2. Write Secrets
    ```go
    err = client.WriteSecret("gateway/config", map[string]interface{}{
      "api_key": "your-api-key-here",
    })
    if err != nil {
      log.Fatalf("failed to write secret: %v", err)
    }
    ```
3. Read Secrets
    ```go
    secret, err := client.ReadSecret("gateway/config")
    if err != nil {
      log.Fatalf("failed to read secret: %v", err)
    }
    log.Printf("Loaded secret: %v", secret)
    ```
4. List Secrets
    ```go
    paths, err := client.ListSecrets("gateway")
    if err != nil {
      log.Fatalf("failed to list secrets: %v", err)
    }
    log.Printf("Secrets under /gateway/: %v", paths)
    ```
5. Delete Secrets
    ```go
    err = client.DeleteSecret("gateway/config")
    if err != nil {
      log.Fatalf("failed to delete secret: %v", err)
    }
    ```

## Secrets Storage Convention

| Purpose | Path | Example |
| ------- | ---- | ------- |
| Service Configuration | service-name/config | gateway/config |
| Tenant-specific Secrets | tenants/{tenant_id}/service-name/secret-name | tenants/tenant123/billing/stripe_keys |
| Global Secrets | global/service-name/secret-name | global/notification/smtp |

### Example for Multi-Tenant:

- /tenants/tenant123/billing/db_credentials
- /tenants/tenant456/gateway/api_keys

## Vault Best Practices (Production)

- Never commit secrets into Git.
- Use dynamic secrets whenever possible (e.g., DB creds that auto-expire).
- Apply least privilege policies using Vault’s RBAC if moving to production.
- Use Vault Agent and Auto-Auth in Kubernetes for seamless service auth.
- Enable audit logging in production Vault.

## Additional Resources

- Vault Official Documentation
- Go Vault API Reference


## Quick Developer Checklist
- Vault running in Docker (make up)
-  Vault client available at pkg/vault/vault.go
-  Secrets written per service or per tenant
-  No raw HTTP calls — clean Go client usage
-  Close to production practices even in local dev