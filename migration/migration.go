package migration

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

// MigrationCommand get migration command
func MigrationCommand(resolver func(driver string) *sqlx.DB, defDriver string, migDir string, seedDir string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "migration"
	cmd.Short = "migrate database"
	cmd.AddCommand(clearCMD(resolver))
	cmd.AddCommand(migrateCmd(resolver))
	cmd.AddCommand(migratedCmd(resolver))
	cmd.AddCommand(seedCmd(resolver))
	cmd.AddCommand(seededCmd(resolver))
	cmd.PersistentFlags().StringP("driver", "d", defDriver, "database driver name")
	cmd.PersistentFlags().StringP("migration_dir", "m", migDir, "migrations path")
	cmd.PersistentFlags().StringP("seed_dir", "s", seedDir, "seeds path")
	return cmd
}
