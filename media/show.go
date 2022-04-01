package media

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/l3uddz/nabarr/media/trakt"
)

func (c *Client) GetShowInfo(item *FeedItem) (*Item, error) {
	// retrieve and validate media provider data
	mdp, mdi := item.GetProviderData()
	if mdp == "" || mdi == "" {
		return nil, fmt.Errorf("trakt: get show: no media provider details found")
	}

	// lookup on trakt
	t, err := c.trakt.GetShow(mdp, mdi)
	if err != nil {
		if errors.Is(err, trakt.ErrItemNotFound) {
			return nil, fmt.Errorf("trakt: get show: show with %sId %q: %w", mdp, mdi, ErrItemNotFound)
		}
		return nil, fmt.Errorf("trakt: get show: show with %sId %q: %w", mdp, mdi, err)
	}

	// transform trakt info to MediaItem
	mi := &Item{
		TvdbId:        strconv.Itoa(t.Ids.Tvdb),
		TmdbId:        strconv.Itoa(t.Ids.Tmdb),
		ImdbId:        t.Ids.Imdb,
		Slug:          t.Ids.Slug,
		Title:         t.Title,
		FeedTitle:     item.Title,
		Summary:       t.Overview,
		Country:       []string{t.Country},
		Network:       t.Network,
		Date:          t.FirstAired,
		Year:          t.FirstAired.Year(),
		Runtime:       t.Runtime,
		Rating:        t.Rating,
		Votes:         t.Votes,
		Status:        t.Status,
		Genres:        t.Genres,
		Languages:     []string{t.Language},
		AiredEpisodes: t.AiredEpisodes,
	}

	// omdb
	if oi, err := c.omdb.GetItem(t.Ids.Imdb); err != nil {
		c.log.Debug().
			Err(err).
			Str("imdb_id", t.Ids.Imdb).
			Msg("Item was not found on omdb")
	} else if oi != nil {
		mi.Omdb = *oi
	}

	// tvdb
	if ti, err := c.tvdb.GetItem(strconv.Itoa(t.Ids.Tvdb)); err != nil {
		c.log.Debug().
			Err(err).
			Int("tvdb_id", t.Ids.Tvdb).
			Msg("Item was not found on tvdb")
	} else if ti != nil {
		mi.Tvdb = *ti

		// merge with trakt data
		if mi.Network == "" {
			mi.Network = ti.Network
		}
	}

	return mi, nil
}
