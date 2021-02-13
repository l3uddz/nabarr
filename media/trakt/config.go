package trakt

type Config struct {
	ClientId string `yaml:"client_id"`

	Verbosity string `yaml:"verbosity,omitempty"`
}
