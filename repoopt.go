package database

import (
	"database/sql"
	"strings"
)

type Executable interface {
	Exec(string, ...any) (sql.Result, error)
}

type Resolver[T any] func(*T) error

type Option[T any] interface {
	WithDriver(Driver) Option[T]
	WithPlaceholder(...string) Option[T]
	WithResolver(Resolver[T]) Option[T]

	getDriver() Driver
	getResolvers() []Resolver[T]

	resolveQuery(string) string
	resolve(string) string
	replacer() *strings.Replacer
	replacerWithFields(string) *strings.Replacer
}

type optionDriver[T any] struct {
	driver       Driver
	placeholders []string
	resolvers    []Resolver[T]
}

// WithDriver set database driver (default Postgres)
func (opt *optionDriver[T]) WithDriver(driver Driver) Option[T] {
	opt.driver = driver
	return opt
}

// WithPlaceholder add placeholder for query to replace before execute
func (opt *optionDriver[T]) WithPlaceholder(oldNew ...string) Option[T] {
	opt.placeholders = append(opt.placeholders, oldNew...)
	return opt
}

// WithResolver add resolver to query
func (opt *optionDriver[T]) WithResolver(resolver Resolver[T]) Option[T] {
	opt.resolvers = append(opt.resolvers, resolver)
	return opt
}

func (opt *optionDriver[T]) getDriver() Driver {
	return opt.driver
}

func (opt *optionDriver[T]) getResolvers() []Resolver[T] {
	return opt.resolvers
}

func (opt *optionDriver[T]) resolveQuery(query string) string {
	var sample T
	fields := strings.Join(structQueryColumns(sample), ",")
	replacer := opt.replacerWithFields(fields)
	query = replacer.Replace(query)
	if opt.driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	return query
}

func (opt *optionDriver[T]) resolve(query string) string {
	replacer := opt.replacer()
	query = replacer.Replace(query)
	if opt.driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	return query
}

func (opt *optionDriver[T]) replacer() *strings.Replacer {
	return strings.NewReplacer(opt.placeholders...)
}

func (opt *optionDriver[T]) replacerWithFields(fields string) *strings.Replacer {
	opt.placeholders = append(opt.placeholders, "@fields", fields)
	return strings.NewReplacer(opt.placeholders...)
}

// NewOption generate new options with default parameters
func NewOption[T any]() Option[T] {
	opt := new(optionDriver[T])
	opt.driver = DriverPostgres
	opt.placeholders = []string{}
	opt.resolvers = []Resolver[T]{}
	return opt
}
