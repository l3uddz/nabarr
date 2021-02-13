package media

import (
	"fmt"
	"github.com/l3uddz/nabarr/logger"
	"github.com/l3uddz/nabarr/media/omdb"
	"github.com/l3uddz/nabarr/media/trakt"
	"github.com/rs/zerolog"
)

type Client struct {
	trakt *trakt.Client
	omdb  *omdb.Client

	log zerolog.Logger
}

func New(cfg *Config) (*Client, error) {
	// trakt
	if cfg.Trakt.ClientId == "" {
		return nil, fmt.Errorf("trakt: no client_id specified")
	}
	t := trakt.New(&cfg.Trakt)

	// omdb
	o := omdb.New(&cfg.Omdb)

	return &Client{
		trakt: t,
		omdb:  o,

		log: logger.New(cfg.Verbosity).With().Logger(),
	}, nil
}
