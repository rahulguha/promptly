package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rahulguha/promptly/internal/routes"
	"github.com/rahulguha/promptly/internal/storage/jsonstore"
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

		// Search config in home directory with name ".promptly" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".promptly")
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
	dataPath := viper.GetString("data")

	// Initialize storage and handler
	store, err := jsonstore.NewFileStorage(dataPath)
	if err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}
	handler := &routes.Handler{Store: store}

	// Setup Gin router
	r := gin.Default()
	routes.RegisterRoutes(r, handler)

	// Start server
	fmt.Printf("Starting Promptly server on port %s\n", port)
	fmt.Printf("Using data file: %s\n", dataPath)

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