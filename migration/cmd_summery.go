package migration

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func summeryCmd(db *sqlx.DB) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "summery"
	cmd.Short = "show migration summery"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		fmt.Println("Migration Summery:")

		if err := createMT(db); err != nil {
			throw(err)
		}

		if records, err := getRunnedMigrate(db); err != nil {
			throw(err)
		} else {
			fmt.Printf("Executed Migration (%d):\n", len(records))
			for _, rec := range records {
				fmt.Println(rec)
			}
			fmt.Println()
		}

		if records, err := getRunnedScript(db); err != nil {
			throw(err)
		} else {
			fmt.Printf("Executed Script (%d):\n", len(records))
			for _, rec := range records {
				fmt.Println(rec)
			}
			fmt.Println()
		}

		if records, err := getRunnedSeed(db); err != nil {
			throw(err)
		} else {
			fmt.Printf("Executed Seed (%d):\n", len(records))
			for _, rec := range records {
				fmt.Println(rec)
			}
			fmt.Println()
		}
	}
	return cmd
}
