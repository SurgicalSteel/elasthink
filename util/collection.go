package util

import ()

//SliceStringToInt64 converts a slice of string into a slice of int64
func SliceStringToInt64(ss []string) []int64 {
	result := make([]int64, len(ss))
	for i := 0; i < len(ss); i++ {
		result[i] = StringToInt64(ss[i])
	}
	return result
}
