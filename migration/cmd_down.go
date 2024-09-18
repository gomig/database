package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func downCmd(db *sqlx.DB, root, ext string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "down"
	cmd.Short = "rollback migrations"
	cmd.Flags().StringP("name", "n", "", "migration name")
	cmd.Flags().StringP("dir", "d", "", "directory path")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := InitMigration(db); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
		} else if files, err := ReadFS(uri(root, flag(cmd, "dir")), ext); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
		} else if files.Len() == 0 {
			Formatter("{m}{I}no migration file found{R}\n")
		} else {
			// prepare
			if name := flag(cmd, "name"); name != "" {
				files = files.Filter(name)
			}

			// execute
			if migrated, err := Migrated(db); err != nil {
				Formatter("{r}FAIL!{R} %s\n", err.Error())
				return
			} else {
				for _, file := range files {
					if ok, err := file.Rollback(db, migrated.Names()...); err != nil {
						Formatter(
							"%s {b}{B}[ROLLBACK]{R} {r}FAIL!{R} %s\n",
							file.Name(), err.Error(),
						)
					} else if ok {
						Formatter(
							"%s {b}{B}[ROLLBACK]{R} {g}OK!{R}\n",
							file.Name(),
						)
					}
				}
			}
		}
	}
	return cmd
}
