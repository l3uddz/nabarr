package nabarr

import (
	"time"
)

type PvrConfig struct {
	Name           string        `yaml:"name"`
	Type           string        `yaml:"type"`
	URL            string        `yaml:"url"`
	ApiKey         string        `yaml:"api_key"`
	QualityProfile string        `yaml:"quality_profile"`
	RootFolder     string        `yaml:"root_folder"`
	Filters        PvrFilters    `yaml:"filters"`
	CacheDuration  time.Duration `yaml:"cache_duration"`
	Verbosity      string        `yaml:"verbosity,omitempty"`
}

type PvrFilters struct {
	Ignores []string
}

type PvrOption func(options *PvrOptions)

type PvrOptions struct {
	// the seriesType returned from the lookup before adding (sonarr)
	LookupType string
}
