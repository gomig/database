package database

import (
	"fmt"
	"strings"
)

// QueryBuilder query manager
type QueryBuilder struct {
	queries []Query
}

// Add add new query
func (qb *QueryBuilder) Add(q Query) {
	if q.Query != "" {
		qb.queries = append(qb.queries, q)
	}
}

// Query get query string
func (qb *QueryBuilder) Query() (res string) {
	res = ""
	for _, q := range qb.queries {
		query := q.Query
		// Compile In Params
		if strings.Contains(query, "@in") {
			params := "IN (?"
			params = params + strings.Repeat(",?", len(q.Params)-1)
			params = params + ")"
			query = strings.Replace(query, "@in", params, 1)
		}
		// Generate subquery
		if q.Closure {
			query = "(" + query + ")"
		}

		if res == "" {
			res = fmt.Sprintf(" %s", query)
		} else {
			res = fmt.Sprintf("%s %s %s", res, q.Type, query)
		}
	}
	return
}

// Params get query parameters
func (qb *QueryBuilder) Params() (vars []any) {
	for _, q := range qb.queries {
		vars = append(vars, q.Params...)
	}
	return
}
