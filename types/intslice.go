package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"strconv"
	"strings"
)

// IntSlice a comma separated int
type IntSlice struct {
	Ints []int
}

// Scan implements the Scanner interface.
func (me *IntSlice) Scan(value any) error {
	if value != nil {
		temp := sql.NullString{}
		err := temp.Scan(value)
		if err != nil {
			return err
		}
		res := strings.Split(strings.TrimSpace(temp.String), ",")
		me.Ints = make([]int, 0)
		for _, item := range res {
			n, err := strconv.Atoi(item)
			if err != nil {
				me.Ints = make([]int, 0)
				return err
			}
			me.Ints = append(me.Ints, n)
		}
	}
	return nil
}

// Value implements the driver Valuer interface.
func (me *IntSlice) Value() (driver.Value, error) {
	res := make([]string, 0)
	for _, n := range me.Ints {
		res = append(res, strconv.Itoa(n))
	}
	return strings.Join(res, ","), nil
}

// Val get nullable value
func (me *IntSlice) Val() any {
	res := make([]string, 0)
	for _, n := range me.Ints {
		res = append(res, strconv.Itoa(n))
	}
	if len(res) == 0 {
		return nil
	}
	return strings.Join(res, ",")
}

// MarshalJSON convert to json
func (me IntSlice) MarshalJSON() ([]byte, error) {
	return json.Marshal(me.Ints)
}

// UnmarshalJSON parse from json
func (me *IntSlice) UnmarshalJSON(data []byte) error {
	var v []int
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	me.Ints = v
	return nil
}
