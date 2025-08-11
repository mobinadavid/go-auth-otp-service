package database

import (
	"github.com/spf13/cobra"
	"go-auth-otp-service/src/database/migrations"
	"log"
)

func init() {
	migrateCmd.AddCommand(
		migrateUpCmd,
		migrateDownCmd,
	)
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "run up method of migrations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running up methods of migrations")
		if err := migrations.Up(); err != nil {
			log.Fatalln(err)
		}
		log.Println("Database has migrated successfully!")
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "run down method of migrations",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("Running down methods of migrations")
		if err := migrations.Down(); err != nil {
			log.Fatalln(err)
		}
		log.Println("Database has migrated successfully!")
	},
}
