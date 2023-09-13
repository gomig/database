package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullUInt8 nullable string keeper
type NullUInt8 struct {
	UInt8 uint8
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullUInt8) Scan(value any) error {
	me.UInt8, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.UInt8, me.Valid = uint8(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullUInt8) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.UInt8, nil
}

// Val get nullable value
func (me *NullUInt8) Val() any {
	if !me.Valid {
		return nil
	}
	return me.UInt8
}

// MarshalJSON convert to json
func (me NullUInt8) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.UInt8)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullUInt8) UnmarshalJSON(data []byte) error {
	var v *uint8
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.UInt8 = *v
	} else {
		me.Valid = false
	}
	return nil
}
