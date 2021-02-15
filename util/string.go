package util

import "regexp"

var (
	reStripNonAlphaNumeric *regexp.Regexp
)

func init() {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	reStripNonAlphaNumeric = reg
}

func StripNonAlphaNumeric(value string) string {
	// credits: https://golangcode.com/how-to-remove-all-non-alphanumerical-characters-from-a-string/
	return reStripNonAlphaNumeric.ReplaceAllString(value, "")
}
