package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
)

// Float64Slice a comma separated float64
type Float64Slice struct {
	Floats []float64
}

// Scan implements the Scanner interface.
func (me *Float64Slice) Scan(value any) error {
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		res := strings.Split(strings.TrimSpace(temp.String), ",")
		me.Floats = make([]float64, 0)
		for _, item := range res {
			f, err := strconv.ParseFloat(item, 64)
			if err != nil {
				me.Floats = make([]float64, 0)
				return err
			}
			me.Floats = append(me.Floats, f)
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *Float64Slice) Value() (driver.Value, error) {
	res := make([]string, 0)
	for _, f := range me.Floats {
		res = append(res, strconv.FormatFloat(f, 'f', -1, 64))
	}
	return strings.Join(res, ","), nil
}

// Val get nullable value
func (me *Float64Slice) Val() any {
	res := make([]string, 0)
	for _, f := range me.Floats {
		res = append(res, strconv.FormatFloat(f, 'f', -1, 64))
	}
	if len(res) == 0 {
		return nil
	}
	return strings.Join(res, ",")
}

// MarshalJSON convert to json
func (me Float64Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Floats)
}

// UnmarshalJSON parse from json
func (me *Float64Slice) UnmarshalJSON(data []byte) error {
	var v []float64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	me.Floats = v
	return nil
}
