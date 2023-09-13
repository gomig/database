package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"time"
)

// NullTime nullable string keeper
type NullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (me *NullTime) Scan(value any) error {
	me.Time, me.Valid = time.Time{}, false
	if value != nil {
		temp := sql.NullTime{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		me.Time, me.Valid = temp.Time, temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *NullTime) Value() (driver.Value, error) {
	if !me.Valid {
		return nil, nil
	}
	return me.Time, nil
}

// Val get nullable value
func (me *NullTime) Val() any {
	if !me.Valid {
		return nil
	}
	return me.Time
}

// MarshalJSON convert to json
func (me NullTime) MarshalJSON() ([]byte, error) {
	if me.Valid {
		return json.Marshal(me.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (me *NullTime) UnmarshalJSON(data []byte) error {
	var v *time.Time
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		me.Valid = true
		me.Time = *v
	} else {
		me.Valid = false
	}
	return nil
}
