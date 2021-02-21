package nabarr

import (
	"time"
)

/* pvr config / filters */

type PvrConfig struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"`
	URL            string `yaml:"url"`
	ApiKey         string `yaml:"api_key"`
	QualityProfile string `yaml:"quality_profile"`
	RootFolder     string `yaml:"root_folder"`
	Options        struct {
		// add options
		AddMonitored  *bool `yaml:"add_monitored"`
		SearchMissing *bool `yaml:"search_missing"`
		// skip objects
		SkipAnime *bool `yaml:"skip_anime"`
	} `yaml:"options"`
	Filters       PvrFilters    `yaml:"filters"`
	CacheDuration time.Duration `yaml:"cache_duration"`
	Verbosity     string        `yaml:"verbosity,omitempty"`
}

type PvrFilters struct {
	Ignores []string
}

/* pvr options */

type PvrOption func(options *PvrOptions)

type PvrOptions struct {
	// seriesType returned from the lookup before adding (sonarr)
	SeriesType string

	AddMonitored  bool
	SearchMissing bool
}

func BuildPvrOptions(opts ...PvrOption) (*PvrOptions, error) {
	os := &PvrOptions{}

	for _, opt := range opts {
		opt(os)
	}

	return os, nil
}

func WithSeriesType(seriesType string) PvrOption {
	return func(opts *PvrOptions) {
		opts.SeriesType = seriesType
	}
}

func WithAddMonitored(monitored bool) PvrOption {
	return func(opts *PvrOptions) {
		opts.AddMonitored = monitored
	}
}

func WithSearchMissing(missing bool) PvrOption {
	return func(opts *PvrOptions) {
		opts.SearchMissing = missing
	}
}
