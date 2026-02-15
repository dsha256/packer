package safeconv

import (
	"strconv"
)

func ParseInt64(s string) int64 {
	if u, err := strconv.ParseInt(s, 10, 64); err == nil {
		return u
	}

	return 0
}

func ParseInt(s string) int {
	if u, err := strconv.Atoi(s); err == nil {
		return u
	}

	return 0
}
