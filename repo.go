package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Find get multiple entity (resolve entity from db struct tag)
//
// You can pass resolver to manipulate record after read
// you can use `q` struct for advanced field select query
func Find[T any](db *sqlx.DB, query string, options ...Option[T]) ([]T, error) {
	option := resolveOptions(options...)
	if rows, err := db.Queryx(option.resolveQuery(query), option.getArgs()...); err == sql.ErrNoRows {
		return []T{}, nil
	} else if err != nil {
		return nil, err
	} else {
		res := make([]T, 0)
		for rows.Next() {
			record := new(T)
			if err := rows.StructScan(record); err != nil {
				return nil, err
			} else {
				if decoder, ok := any(record).(IDecoder); ok {
					if err := decoder.Decode(); err != nil {
						return nil, err
					}
				}

				for _, resolver := range option.getResolvers() {
					if err := resolver(record); err != nil {
						return nil, err
					}
				}

				res = append(res, *record)
			}
		}
		return res, nil
	}
}

// FindOne get single entity
//
// You can pass resolver to manipulate record after read
// you can use `q` or `db` struct tag to map field to database column
func FindOne[T any](db *sqlx.DB, query string, options ...Option[T]) (*T, error) {
	option := resolveOptions(options...)
	// handle options
	record := new(T)
	if err := db.Get(record, option.resolveQuery(query), option.getArgs()...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		if decoder, ok := any(record).(IDecoder); ok {
			if err := decoder.Decode(); err != nil {
				return record, err
			}
		}

		for _, resolver := range option.getResolvers() {
			if err := resolver(record); err != nil {
				return nil, err
			}
		}

		return record, nil
	}
}

// Count get count of records
func Count[T any](db *sqlx.DB, query string, options ...Option[T]) (int64, error) {
	option := resolveOptions(options...)
	// handle options
	var count int64
	if err := db.Get(&count, option.resolve(query), option.getArgs()...); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}

// Insert struct to database
func Insert[T any](db Executable, entity T, table string, options ...Option[T]) (sql.Result, error) {
	option := resolveOptions(options...)
	cmd, args := ResolveInsert(entity, table, option.getDriver())
	return db.Exec(cmd, args...)
}

// Update update struct in database
func Update[T any](db Executable, entity T, table string, condition string, options ...Option[T]) (sql.Result, error) {
	option := resolveOptions(options...)
	cmd, args := ResolveUpdate(entity, table, option.getDriver(), condition, option.getArgs()...)
	return db.Exec(cmd, args...)
}
