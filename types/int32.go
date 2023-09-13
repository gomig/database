package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullInt32 nullable string keeper
type NullInt32 struct {
	Int32 int32
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullInt32) Scan(value any) error {
	me.Int32, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Int32, me.Valid = int32(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullInt32) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Int32, nil
}

// Val get nullable value
func (me *NullInt32) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Int32
}

// MarshalJSON convert to json
func (me NullInt32) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullInt32) UnmarshalJSON(data []byte) error {
	var v *int32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Int32 = *v
	} else {
		me.Valid = false
	}
	return nil
}
