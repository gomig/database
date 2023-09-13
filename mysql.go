package database

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gomig/utils"
	"github.com/jmoiron/sqlx"
)

// NewMySQLConnector create new mysql connection
func NewMySQLConnector(host string, username string, password string, database string) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@%s/%s?charset=utf8mb4&parseTime=true", username, password, host, database))
	if err != nil {
		return nil, utils.TaggedError([]string{"MySQLDriver"}, err.Error())
	}
	if err := db.Ping(); err != nil {
		return nil, utils.TaggedError([]string{"MySQLDriver"}, err.Error())
	}
	return db, nil
}
