package sonarr

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

func seriesType(seriesType string) nabarr.PvrOption {
	return func(opts *nabarr.PvrOptions) {
		opts.LookupType = seriesType
	}
}
