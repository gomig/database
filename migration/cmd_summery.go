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

		if err := InitMigration(db); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		}

		if summery, err := Migrated(db); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		} else {
			fmt.Println("Migration Summery:")

			for stage, files := range summery.GroupByStage() {
				Formatter(
					"Migrate {b}{B}[%s]{R}: {B}(%d){R}\n",
					stage, len(files),
				)

				for _, file := range files {
					fmt.Printf("    %s\n", file)
				}

				fmt.Println()
			}
		}
	}
	return cmd
}
