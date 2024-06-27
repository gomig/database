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
	And(cond string, args ...any)
	// Or add new simple condition to query with OR
	Or(cond string, args ...any)
	// AndClosure add new condition to query with AND in nested ()
	AndClosure(cond string, args ...any)
	// OrClosure add new condition to query with OR in nested ()
	OrClosure(cond string, args ...any)
	// ToSQL generate query with placeholder based on counter
	ToSQL(counter int) string
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
