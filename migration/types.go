package migration

import (
	"regexp"
	"strconv"

	"github.com/gomig/utils"
)

// MigrationFile migration file with content
type MigrationFile struct {
	Name    string
	Content string
}

// Is check migration name without dash, timestamp and extension
func (m MigrationFile) Is(name string) bool {
	return m.name() == regexp.
		MustCompile(`^(\d+-)|(\.sql)$`).
		ReplaceAllString(name, "")
}

func (m MigrationFile) timestamp() int64 {
	if res, err := strconv.ParseInt(utils.ExtractNumbers(m.Name), 10, 64); err == nil {
		return res
	} else {
		return 0
	}
}

func (m MigrationFile) name() string {
	return regexp.
		MustCompile(`^(\d+-)|(\.sql)$`).
		ReplaceAllString(m.Name, "")
}

func (m MigrationFile) in(skips ...string) bool {
	for _, skip := range skips {
		if m.Name == skip {
			return true
		}
	}
	return false
}

// MigrationFiles migration file arrays
type MigrationFiles []MigrationFile

func (m MigrationFiles) Len() int {
	return len(m)
}

func (m MigrationFiles) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MigrationFiles) Less(i, j int) bool {
	return m[i].timestamp() < m[j].timestamp()
}

// Filter filter migrations by name
func (m MigrationFiles) Filter(name string) MigrationFiles {
	if name == "" {
		return m
	} else {
		res := make(MigrationFiles, 0)

		for _, migration := range m {
			if migration.Is(utils.Slugify(name)) {
				res = append(res, migration)
			}
		}

		return res
	}
}

// Reverse reverse array order
func (m MigrationFiles) Reverse() {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}
}
