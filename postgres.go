package database

import (
	"strings"

	"github.com/gomig/utils"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// NewPostgresConnector create new POSTGRES connection
func NewPostgresConnector(host string, port string, user string, password string, database string) (*sqlx.DB, error) {
	params := []string{}
	params = append(params, "host="+host)
	params = append(params, "user="+user)
	params = append(params, "sslmode=disable")
	if port != "" {
		params = append(params, "port="+port)
	}
	if password != "" {
		params = append(params, "password="+password)
	}
	if database != "" {
		params = append(params, "dbname="+database)
	}

	db, err := sqlx.Open("postgres", strings.Join(params, " "))
	if err != nil {
		return nil, utils.TaggedError([]string{"PostgresDriver"}, err.Error())
	}
	if err := db.Ping(); err != nil {
		return nil, utils.TaggedError([]string{"PostgresDriver"}, err.Error())
	}
	return db, nil
}
