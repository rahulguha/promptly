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
	DynamoDBRegion      string
	DynamoDBTableName   string
	DynamoDBActivityTableName string
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
		DynamoDBRegion:      viper.GetString("DYNAMODB_REGION"),
		DynamoDBTableName:   viper.GetString("DYNAMODB_TABLE_NAME"),
		DynamoDBActivityTableName: viper.GetString("DYNAMODB_ACTIVITY_TABLE_NAME"),
	}

	// --- Critical Debugging Step ---
	// Print out the loaded configuration to be 100% sure.
	fmt.Println("--- Loaded Configuration ---")
	fmt.Printf("COGNITO_DOMAIN: %s\n", cfg.CognitoDomain)
	fmt.Printf("COGNITO_CLIENT_ID: %s\n", cfg.CognitoClientID)
	fmt.Printf("COGNITO_REDIRECT_URI: %s\n", cfg.CognitoRedirectURI)
	fmt.Printf("PORT: %s\n", cfg.Port)
	fmt.Printf("DYNAMODB_REGION: %s\n", cfg.DynamoDBRegion)
	fmt.Printf("DYNAMODB_TABLE_NAME: %s\n", cfg.DynamoDBTableName)
	fmt.Printf("DYNAMODB_ACTIVITY_TABLE_NAME: %s\n", cfg.DynamoDBActivityTableName)
	fmt.Println("--------------------------")

	if cfg.CognitoDomain == "" {
		return nil, fmt.Errorf("FATAL: COGNITO_DOMAIN is not set")
	}
	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("FATAL: SESSION_SECRET is not set")
	}
	if cfg.DynamoDBTableName == "" {
		return nil, fmt.Errorf("FATAL: DYNAMODB_TABLE_NAME is not set")
	}
	if cfg.DynamoDBActivityTableName == "" {
		return nil, fmt.Errorf("FATAL: DYNAMODB_ACTIVITY_TABLE_NAME is not set")
	}

	return cfg, nil
}

