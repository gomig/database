package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
)

// NullBool nullable string keeper
type NullBool struct {
	Bool  bool
	Valid bool
}

// Scan implements the Scanner interface.
func (nb *NullBool) Scan(value any) error {
	nb.Bool, nb.Valid = false, false
	if value != nil {
		temp := sql.NullBool{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		nb.Bool, nb.Valid = temp.Bool, temp.Valid
	}
	return nil
}

// Value implements the driver Valuer interface.
func (nb *NullBool) Value() (driver.Value, error) {
	if !nb.Valid {
		return nil, nil
	}
	return nb.Bool, nil
}

// Val get nullable value
func (nb *NullBool) Val() any {
	if !nb.Valid {
		return nil
	}
	return nb.Bool
}

// MarshalJSON convert to json
func (nb NullBool) MarshalJSON() ([]byte, error) {
	if nb.Valid {
		return json.Marshal(nb.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON parse from json
func (nb *NullBool) UnmarshalJSON(data []byte) error {
	var v *bool
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	if v != nil {
		nb.Valid = true
		nb.Bool = *v
	} else {
		nb.Valid = false
	}
	return nil
}
