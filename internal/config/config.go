package config

import (
	"encoding/json"
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
	LLMPrompts          *LLMPrompts
}

// LLMPrompts holds the system prompts for different LLM providers
type LLMPrompts struct {
	Claude string `json:"claude"`
	OpenAI string `json:"openai"`
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

	// Load LLM prompts from llm_config.json
	llmPrompts, err := loadLLMPrompts()
	if err != nil {
		return nil, fmt.Errorf("FATAL: Failed to load LLM prompts: %w", err)
	}
	cfg.LLMPrompts = llmPrompts

	return cfg, nil
}

// loadLLMPrompts loads the LLM prompts from llm_config.json
func loadLLMPrompts() (*LLMPrompts, error) {
	data, err := os.ReadFile("llm_config.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read llm_config.json: %w", err)
	}

	var prompts LLMPrompts
	err = json.Unmarshal(data, &prompts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal llm_config.json: %w", err)
	}

	return &prompts, nil
}

