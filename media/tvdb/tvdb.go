package tvdb

import (
	"fmt"
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

	apiKey     string
	apiURL     string
	apiHeaders map[string]string
}

func New(cfg *Config) *Client {
	l := logger.Child(logger.WithLevel(cfg.Verbosity)).With().
		Str("media", "tvdb").Logger()

	return &Client{
		log:  l,
		http: util.NewRetryableHttpClient(30*time.Second, ratelimit.New(1, ratelimit.WithoutSlack), &l),

		apiKey: cfg.ApiKey,
		apiURL: "https://api.thetvdb.com",
		apiHeaders: map[string]string{
			"Authorization": fmt.Sprintf("Bearer %s", cfg.ApiKey),
		},
	}
}
