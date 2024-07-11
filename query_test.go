package database_test

import (
	"testing"

	"github.com/gomig/database/v2"
)

func TestQueryBuilder(t *testing.T) {
	query := database.NewQuery(database.DriverPostgres)
	query.And("firstname LIKE '%?%'", "John").
		And("role @in", "admin", "support", "user").
		OrClosure("age > ? AND age < ?", 15, 30)
	if query.ToSQL(1) != " firstname LIKE '%$1%' AND role IN ($2,$3,$4) OR (age > $5 AND age < $6)" {
		t.Error("ToSQL failed")
	}

	if len(query.Params()) != 6 {
		t.Error("Params resolve failed")
	}

	query = database.NewQuery(database.DriverPostgres).
		And("name = ?", "John Doe").
		And("id = ?", 3)
	sql := query.ToString("SELECT * FROM users WHERE @q ORDER BY %s %s;", 1, "name", "asc")
	if sql != `SELECT * FROM users WHERE  name = $1 AND id = $2 ORDER BY name asc;` {
		t.Log(sql)
		t.Error("ToString failed")
	}

}
