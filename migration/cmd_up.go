package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func upCmd(db *sqlx.DB, root string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "up"
	cmd.Short = "run migrations"
	cmd.Flags().StringP("name", "n", "", "migration name")
	cmd.Flags().StringP("dir", "d", "", "directory path")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if files, err := ReadDirectory(uri(root, flag(cmd, "dir"))); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
		} else if files.Len() == 0 {
			Formatter("{m}{I}no migration found{R}\n")
		} else {
			for _, file := range files {
				if res, err := Migrate(db, flag(cmd, "name"), file); err != nil {
					Formatter("{m}{I}%s{R} migrate: {r}FAIL!{R} %s\n", file.Name, err.Error())
				} else if len(res) > 0 {
					Formatter("{m}{I}%s{R} migrate: {g}OK!{R}\n", file.Name)
				}
			}
		}
	}
	return cmd
}
