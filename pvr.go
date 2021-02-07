package nabarr

type PvrConfig struct {
	Name           string     `yaml:"name"`
	Type           string     `yaml:"type"`
	URL            string     `yaml:"url"`
	ApiKey         string     `yaml:"api_key"`
	QualityProfile string     `yaml:"quality_profile"`
	RootFolder     string     `yaml:"root_folder"`
	Filters        PvrFilters `yaml:"filters"`

	Verbosity string `yaml:"verbosity,omitempty"`
}

type PvrFilters struct {
	Ignores []string
}
