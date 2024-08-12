package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

// MigrationCommand get migration command
func MigrationCommand(db *sqlx.DB, root string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "migration"
	cmd.Short = "migrate database"
	cmd.AddCommand(newCMD(root))
	cmd.AddCommand(summeryCmd(db))
	cmd.AddCommand(runCmd(db, root))
	cmd.AddCommand(upCmd(db, root))
	cmd.AddCommand(scriptCmd(db, root))
	cmd.AddCommand(seedCmd(db, root))
	cmd.AddCommand(downCmd(db, root))
	return cmd
}
