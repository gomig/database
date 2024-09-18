package database

import (
	"database/sql"
	"strings"
)

type Updater[T any] interface {
	// NumericArgs specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder
	NumericArgs(isNumeric bool) Updater[T]
	// Table table name
	Table(table string) Updater[T]
	// Where update condition
	Where(cond string, args ...any) Updater[T]
	// Update update and return result
	Update(entity T) (sql.Result, error)
}

func NewUpdater[T any](db Executable) Updater[T] {
	updater := new(updaterDriver[T])
	updater.db = db
	updater.numeric = true
	return updater
}

type updaterDriver[T any] struct {
	db        Executable
	numeric   bool
	table     string
	condition string
	args      []any
}

func (updater *updaterDriver[T]) NumericArgs(numeric bool) Updater[T] {
	updater.numeric = numeric
	return updater
}

func (updater *updaterDriver[T]) Table(table string) Updater[T] {
	updater.table = table
	return updater
}

func (updater *updaterDriver[T]) Where(cond string, args ...any) Updater[T] {
	updater.condition = cond
	updater.args = args
	return updater
}

func (updater *updaterDriver[T]) Update(entity T) (sql.Result, error) {
	fields := structColumns(entity)
	for i, v := range fields {
		fields[i] = v + " = ?"
	}

	sql := strings.NewReplacer(
		"@table", updater.table,
		"@cond", updater.condition,
		"@fields", strings.Join(fields, " ,"),
	).Replace("UPDATE @table SET @fields WHERE @cond;")

	if updater.numeric {
		sql = numericArgs(sql, 1)
	}

	return updater.db.Exec(sql, append(structValues(entity), updater.args...)...)
}
