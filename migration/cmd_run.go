package migration

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func runCmd(db *sqlx.DB, root string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "run"
	cmd.Short = "run [up], [script] and [seed] migrations"
	cmd.Flags().StringP("name", "n", "", "migration name")
	cmd.Flags().StringP("dir", "d", "", "directory path")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if files, err := ReadDirectory(uri(root, flag(cmd, "dir"))); err != nil {
			throw(err)
		} else if files.Len() == 0 {
			throw(errors.New("no migration found"))
		} else {
			name := flag(cmd, "name")
			for _, file := range files {
				if res, err := Migrate(db, name, file); err != nil {
					fmt.Printf("%s migrate failed: %s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s migrate done\n", file.Name)
				}

				if res, err := Script(db, name, file); err != nil {
					fmt.Printf("%s script failed: %s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s script done\n", file.Name)
				}

				if res, err := Seed(db, name, file); err != nil {
					fmt.Printf("%s seed failed: %s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s seed done\n", file.Name)
				}
			}
		}
	}
	return cmd
}
