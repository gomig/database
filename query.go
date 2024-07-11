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
	// ToSQL generate query with placeholder based on counter
	ToSQL(counter int) string
	// ToString generate query string with
	// This method replace @q with query to sql
	ToString(pattern string, counter int, params ...any) string
	// Params get list of query parameters
	Params() []any
}

// NewQuery generate new query builder
func NewQuery(driver Driver) QueryBuilder {
	res := new(qBuilder)
	res.driver = driver
	res.queries = make([]Query, 0)
	return res
}
