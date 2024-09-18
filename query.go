package database

// Query builder
type QueryBuilder interface {
	// And add new simple condition to query with AND
	And(query string, args ...any) QueryBuilder
	// AndIf add new And condition if first parameter is true
	AndIf(cond bool, query string, args ...any) QueryBuilder
	// Or add new simple condition to query with OR
	Or(query string, args ...any) QueryBuilder
	// OrIf add new Or condition if first parameter is true
	OrIf(cond bool, query string, args ...any) QueryBuilder
	// AndClosure add new condition to query with AND in nested ()
	AndClosure(query string, args ...any) QueryBuilder
	// AndClosureIf add new AndClosure condition if first parameter is true
	AndClosureIf(cond bool, query string, args ...any) QueryBuilder
	// OrClosure add new condition to query with OR in nested ()
	OrClosure(query string, args ...any) QueryBuilder
	// OrClosureIf add new AndClosure condition if first parameter is true
	OrClosureIf(cond bool, query string, args ...any) QueryBuilder
	// NumericArgs specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder
	NumericArgs(bool) QueryBuilder
	// NumericStart set numeric argument start for numeric args mode
	NumericStart(int) QueryBuilder
	// Replace replace phrase in query string before run
	Replace(string, string) QueryBuilder
	// Raw get raw generated query
	Raw() string
	// SQL use generated query in part of sql command
	// automatically replace @query with Raw() value
	// automatically replace @where with `WHERE Raw()` value
	SQL(query string) string
	// Args get list of arguments
	Args() []any
}

// NewQuery generate new query builder
func NewQuery() QueryBuilder {
	res := new(qBuilder)
	res.numeric = true
	res.start = 1
	return res
}
