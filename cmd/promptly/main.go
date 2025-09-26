package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rahulguha/promptly/internal/api"
	"github.com/rahulguha/promptly/internal/config"
	"github.com/rahulguha/promptly/internal/llm"
	"github.com/rahulguha/promptly/internal/routes"
	"github.com/rahulguha/promptly/internal/storage"
	"github.com/rahulguha/promptly/internal/tracking"
	"github.com/spf13/cobra"
)

var (
	port        int
	llmProvider string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "promptly",
	Short: "Promptly is a prompt management tool.",
	Long:  `A versatile tool for managing and generating prompts for various AI models.`,
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the Promptly API server.",
	Long:  `Starts the API server to manage personas, templates, and prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.New()
		if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
		}

		// Override config with flags if provided (only if flag was explicitly set)
		if cmd.Flags().Changed("port") {
			cfg.Port = fmt.Sprintf("%d", port)
		}

		// Initialize the database manager
		dbManager := storage.NewDBManager()

		// Initialize the tracker
		tracker, err := tracking.NewDynamoDBTracker(cfg.DynamoDBRegion, cfg.DynamoDBTableName, cfg.DynamoDBActivityTableName)
		if err != nil {
			log.Fatalf("Failed to initialize tracker: %v", err)
		}

		// Initialize the llm evaluator
		var evaluator llm.Evaluator
		switch llmProvider {
		case "openai":
			evaluator, err = llm.NewOpenAIEvaluator(cfg.LLMPrompts.OpenAI)
		case "claude":
			evaluator, err = llm.NewClaudeEvaluator(cfg.LLMPrompts.Claude)
		default:
			log.Fatalf("Invalid llm-provider: %s", llmProvider)
		}
		if err != nil {
			log.Fatalf("Failed to initialize LLM evaluator: %v", err)
		}

		// Create user tracking handler wrapper
		userTrackingHandler := api.NewUserTrackingHandler(tracker)

		// Initialize the handler with the config and tracker
		handler := routes.NewHandler(cfg, dbManager, userTrackingHandler, evaluator)

		// Initialize the router
		router := routes.NewRouter(handler)

		// Start the server
		address := fmt.Sprintf(":%s", cfg.Port)
		log.Printf("Server starting on %s", address)
		if err := router.Run(address); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	},
}

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	rootCmd.AddCommand(serveCmd)
	
	// Define flags for the serve command
	serveCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to run the server on")
	serveCmd.Flags().StringVar(&llmProvider, "llm-provider", "openai", "LLM provider to use (openai or claude)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}