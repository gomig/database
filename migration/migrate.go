package migration

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func migrateCmd(resolver func(driver string) *sqlx.DB) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "migrate"
	cmd.Short = "migrate database"
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

		// Create migrations table if not exists
		createMigrationTable(db)

		// Get available migrations
		files, paths := getMigrationFiles(getMigratedFiles(db, false), migrationsDir)

		// Run migrations
		for i, mFile := range files {

			fmt.Printf("\nmigrating: %s\n", mFile)

			// Read file
			content, err := os.ReadFile(paths[i])
			if err != nil {
				fmt.Printf("failed: %s\n", err.Error())
				return
			}

			commands := getValidLines(string(content))
			// Validate commands
			for _, cmd := range commands {
				if err := validateStatement(cmd, db); err != nil {
					fmt.Printf("\nmigration failed\n%s\n%s\n", mFile, err.Error())
					return
				}
			}

			// Run Migration
			for _, cmd := range commands {
				if _, err = db.Exec(cmd); err != nil {
					fmt.Printf("\nmigration failed\n%s\n%s\n", mFile, err.Error())
					return
				}
			}

			// Add to migrated table
			_, err = db.Exec(fmt.Sprintf("INSERT INTO migrations VALUES(%s, 0);", mFile))
			if err != nil {
				fmt.Printf("\nmigration failed\n%s\n%s\n", mFile, err.Error())
				return
			}

			fmt.Printf("migrated!\n\n")
		}
	}
	return cmd
}
