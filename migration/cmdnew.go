package migration

import (
	"github.com/spf13/cobra"
)

func newCMD(driver Migration, autoExec []string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "new [name]"
	cmd.Short = "create new migration file"
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Run = func(cmd *cobra.Command, args []string) {
		name := args[0]
		base := normalizePath(driver.Root())

		if err := NewMigrationFile(base, name, driver.Extension(), autoExec...); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
		} else {
			Formatter("{m}{I}%s{R}: {g}CREATED!{R}\n", name)
		}
	}
	return cmd
}
