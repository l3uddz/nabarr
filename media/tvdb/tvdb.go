package tvdb

import (
	"fmt"
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
	apiHeaders map[string]string
	apiTimeout time.Duration
}

func New(cfg *Config) *Client {
	return &Client{
		apiKey: cfg.ApiKey,
		log:    logger.New(cfg.Verbosity).With().Logger(),
		rl:     ratelimit.New(1, ratelimit.WithoutSlack),

		apiURL: "https://api.thetvdb.com",
		apiHeaders: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cfg.ApiKey),
		},
		apiTimeout: 30 * time.Second,
	}
}
