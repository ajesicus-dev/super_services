package vault

import (
	"context"
	"fmt"
	"github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *api.Client
}

// NewVaultClient creates a new VaultClient instance.
func NewVaultClient(address, token string) (*VaultClient, error) {
	cfg := api.DefaultConfig()
	cfg.Address = address

	client, err := api.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	client.SetToken(token)

	return &VaultClient{client: client}, nil
}

// WriteSecret writes key-value data to a given path.
func (vc *VaultClient) WriteSecret(path string, data map[string]interface{}) error {
	_, err := vc.client.KVv2("secret").Put(context.Background(), path, data)
	if err != nil {
		return fmt.Errorf("failed to write secret: %w", err)
	}
	return nil
}

// ReadSecret reads key-value data from a given path.
func (vc *VaultClient) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := vc.client.KVv2("secret").Get(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret at path %s not found", path)
	}
	return secret.Data, nil
}

// DeleteSecret deletes a secret at a given path.
func (vc *VaultClient) DeleteSecret(path string) error {
	err := vc.client.KVv2("secret").Delete(context.Background(), path)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}
	return nil
}

// ListSecrets lists all secrets under a given path.
func (vc *VaultClient) ListSecrets(path string) ([]string, error) {
	list, err := vc.client.KVv2("secret").List(context.Background(), path)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}
	if list == nil || list.Data == nil {
		return nil, fmt.Errorf("no secrets found at path %s", path)
	}

	keysRaw, ok := list.Data["keys"]
	if !ok {
		return nil, fmt.Errorf("no keys found in secret data")
	}

	keys, ok := keysRaw.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format for keys")
	}

	var result []string
	for _, key := range keys {
		if strKey, ok := key.(string); ok {
			result = append(result, strKey)
		}
	}

	return result, nil
}
