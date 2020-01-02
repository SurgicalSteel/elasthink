package util

import (
	"strconv"
)

//StringToInt64 convert a string into int64
func StringToInt64(s string) int64 {
	result, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return int64(0)
	}
	return result
}
