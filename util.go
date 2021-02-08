package nabarr

import (
	"crypto/sha256"
	"fmt"
	"net/url"
	"path"
	"strconv"
	"strings"
)

func JoinURL(base string, paths ...string) string {
	// credits: https://stackoverflow.com/a/57220413
	p := path.Join(paths...)
	return fmt.Sprintf("%s/%s", strings.TrimRight(base, "/"), strings.TrimLeft(p, "/"))
}

func URLWithQuery(base string, q url.Values) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", fmt.Errorf("url parse: %w", err)
	}

	u.RawQuery = q.Encode()
	return u.String(), nil
}

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
