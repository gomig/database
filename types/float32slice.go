package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
)

// Float32Slice a comma separated float32
type Float32Slice struct {
	Floats []float32
}

// Scan implements the Scanner interface.
func (me *Float32Slice) Scan(value any) error {
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		res := strings.Split(strings.TrimSpace(temp.String), ",")
		me.Floats = make([]float32, 0)
		for _, item := range res {
			f, err := strconv.ParseFloat(item, 32)
			if err != nil {
				me.Floats = make([]float32, 0)
				return err
			}
			me.Floats = append(me.Floats, float32(f))
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *Float32Slice) Value() (driver.Value, error) {
	res := make([]string, 0)
	for _, f := range me.Floats {
		res = append(res, strconv.FormatFloat(float64(f), 'f', -1, 32))
	}
	return strings.Join(res, ","), nil
}

// Val get nullable value
func (me *Float32Slice) Val() any {
	res := make([]string, 0)
	for _, f := range me.Floats {
		res = append(res, strconv.FormatFloat(float64(f), 'f', -1, 32))
	}
	if len(res) == 0 {
		return nil
	}
	return strings.Join(res, ",")
}

// MarshalJSON convert to json
func (me Float32Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Floats)
}

// UnmarshalJSON parse from json
func (me *Float32Slice) UnmarshalJSON(data []byte) error {
	var v []float32
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	me.Floats = v
	return nil
}
