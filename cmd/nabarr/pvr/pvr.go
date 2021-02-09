package pvr

import (
	"errors"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/radarr"
	"github.com/l3uddz/nabarr/sonarr"
	"github.com/l3uddz/nabarr/trakt"
	"strings"
)

type PVR interface {
	Type() string
	AddMediaItem(*nabarr.MediaItem) error
	ShouldIgnore(*nabarr.MediaItem) (bool, string, error)
	Start()
	Stop()
	QueueFeedItem(*nabarr.FeedItem)
}

func NewPVR(c nabarr.PvrConfig, t *trakt.Client, cc *cache.Client) (PVR, error) {
	// return pvr object
	switch strings.ToLower(c.Type) {
	case "sonarr":
		return sonarr.New(c, t, cc)
	case "radarr":
		return radarr.New(c, t, cc)
	}

	return nil, errors.New("unknown pvr")
}
