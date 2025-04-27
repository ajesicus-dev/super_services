package config

import (
	"fmt"
	"os"
	"strings"

	"gitlab.com/ajesicus/super_services/pkg/vault"

	"github.com/spf13/viper"
)

type Config struct {
	VaultEnabled bool
	VaultClient  *vault.VaultClient
	ServiceName  string
}

var config *Config

func Init(serviceName string) error {
	vaultEnabled := os.Getenv("VAULT_ENABLED") == "true"

	var vaultClient *vault.VaultClient
	if vaultEnabled {
		client, err := vault.NewVaultClient()
		if err != nil {
			return fmt.Errorf("vault client error: %w", err)
		}
		vaultClient = client
	}

	viper.AutomaticEnv() // Read env variables automatically

	config = &Config{
		VaultEnabled: vaultEnabled,
		VaultClient:  vaultClient,
		ServiceName:  serviceName,
	}

	return nil
}

// GetConfig loads a value with fallback
func GetConfig(key string) (string, error) {
	// 1. Try Vault first
	if config.VaultEnabled && config.VaultClient != nil {
		secretPath := fmt.Sprintf("%s/config", config.ServiceName) // e.g., auth/config
		secrets, err := config.VaultClient.ReadSecret(secretPath)
		if err == nil {
			// Normalize key
			k := strings.ToLower(key)
			if val, ok := secrets[k]; ok {
				return fmt.Sprintf("%v", val), nil
			}
		}
	}

	// 2. Try Viper (env)
	val := viper.GetString(key)
	if val == "" {
		return "", fmt.Errorf("config key not found: %s", key)
	}

	return val, nil
}
