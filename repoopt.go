package database

import (
	"database/sql"
	"strings"
)

// types
type RepositoryOpt[T any] func(*Option[T])
type Resolver[T any] func(*T) error
type Option[T any] struct {
	driver    Driver
	table     string
	args      []any
	oldNew    []string
	resolvers []Resolver[T]
}
type Executable interface {
	Exec(string, ...any) (sql.Result, error)
}

// newOption generate new option with defaults
func newOption[T any](opts ...RepositoryOpt[T]) *Option[T] {
	opt := new(Option[T])
	opt.args = []any{}
	opt.oldNew = []string{}
	opt.resolvers = []Resolver[T]{}
	opt.driver = DriverPostgres
	for _, o := range opts {
		o(opt)
	}
	return opt
}

// Methods
// resolveQ normalize query for select sql commands
func (opt Option[T]) resolveQ(query string) string {
	var sample T
	fields := strings.Join(structQueryColumns(sample), ",")
	replacer := opt.ReplacerWithFields(fields)
	query = replacer.Replace(query)
	if opt.driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	return query
}

// resolve normalize query for non-select sql commands
func (opt Option[T]) resolve(query string) string {
	replacer := opt.Replacer()
	query = replacer.Replace(query)
	if opt.driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	return query
}

// Replacer generate new placeholder replacer for query
func (opt Option[T]) Replacer() *strings.Replacer {
	return strings.NewReplacer(opt.oldNew...)
}

// Replacer generate new placeholder replacer for query with @fields list
func (opt Option[T]) ReplacerWithFields(fields string) *strings.Replacer {
	opt.oldNew = append(opt.oldNew, "@fields", fields)
	return strings.NewReplacer(opt.oldNew...)
}

// Options
// WithDriver is a functional option to set database driver
func WithDriver[T any](driver Driver) RepositoryOpt[T] {
	return func(opt *Option[T]) {
		opt.driver = driver
	}
}

// WithTable is a functional option to set database table for insert and update query
func WithTable[T any](table string) RepositoryOpt[T] {
	return func(opt *Option[T]) {
		opt.table = table
	}
}

// WithArgs is a functional option to pass args to query
func WithArgs[T any](args ...any) RepositoryOpt[T] {
	return func(opt *Option[T]) {
		opt.args = append(opt.args, args...)
	}
}

// WithPlaceholder is a functional option to fill query placeholders
func WithPlaceholder[T any](oldNew ...string) RepositoryOpt[T] {
	return func(opt *Option[T]) {
		opt.oldNew = append(opt.oldNew, oldNew...)
	}
}

// WithResolver is a functional option to add resolver to query
func WithResolver[T any](resolver ...Resolver[T]) RepositoryOpt[T] {
	return func(opt *Option[T]) {
		opt.resolvers = append(opt.resolvers, resolver...)
	}
}
