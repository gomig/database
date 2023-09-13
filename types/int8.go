package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullInt8 nullable string keeper
type NullInt8 struct {
	Int8  int8
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullInt8) Scan(value any) error {
	me.Int8, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Int8, me.Valid = int8(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullInt8) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Int8, nil
}

// Val get nullable value
func (me *NullInt8) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Int8
}

// MarshalJSON convert to json
func (me NullInt8) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Int8)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullInt8) UnmarshalJSON(data []byte) error {
	var v *int8
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Int8 = *v
	} else {
		me.Valid = false
	}
	return nil
}
