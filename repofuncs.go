package database

import (
	"fmt"
	"strings"
)

// ResolveQuery get list of fields from struct `q` and `db` tag and replace with `SELECT ...;` keyword in query
func ResolveQuery[T any](query string, driver Driver) string {
	var sample T
	if strings.Contains(strings.ToLower(query), "select ...") {
		query = strings.Replace(query, "...", strings.Join(structQueryColumns(sample), ","), 1)
	}

	if driver == DriverPostgres {
		query = numericArgs(query, 1)
	}

	return query
}

// ResolveInsert create insert cmd for table
//
// @returns insert command and params as result
func ResolveInsert[T any](entity T, table string, driver Driver) (string, []any) {
	fields := structColumns(entity)
	placeholders := make([]string, 0)
	for range fields {
		placeholders = append(placeholders, "?")
	}

	command := fmt.Sprintf(
		"INSERT INTO %s(%s) VALUES(%s);",
		table,
		strings.Join(fields, ","),
		strings.Join(placeholders, ","),
	)

	if driver == DriverPostgres {
		command = numericArgs(command, 1)
	}

	return command, structValues(entity)
}

// ResolveUpdate create update cmd for table and
//
// You must pass condition argument with ?
// @returns query and params as result
func ResolveUpdate[T any](entity T, table string, driver Driver, condition string, args ...any) (string, []any) {
	fields := structColumns(entity)
	for i, v := range fields {
		fields[i] = fmt.Sprintf("%s = ?", v)
	}

	command := fmt.Sprintf(
		"UPDATE %s SET %s WHERE %s;",
		table,
		strings.Join(fields, ","),
		condition,
	)

	if driver == DriverPostgres {
		command = numericArgs(command, 1)
	}

	return command, append(structValues(entity), args...)
}
