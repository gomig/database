package database

import (
	"math"
	"strings"
)

type qItem struct {
	Type    string
	Query   string
	Args    []any
	Closure bool
}

type qBuilder struct {
	numeric      bool
	start        int
	queries      []qItem
	replacements []string
}

func (builder *qBuilder) addItem(q string, and, closure bool, args ...any) {
	item := qItem{}
	if and {
		item.Type = "AND"
	} else {
		item.Type = "OR"
	}
	item.Query = q
	item.Closure = closure
	item.Args = args
	builder.queries = append(builder.queries, item)
}

func (builder *qBuilder) And(query string, args ...any) QueryBuilder {
	if query != "" {
		builder.addItem(query, true, false, args...)
	}
	return builder
}

func (builder *qBuilder) AndIf(cond bool, query string, args ...any) QueryBuilder {
	if cond && query != "" {
		builder.addItem(query, true, false, args...)
	}
	return builder
}

func (builder *qBuilder) Or(query string, args ...any) QueryBuilder {
	if query != "" {
		builder.addItem(query, false, false, args...)
	}
	return builder
}

func (builder *qBuilder) OrIf(cond bool, query string, args ...any) QueryBuilder {
	if cond && query != "" {
		builder.addItem(query, false, false, args...)
	}
	return builder
}

func (builder *qBuilder) AndClosure(query string, args ...any) QueryBuilder {
	if query != "" {
		builder.addItem(query, true, true, args...)
	}
	return builder
}

func (builder *qBuilder) AndClosureIf(cond bool, query string, args ...any) QueryBuilder {
	if cond && query != "" {
		builder.addItem(query, true, true, args...)
	}
	return builder
}

func (builder *qBuilder) OrClosure(query string, args ...any) QueryBuilder {
	if query != "" {
		builder.addItem(query, false, true, args...)
	}
	return builder
}

func (builder *qBuilder) OrClosureIf(cond bool, query string, args ...any) QueryBuilder {
	if cond && query != "" {
		builder.addItem(query, false, true, args...)
	}
	return builder
}

func (builder *qBuilder) NumericArgs(numeric bool) QueryBuilder {
	builder.numeric = numeric
	return builder
}

func (builder *qBuilder) NumericStart(start int) QueryBuilder {
	builder.start = start
	return builder
}

func (builder *qBuilder) Replace(old, new string) QueryBuilder {
	builder.replacements = append(builder.replacements, old, new)
	return builder
}

func (builder *qBuilder) Raw() string {
	command := ""
	for _, q := range builder.queries {
		query := q.Query

		// generate @in
		if strings.Contains(query, "@in") {
			placeholders := strings.TrimLeft(
				strings.Repeat(", ?", len(q.Args)),
				", ",
			)
			query = strings.Replace(query, "@in", "IN ("+placeholders+")", 1)
		}

		// generate subquery
		if q.Closure {
			query = "(" + query + ")"
		}

		if command == "" {
			command = query
		} else {
			command = command + " " + q.Type + " " + query
		}
	}

	if builder.numeric {
		command = numericArgs(command, int(math.Max(float64(builder.start), 1)))
	}

	return command
}

func (builder *qBuilder) SQL(query string) string {
	if raw := builder.Raw(); raw == "" {
		return strings.NewReplacer(
			append(
				builder.replacements,
				"@query", raw,
				"@where", "",
			)...,
		).Replace(query)
	} else {
		return strings.NewReplacer(
			append(
				builder.replacements,
				"@query", raw,
				"@where", "WHERE "+raw,
			)...,
		).Replace(query)
	}
}

func (builder *qBuilder) Args() []any {
	args := make([]any, 0)
	for _, q := range builder.queries {
		args = append(args, q.Args...)
	}
	return args
}
