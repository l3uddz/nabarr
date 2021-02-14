package media

import (
	"github.com/l3uddz/nabarr/media/omdb"
	"github.com/l3uddz/nabarr/media/trakt"
)

type Config struct {
	Trakt trakt.Config `yaml:"trakt"`
	Omdb  omdb.Config  `yaml:"omdb"`

	Verbosity string `yaml:"verbosity,omitempty"`
}
