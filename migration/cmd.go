package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

// MigrationCommand get migration command
func MigrationCommand(db *sqlx.DB, root, ext string, autExec ...string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "migration"
	cmd.Short = "migrate database"
	cmd.AddCommand(newCMD(root, ext, autExec))
	cmd.AddCommand(summeryCmd(db))
	cmd.AddCommand(runCmd(db, root, ext, autExec))
	cmd.AddCommand(downCmd(db, root, ext))
	return cmd
}
