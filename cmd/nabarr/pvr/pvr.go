package pvr

import (
	"errors"
	"strings"

	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/cache"
	"github.com/l3uddz/nabarr/media"
	"github.com/l3uddz/nabarr/radarr"
	"github.com/l3uddz/nabarr/sonarr"
	"github.com/lefelys/state"
)

type PVR interface {
	Type() string
	GetFiltersHash() string
	AddMediaItem(*media.Item, ...nabarr.PvrOption) error
	ShouldIgnore(*media.Item) (bool, string, error)
	Start() state.State
	QueueFeedItem(*media.FeedItem)
}

func NewPVR(c nabarr.PvrConfig, mode string, m *media.Client, cc *cache.Client) (PVR, error) {
	// return pvr object
	switch strings.ToLower(c.Type) {
	case "sonarr":
		return sonarr.New(c, mode, m, cc)
	case "radarr":
		return radarr.New(c, mode, m, cc)
	}

	return nil, errors.New("unknown pvr")
}
