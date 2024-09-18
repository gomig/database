package database

import (
	"database/sql"
	"strings"
)

type Commander interface {
	// NumericArgs specifies whether to use numeric ($1, $2) or normal (?, ?) placeholder
	NumericArgs(isNumeric bool) Commander
	// Command set sql comman
	Command(cmd string) Commander
	// Replace replace phrase in query string before ru
	Replace(old string, new string) Commander
	// Exec normalize command and exe
	Exec(args ...any) (sql.Result, error)
}

func NewCMD(db Executable) Commander {
	cmd := new(cmdDriver)
	cmd.db = db
	cmd.numeric = true
	return cmd
}

type cmdDriver struct {
	db           Executable
	numeric      bool
	command      string
	replacements []string
}

func (cmd *cmdDriver) sql() string {
	if cmd.numeric {
		return numericArgs(
			strings.
				NewReplacer(cmd.replacements...).
				Replace(cmd.command),
			1,
		)
	} else {
		return strings.
			NewReplacer(cmd.replacements...).
			Replace(cmd.command)
	}
}

func (cmd *cmdDriver) NumericArgs(numeric bool) Commander {
	cmd.numeric = numeric
	return cmd
}

func (cmd *cmdDriver) Command(sql string) Commander {
	cmd.command = sql
	return cmd
}

func (cmd *cmdDriver) Replace(o string, n string) Commander {
	cmd.replacements = append(cmd.replacements, o, n)
	return cmd
}

func (cmd *cmdDriver) Exec(args ...any) (sql.Result, error) {
	return cmd.db.Exec(cmd.sql(), args...)
}
