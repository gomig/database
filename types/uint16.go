package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullUInt16 nullable string keeper
type NullUInt16 struct {
	UInt16 uint16
	Valid  bool
}

// Scan implements the Scanner interface.
func (me *NullUInt16) Scan(value any) error {
	me.UInt16, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.UInt16, me.Valid = uint16(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullUInt16) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.UInt16, nil
}

// Val get nullable value
func (me *NullUInt16) Val() any {
	if !me.Valid {
		return nil
	}
	return me.UInt16
}

// MarshalJSON convert to json
func (me NullUInt16) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.UInt16)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullUInt16) UnmarshalJSON(data []byte) error {
	var v *uint16
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.UInt16 = *v
	} else {
		me.Valid = false
	}
	return nil
}
