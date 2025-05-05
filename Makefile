# General Variables
VAULT_ADDR ?= http://127.0.0.1:8200
VAULT_ROOT_TOKEN ?= (your-root-token-here after init)

# Directories and Compose Files
COMPOSE_DIR=deployments/docker-compose
COMPOSE_FILES=-f $(COMPOSE_DIR)/common.yml -f $(COMPOSE_DIR)/vault.yml  -f $(COMPOSE_DIR)/stack.yml -f $(COMPOSE_DIR)/services.yml
VAULT_COMPOSE_FILES=-f deployments/vault/docker-compose.yml -f deployments/vault/docker-compose.dev.yml

.PHONY: build build-cli run-gateway run-auth test lint up down restart logs clean \
        vault-up vault-log vault-down vault-status vault-init vault-unseal vault-login vault-put vault-get vault-root-token vault-secrets

# Build all Go binaries
build:
	go build ./...

# build specificity: cli
build-cli:
	go build -o ./bin/cli ./cmd/cli
	chmod +x ./bin/cli

# Run specific services
run-gateway:
	go run ./cmd/gateway

run-auth:
	go run ./cmd/auth

# Run tests for all packages
test:
	go test ./... -v

# Run linter using golangci-lint
lint:
	golangci-lint run

# Start all services using Docker Compose
up:
	docker compose $(COMPOSE_FILES) up -d

# Stop all services
down:
	docker compose $(COMPOSE_FILES)

# Restart all services
restart: down up

# View logs for all services
logs:
	docker compose $(COMPOSE_FILES) logs -f

# Clean up Docker resources
clean:
	docker compose $(COMPOSE_FILES) down --volumes --remove-orphans

# Vault-specific commands

vault-up:
	docker compose $(VAULT_COMPOSE_FILES) up -d

vault-log:
	docker compose $(VAULT_COMPOSE_FILES) logs -f

vault-down:
	docker compose $(VAULT_COMPOSE_FILES) down -v

vault-status:
	docker exec -it vault-node1 vault operator raft list-peers
	docker exec -it vault-node2 vault status

vault-init:
	docker compose run --rm vault-init

vault-unseal:
	docker compose run --rm vault-unseal

vault-preload:
	docker compose run --rm vault-preload

vault-root-token:
	cat ./deployments/vault/root.token

vault-secrets:
	curl --header "X-Vault-Token: $$(cat ./vault/init-outputs/root.token)" $(VAULT_ADDR)/v1/secret/data/database

vault-login:
	vault login $(VAULT_TOKEN)

vault-put:
	curl --header "X-Vault-Token: $(VAULT_TOKEN)" --request POST --data '{"data": {"key":"value"}}' $(VAULT_ADDR)/v1/secret/data/super_services/config

vault-get:
	curl --header "X-Vault-Token: $(VAULT_TOKEN)" $(VAULT_ADDR)/v1/secret/data/super_services/config

