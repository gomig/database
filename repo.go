package database

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Find get multiple entity (resolve entity from db struct tag)
//
// You can pass resolver to manipulate record after read
// you can use `q` struct for advanced field select query
func Find[T any](db *sqlx.DB, driver Driver, query string, resolver func(*T), args ...any) ([]T, error) {
	if rows, err := db.Queryx(ResolveQuery[T](query, driver), args...); err == sql.ErrNoRows {
		return []T{}, nil
	} else if err != nil {
		return nil, err
	} else {
		res := make([]T, 0)
		for rows.Next() {
			record := new(T)
			if err := rows.StructScan(&record); err != nil {
				return nil, err
			}

			if decoder, ok := any(record).(IDecoder); ok {
				if err := decoder.Decode(); err != nil {
					return nil, err
				}
			}

			if resolver != nil {
				resolver(record)
			}

			if record != nil {
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
func FindOne[T any](db *sqlx.DB, driver Driver, query string, resolver func(*T), args ...any) (*T, error) {
	res := new(T)
	if err := db.Get(res, ResolveQuery[T](query, driver), args...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	} else {
		if decoder, ok := any(res).(IDecoder); ok {
			if err := decoder.Decode(); err != nil {
				return res, err
			}
		}

		if resolver != nil {
			resolver(res)
		}

		return res, nil
	}
}

// Count get count of records
func Count(db *sqlx.DB, driver Driver, query string, args ...any) (int64, error) {
	var res int64
	if driver == DriverPostgres {
		query = numericArgs(query, 1)
	}
	if err := db.Get(&res, query, args...); err != nil {
		return 0, err
	} else {
		return res, nil
	}
}

// Insert struct to database
func Insert(db *sqlx.DB, entity any, table string, driver Driver) (sql.Result, error) {
	cmd, args := ResolveInsert(entity, table, driver)
	return db.Exec(cmd, args...)
}

// Update update struct in database
func Update(db *sqlx.DB, entity any, table string, driver Driver, condition string, args ...any) (sql.Result, error) {
	cmd, args := ResolveUpdate(entity, table, driver, condition, args...)
	return db.Exec(cmd, args...)
}
