package migration

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func summaryCmd(driver Migration) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "summary"
	cmd.Short = "show migration summary"
	cmd.Run = func(cmd *cobra.Command, args []string) {
		if err := driver.Init(); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		}

		summary, err := driver.Summary()
		if err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		}

		if len(summary) == 0 {
			Formatter("{m}{I}Nothing migrated!{R}\n")
			return
		}

		fmt.Println("Migration Summery:")
		fmt.Println("")
		for stage, files := range summary.GroupByStage() {
			Formatter(
				"Stage {b}{B}%s{R}: {B}(%d){R}\n",
				stage, len(files),
			)

			for _, file := range files {
				f := File{name: file, content: "", ext: driver.Extension()}
				fmt.Printf("    %s\n", strings.ReplaceAll(f.HumanizedName(), "-", " "))
			}

			fmt.Println()
		}
	}
	return cmd
}
