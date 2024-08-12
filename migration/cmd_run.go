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
					fmt.Printf("%s [MIGRATE] failed!\n\t%s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s [MIGRATE] ok!\n", file.Name)
				}

				if res, err := Script(db, name, file); err != nil {
					fmt.Printf("%s [SCRIPT] failed!\n\t%s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s [SCRIPT] ok!\n", file.Name)
				}

				if res, err := Seed(db, name, file); err != nil {
					fmt.Printf("%s [SEED] failed!\n\t%s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s [SEED] ok!\n", file.Name)
				}
			}
		}
	}
	return cmd
}
