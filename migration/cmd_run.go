package migration

import (
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func runCmd(db *sqlx.DB, root, ext string, defaults []string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "run [stage1, stage2, ...]"
	cmd.Short = "run stages script"
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
			if len(args) == 0 {
				args = defaults
			}
			if name := flag(cmd, "name"); name != "" {
				files = files.Filter(name)
			}

			// execute
			for _, stage := range args {
				if migrated, err := StageMigrated(db, stage); err != nil {
					Formatter("{r}FAIL!{R} %s\n", err.Error())
					return
				} else {
					for _, file := range files {
						if ok, err := file.Migrate(db, stage, migrated.Names()...); err != nil {
							Formatter(
								"%s {b}{B}[%s]{R} {r}FAIL!{R} %s\n",
								file.Name(), strings.ToUpper(stage), err.Error(),
							)
						} else if ok {
							Formatter(
								"%s {b}{B}[%s]{R} {g}OK!{R}\n",
								file.Name(), strings.ToUpper(stage),
							)
						}
					}
				}
			}
		}
	}
	return cmd
}
