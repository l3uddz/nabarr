package pvr

import (
	"errors"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/sonarr"
	"github.com/l3uddz/nabarr/trakt"
	"strings"
)

type PVR interface {
	Type() string
	AddMediaItem(*nabarr.MediaItem) error
	ShouldIgnore(*nabarr.MediaItem) (bool, error)
	QueueFeedItem(*nabarr.FeedItem)
}

func NewPVR(c nabarr.PvrConfig, t *trakt.Client) (PVR, error) {
	// return pvr object
	switch strings.ToLower(c.Type) {
	case "sonarr":
		return sonarr.New(c, t)
		//case "radarr":
		//	return radarr.New(c)
	}

	return nil, errors.New("unknown pvr")
}