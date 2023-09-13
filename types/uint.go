package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullUInt nullable string keeper
type NullUInt struct {
	UInt  uint
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullUInt) Scan(value any) error {
	me.UInt, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.UInt, me.Valid = uint(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullUInt) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.UInt, nil
}

// Val get nullable value
func (me *NullUInt) Val() any {
	if !me.Valid {
		return nil
	}
	return me.UInt
}

// MarshalJSON convert to json
func (me NullUInt) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.UInt)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullUInt) UnmarshalJSON(data []byte) error {
	var v *uint
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.UInt = *v
	} else {
		me.Valid = false
	}
	return nil
}
