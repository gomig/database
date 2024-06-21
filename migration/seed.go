package migration

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func seedCmd(resolver func(driver string) *sqlx.DB) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "seed"
	cmd.Short = "seed database"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		var err error

		driver, err := cmd.Flags().GetString("driver")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		seedsDir, err := cmd.Flags().GetString("seed_dir")
		if err != nil {
			fmt.Printf("failed: %s\n", err.Error())
			return
		}

		db := resolver(driver)
		if db == nil {
			fmt.Printf("failed: %s database driver not found\n", driver)
			os.Exit(1)
			return
		}

		// Create migrations table if not exists
		createMigrationTable(db)

		// Get available migrations
		files, paths := getMigrationFiles(getMigratedFiles(db, true), seedsDir)

		// Run migrations
		for i, mFile := range files {
			fmt.Printf("\nseeding: %s\n", mFile)

			// Read file
			content, err := ioutil.ReadFile(paths[i])
			if err != nil {
				fmt.Printf("failed: %s\n", err.Error())
				return
			}

			commands := getValidLines(string(content))
			// Validate commands
			for _, cmd := range commands {
				if err := validateStatement(cmd, db); err != nil {
					fmt.Printf("\nseeding failed\n%s\n%s\n", mFile, err.Error())
					os.Exit(1)
					return
				}
			}

			// Run Seeds
			for _, cmd := range commands {
				if _, err = db.Exec(cmd); err != nil {
					fmt.Printf("\nseeding failed\n%s\n%s\n", mFile, err.Error())
					os.Exit(1)
					return
				}
			}

			// Add to seed table
			_, err = db.Exec(fmt.Sprintf("INSERT INTO migrations VALUES(%s, 1);", mFile))
			if err != nil {
				fmt.Printf("\nseeding failed\n%s\n%s\n", mFile, err.Error())
				os.Exit(1)
				return
			}

			fmt.Printf("seeded\n\n")
		}
	}
	return cmd
}
