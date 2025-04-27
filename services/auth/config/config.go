package service

import (
	"gitlab.com/ajesicus/super_services/pkg/config"
)

type AuthConfig struct {
	JWTSecret         string
	OAuthClientID     string
	OAuthClientSecret string
}

func LoadAuthConfig() (*AuthConfig, error) {
	jwtSecret, err := config.GetConfig("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	oauthClientID, err := config.GetConfig("OAUTH_CLIENT_ID")
	if err != nil {
		return nil, err
	}
	oauthClientSecret, err := config.GetConfig("OAUTH_CLIENT_SECRET")
	if err != nil {
		return nil, err
	}

	return &AuthConfig{
		JWTSecret:         jwtSecret,
		OAuthClientID:     oauthClientID,
		OAuthClientSecret: oauthClientSecret,
	}, nil
}
