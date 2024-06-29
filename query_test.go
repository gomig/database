package database_test

import (
	"testing"

	"github.com/gomig/database/v2"
)

func TestQueryBuilder(t *testing.T) {
	query := database.NewQuery(database.DriverPostgres)
	query.And("firstname LIKE '%?%'", "John")
	query.And("role @in", "admin", "support", "user")
	query.OrClosure("age > ? AND age < ?", 15, 30)
	if query.ToSQL(1) != " firstname LIKE '%$1%' AND role IN ($2,$3,$4) OR (age > $5 AND age < $6)" {
		t.Error("ToSQL failed")
	}

	if len(query.Params()) != 6 {
		t.Error("Params resolve failed")
	}
}
