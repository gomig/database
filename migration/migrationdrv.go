package migration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type migration struct {
	root string
	ext  string
	db   *sqlx.DB
	fs   FS
}

func (driver migration) Root() string {
	return driver.root
}

func (driver migration) Extension() string {
	return driver.ext
}

func (driver migration) Path(name string) string {
	return normalizePath(driver.root, name)
}

func (driver migration) Init() error {
	if driver.db == nil {
		return errors.New("database driver is nil")
	} else if stmt, err := driver.db.Prepare(`
		CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(100) NOT NULL,
			stage VARCHAR(30) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY(name, stage)
		);
	`); err != nil {
		return err
	} else if _, err = stmt.Exec(); err != nil {
		return err
	} else {
		return nil
	}
}

func (driver migration) Summary() (Summary, error) {
	result := make(Summary, 0)
	if driver.db == nil {
		return nil, errors.New("database driver is nil")
	} else if err := driver.db.Select(&result, `
		SELECT name, stage
		FROM migrations
		ORDER BY created_at ASC;
	`); err == sql.ErrNoRows {
		return Summary{}, nil
	} else if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (driver migration) StageSummary(stage string) (Summary, error) {
	result := make(Summary, 0)
	if driver.db == nil {
		return nil, errors.New("database driver is nil")
	} else if err := driver.db.Select(
		&result,
		fmt.Sprintf(`
			SELECT name, stage
			FROM migrations
			WHERE stage = '%s'
			ORDER BY created_at ASC;
		`, stage),
	); err == sql.ErrNoRows {
		return Summary{}, nil
	} else if err != nil {
		return nil, err
	} else {
		return result, nil
	}
}

func (driver migration) Up(stage string, only ...string) ([]string, error) {
	errOf := func(file, mode string, err error) error {
		return fmt.Errorf(
			"[%s] (%s %s): %s",
			file, mode, stage, err.Error(),
		)
	}

	if driver.db == nil {
		return nil, errors.New("database driver is nil")
	}

	// Get old migrated files
	migrated, err := driver.StageSummary(stage)
	if err != nil {
		return nil, err
	}

	// Parse and load files
	files := driver.fs.Filter(only...).ExcludeMigrated(migrated.Names()...)
	if files.Len() == 0 {
		return nil, nil
	}

	// run migrations
	result := make([]string, 0)
	onFail := func(err error) ([]string, error) {
		return nil, err
	}

	tx, err := driver.db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		for _, file := range files {
			if scripts, err := file.UpScripts(stage); err != nil {
				return onFail(errOf(file.Name(), "PARSE UP", err))
			} else if len(scripts) == 0 {
				continue
			} else {
				if _, err := tx.Exec(scripts); err != nil {
					return onFail(errOf(file.Name(), "UP", err))
				}

				if _, err := tx.Exec(fmt.Sprintf(
					`INSERT INTO migrations (name, stage) VALUES('%s', '%s');`,
					file.Name(), stage,
				)); err != nil {
					return onFail(errOf(file.Name(), "UP", err))
				} else {
					result = append(result, file.Name())
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return onFail(err)
		} else {
			return result, nil
		}
	}
}

func (driver migration) Down(stage string, only ...string) ([]string, error) {
	errOf := func(file, mode string, err error) error {
		return fmt.Errorf(
			"[%s] (%s %s): %s",
			file, mode, stage, err.Error(),
		)
	}

	if driver.db == nil {
		return nil, errors.New("database driver is nil")
	}

	// Get old migrated files
	migrated, err := driver.StageSummary(stage)
	if err != nil {
		return nil, err
	}

	// Parse and load files
	files := driver.fs.Reverse().Filter(only...).FilterMigrated(migrated.Names()...)
	if files.Len() == 0 {
		return nil, nil
	}

	// run migrations
	result := make([]string, 0)
	onFail := func(err error) ([]string, error) {
		return nil, err
	}

	tx, err := driver.db.BeginTx(context.Background(), nil)
	if err != nil {
		tx.Rollback()
		return nil, err
	} else {
		for _, file := range files {
			if scripts, err := file.DownScripts(stage); err != nil {
				return onFail(errOf(file.Name(), "PARSE DOWN", err))
			} else if len(scripts) == 0 {
				continue
			} else {
				if _, err := tx.Exec(scripts); err != nil {
					return onFail(errOf(file.Name(), "DOWN", err))
				}

				if _, err := tx.Exec(fmt.Sprintf(
					`DELETE FROM migrations WHERE name = '%s' AND stage = '%s';`,
					file.Name(), stage,
				)); err != nil {
					return onFail(errOf(file.Name(), "DOWN", err))
				} else {
					result = append(result, file.Name())
				}
			}
		}

		if err := tx.Commit(); err != nil {
			return onFail(err)
		} else {
			return result, nil
		}
	}
}
