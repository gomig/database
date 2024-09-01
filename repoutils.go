package database

import (
	"fmt"
	"reflect"
	"strings"
)

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
