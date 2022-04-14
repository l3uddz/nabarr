package media

import (
	"fmt"

	"github.com/rs/zerolog"

	"github.com/l3uddz/nabarr/logger"
	"github.com/l3uddz/nabarr/media/omdb"
	"github.com/l3uddz/nabarr/media/trakt"
	"github.com/l3uddz/nabarr/media/tvdb"
)

type Client struct {
	trakt *trakt.Client
	omdb  *omdb.Client
	tvdb  *tvdb.Client

	log zerolog.Logger
}

func New(cfg *Config) (*Client, error) {
	// validate trakt configured (it is mandatory)
	if cfg.Trakt.ClientId == "" {
		return nil, fmt.Errorf("trakt: no client_id specified")
	}

	return &Client{
		trakt: trakt.New(&cfg.Trakt),
		omdb:  omdb.New(&cfg.Omdb),
		tvdb:  tvdb.New(&cfg.Tvdb),

		log: logger.Child(logger.WithLevel(cfg.Verbosity)),
	}, nil
}
