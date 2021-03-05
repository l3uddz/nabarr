package util

import "strings"

func StringSliceContains(slice []string, val string) bool {
	for _, s := range slice {
		if strings.EqualFold(s, val) {
			return true
		}
	}
	return false
}

func StringSliceContainsAny(slice []string, vals []string) bool {
	for _, s := range slice {
		for _, v := range vals {
			if strings.EqualFold(s, v) {
				return true
			}
		}
	}
	return false
}

func StringSliceMergeUnique(existingSlice []string, mergeSlice []string) []string {
	// add existing
	data := make([]string, 0)
	for _, es := range existingSlice {
		if es == "" {
			continue
		}
		data = append(data, es)
	}

	// add merge items (unique)
	for _, ms := range mergeSlice {
		if ms == "" {
			continue
		}

		merge := true
		for _, es := range data {
			if strings.EqualFold(es, ms) {
				merge = false
				break
			}
		}

		if merge {
			data = append(data, ms)
		}
	}
	return data
}
