package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullInt64 nullable string keeper
type NullInt64 struct {
	Int64 int64
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullInt64) Scan(value any) error {
	me.Int64, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Int64, me.Valid = temp.Int64, temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullInt64) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Int64, nil
}

// Val get nullable value
func (me *NullInt64) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Int64
}

// MarshalJSON convert to json
func (me NullInt64) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Int64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullInt64) UnmarshalJSON(data []byte) error {
	var v *int64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Int64 = *v
	} else {
		me.Valid = false
	}
	return nil
}
