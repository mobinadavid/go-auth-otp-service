package database

import (
	"github.com/spf13/cobra"
	"log"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Manage database seeders",
}

var seedRunCmd = &cobra.Command{
	Use:   "run",
	Short: "run all seeders",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running seeders")
		//seeders.SeedAuthorization()

		log.Println("Database has seeded successfully!")
	},
}

func init() {
	seedCmd.AddCommand(
		seedRunCmd,
	)
}
