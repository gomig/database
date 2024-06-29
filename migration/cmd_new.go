package migration

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gomig/utils"
	"github.com/spf13/cobra"
)

func newCMD(root string) *cobra.Command {
	var cmd = new(cobra.Command)
	cmd.Use = "new [name]"
	cmd.Short = "create new migration file"
	cmd.Args = cobra.MinimumNArgs(1)
	cmd.Flags().StringP("dir", "d", "", "directory path")
	cmd.Run = func(cmd *cobra.Command, args []string) {
		// Make directory
		base := uri(root, flag(cmd, "dir"))
		if err := deepMK(base); err != nil {
			throw(err)
		}

		// Generate name
		var name string
		if slug := utils.Slugify(args[0]); slug == "" {
			throw(errors.New("invalid migration file name"))
		} else {
			name = fmt.Sprintf(
				"%s-%s.sql",
				strconv.FormatInt(time.Now().Unix(), 10),
				slug,
			)
		}

		// Write template
		content := strings.Join([]string{
			"-- [SECTION UP] migrate",
			"\r\n",
			"\r\n",
			"-- [SECTION SCRIPT] extra script, triggers, etc.",
			"\r\n",
			"\r\n",
			"-- [SECTION SEED] seed",
			"\r\n",
			"\r\n",
			"-- [SECTION DOWN] rollback",
			"\r\n",
		}, "")
		if err := os.WriteFile(
			uri(base, name),
			[]byte(content),
			0644,
		); err != nil {
			throw(err)
		} else {
			fmt.Printf("%s migration file created\n", uri(base, name))
		}
	}
	return cmd
}
