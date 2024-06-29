package migration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"github.com/jmoiron/sqlx"
)

// ReadDirectory read migration from file system
func ReadDirectory(dir string) (MigrationsT, error) {
	result := make(MigrationsT, 0)
	if err := filepath.Walk(uri(dir), func(path string, f os.FileInfo, _ error) error {
		if ok, err := regexp.MatchString(`^([0-9])(.+)(\.sql)$`, f.Name()); err != nil {
			return err
		} else if ok && !f.IsDir() {
			if content, err := os.ReadFile(path); err != nil {
				return err
			} else {
				result = append(result, MigrationT{
					Name:    f.Name(),
					Content: string(content),
				})
			}
		}
		return nil
	}); err != nil {
		return nil, err
	} else {
		sort.Sort(result)
		return result, nil
	}
}

// Migrate run migration up on database
//
// pass migration name to run migrate on specific migration
func Migrate(db *sqlx.DB, migrations MigrationsT, name string) ([]string, error) {
	if err := createMT(db); err != nil {
		return nil, err
	} else if migrated, err := getRunnedMigrate(db); err != nil {
		return nil, err
	} else {
		migrations := migrations.Filter(name)
		res := make([]string, 0)
		for _, migration := range migrations {
			if migration.Name == "" {
				return nil, errors.New("migration name is empty")
			} else if migration.In(migrated...) {
				continue
			} else if scripts, err := readScripts(migration.Content, "up"); err != nil {
				return nil, err
			} else {
				for _, script := range scripts {
					if _, err := db.Exec(script); err != nil {
						return nil, err
					}
				}

				if _, err := db.Exec(fmt.Sprintf("INSERT INTO migrations(name, section) VALUES('%s', 'migration');", migration.Name)); err != nil {
					return nil, err
				}

				res = append(res, migration.Name)
			}
		}
		return res, nil
	}
}

// Script run script migration on database
//
// pass migration name to run script on specific migration
func Script(db *sqlx.DB, migrations MigrationsT, name string) ([]string, error) {
	if err := createMT(db); err != nil {
		return nil, err
	} else if migrated, err := getRunnedScript(db); err != nil {
		return nil, err
	} else {
		migrations := migrations.Filter(name)
		res := make([]string, 0)
		for _, migration := range migrations {
			if migration.Name == "" {
				return nil, errors.New("migration name is empty")
			} else if migration.In(migrated...) {
				continue
			} else if scripts, err := readScripts(migration.Content, "script"); err != nil {
				return nil, err
			} else {
				for _, script := range scripts {
					if _, err := db.Exec(script); err != nil {
						return nil, err
					}
				}

				if _, err := db.Exec(fmt.Sprintf("INSERT INTO migrations(name, section) VALUES('%s', 'script');", migration.Name)); err != nil {
					return nil, err
				}
				res = append(res, migration.Name)
			}
		}
		return res, nil
	}
}

// Seed run seed on database
//
// pass migration name to run seed on specific migration
func Seed(db *sqlx.DB, migrations MigrationsT, name string) ([]string, error) {
	if err := createMT(db); err != nil {
		return nil, err
	} else if seeded, err := getRunnedSeed(db); err != nil {
		return nil, err
	} else {
		migrations := migrations.Filter(name)
		res := make([]string, 0)
		for _, migration := range migrations {
			if migration.Name == "" {
				return nil, errors.New("migration name is empty")
			} else if migration.In(seeded...) {
				continue
			} else if scripts, err := readScripts(migration.Content, "seed"); err != nil {
				return nil, err
			} else {
				for _, script := range scripts {
					if _, err := db.Exec(script); err != nil {
						return nil, err
					}
				}

				if _, err := db.Exec(fmt.Sprintf("INSERT INTO migrations(name, section) VALUES('%s', 'seed');", migration.Name)); err != nil {
					return nil, err
				}
				res = append(res, migration.Name)
			}
		}
		return res, nil
	}
}

// Rollback run migration down on database
//
// pass migration name to run rollback on specific migration
func Rollback(db *sqlx.DB, migrations MigrationsT, name string) ([]string, error) {
	if err := createMT(db); err != nil {
		return nil, err
	} else if migrated, err := getRunnedMigrate(db); err != nil {
		return nil, err
	} else {
		migrations := migrations.Filter(name)
		res := make([]string, 0)
		for _, migration := range migrations {
			if migration.Name == "" {
				return nil, errors.New("migration name is empty")
			} else if !migration.In(migrated...) {
				continue
			} else if scripts, err := readScripts(migration.Content, "down"); err != nil {
				return nil, err
			} else {
				for _, script := range scripts {
					if _, err := db.Exec(script); err != nil {
						return nil, err
					}
				}

				if _, err := db.Exec(fmt.Sprintf("DELETE FROM migrations WHERE name = '%s';", migration.Name)); err != nil {
					return nil, err
				}

				res = append(res, migration.Name)
			}
		}
		return res, nil
	}
}
