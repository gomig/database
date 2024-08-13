package migration

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/jmoiron/sqlx"
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
	return regexp.
		MustCompile(`\/+`).
		ReplaceAllString(filepath.ToSlash(path.Join(p...)), "/")
}

// deepMK make nested directory if not exists
func deepMK(dir string) error {
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

// readLines read valid lines for section
func readLines(content string, section string) ([]string, error) {
	trim := func(str string) string {
		return strings.ToUpper(strings.ReplaceAll(str, " ", ""))
	}
	normalize := func(str string) string {
		str = strings.ReplaceAll(str, "-- [br]", "--[br]")
		return strings.ReplaceAll(str, "--[br]", "--[BR]")
	}
	isSection := func(str, section string) bool {
		if section == "" {
			return strings.HasPrefix(trim(str), "--[SECTION")
		} else {
			return strings.HasPrefix(trim(str), trim("--[SECTION"+section+"]"))
		}
	}

	lines := make([]string, 0)
	founded := false
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if founded {
			if isSection(line, "") {
				break
			} else if trim(line) != "" {
				lines = append(lines, normalize(line))
			}
		} else if isSection(line, section) {
			founded = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	} else {
		return lines, nil
	}
}

// readScripts read scripts splitted bt -- [br] for sections
func readScripts(content string, section string) ([]string, error) {
	if lines, err := readLines(content, section); err != nil {
		return nil, err
	} else if len(lines) == 0 {
		return []string{}, nil
	} else {
		return strings.Split(strings.Join(lines, "\r\n"), "--[BR]"), nil
	}
}

// createMT create migration table
func createMT(db *sqlx.DB) error {
	cmd := `CREATE TABLE IF NOT EXISTS migrations(
        name VARCHAR(100) NOT NULL,
        section VARCHAR(10) NOT NULL,
		PRIMARY KEY(name, section)
    );`
	if stmt, err := db.Prepare(cmd); err != nil {
		return err
	} else if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}

// getRunnedMigrate get getRunnedMigrate file list
func getRunnedMigrate(db *sqlx.DB) ([]string, error) {
	var migrated []string
	if err := db.Select(&migrated, "select name from migrations WHERE section = 'migration';"); err != nil {
		return nil, err
	} else {
		return migrated, nil
	}
}

// getRunnedScript get migrated getRunnedScript file list
func getRunnedScript(db *sqlx.DB) ([]string, error) {
	var scripts []string
	if err := db.Select(&scripts, "select name from migrations WHERE section = 'script';"); err != nil {
		return nil, err
	} else {
		return scripts, nil
	}
}

// getRunnedSeed get getRunnedSeed file list
func getRunnedSeed(db *sqlx.DB) ([]string, error) {
	var seeded []string
	if err := db.Select(&seeded, "select name from migrations WHERE section = 'seed';"); err != nil {
		return nil, err
	} else {
		return seeded, nil
	}
}
