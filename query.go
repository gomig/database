package database

// Query object
type Query struct {
	Type    string
	Query   string
	Params  []any
	Closure bool
}
