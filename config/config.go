package config

import (
	"os"

	"github.com/Nerzal/gocloak/v13"
)

type Config struct {
	KeycloakURL   string
	Realm         string
	ClientID      string
	ClientSecret  string
	GocloakClient *gocloak.GoCloak
}
// this is a function that set config
func NewConfig() *Config {
	keycloakURL := getEnv("KEYCLOAK_URL", "http://localhost:8080")
	return &Config{
		KeycloakURL:   keycloakURL,
		Realm:         getEnv("REALM", "mani"),
		ClientID:      getEnv("CLIENT_ID", "go"),
		ClientSecret:  getEnv("CLIENT_SECRET", "aZbUGUvcXxETMznmtsFB1Mm8TEWK4Hpz"),
		GocloakClient: gocloak.NewClient(keycloakURL),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
