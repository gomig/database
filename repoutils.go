package database

import (
	"fmt"
	"reflect"
	"strings"
)

// ResolveQuery get list of fields from struct `q` and `db` tag and replace with `SELECT @fields;` keyword in query
func ResolveQuery[T any](query string, driver Driver) string {
	var sample T
	if strings.Contains(strings.ToLower(query), "@fields") {
		query = strings.Replace(query, "@fields", strings.Join(structQueryColumns(sample), ","), 1)
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

// structQueryColumns get columns list from `q` or `db` struct tag
func structQueryColumns(v any) []string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() || !val.Elem().CanInterface() {
			return []string{}
		} else {
			return structQueryColumns(val.Elem().Interface())
		}
	} else if val.Kind() != reflect.Struct {
		return []string{}
	} else {
		res := make([]string, 0)
		typ := reflect.TypeOf(v)
		for i := 0; i < typ.NumField(); i++ {
			if typ.Field(i).IsExported() {
				if typ.Field(i).Anonymous {
					res = append(res, structQueryColumns(val.Field(i).Interface())...)
				} else {
					if q, ok := typ.Field(i).Tag.Lookup("q"); ok {
						if q != "-" && q != "" {
							res = append(res, q)
						}
					} else if tag, ok := typ.Field(i).Tag.Lookup("db"); ok && tag != "-" && tag != "" {
						res = append(res, tag)
					}
				}
			}
		}
		return res
	}
}

// structColumns get columns list from `db` struct tag
func structColumns(v any) []string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() || !val.Elem().CanInterface() {
			return []string{}
		} else {
			return structColumns(val.Elem().Interface())
		}
	} else if val.Kind() != reflect.Struct {
		return []string{}
	} else {
		res := make([]string, 0)
		typ := reflect.TypeOf(v)
		for i := 0; i < typ.NumField(); i++ {
			if typ.Field(i).IsExported() {
				if tag, ok := typ.Field(i).Tag.Lookup("db"); ok && tag != "-" && tag != "" {
					res = append(res, tag)
				}
			}
		}
		return res
	}
}

// structValues get struct value where `db` tag not - or empty
func structValues(v any) []any {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() || !val.Elem().CanInterface() {
			return []any{}
		} else {
			return structValues(val.Elem().Interface())
		}
	} else if val.Kind() != reflect.Struct {
		return []any{}
	} else {
		res := make([]any, 0)
		typ := reflect.TypeOf(v)
		for i := 0; i < typ.NumField(); i++ {
			if typ.Field(i).IsExported() {
				if tag, ok := typ.Field(i).Tag.Lookup("db"); ok && tag != "-" && tag != "" {
					res = append(res, val.Field(i).Interface())
				}
			}
		}
		return res
	}
}

// numericArgs convert ? placeholder to numeric $1 placeholder
func numericArgs(query string, counter int) string {
	if counter <= 0 {
		counter = 1
	}
	for {
		if strings.Contains(query, "?") {
			query = strings.Replace(query, "?", fmt.Sprintf("$%d", counter), 1)
			counter++
		} else {
			break
		}
	}
	return query
}
