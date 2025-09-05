package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// Config stores all configuration for the application.
type Config struct {
	CognitoDomain       string
	CognitoClientID     string
	CognitoClientSecret string
	CognitoRedirectURI  string
	SessionSecret       string
	Port                string
	FrontendURL         string
}

// New loads configuration from environment variables and .env file.
func New() (*Config, error) {
	// Load .env file. This is not fatal.
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "No .env file found, reading from environment")
	}

	// Set up Viper to read environment variables
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8082")

	cfg := &Config{
		CognitoDomain:       viper.GetString("COGNITO_DOMAIN"),
		CognitoClientID:     viper.GetString("COGNITO_CLIENT_ID"),
		CognitoClientSecret: viper.GetString("COGNITO_CLIENT_SECRET"),
		CognitoRedirectURI:  viper.GetString("COGNITO_REDIRECT_URI"),
		SessionSecret:       viper.GetString("SESSION_SECRET"),
		Port:                viper.GetString("PORT"),
		FrontendURL:         viper.GetString("FRONTEND_URL"),
	}

	// --- Critical Debugging Step ---
	// Print out the loaded configuration to be 100% sure.
	fmt.Println("--- Loaded Configuration ---")
	fmt.Printf("COGNITO_DOMAIN: %s\n", cfg.CognitoDomain)
	fmt.Printf("COGNITO_CLIENT_ID: %s\n", cfg.CognitoClientID)
	fmt.Printf("COGNITO_REDIRECT_URI: %s\n", cfg.CognitoRedirectURI)
	fmt.Printf("PORT: %s\n", cfg.Port)
	fmt.Println("--------------------------")

	if cfg.CognitoDomain == "" {
		return nil, fmt.Errorf("FATAL: COGNITO_DOMAIN is not set")
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("FATAL: SESSION_SECRET is not set")
	}

	return cfg, nil
}
