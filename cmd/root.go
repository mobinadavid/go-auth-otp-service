package cmd

import (
	"github.com/spf13/cobra"
	"go-auth-otp-service/cmd/app"
	"go-auth-otp-service/cmd/database"
	"go-auth-otp-service/src/cache"
	"go-auth-otp-service/src/config"
	"log"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fund-app",
	Short: "Fund app backend API project.",
	Long:  ``,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(
		initConfig,
		initCache,
	)

	rootCmd.AddCommand(
		app.AppCmd,
		database.DatabaseCmd,
	)
}

func initConfig() {
	// Initialize configuration
	err := config.Init()
	if err != nil {
		log.Fatalf("Config Service: Failed to Initialize. %v", err)
	}
	log.Println("Config Service: Initialized Successfully.")
}

func initCache() {
	err := cache.Init()
	if err != nil {
		log.Fatalf("Cache Service: Failed to Initialize. %v", err)
	}
}
