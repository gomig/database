package migration_test

import (
	"testing"

	"github.com/gomig/database/v2/migration"
)

func TestFormatter(t *testing.T) {
	migration.Formatter("{m}{I}Executed Script{R}: {b}{B}(%d){R}\n", 10)
	migration.Formatter("{r}FAIL!{R} invalid migration file name\n")
	migration.Formatter("{g}OK!{R} invalid migration file name\n")
}
