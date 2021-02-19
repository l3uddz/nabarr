package media

import (
	"github.com/l3uddz/nabarr/media/omdb"
	"github.com/l3uddz/nabarr/media/trakt"
	"github.com/l3uddz/nabarr/media/tvdb"
)

type Config struct {
	Trakt trakt.Config `yaml:"trakt"`
	Omdb  omdb.Config  `yaml:"omdb"`
	Tvdb  tvdb.Config  `yaml:"tvdb"`

	Verbosity string `yaml:"verbosity,omitempty"`
}
