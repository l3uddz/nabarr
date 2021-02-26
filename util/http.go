package util

import (
	"github.com/hashicorp/go-retryablehttp"
	"github.com/l3uddz/nabarr/build"
	"github.com/rs/zerolog"
	"go.uber.org/ratelimit"
	"net/http"
	"time"
)

func NewRetryableHttpClient(timeout time.Duration, rl ratelimit.Limiter, log *zerolog.Logger) *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10
	retryClient.RetryWaitMin = 1 * time.Second
	retryClient.RetryWaitMax = 10 * time.Second
	retryClient.RequestLogHook = func(l retryablehttp.Logger, request *http.Request, i int) {
		// set user-agent
		if request != nil {
			request.Header.Set("User-Agent", "nabarr/"+build.Version)
		}

		// rate limit
		if rl != nil {
			rl.Take()
		}

		// log
		if log != nil && request != nil && request.URL != nil {
			switch i {
			case 0:
				// first
				log.Trace().
					Str("url", request.URL.String()).
					Msg("Sending request")
			default:
				// retry
				log.Debug().
					Str("url", request.URL.String()).
					Int("attempt", i).
					Msg("Retrying failed request")
			}
		}
	}
	retryClient.HTTPClient.Timeout = timeout
	retryClient.Logger = nil
	return retryClient.StandardClient()
}
