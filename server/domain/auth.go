package domain

import (
	"log"
	"os"

	"golang.org/x/oauth2"
)

type GoogleAuthConfig struct {
	ClientID          string
	ClientSecret      string
	AuthorizeEndpoint string
	TokenEndpoint     string
}

var authConfig GoogleAuthConfig

func InitAuthConfig(logger *log.Logger) {
	authConfig.ClientID = os.Getenv("GOOGLE_AUTH_CLIENT_ID")
	if authConfig.ClientID == "" {
		logger.Panic("Failed to import env: GOOGLE_AUTH_CLIENT_ID")
	}
	authConfig.ClientSecret = os.Getenv("GOOGLE_AUTH_CLIENT_SECRET")
	if authConfig.ClientSecret == "" {
		logger.Panic("Failed to import env: GOOGLE_AUTH_CLIENT_SECRET")
	}
	authConfig.AuthorizeEndpoint = "https://accounts.google.com/o/oauth2/v2/auth"
	authConfig.TokenEndpoint = "https://www.googleapis.com/oauth2/v4/token"
}

func GetConfig() *oauth2.Config {
	config := &oauth2.Config{
		ClientID:     authConfig.ClientID,
		ClientSecret: authConfig.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authConfig.AuthorizeEndpoint,
			TokenURL: authConfig.TokenEndpoint,
		},
		Scopes:      []string{"openid", "email", "profile"},
		RedirectURL: "http://localhost:8080/callback",
	}

	return config
}
