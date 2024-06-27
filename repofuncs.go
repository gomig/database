package database

import (
	"fmt"
	"strings"
)

// ResolveQuery get list of fields from struct `q` or `db` tag and replace with ? keyword in query
func ResolveQuery[T any](query string) string {
	var sample T
	if strings.Contains(query, "?") {
		return strings.ReplaceAll(query, "?", strings.Join(structQueryColumns(sample), ","))
	} else {
		return query
	}
}

// ResolveInsert create insert cmd for table
//
// @returns insert command and params as result
func ResolveInsert[T any](entity T, table string, driver Driver) (string, []any) {
	tags := structColumns(entity)
	placeholders := make([]string, 0)

	if driver == DriverMySQL {
		for range tags {
			placeholders = append(placeholders, "?")
		}
	} else {
		for i := range tags {
			placeholders = append(placeholders, fmt.Sprintf("$%d", i+1))
		}
	}

	return fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES(%s);",
			table,
			strings.Join(tags, ","),
			strings.Join(placeholders, ","),
		),
		structValues(entity)
}

// ResolveUpdate create update cmd for table and
//
// You must pass condition argument with ?
// @returns query and params as result
func ResolveUpdate[T any](entity T, table string, driver Driver, condition string, args ...any) (string, []any) {
	tags := structColumns(entity)

	if driver == DriverMySQL {
		for i, v := range tags {
			tags[i] = fmt.Sprintf("%s = ?", v)
		}
	} else {
		counter := 0
		for i, v := range tags {
			tags[i] = fmt.Sprintf("%s = $%d", v, i+1)
			counter = i
		}
		counter++
		for {
			if strings.Contains(condition, "?") {
				condition = strings.Replace(condition, "?", fmt.Sprintf("$%d", counter+1), 1)
				counter++
			} else {
				break
			}
		}
	}

	return fmt.Sprintf(
			"UPDATE %s SET %s WHERE %s;",
			table,
			strings.Join(tags, ","),
			condition,
		),
		append(structValues(entity), args...)
}
