package radarr

import (
	"github.com/l3uddz/nabarr"
)

func buildOptions(opts ...nabarr.PvrOption) (*nabarr.PvrOptions, error) {
	os := &nabarr.PvrOptions{}

	for _, opt := range opts {
		opt(os)
	}

	return os, nil
}
