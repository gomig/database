package database_test

import (
	"testing"

	"github.com/gomig/database/v2"
)

func TestQueryBuilder(t *testing.T) {
	rawExp := `firstname LIKE '%$5%' AND role IN ($6, $7, $8) OR (age > $9 AND age < $10)`
	rawQ := database.NewQuery().
		And("firstname LIKE '%?%'", "John").
		And("role @in", "admin", "support", "user").
		OrClosure("age > ? AND age < ?", 15, 30).
		NumericStart(5)

	if raw := rawQ.Raw(); raw != rawExp {
		t.Logf("Expected: %s\nReturns: %s\n", rawExp, raw)
		t.Error("Raw() failed")
	}

	if len(rawQ.Args()) != 6 {
		t.Error("Args() failed")
	}

	sqlExp := `SELECT * FROM users WHERE name = ? AND id = ? ORDER BY name asc;`
	sqlQ := database.NewQuery().
		And("name = ?", "John Doe").
		And("id = ?", 3).
		NumericArgs(false).
		Replace("@sort", "name").
		Replace("@order", "asc")
	if sql := sqlQ.SQL(`SELECT * FROM users @where ORDER BY @sort @order;`); sql != sqlExp {
		t.Logf("Expected: %s\nReturns: %s\n", sqlExp, sql)
		t.Error("SQL() failed")
	}
}
