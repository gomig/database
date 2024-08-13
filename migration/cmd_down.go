package migration

import (
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
			Formatter("{r}FAIL!{R} %s\n", err.Error())
		} else if files.Len() == 0 {
			Formatter("{m}{I}no migration found{R}\n")
		} else {
			for _, file := range files {
				if res, err := Rollback(db, flag(cmd, "name"), file); err != nil {
					Formatter("{m}{I}%s{R} rollback: {r}FAIL!{R} %s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					Formatter("{m}{I}%s{R} rollback: {g}OK!{R}\n", file.Name)
				}
			}
		}
	}
	return cmd
}
