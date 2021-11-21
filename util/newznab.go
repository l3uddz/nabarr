package util

import "strconv"

func ContainsTvCategory(cats []string) bool {
	for _, cat := range cats {
		cn, err := strconv.Atoi(cat)
		if err != nil {
			continue
		}

		if cn >= 5000 && cn <= 5999 {
			return true
		}
	}

	return false
}

func ContainsMovieCategory(cats []string) bool {
	for _, cat := range cats {
		cn, err := strconv.Atoi(cat)
		if err != nil {
			continue
		}

		if cn >= 2000 && cn <= 2999 {
			return true
		}
	}

	return false
}
