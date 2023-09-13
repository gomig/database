package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strings"
)

// StringSlice a comma separated string
type StringSlice struct {
	Strings []string
}

// Scan implements the Scanner interface.
func (me *StringSlice) Scan(value any) error {
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		res := strings.Split(strings.TrimSpace(temp.String), ",")
		me.Strings = make([]string, 0)
		for _, item := range res {
			me.Strings = append(me.Strings, item)
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *StringSlice) Value() (driver.Value, error) {
	return strings.Join(me.Strings, ","), nil
}

// Val get nullable value
func (me *StringSlice) Val() any {
	if len(me.Strings) == 0 {
		return nil
	}
	return strings.Join(me.Strings, ",")
}

// MarshalJSON convert to json
func (me StringSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Strings)
}

// UnmarshalJSON parse from json
func (me *StringSlice) UnmarshalJSON(data []byte) error {
	var v []string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	me.Strings = v
	return nil
}
