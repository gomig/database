package migration

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomig/utils"
	"github.com/spf13/cobra"
)

func newCMD(root, ext string, autoExec []string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "new [name]"
	cmd.Short = "create new migration file"
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Flags().StringP("dir", "d", "", "directory path")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// make directory
		base := uri(root, flag(cmd, "dir"))
		if err := makeDir(base); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
			return
		}

		// generate name
		var name string
		if slug := utils.Slugify(args[0]); slug == "" {
			Formatter("{r}FAIL!{R} invalid migration file name\n")
			return
		} else {
			name = fmt.Sprintf(
				"%s-%s."+ext,
				strconv.FormatInt(time.Now().Unix(), 10),
				slug,
			)
		}

		if len(autoExec) == 0 {
			autoExec = []string{"UP", "DOWN"}
		}

		// generate template
		content := make([]string, 0)
		for _, stage := range autoExec {
			stage = strings.ToUpper(stage)
			if stage != "DOWN" {
				content = append(content, fmt.Sprintf("-- [STAGE %s]\r\n\r\n", stage))
			}
		}
		content = append(content, "-- [STAGE DOWN] rollback\r\n\r\n")

		// write file
		if err := os.WriteFile(
			uri(base, name),
			[]byte(strings.Join(content, "")),
			0644,
		); err != nil {
			Formatter("{r}FAIL!{R} %s\n", err.Error())
		} else {
			Formatter("{m}{I}%s{R}: {g}CREATED!{R}\n", uri(base, name))
		}
	}
	return cmd
}
