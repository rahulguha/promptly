package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rahulguha/promptly/internal/routes"
	"github.com/rahulguha/promptly/internal/storage"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "promptly",
	Short: "A prompt management and API server",
	Long:  `Promptly is a CLI application for managing prompts with a built-in API server.`,
}

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the HTTP API server for managing prompts.`,
	Run: func(cmd *cobra.Command, args []string) {
		startServer()
	},
}

func init() {
	cobra.OnInitialize(initConfig)

	// Add serve command to root
	rootCmd.AddCommand(serveCmd)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

	// Serve command flags
	serveCmd.Flags().StringP("port", "p", "8080", "Port to run the server on")
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))

	serveCmd.Flags().StringP("data", "d", "data/prompts.json", "Path to data file")
	viper.BindPFlag("data", serveCmd.Flags().Lookup("data"))

	serveCmd.Flags().StringP("storage", "s", "json", "Storage backend type (json or sqlite)")
	viper.BindPFlag("storage", serveCmd.Flags().Lookup("storage"))

	serveCmd.Flags().String("db", "data/promptly.db", "Path to SQLite database file (used when --storage=sqlite)")
	viper.BindPFlag("db", serveCmd.Flags().Lookup("db"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory and current directory
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".promptly")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func startServer() {
	// Get configuration from flags/config
	port := viper.GetString("port")
	storageTypeStr := viper.GetString("storage")
	dataPath := viper.GetString("data")
	dbPath := viper.GetString("db")

	// Validate and get storage type
	storageType, err := storage.ValidateStorageType(storageTypeStr)
	if err != nil {
		log.Fatalf("Invalid storage type: %v", err)
	}

	// Create storage configuration
	var config storage.StorageConfig
	switch storageType {
	case storage.StorageTypeJSON:
		config = storage.StorageConfig{
			Type:     storage.StorageTypeJSON,
			JSONPath: dataPath,
		}
	case storage.StorageTypeSQLite:
		config = storage.StorageConfig{
			Type:   storage.StorageTypeSQLite,
			DBPath: dbPath,
		}
	}

	// Initialize storage
	store, err := storage.NewStorage(config)
	if err != nil {
		log.Fatalf("Failed to initialize %s storage: %v", storageType, err)
	}
	defer store.Close()

	handler := &routes.Handler{Store: store}

	// Setup Gin router
	r := gin.Default()
	routes.RegisterRoutes(r, handler)

	// Start server
	fmt.Printf("Starting Promptly server on port %s\n", port)
	fmt.Printf("Using %s storage\n", storageType)
	if storageType == storage.StorageTypeJSON {
		fmt.Printf("Data file: %s\n", dataPath)
	} else {
		fmt.Printf("Database: %s\n", dbPath)
	}

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}