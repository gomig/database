package migration

import (
	"errors"
	"fmt"

	"github.com/gomig/utils"
	"github.com/jmoiron/sqlx"
)

type MigrationScript struct {
	Name   string
	CMD    string
	IsSeed bool
}

// ExecuteScripts execute named migration scripts
func ExecuteScripts(db *sqlx.DB, commands []MigrationScript) error {
	var migrated []string
	var seeded []string

	// Prepare
	if _, err := db.Exec(`CREATE TABLE IF NOT EXISTS migrations(name VARCHAR(100) NOT NULL PRIMARY KEY,mode VARCHAR(1) NOT NULL DEFAULT 'M');`); err != nil {
		return err
	}

	if err := db.Select(&migrated, "select name from migrations WHERE mode = 'M';"); err != nil {
		return err
	}

	if err := db.Select(&seeded, "select name from migrations WHERE mode = 'S';"); err != nil {
		return err
	}

	for _, script := range commands {
		if script.Name == "" {
			return errors.New("script name must pass")
		}

		if (script.IsSeed && utils.Contains(seeded, script.Name)) ||
			(!script.IsSeed && utils.Contains(migrated, script.Name)) {
			continue
		}

		for _, cmd := range getValidLines(string(script.CMD)) {
			if _, err := db.Exec(cmd); err != nil {
				return err
			}
		}

		if script.IsSeed {
			if _, err := db.Exec(fmt.Sprintf("INSERT INTO migrations(name, mode) VALUES('%s', 'S');", script.Name)); err != nil {
				return err
			}
		} else {
			if _, err := db.Exec(fmt.Sprintf("INSERT INTO migrations(name, mode) VALUES('%s', 'M');", script.Name)); err != nil {
				return err
			}
		}

	}

	return nil
}
