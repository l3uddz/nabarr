package omdb

type Config struct {
	ApiKey string `yaml:"api_key"`

	Verbosity string `yaml:"verbosity,omitempty"`
}
