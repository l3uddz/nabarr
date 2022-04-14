package radarr

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/logger"
	"github.com/l3uddz/nabarr/media"
	"github.com/l3uddz/nabarr/util"
)

type Client struct {
	pvrType     string
	name        string
	testMode    bool
	testModeAdd bool

	rootFolder       string
	qualityProfileId int

	// options
	searchMissing bool
	addMonitored  bool

	apiURL     string
	apiHeaders map[string]string

	cache             *cache.Client
	cacheTempDuration time.Duration
	cacheFiltersHash  string

	queue chan *media.FeedItem

	m           *media.Client
	http        *http.Client
	log         zerolog.Logger
	ignoresExpr []*nabarr.ExprProgram
}

func New(c nabarr.PvrConfig, mode string, m *media.Client, cc *cache.Client) (*Client, error) {
	l := logger.Child(logger.WithLevel(c.Verbosity)).With().
		Str("pvr_name", c.Name).
		Str("pvr_type", c.Type).
		Logger()

	// set config defaults (if not set)
	if c.CacheDuration == 0 {
		c.CacheDuration = 24 * time.Hour
	}

	// set api url
	apiURL := ""
	if strings.Contains(strings.ToLower(c.URL), "/api") {
		apiURL = c.URL
	} else {
		apiURL = util.JoinURL(c.URL, "api", "v3")
	}

	// set api headers
	apiHeaders := map[string]string{
		"X-Api-Key": c.ApiKey,
	}

	// create client
	cl := &Client{
		pvrType:     "radarr",
		name:        strings.ToLower(c.Name),
		testMode:    util.StringSliceContains([]string{"test", "test-add"}, mode),
		testModeAdd: strings.EqualFold(mode, "test-add"),

		rootFolder:    c.RootFolder,
		searchMissing: util.BoolOrDefault(c.Options.SearchMissing, true),
		addMonitored:  util.BoolOrDefault(c.Options.AddMonitored, true),

		cache:             cc,
		cacheTempDuration: c.CacheDuration,
		cacheFiltersHash:  util.AsSHA256(c.Filters),

		queue: make(chan *media.FeedItem, 1024),

		apiURL:     apiURL,
		apiHeaders: apiHeaders,

		m:    m,
		http: util.NewRetryableHttpClient(60*time.Second, nil, &l),
		log:  l,
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

func (c *Client) GetFiltersHash() string {
	return c.cacheFiltersHash
}
