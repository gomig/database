package migration

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// flag get string flag or return empty string
func flag(cmd *cobra.Command, name string) string {
	if v, err := cmd.Flags().GetString(name); err == nil {
		return v
	}
	return ""
}

// uri get normalized path
func uri(p ...string) string {
	return filepath.Clean(path.Join(p...))
}

// makeDir make nested directory if not exists
func makeDir(dir string) error {
	dir = uri(dir)
	if stat, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModeDir|0755)
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return fmt.Errorf("%s is not directory", dir)
	}
	return nil
}

// readLines read valid lines for stage
func readLines(content string, stage string) ([]string, error) {
	trim := func(str string) string {
		return strings.ToUpper(strings.ReplaceAll(str, " ", ""))
	}
	normalize := func(str string) string {
		return strings.NewReplacer("-- [end]", "--[END]", "--[end]", "--[END]").Replace(str)
	}
	isStage := func(str, stage string) bool {
		if stage == "" {
			return strings.HasPrefix(trim(str), "--[STAGE")
		} else {
			return strings.HasPrefix(trim(str), trim("--[STAGE"+stage+"]"))
		}
	}

	lines := make([]string, 0)
	founded := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if founded {
			if isStage(line, "") {
				break
			} else if trim(line) != "" {
				lines = append(lines, normalize(line))
			}
		} else if isStage(line, stage) {
			founded = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	} else {
		return lines, nil
	}
}

// scriptsOf read scripts splitted bt -- [end] for stage
func scriptsOf(content string, stage string) ([]string, error) {
	if lines, err := readLines(content, stage); err != nil {
		return nil, err
	} else if len(lines) == 0 {
		return []string{}, nil
	} else {
		return strings.Split(strings.Join(lines, "\r\n"), "--[END]"), nil
	}
}
