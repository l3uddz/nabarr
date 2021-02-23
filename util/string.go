package util

import "regexp"

var (
	reStripNonAlphaNumeric *regexp.Regexp
	reStripNonNumeric      *regexp.Regexp
)

func init() {
	reg1, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		panic(err)
	}
	reStripNonAlphaNumeric = reg1

	reg2, err := regexp.Compile("[^0-9]+")
	if err != nil {
		panic(err)
	}
	reStripNonNumeric = reg2
}

func StripNonAlphaNumeric(value string) string {
	// credits: https://golangcode.com/how-to-remove-all-non-alphanumerical-characters-from-a-string/
	return reStripNonAlphaNumeric.ReplaceAllString(value, "")
}

func StripNonNumeric(value string) string {
	return reStripNonNumeric.ReplaceAllString(value, "")
}
