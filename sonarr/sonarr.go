package sonarr

import (
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/trakt"
	"github.com/rs/zerolog"
	"strings"
	"sync"
	"time"
)

type Client struct {
	pvrType string
	name    string

	rootFolder       string
	qualityProfileId int

	apiURL     string
	apiHeaders map[string]string
	apiTimeout time.Duration

	cache             *cache.Client
	cacheTempDuration time.Duration
	cacheFiltersHash  string

	queue chan *nabarr.FeedItem
	wg    *sync.WaitGroup

	t           *trakt.Client
	log         zerolog.Logger
	ignoresExpr []*nabarr.ExprProgram
}

func New(c nabarr.PvrConfig, t *trakt.Client, cc *cache.Client) (*Client, error) {
	l := nabarr.GetLogger(c.Verbosity).With().
		Str("pvr_name", c.Name).
		Str("pvr_type", c.Type).
		Logger()

	// set config defaults (if not set)
	if c.Cache.TemporaryDuration == 0 {
		c.Cache.TemporaryDuration = 24 * time.Hour
	}

	// set api url
	apiURL := ""
	if strings.Contains(strings.ToLower(c.URL), "/api") {
		apiURL = c.URL
	} else {
		apiURL = nabarr.JoinURL(c.URL, "/api")
	}

	// set api headers
	apiHeaders := map[string]string{
		"X-Api-Key": c.ApiKey,
	}

	// create client
	cl := &Client{
		pvrType: "sonarr",
		name:    strings.ToLower(c.Name),

		rootFolder: c.RootFolder,

		cache:             cc,
		cacheTempDuration: c.Cache.TemporaryDuration,
		cacheFiltersHash:  nabarr.AsSHA256(c.Filters),

		queue: make(chan *nabarr.FeedItem, 1024),
		wg:    &sync.WaitGroup{},

		apiURL:     apiURL,
		apiHeaders: apiHeaders,
		apiTimeout: 60 * time.Second,

		t:   t,
		log: l,
	}

	// compile expressions
	if err := cl.compileExpressions(c.Filters); err != nil {
		return nil, fmt.Errorf("compile expressions: %w", err)
	}

	// validate api access
	ss, err := cl.getSystemStatus()
	if err != nil {
		return nil, fmt.Errorf("validate api: %w", err)
	}

	// get quality profile
	if qid, err := cl.getQualityProfileId(c.QualityProfile); err != nil {
		return nil, fmt.Errorf("get quality profile: %v: %w", c.QualityProfile, err)
	} else {
		cl.qualityProfileId = qid
	}

	cl.log.Info().
		Str("pvr_version", ss.Version).
		Msg("Initialised")
	return cl, nil
}

func (c *Client) Type() string {
	return c.pvrType
}
