package util

import (
	"strconv"
	"strings"
)

func Atoi(val string, defaultVal int) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func Atof64(val string, defaultVal float64) float64 {
	n, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
	if err != nil {
		return defaultVal
	}
	return n
}

func StringOrDefault(val string, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}

func BoolOrDefault(val *bool, defaultVal bool) bool {
	if val == nil {
		return defaultVal
	}

	return *val
}
