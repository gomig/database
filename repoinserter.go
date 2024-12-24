package database

import (
	"database/sql"
	"strings"
)

type Inserter[T any] interface {
	// NumericArgs specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder
	NumericArgs(isNumeric bool) Inserter[T]
	// QuoteFields specifies whether to use quoted field name ("id", "name") or not
	QuoteFields(quoted bool) Inserter[T]
	// Table table name
	Table(table string) Inserter[T]
	// Insert insert and return result
	Insert(entity T) (sql.Result, error)
}

func NewInserter[T any](db Executable) Inserter[T] {
	inserter := new(insertDriver[T])
	inserter.db = db
	inserter.numeric = true
	inserter.quoted = true
	return inserter
}

type insertDriver[T any] struct {
	db      Executable
	numeric bool
	quoted  bool
	table   string
}

func (inserter *insertDriver[T]) NumericArgs(numeric bool) Inserter[T] {
	inserter.numeric = numeric
	return inserter
}

func (inserter *insertDriver[T]) QuoteFields(quoted bool) Inserter[T] {
	inserter.quoted = quoted
	return inserter
}

func (inserter *insertDriver[T]) Table(table string) Inserter[T] {
	inserter.table = table
	return inserter
}

func (inserter *insertDriver[T]) Insert(entity T) (sql.Result, error) {
	fields := structColumns(entity, inserter.quoted)
	placeholders := make([]string, 0)
	for range fields {
		placeholders = append(placeholders, "?")
	}

	sql := strings.NewReplacer(
		"@table", inserter.table,
		"@fields", strings.Join(fields, " ,"),
		"@values", strings.Join(placeholders, " ,"),
	).Replace("INSERT INTO @table (@fields) VALUES(@values);")

	if inserter.numeric {
		sql = numericArgs(sql, 1)
	}

	return inserter.db.Exec(sql, structValues(entity)...)
}
