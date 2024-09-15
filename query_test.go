package database_test

import (
	"testing"

	"github.com/gomig/database/v2"
)

func TestQueryBuilder(t *testing.T) {
	query := database.NewQuery()
	query.And("firstname LIKE '%?%'", "John").
		And("role @in", "admin", "support", "user").
		OrClosure("age > ? AND age < ?", 15, 30)
	if raw := query.RawPostgres(1); raw != `firstname LIKE '%$1%' AND role IN ($2,$3,$4) OR (age > $5 AND age < $6)` {
		t.Log(raw)
		t.Error("RawPostgres failed")
	}

	if len(query.Params()) != 6 {
		t.Error("Params resolve failed")
	}

	query = database.NewQuery().
		And("name = ?", "John Doe").
		And("id = ?", 3)
	if sql := query.ToPostgres(
		`SELECT * FROM users @where ORDER BY @sort @order;`,
		1, "@sort", "name", "@order", "asc",
	); sql != `SELECT * FROM users WHERE name = $1 AND id = $2 ORDER BY name asc;` {
		t.Log(sql)
		t.Error("ToPostgres failed")
	}

}
