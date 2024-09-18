package migration

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strconv"

	"github.com/gomig/utils"
	"github.com/jmoiron/sqlx"
)

// ReadFS read migration from file system
func ReadFS(dir, ext string) (Files, error) {
	result := make(Files, 0)
	if err := filepath.Walk(uri(dir), func(path string, f os.FileInfo, _ error) error {
		if ok, err := regexp.MatchString(`^([0-9])(.+)(\.`+ext+`)$`, f.Name()); err != nil {
			return err
		} else if ok && !f.IsDir() {
			if content, err := os.ReadFile(path); err != nil {
				return err
			} else {
				result = append(result, NewFile(f.Name(), ext, string(content)))
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

// NewFile create new migration file from content
//
// file name must full with extension
func NewFile(name, ext, content string) File {
	return File{name: name, ext: ext, content: content}
}

// File migration file
type File struct {
	name    string
	content string
	ext     string
}

// RealName get file name
func (file File) RealName() string {
	return file.name
}

// Content get file content
func (file File) Content() string {
	return file.content
}

// Get timestamp part of filename
func (file File) Timestamp() int64 {
	if res, err := strconv.ParseInt(utils.ExtractNumbers(file.name), 10, 64); err == nil {
		return res
	} else {
		return 0
	}
}

// Name get readable file name
func (file File) Name() string {
	return regexp.
		MustCompile(`^(\d+-)|(\.`+file.ext+`)$`).
		ReplaceAllString(file.name, "")
}

// Extension get file extension
func (file File) Extension() string {
	return file.ext
}

// Is check migration name without dash, timestamp and extension
func (file File) Is(name string) bool {
	return file.Name() == regexp.
		MustCompile(`^(\d+-)|(\.`+file.ext+`)$`).
		ReplaceAllString(name, "")
}

// MustSkip check if file must skipped
func (file File) MustSkip(skips ...string) bool {
	return slices.Contains(skips, file.name)
}

// Scripts get file scripts list for stage
func (file File) Scripts(stage string) ([]string, error) {
	return scriptsOf(file.content, stage)
}

// Migrate migrate file stage, return true if migrate or false if file already migrated
func (file File) Migrate(db *sqlx.DB, stage string, migrated ...string) (bool, error) {
	if file.name == "" {
		return false, errors.New("migration name is empty")
	} else if file.MustSkip(migrated...) {
		return false, nil
	} else if scripts, err := file.Scripts(stage); err != nil {
		return false, err
	} else if tx, err := db.BeginTx(context.Background(), nil); err != nil {
		return false, err
	} else if len(scripts) == 0 {
		return false, nil
	} else {
		for _, script := range scripts {
			if _, err := tx.Exec(script); err != nil {
				tx.Rollback()
				return false, err
			}

			if _, err := tx.Exec(fmt.Sprintf(
				`INSERT INTO migrations (name, stage) VALUES('%s', '%s');`,
				file.name, stage,
			)); err != nil {
				tx.Rollback()
				return false, err
			}
		}

		if err := tx.Commit(); err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

// Rollback run migration down for file
func (file File) Rollback(db *sqlx.DB, migrated ...string) (bool, error) {
	if file.name == "" {
		return false, errors.New("migration name is empty")
	} else if !file.MustSkip(migrated...) {
		return false, nil
	} else if scripts, err := file.Scripts("DOWN"); err != nil {
		return false, err
	} else if tx, err := db.BeginTx(context.Background(), nil); err != nil {
		return false, err
	} else if len(scripts) == 0 {
		return false, nil
	} else {
		for _, script := range scripts {
			if _, err := tx.Exec(script); err != nil {
				tx.Rollback()
				return false, err
			}

			if _, err := tx.Exec(fmt.Sprintf(
				`DELETE FROM migrations WHERE name = '%s';`,
				file.name,
			)); err != nil {
				tx.Rollback()
				return false, err
			}
		}

		if err := tx.Commit(); err != nil {
			return false, err
		} else {
			return true, nil
		}
	}
}

// Files migration file list
type Files []File

// Len get files length
func (files Files) Len() int {
	return len(files)
}

// Swap swap item i and j
func (files Files) Swap(i, j int) {
	files[i], files[j] = files[j], files[i]
}

// Less check if name timestamp is smaller
func (files Files) Less(i, j int) bool {
	return files[i].Timestamp() < files[j].Timestamp()
}

// Filter filter files by name
func (files Files) Filter(name string) Files {
	if name == "" {
		return files
	} else {
		res := make(Files, 0)

		for _, file := range files {
			if file.Is(utils.Slugify(name)) {
				res = append(res, file)
			}
		}

		return res
	}
}

// Reverse reverse array order
func (files Files) Reverse() {
	for i, j := 0, len(files)-1; i < j; i, j = i+1, j-1 {
		files[i], files[j] = files[j], files[i]
	}
}
