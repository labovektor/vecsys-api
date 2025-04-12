package util

import "strconv"

func StringToInteger(source string, fallback ...int) int {
	val, err := strconv.Atoi(source)
	if err != nil {
		if len(fallback) > 0 {
			return fallback[0]
		}
		return 0
	}
	return val
}
