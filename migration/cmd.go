package migration

import (
	"github.com/spf13/cobra"
)

// MigrationCommand migration cli commands
func MigrationCommand(driver Migration, autExec ...string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "migration"
	cmd.Short = "migrate database"
	cmd.AddCommand(newCMD(driver, autExec))
	cmd.AddCommand(summaryCmd(driver))
	cmd.AddCommand(upCmd(driver, autExec))
	cmd.AddCommand(downCmd(driver, autExec))
	return cmd
}
