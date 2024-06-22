package migration

import (
	"strconv"

	"github.com/gomig/utils"
)

type byNumber []string

func getCode(str string) int {
	if res, err := strconv.Atoi(utils.ExtractNumbers(str)); err == nil {
		return res
	} else {
		return 0
	}
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
