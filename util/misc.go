package util

import (
	"crypto/sha256"
	"fmt"
	"strconv"
)

func Atoi(val string, defaultVal int) int {
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func AsSHA256(o interface{}) string {
	// credits: https://blog.8bitzen.com/posts/22-08-2019-how-to-hash-a-struct-in-go
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))

	return fmt.Sprintf("%x", h.Sum(nil))
}
