package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullFloat32 nullable string keeper
type NullFloat32 struct {
	Float32 float32
	Valid   bool
}

// Scan implements the Scanner interface.
func (me *NullFloat32) Scan(value any) error {
	me.Float32, me.Valid = 0.0, false
	if value != nil {
		temp := sql.NullFloat64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Float32, me.Valid = float32(temp.Float64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullFloat32) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Float32, nil
}

// Val get nullable value
func (me *NullFloat32) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Float32
}

// MarshalJSON convert to json
func (me NullFloat32) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Float32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullFloat32) UnmarshalJSON(data []byte) error {
	var v *float32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Float32 = *v
	} else {
		me.Valid = false
	}
	return nil
}
