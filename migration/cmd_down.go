package migration

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func downCmd(db *sqlx.DB, root string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "down"
	cmd.Short = "rollback migrations"
	cmd.Flags().StringP("name", "n", "", "migration name")
	cmd.Flags().StringP("dir", "d", "", "directory path")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if files, err := ReadDirectory(uri(root, flag(cmd, "dir"))); err != nil {
			throw(err)
		} else if len(files) == 0 {
			throw(errors.New("no migration found"))
		} else {
			for _, file := range files {
				if res, err := Rollback(db, flag(cmd, "name"), file); err != nil {
					fmt.Printf("%s [ROLLBACK] failed!\n\t%s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					fmt.Printf("%s [ROLLBACK] ok!\n", file.Name)
				}
			}
		}
	}
	return cmd
}
