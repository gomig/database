package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullInt16 nullable string keeper
type NullInt16 struct {
	Int16 int16
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullInt16) Scan(value any) error {
	me.Int16, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Int16, me.Valid = int16(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullInt16) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Int16, nil
}

// Val get nullable value
func (me *NullInt16) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Int16
}

// MarshalJSON convert to json
func (me NullInt16) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Int16)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullInt16) UnmarshalJSON(data []byte) error {
	var v *int16
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Int16 = *v
	} else {
		me.Valid = false
	}
	return nil
}
