package omdb

import (
	"net/http"
	"time"

	"github.com/l3uddz/nabarr/logger"
	"github.com/l3uddz/nabarr/util"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
)

type Client struct {
	log  zerolog.Logger
	http *http.Client

	apiKey string
	apiURL string
}

func New(cfg *Config) *Client {
	l := logger.New(cfg.Verbosity).With().
		Str("media", "omdb").
		Logger()

	return &Client{
		log:  l,
		http: util.NewRetryableHttpClient(30*time.Second, ratelimit.New(1, ratelimit.WithoutSlack), &l),

		apiKey: cfg.ApiKey,
		apiURL: "https://www.omdbapi.com",
	}
}
