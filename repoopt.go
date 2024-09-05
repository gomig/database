package database

import (
	"database/sql"
	"strings"
)

// Types
type Executable interface {
	Exec(string, ...any) (sql.Result, error)
}

type Resolver[T any] func(*T) error

type RepoOption[T any] struct {
	driver    Driver
	args      []any
	oldNew    []string
	resolvers []Resolver[T]
}

// Methods
// NewOption generate new options with default parameters
func NewOption[T any]() *RepoOption[T] {
	opt := new(RepoOption[T])
	opt.driver = DriverPostgres
	opt.args = []any{}
	opt.oldNew = []string{}
	opt.resolvers = []Resolver[T]{}
	return opt
}

// WithDriver set database driver (default Postgres)
func (opt *RepoOption[T]) WithDriver(driver Driver) *RepoOption[T] {
	opt.driver = driver
	return opt
}

// WithArgs add args to query
func (opt *RepoOption[T]) WithArgs(args ...any) *RepoOption[T] {
	opt.args = append(opt.args, args...)
	return opt
}

// WithPlaceholder add placeholder for query to replace before execute
func (opt *RepoOption[T]) WithPlaceholder(oldNew ...string) *RepoOption[T] {
	opt.oldNew = append(opt.oldNew, oldNew...)
	return opt
}

// WithResolver add resolver to query
func (opt *RepoOption[T]) WithResolver(resolver Resolver[T]) *RepoOption[T] {
	opt.resolvers = append(opt.resolvers, resolver)
	return opt
}

// resolveOptions get first option or return default one
func resolveOptions[T any](options ...RepoOption[T]) *RepoOption[T] {
	if len(options) > 0 {
		return &options[0]
	} else {
		return NewOption[T]()
	}
}

// resolveQ normalize query for select sql commands
func (opt *RepoOption[T]) resolveQ(query string) string {
	var sample T
	fields := strings.Join(structQueryColumns(sample), ",")
	replacer := opt.replacerWithFields(fields)
	query = replacer.Replace(query)
	if opt.driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	return query
}

// resolve normalize query for non-select sql commands
func (opt *RepoOption[T]) resolve(query string) string {
	replacer := opt.replacer()
	query = replacer.Replace(query)
	if opt.driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	return query
}

// replacer generate new placeholder replacer for query
func (opt *RepoOption[T]) replacer() *strings.Replacer {
	return strings.NewReplacer(opt.oldNew...)
}

// replacerWithFields generate new placeholder replacer for query with @fields list
func (opt *RepoOption[T]) replacerWithFields(fields string) *strings.Replacer {
	opt.oldNew = append(opt.oldNew, "@fields", fields)
	return strings.NewReplacer(opt.oldNew...)
}
