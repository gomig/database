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
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		}

		if records, err := getRunnedMigrate(db); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		} else {
			Formatter("{m}{I}Executed Migration{R}: {b}{B}(%d){R}\n", len(records))
			for _, rec := range records {
				fmt.Printf("\t%s", rec)
			}
			fmt.Println()
		}

		if records, err := getRunnedScript(db); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		} else {
			Formatter("{m}{I}Executed Script{R}: {b}{B}(%d){R}\n", len(records))
			for _, rec := range records {
				fmt.Printf("\t%s", rec)
			}
			fmt.Println()
		}

		if records, err := getRunnedSeed(db); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		} else {
			Formatter("{m}{I}Executed Seed{R}: {b}{B}(%d){R}\n", len(records))
			for _, rec := range records {
				fmt.Printf("\t%s", rec)
			}
			fmt.Println()
		}
	}
	return cmd
}
