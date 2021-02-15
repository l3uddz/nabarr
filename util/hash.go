package util

import (
	"crypto/sha256"
	"fmt"
)

func AsSHA256(o interface{}) string {
	// credits: https://blog.8bitzen.com/posts/22-08-2019-how-to-hash-a-struct-in-go
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v", o)))
	return fmt.Sprintf("%x", h.Sum(nil))
}
