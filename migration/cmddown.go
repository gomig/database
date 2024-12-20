package migration

import (
	"slices"

	"github.com/spf13/cobra"
)

func downCmd(driver Migration, autoExec []string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "down [stage1, stage2, ...]"
	cmd.Short = "rollback migrations"
	cmd.Flags().StringP("name", "n", "", "migration name")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := driver.Init(); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		}

		var stages []string
		if len(args) > 0 {
			stages = append([]string{}, args...)
		} else {
			stages = append([]string{}, autoExec...)
			slices.Reverse(stages)
		}

		names := []string{}
		if name := flag(cmd, "name"); name != "" {
			names = append(names, name)
		}

		total := 0
		for _, stage := range stages {
			if files, err := driver.Down(stage, names...); err != nil {
				Formatter("Stage {b}{B}%s{R}:\n", stage)
				Formatter("    {r}FAIL!{R} %s\n", err.Error())
				return
			} else if len(files) > 0 {
				total += len(files)
				Formatter("Stage {b}{B}%s{R}:\n", stage)
				for _, file := range files {
					Formatter("    %s {g}RolledBack!{R}\n", file)
				}
			}
		}

		if total == 0 {
			Formatter("{m}{I}Nothing to rollback!{R}\n")
		}
	}
	return cmd
}
