package database

import (
	"database/sql"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Finder[T any] interface {
	// NumericArgs specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder
	NumericArgs(isNumeric bool) Finder[T]
	// Query set sql query
	Query(query string) Finder[T]
	// Replace replace phrase in query string before run
	Replace(old string, new string) Finder[T]
	// Resolve reginster new resolver to run on record after read
	Resolve(resolver func(*T) error) Finder[T]
	// Single get first result
	Single(args ...any) (*T, error)
	// Result get multiple result
	Result(args ...any) ([]T, error)
}

func NewFinder[T any](db *sqlx.DB) Finder[T] {
	finder := new(finderDriver[T])
	finder.db = db
	finder.numeric = true
	return finder
}

type finderDriver[T any] struct {
	db           *sqlx.DB
	numeric      bool
	query        string
	replacements []string
	resolvers    []func(*T) error
}

func (finder *finderDriver[T]) sql() string {
	if strings.Contains(finder.query, "@fields") {
		var sample T
		finder.replacements = append(
			finder.replacements,
			"@fields",
			strings.Join(structQueryColumns(sample), " ,"),
		)
	}

	if finder.numeric {
		return numericArgs(
			strings.
				NewReplacer(finder.replacements...).
				Replace(finder.query),
			1,
		)
	} else {
		return strings.
			NewReplacer(finder.replacements...).
			Replace(finder.query)
	}
}

func (finder *finderDriver[T]) NumericArgs(numeric bool) Finder[T] {
	finder.numeric = numeric
	return finder
}

func (finder *finderDriver[T]) Query(query string) Finder[T] {
	finder.query = query
	return finder
}

func (finder *finderDriver[T]) Replace(old, new string) Finder[T] {
	finder.replacements = append(finder.replacements, old, new)
	return finder
}

func (finder *finderDriver[T]) Resolve(resolver func(*T) error) Finder[T] {
	finder.resolvers = append(finder.resolvers, resolver)
	return finder
}

func (finder *finderDriver[T]) Single(args ...any) (*T, error) {
	if cursor, err := finder.db.Queryx(finder.sql(), args...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		defer cursor.Close()
		for cursor.Next() {
			record := new(T)
			if err := cursor.StructScan(record); err != nil {
				return nil, err
			} else {
				if decoder, ok := any(record).(IDecoder); ok {
					if err := decoder.Decode(); err != nil {
						return nil, err
					}
				}

				for _, resolver := range finder.resolvers {
					if err := resolver(record); err != nil {
						return nil, err
					}
				}

				return record, nil
			}
		}
		return nil, nil
	}
}

func (finder *finderDriver[T]) Result(args ...any) ([]T, error) {
	if cursor, err := finder.db.Queryx(finder.sql(), args...); err == sql.ErrNoRows {
		return []T{}, nil
	} else if err != nil {
		return nil, err
	} else {
		results := make([]T, 0)
		defer cursor.Close()
		for cursor.Next() {
			record := new(T)
			if err := cursor.StructScan(record); err != nil {
				return nil, err
			} else {
				if decoder, ok := any(record).(IDecoder); ok {
					if err := decoder.Decode(); err != nil {
						return nil, err
					}
				}

				for _, resolver := range finder.resolvers {
					if err := resolver(record); err != nil {
						return nil, err
					}
				}

				results = append(results, *record)
			}
		}
		return results, nil
	}
}
