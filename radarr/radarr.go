package radarr

import (
	"fmt"
	"github.com/ReneKroon/ttlcache/v2"
	"github.com/antonmedv/expr/vm"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/trakt"
	"github.com/rs/zerolog"
	"strings"
	"time"
)

type Client struct {
	pvrType          string
	rootFolder       string
	qualityProfileId int

	apiURL     string
	apiHeaders map[string]string
	apiTimeout time.Duration

	cacheTemp *ttlcache.Cache
	cachePerm map[string]int

	queue chan *nabarr.FeedItem

	t           *trakt.Client
	log         zerolog.Logger
	ignoresExpr []*vm.Program
}

func New(c nabarr.PvrConfig, t *trakt.Client) (*Client, error) {
	l := nabarr.GetLogger(c.Verbosity).With().
		Str("pvr_name", c.Name).
		Str("pvr_type", c.Type).
		Logger()

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
		pvrType:    "radarr",
		rootFolder: c.RootFolder,

		cacheTemp: ttlcache.NewCache(),
		cachePerm: make(map[string]int, 0),

		queue: make(chan *nabarr.FeedItem, 1024),

		apiURL:     apiURL,
		apiHeaders: apiHeaders,
		apiTimeout: 60 * time.Second,

		t:   t,
		log: l,
	}

	// setup cache
	if err := cl.cacheTemp.SetTTL(24 * time.Hour); err != nil {
		return nil, fmt.Errorf("set cache ttl: %w", err)
	}
	cl.cacheTemp.SkipTTLExtensionOnHit(false)

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

	// start queue queueProcessor
	go cl.queueProcessor()

	cl.log.Info().
		Str("pvr_version", ss.Version).
		Msg("Initialised")
	return cl, nil
}

func (c *Client) Type() string {
	return c.pvrType
}
