package database

// Query object
type Query struct {
	Type    string
	Query   string
	Params  []any
	Closure bool
}

type QueryBuilder interface {
	// And add new simple condition to query with AND
	And(cond string, args ...any) QueryBuilder
	// AndIf add new And condition if first parameter is true
	AndIf(ifCond bool, cond string, args ...any) QueryBuilder
	// Or add new simple condition to query with OR
	Or(cond string, args ...any) QueryBuilder
	// OrIf add new Or condition if first parameter is true
	OrIf(ifCond bool, cond string, args ...any) QueryBuilder
	// AndClosure add new condition to query with AND in nested ()
	AndClosure(cond string, args ...any) QueryBuilder
	// AndClosureIf add new AndClosure condition if first parameter is true
	AndClosureIf(ifCond bool, cond string, args ...any) QueryBuilder
	// OrClosure add new condition to query with OR in nested ()
	OrClosure(cond string, args ...any) QueryBuilder
	// OrClosureIf add new AndClosure condition if first parameter is true
	OrClosureIf(ifCond bool, cond string, args ...any) QueryBuilder
	// RawPostgres get raw generated query for postgres
	RawPostgres(counter int) string
	// RawSQL get raw generated query for mysql
	RawSQL() string
	// ToPostgres generate query string for postgres with numeric arguments based on counter
	// this method replace @query with Raw() value
	// this method replace @where with `WHERE Raw()` method
	// you can use @[key], value to replacement in query string
	ToPostgres(query string, counter int, replacements ...string) string
	// ToSQL generate query string for sql with ? arguments
	// this method replace @query with Raw() value
	// this method replace @where with `WHERE Raw()` method
	// you can use @[key], value to replacement in query string
	ToSQL(query string, replacements ...string) string
	// Params get list of query parameters
	Params() []any
}

// NewQuery generate new query builder
func NewQuery() QueryBuilder {
	res := new(qBuilder)
	res.queries = make([]Query, 0)
	return res
}
