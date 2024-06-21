package migration

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

func getValidLines(content string) []string {
	var res []string
	list := strings.Split(content, ";")
	for _, item := range list {
		item = strings.Trim(item, " ")
		if item != "" && item != "\r\n" {
			res = append(res, item)
		}
	}
	return res
}

func validateStatement(statement string, db *sqlx.DB) error {
	stmt, err := db.Prepare(statement)
	defer stmt.Close()
	return err
}

func createMigrationTable(db *sqlx.DB) {
	cmd := `CREATE TABLE IF NOT EXISTS migrations(
        name VARCHAR(100) PRIMARY KEY,
        is_seed BOOLEAN NOT NULL DEFAULT FALSE
    );`
	stmt, err := db.Prepare(cmd)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getMigratedFiles(db *sqlx.DB, isSeed bool) []string {
	var migrated []string
	rows, err := db.Query("select name from migrations WHERE is_seed = ?;", isSeed)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for rows.Next() {
		var migration string
		err := rows.Scan(&migration)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		migrated = append(migrated, migration)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return migrated
}

func getMigrationFiles(migrated []string, dir string) ([]string, []string) {
	temp := make(map[string]string)
	filepath.Walk(dir, func(path string, f os.FileInfo, _ error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".sql", f.Name())

			if err == nil && r && !strings.HasPrefix(f.Name(), "_") {
				if len(migrated) == 0 {
					temp[f.Name()] = path
				} else {
					found := false
					for _, mtd := range migrated {
						if f.Name() == mtd {
							found = true
							break
						}
					}
					if !found {
						temp[f.Name()] = path
					}
				}
			}
		}
		return nil
	})

	keys := make([]string, 0, len(temp))
	vals := make([]string, 0, len(temp))
	for k := range temp {
		keys = append(keys, k)
	}
	sort.Sort(byNumber(keys))
	for _, k := range keys {
		vals = append(vals, temp[k])
	}

	return keys, vals
}
