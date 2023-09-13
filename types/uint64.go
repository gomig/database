package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullUInt64 nullable string keeper
type NullUInt64 struct {
	UInt64 uint64
	Valid  bool
}

// Scan implements the Scanner interface.
func (me *NullUInt64) Scan(value any) error {
	me.UInt64, me.Valid = 0, false
	if value != nil {
		temp := sql.NullInt64{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.UInt64, me.Valid = uint64(temp.Int64), temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullUInt64) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.UInt64, nil
}

// Val get nullable value
func (me *NullUInt64) Val() any {
	if !me.Valid {
		return nil
	}
	return me.UInt64
}

// MarshalJSON convert to json
func (me NullUInt64) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.UInt64)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullUInt64) UnmarshalJSON(data []byte) error {
	var v *uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.UInt64 = *v
	} else {
		me.Valid = false
	}
	return nil
}
