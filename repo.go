package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Find get multiple entity (resolve entity from db struct tag)
//
// You can pass resolver to manipulate record after read
// you can use `q` struct for advanced field select query
func FindOpt[T any](db *sqlx.DB, query string, option Option[T], args ...any) ([]T, error) {
	if rows, err := db.Queryx(option.resolveQuery(query), args...); err == sql.ErrNoRows {
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
func Find[T any](db *sqlx.DB, query string, args ...any) ([]T, error) {
	return FindOpt(db, query, NewOption[T](), args...)
}

// FindOne get single entity
//
// You can pass resolver to manipulate record after read
// you can use `q` or `db` struct tag to map field to database column
func FindOneOpt[T any](db *sqlx.DB, query string, option Option[T], args ...any) (*T, error) {
	// handle options
	record := new(T)
	if err := db.Get(record, option.resolveQuery(query), args...); err == sql.ErrNoRows {
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
func FindOne[T any](db *sqlx.DB, query string, args ...any) (*T, error) {
	return FindOneOpt(db, query, NewOption[T](), args...)
}

// Count get count of records
func CountOpt(db *sqlx.DB, query string, option Option[int64], args ...any) (int64, error) {
	var count int64
	if err := db.Get(&count, option.resolve(query), args...); err != nil {
		return 0, err
	} else {
		return count, nil
	}
}
func Count(db *sqlx.DB, query string, args ...any) (int64, error) {
	return CountOpt(db, query, NewOption[int64](), args...)
}

// Insert struct to database
func InsertOpt[T any](db Executable, entity T, table string, option Option[T]) (sql.Result, error) {
	cmd, args := ResolveInsert(entity, table, option.getDriver())
	return db.Exec(cmd, args...)
}
func Insert[T any](db Executable, entity T, table string) (sql.Result, error) {
	return InsertOpt(db, entity, table, NewOption[T]())
}

// Update update struct in database
func UpdateOpt[T any](db Executable, entity T, table string, condition string, option Option[T], args ...any) (sql.Result, error) {
	cmd, args := ResolveUpdate(entity, table, option.getDriver(), condition, args...)
	return db.Exec(cmd, args...)
}
func Update[T any](db Executable, entity T, table string, condition string, args ...any) (sql.Result, error) {
	return UpdateOpt(db, entity, table, condition, NewOption[T](), args...)
}
