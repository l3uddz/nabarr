package radarr

import (
	"github.com/l3uddz/nabarr"
	"github.com/rs/zerolog"
)

type Client struct {
	pvrType string

	log zerolog.Logger
}

func New(c *nabarr.PvrConfig) (*Client, error) {
	l := nabarr.GetLogger(c.Verbosity).With().
		Str("name", c.Name).
		Str("type", c.Type).
		Logger()

	return &Client{
		pvrType: "radarr",

		log: l,
	}, nil
}

func (c *Client) Type() string {
	return c.pvrType
}
