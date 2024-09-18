package migration

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Migration struct {
	Name  string `db:"name"`
	Stage string `db:"stage"`
}

type Migrations []Migration

// Get migration file names
func (migrations Migrations) Names() []string {
	result := make([]string, 0)
	for _, migration := range migrations {
		result = append(result, migration.Name)
	}
	return result
}

// Group migration by stage
func (migrations Migrations) GroupByStage() map[string][]string {
	result := make(map[string][]string)
	for _, file := range migrations {
		result[file.Stage] = append(result[file.Stage], file.Name)
	}
	return result
}

// Group migration by file
func (migrations Migrations) GroupByFile() map[string][]string {
	result := make(map[string][]string)
	for _, file := range migrations {
		result[file.Name] = append(result[file.Name], file.Stage)
	}
	return result
}

// StageMigrated get migrated items for stage
func StageMigrated(db *sqlx.DB, stage string) (Migrations, error) {
	result := make(Migrations, 0)
	if err := db.Select(
		&result,
		fmt.Sprintf(
			`SELECT name, stage FROM migrations WHERE stage = '%s' ORDER BY created_at ASC;`,
			stage,
		),
	); err == sql.ErrNoRows {
		return Migrations{}, nil
	} else if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

// Migrated get migrated items
func Migrated(db *sqlx.DB) (Migrations, error) {
	result := make(Migrations, 0)
	if err := db.Select(
		&result,
		`SELECT name, stage FROM migrations ORDER BY created_at ASC;`,
	); err == sql.ErrNoRows {
		return Migrations{}, nil
	} else if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}
