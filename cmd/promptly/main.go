package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rahulguha/promptly/internal/api"
	"github.com/rahulguha/promptly/internal/storage"
	"github.com/rahulguha/promptly/internal/storage/jsonstore"
)

var cfgFile string

func main() {
	var rootCmd = &cobra.Command{
		Use:   "promptly",
		Short: "Promptly - a prompt management API",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is config.yaml)")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config: %v", err)
	}
}

func startServer() {
	var store storage.StorageProvider

	switch viper.GetString("storage.provider") {
	case "json":
		path := viper.GetString("storage.path")
		s, err := jsonstore.New(path)
		if err != nil {
			log.Fatalf("failed to init storage: %v", err)
		}
		store = s
	default:
		log.Fatal("unknown storage provider")
	}

	port := viper.GetInt("server.port")
	router := api.NewRouter(store)
	log.Printf("Promptly API running on :%d", port)
	log.Fatal(router.Listen(port))
}
