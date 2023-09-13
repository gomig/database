package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullUInt32 nullable string keeper
type NullUInt32 struct {
	UInt32 uint32
	Valid  bool
}

// Scan implements the Scanner interface.
func (me *NullUInt32) Scan(value any) error {
	me.UInt32, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.UInt32, me.Valid = uint32(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullUInt32) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.UInt32, nil
}

// Val get nullable value
func (me *NullUInt32) Val() any {
	if !me.Valid {
		return nil
	}
	return me.UInt32
}

// MarshalJSON convert to json
func (me NullUInt32) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.UInt32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullUInt32) UnmarshalJSON(data []byte) error {
	var v *uint32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.UInt32 = *v
	} else {
		me.Valid = false
	}
	return nil
}
