package pvr

import (
	"errors"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/radarr"
	"github.com/l3uddz/nabarr/sonarr"
	"github.com/l3uddz/nabarr/trakt"
	"github.com/lefelys/state"
	"strings"
)

type PVR interface {
	Type() string
	AddMediaItem(*nabarr.MediaItem) error
	ShouldIgnore(*nabarr.MediaItem) (bool, string, error)
	Start() state.State
	QueueFeedItem(*nabarr.FeedItem)
}

func NewPVR(c nabarr.PvrConfig, mode string, t *trakt.Client, cc *cache.Client) (PVR, error) {
	// return pvr object
	switch strings.ToLower(c.Type) {
	case "sonarr":
		return sonarr.New(c, mode, t, cc)
	case "radarr":
		return radarr.New(c, mode, t, cc)
	}

	return nil, errors.New("unknown pvr")
}
