package migration

import (
	"github.com/jmoiron/sqlx"
)

// InitMigration prepare database to run migrations
func InitMigration(db *sqlx.DB) error {
	if stmt, err := db.Prepare(
		`CREATE TABLE IF NOT EXISTS migrations (
			name VARCHAR(100) NOT NULL,
			stage VARCHAR(30) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			PRIMARY KEY(name, stage)
		);`,
	); err != nil {
		return err
	} else if _, err = stmt.Exec(); err != nil {
		return err
	}
	return nil
}

// Migrate run migration on database
func Migrate(db *sqlx.DB, stage string, files ...File) ([]string, error) {
	if err := InitMigration(db); err != nil {
		return nil, err
	} else if migrated, err := StageMigrated(db, stage); err != nil {
		return nil, err
	} else {
		res := make([]string, 0)

		for _, file := range files {
			if ok, err := file.Migrate(db, stage, migrated.Names()...); err != nil {
				return nil, err
			} else if ok {
				res = append(res, file.RealName())
			}
		}

		return res, nil
	}
}

// Rollback run migration down on database
func Rollback(db *sqlx.DB, files ...File) ([]string, error) {
	if err := InitMigration(db); err != nil {
		return nil, err
	} else if migrated, err := StageMigrated(db, "DOWN"); err != nil {
		return nil, err
	} else {
		res := make([]string, 0)

		for _, file := range files {
			if ok, err := file.Rollback(db, migrated.Names()...); err != nil {
				return nil, err
			} else if ok {
				res = append(res, file.RealName())
			}
		}

		return res, nil
	}
}
