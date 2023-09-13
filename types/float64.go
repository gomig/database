package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullFloat64 nullable string keeper
type NullFloat64 struct {
	Float64 float64
	Valid   bool
}

// Scan implements the Scanner interface.
func (me *NullFloat64) Scan(value any) error {
	me.Float64, me.Valid = 0.0, false
	if value != nil {
		temp := sql.NullFloat64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Float64, me.Valid = temp.Float64, temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullFloat64) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Float64, nil
}

// Val get nullable value
func (me *NullFloat64) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Float64
}

// MarshalJSON convert to json
func (me NullFloat64) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Float64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullFloat64) UnmarshalJSON(data []byte) error {
	var v *float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Float64 = *v
	} else {
		me.Valid = false
	}
	return nil
}
