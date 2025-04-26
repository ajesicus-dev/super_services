package vault

import (
	"fmt"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
)

type VaultClient struct {
	client *vaultapi.Client
}

func NewVaultClient() (*VaultClient, error) {
	addr := os.Getenv("VAULT_ADDR")
	if addr == "" {
		return nil, fmt.Errorf("VAULT_ADDR not set")
	}
	token := os.Getenv("VAULT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("VAULT_TOKEN not set")
	}

	config := vaultapi.DefaultConfig()
	config.Address = addr

	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	client.SetToken(token)

	return &VaultClient{client: client}, nil
}

// ReadSecret fetches a secret at a given Vault path.
func (v *VaultClient) ReadSecret(path string) (map[string]interface{}, error) {
	secret, err := v.client.Logical().Read(fmt.Sprintf("secret/data/%s", path))
	if err != nil {
		return nil, fmt.Errorf("failed to read secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("no data found at path: %s", path)
	}

	// Vault v2 KV secrets: data is inside a "data" field
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format at path: %s", path)
	}

	return data, nil
}
