package trakt

import (
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"

	"github.com/l3uddz/nabarr/logger"
	"github.com/l3uddz/nabarr/util"
)

type Client struct {
	log  zerolog.Logger
	http *http.Client

	apiKey string
	apiURL string
}

func New(cfg *Config) *Client {
	l := logger.Child(logger.WithLevel(cfg.Verbosity)).With().
		Str("media", "trakt").Logger()

	return &Client{
		log:  l,
		http: util.NewRetryableHttpClient(30*time.Second, ratelimit.New(1, ratelimit.WithoutSlack), &l),

		apiKey: cfg.ClientId,
		apiURL: "https://api.trakt.tv",
	}
}

func (c *Client) getAuthHeaders() map[string]string {
	return map[string]string{
		"trakt-api-key":     c.apiKey,
		"trakt-api-version": "2",
	}
}
