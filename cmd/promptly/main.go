package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/rahulguha/promptly/internal/config"
	"github.com/rahulguha/promptly/internal/routes"
	"github.com/rahulguha/promptly/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		fmt.Fprintln(os.Stderr, "No .env file found, reading from environment")
	}

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
	// Initialize configuration
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the DB manager
	dbManager := storage.NewDBManager()

	// The handler now gets the DBManager instead of a specific storage instance
	handler := &routes.Handler{DBManager: dbManager, Cfg: cfg}

	// Setup Gin router
	r := gin.Default()
	routes.RegisterRoutes(r, handler)

	// Start server
	fmt.Printf("Starting Promptly server on port %s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
