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
		if migrations, err := ReadDirectory(uri(root, flag(cmd, "dir"))); err != nil {
			throw(err)
		} else if len(migrations) == 0 {
			throw(errors.New("no migration found"))
		} else if done, err := Migrate(db, migrations, flag(cmd, "name")); err != nil {
			throw(err)
		} else {
			for _, d := range done {
				fmt.Printf("%s done\n", d)
			}
		}
	}
	return cmd
}
