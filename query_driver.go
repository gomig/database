package database

import (
	"fmt"
	"strings"
)

// qBuilder query manager
type qBuilder struct {
	queries []Query
	driver  Driver
}

func (q *qBuilder) And(cond string, args ...any) QueryBuilder {
	if cond != "" {
		q.queries = append(q.queries, Query{
			Type:    "AND",
			Query:   cond,
			Params:  args,
			Closure: false,
		})
	}
	return q
}
func (q *qBuilder) Or(cond string, args ...any) QueryBuilder {
	if cond != "" {
		q.queries = append(q.queries, Query{
			Type:    "OR",
			Query:   cond,
			Params:  args,
			Closure: false,
		})
	}
	return q
}
func (q *qBuilder) AndClosure(cond string, args ...any) QueryBuilder {
	if cond != "" {
		q.queries = append(q.queries, Query{
			Type:    "AND",
			Query:   cond,
			Params:  args,
			Closure: true,
		})
	}
	return q
}
func (q *qBuilder) OrClosure(cond string, args ...any) QueryBuilder {
	if cond != "" {
		q.queries = append(q.queries, Query{
			Type:    "OR",
			Query:   cond,
			Params:  args,
			Closure: true,
		})
	}
	return q
}
func (q qBuilder) ToSQL(counter int) string {
	command := ""
	for _, q := range q.queries {
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

		if command == "" {
			command = fmt.Sprintf(" %s", query)
		} else {
			command = fmt.Sprintf("%s %s %s", command, q.Type, query)
		}
	}
	if q.driver == DriverPostgres {
		command = numericArgs(command, counter)
	}
	return command
}

func (q qBuilder) Params() []any {
	args := make([]any, 0)
	for _, q := range q.queries {
		args = append(args, q.Params...)
	}
	return args
}
