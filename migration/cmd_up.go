package migration

import (
	"errors"
	"fmt"

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
			throw(err)
		} else if len(files) == 0 {
			throw(errors.New("no migration found"))
		} else {
			for _, file := range files {
				fmt.Printf("MIGRATE %s: ", file.Name)
				if res, err := Migrate(db, flag(cmd, "name"), file); err != nil {
					fmt.Printf("FAIL! %s\n", err.Error())
				} else if len(res) > 0 {
					fmt.Printf("OK!\n")
				}
			}
		}
	}
	return cmd
}
