package database

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

type Counter interface {
	// NumericArgs specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder
	NumericArgs(isNumeric bool) Counter
	// Query set sql query
	Query(query string) Counter
	// Replace replace phrase in query string before run
	Replace(old string, new string) Counter
	// Result get count, returns -1 on error
	Result(args ...any) (int64, error)
}

func NewCounter(db *sqlx.DB) Counter {
	counter := new(counterDriver)
	counter.db = db
	counter.numeric = true
	return counter
}

type counterDriver struct {
	db           *sqlx.DB
	numeric      bool
	query        string
	replacements []string
}

func (counter *counterDriver) sql() string {
	if counter.numeric {
		return numericArgs(
			strings.
				NewReplacer(counter.replacements...).
				Replace(counter.query),
			1,
		)
	} else {
		return strings.
			NewReplacer(counter.replacements...).
			Replace(counter.query)
	}
}

func (counter *counterDriver) NumericArgs(numeric bool) Counter {
	counter.numeric = numeric
	return counter
}

func (counter *counterDriver) Query(query string) Counter {
	counter.query = query
	return counter
}

func (counter *counterDriver) Replace(old, new string) Counter {
	counter.replacements = append(counter.replacements, old, new)
	return counter
}

func (counter *counterDriver) Result(args ...any) (int64, error) {
	var count int64
	if err := counter.db.Get(&count, counter.sql(), args...); err != nil {
		return -1, err
	} else {
		return count, nil
	}
}
