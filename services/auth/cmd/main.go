package main

import (
	"log"
	"os"

	"gitlab.com/ajesicus/super_services/pkg/config"
	"gitlab.com/ajesicus/super_services/services/auth/config"
	"gitlab.com/ajesicus/super_services/services/auth/internal/api"
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
	port := getPort()
	r := api.NewServer()

	log.Printf("Starting auth service on port %s...\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}
