package trakt

import (
	"github.com/l3uddz/nabarr"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
	"time"
)

type Client struct {
	clientId string
	log      zerolog.Logger
	rl       ratelimit.Limiter

	apiURL     string
	apiTimeout time.Duration
}

func New(cfg *Config) *Client {
	return &Client{
		clientId: cfg.ClientId,
		log:      nabarr.GetLogger(cfg.Verbosity).With().Logger(),
		rl:       ratelimit.New(1, ratelimit.WithoutSlack),

		apiURL:     "https://api.trakt.tv",
		apiTimeout: 30 * time.Second,
	}
}

func (c *Client) getAuthHeaders() map[string]string {
	return map[string]string{
		"trakt-api-key":     c.clientId,
		"trakt-api-version": "2",
	}
}
