package utils

import "strconv"

func ParseInt(num string) int64 {
	val, _ := strconv.ParseInt(num, 10, 64)
	return val
}
