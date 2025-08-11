package app

import (
	"github.com/spf13/cobra"
	"go-auth-otp-service/src/bootstrap"
	"log"
)

var bootstrapCmd = &cobra.Command{
	Use:   "bootstrap",
	Short: "Bootstraps the application and it's related services.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := bootstrap.Init(); err != nil {
			log.Fatalf("Bootstrap Service: Failed to Initialize. %v", err)
		}
	},
}
