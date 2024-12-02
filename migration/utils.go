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

// normalizePath get normalized path
func normalizePath(p ...string) string {
	return filepath.Clean(path.Join(p...))
}

// makeDir make nested directory if not exists
func makeDir(dir string) error {
	dir = normalizePath(dir)
	if stat, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModeDir|0755)
	} else if err != nil {
		return err
	} else if !stat.IsDir() {
		return fmt.Errorf("%s is not directory", dir)
	}
	return nil
}

// hardTrim remove all space and control character
func hardTrim(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, "\t", "")
	return str
}

// readStage read valid lines for stage
func readStage(content, section, stage string) (string, error) {
	normalize := func(str string) string {
		return strings.ToLower(hardTrim(str))
	}
	isNewStage := func(str string) bool {
		normalized := normalize(str)
		return strings.HasPrefix(normalized, "--{up:") || strings.HasPrefix(normalized, "--{down:")
	}
	isComment := func(str string) bool {
		return strings.HasPrefix(normalize(str), "--")
	}
	isPreferStage := func(str string) bool {
		normalized := normalize(str)
		return strings.HasPrefix(normalized, normalize("--{"+section+":"+stage))
	}

	founded := false
	lines := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if founded {
			if isNewStage(line) {
				break
			} else if hardTrim(line) != "" && !isComment(line) {
				lines = append(lines, line)
			}
		} else if isPreferStage(line) {
			founded = true
		}
	}

	if err := scanner.Err(); err != nil {
		return "", err
	} else {
		return strings.Join(lines, "\n"), nil
	}
}

// upScripts read scripts for up section of script
func upScripts(content string, stage string) (string, error) {
	if code, err := readStage(content, "up", stage); err != nil {
		return "", err
	} else {
		return code, nil
	}
}

// downScripts read scripts for down section of script
func downScripts(content string, stage string) (string, error) {
	if code, err := readStage(content, "down", stage); err != nil {
		return "", err
	} else {
		return code, nil
	}
}
