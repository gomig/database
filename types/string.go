package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullString nullable string keeper
type NullString struct {
	String string
	Valid  bool
}

// Scan implements the Scanner interface.
func (me *NullString) Scan(value any) error {
	me.String, me.Valid = "", false
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.String, me.Valid = temp.String, temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullString) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.String, nil
}

// Val get nullable value
func (me *NullString) Val() any {
	if !me.Valid {
		return nil
	}
	return me.String
}

// MarshalJSON convert to json
func (me NullString) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullString) UnmarshalJSON(data []byte) error {
	var v *string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.String = *v
	} else {
		me.Valid = false
	}
	return nil
}
