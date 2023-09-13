package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
)

// UInt64Slice a comma separated uint64
type UInt64Slice struct {
	Ints []uint64
}

// Scan implements the Scanner interface.
func (me *UInt64Slice) Scan(value any) error {
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		res := strings.Split(strings.TrimSpace(temp.String), ",")
		me.Ints = make([]uint64, 0)
		for _, item := range res {
			n, err := strconv.ParseInt(item, 10, 64)
			if err != nil {
				me.Ints = make([]uint64, 0)
				return err
			}
			me.Ints = append(me.Ints, uint64(n))
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *UInt64Slice) Value() (driver.Value, error) {
	res := make([]string, 0)
	for _, n := range me.Ints {
		res = append(res, strconv.Itoa(int(n)))
	}
	return strings.Join(res, ","), nil
}

// Val get nullable value
func (me *UInt64Slice) Val() any {
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
func (me UInt64Slice) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Ints)
}

// UnmarshalJSON parse from json
func (me *UInt64Slice) UnmarshalJSON(data []byte) error {
	var v []uint64
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	me.Ints = v
	return nil
}
