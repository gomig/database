package migration

import (
	"regexp"
	"strconv"

	"github.com/gomig/utils"
)

type MigrationT struct {
	Name    string
	Content string
}

func (m MigrationT) In(skips ...string) bool {
	for _, skip := range skips {
		if m.Name == skip {
			return true
		}
	}
	return false
}

type MigrationsT []MigrationT

func getTimestamp(str string) int {
	if res, err := strconv.Atoi(utils.ExtractNumbers(str)); err == nil {
		return res
	} else {
		return 0
	}
}

func (m MigrationsT) Len() int {
	return len(m)
}

func (m MigrationsT) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MigrationsT) Less(i, j int) bool {
	return getTimestamp(m[i].Name) < getTimestamp(m[j].Name)
}

// Filter filter migrations by name
func (m MigrationsT) Filter(name string) MigrationsT {
	clear := func(str string) string {
		return regexp.
			MustCompile(`^(\d+-)|(\.sql)$`).
			ReplaceAllString(str, "")
	}

	if name == "" {
		return m
	} else {
		res := make(MigrationsT, 0)

		for _, migration := range m {
			if clear(migration.Name) == clear(utils.Slugify(name)) {
				res = append(res, migration)
			}
		}

		return res
	}
}

// Reverse reverse array order
func (m MigrationsT) Reverse() {
	for i, j := 0, len(m)-1; i < j; i, j = i+1, j-1 {
		m[i], m[j] = m[j], m[i]
	}
}
