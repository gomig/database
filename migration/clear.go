package migration

import (
	"fmt"
	"os"
	"path"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func clearCMD(resolver func(driver string) *sqlx.DB) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "clear"
	cmd.Short = "delete all database table"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var err error
		driver, err := cmd.Flags().GetString("driver")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		migrationsDir, err := cmd.Flags().GetString("migration_dir")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		db := resolver(driver)
		if db == nil {
			fmt.Printf("failed: %s database driver not found\n", driver)
			return
		}

		// Read file
		content, err := os.ReadFile(path.Join(migrationsDir, "clean.sql"))
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		commands := getValidLines(string(content))

		// Validate commands
		for _, cmd := range commands {
			if err := validateStatement(cmd, db); err != nil {
				fmt.Printf("\nclear failed\n%s\n", err.Error())
				return
			}
		}

		// Run commands
		for _, cmd := range commands {
			if _, err = db.Exec(cmd); err != nil {
				fmt.Printf("\nclear failed\n%s\n", err.Error())
				return
			}
		}

		fmt.Println("DONE")
	}
	return cmd
}
