package migration

import (
	"strconv"
	"strings"
)

type byNumber []string

func getCode(str string) int {
	res := 0
	parts := strings.Split(str, "-")
	if len(parts) > 1 {
		if res, err := strconv.Atoi(parts[0]); err == nil {
			return res
		}
	}
	return res
}

func (s byNumber) Len() int {
	return len(s)
}
func (s byNumber) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byNumber) Less(i, j int) bool {
	return getCode(s[i]) < getCode(s[j])
}
