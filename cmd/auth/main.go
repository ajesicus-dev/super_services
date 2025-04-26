package main

import (
	"log"

	"super_services/internal/auth/service"
	"super_services/pkg/config"
)

func main() {
	err := config.Init("auth") // service name for vault path
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	cfg, err := service.LoadAuthConfig()
	if err != nil {
		log.Fatalf("failed to load service config: %v", err)
	}

	log.Printf("Auth service started with client_id: %s", cfg.OAuthClientID)

	// Boot your server...
}