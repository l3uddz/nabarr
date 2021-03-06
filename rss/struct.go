package rss

import (
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/cmd/nabarr/pvr"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"net/http"
	"time"
)

type feedItem struct {
	Name          string        `yaml:"name"`
	URL           string        `yaml:"url"`
	Cron          string        `yaml:"cron"`
	CacheDuration time.Duration `yaml:"cache_duration"`
	Pvrs          []string      `yaml:"pvrs"`
}

type Config struct {
	Feeds []feedItem `yaml:"feeds"`

	Verbosity string `yaml:"verbosity,omitempty"`
}

type rssJob struct {
	name string
	log  zerolog.Logger
	http *http.Client

	url  string
	pvrs map[string]pvr.PVR

	attempts int
	errors   []error

	cron             *cron.Cron
	cache            *cache.Client
	cacheDuration    time.Duration
	cacheFiltersHash string
	jobID            cron.EntryID
}
