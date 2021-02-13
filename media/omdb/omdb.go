package omdb

import (
	"github.com/l3uddz/nabarr/logger"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
	"time"
)

type Client struct {
	apiKey string
	log    zerolog.Logger
	rl     ratelimit.Limiter

	apiURL     string
	apiTimeout time.Duration
}

func New(cfg *Config) *Client {
	return &Client{
		apiKey: cfg.ApiKey,
		log:    logger.New(cfg.Verbosity).With().Logger(),
		rl:     ratelimit.New(1, ratelimit.WithoutSlack),

		apiURL:     "https://www.omdbapi.com",
		apiTimeout: 30 * time.Second,
	}
}
