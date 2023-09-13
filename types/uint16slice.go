package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
)

// UInt16Slice a comma separated uint16
type UInt16Slice struct {
	Ints []uint16
}

// Scan implements the Scanner interface.
func (me *UInt16Slice) Scan(value any) error {
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		res := strings.Split(strings.TrimSpace(temp.String), ",")
		me.Ints = make([]uint16, 0)
		for _, item := range res {
			n, err := strconv.Atoi(item)
			if err != nil {
				me.Ints = make([]uint16, 0)
				return err
			}
			me.Ints = append(me.Ints, uint16(n))
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *UInt16Slice) Value() (driver.Value, error) {
	res := make([]string, 0)
	for _, n := range me.Ints {
		res = append(res, strconv.Itoa(int(n)))
	}
	return strings.Join(res, ","), nil
}

// Val get nullable value
func (me *UInt16Slice) Val() any {
	res := make([]string, 0)
	for _, n := range me.Ints {
		res = append(res, strconv.Itoa(int(n)))
	}
	if len(res) == 0 {
		return nil
	}
	return strings.Join(res, ",")
}

// MarshalJSON convert to json
func (me UInt16Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Ints)
}

// UnmarshalJSON parse from json
func (me *UInt16Slice) UnmarshalJSON(data []byte) error {
	var v []uint16
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	me.Ints = v
	return nil
}
