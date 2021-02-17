package media

import (
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr/media/trakt"
	"strconv"
	"strings"
)

func (c *Client) GetShowInfo(item *FeedItem) (*Item, error) {
	// lookup on trakt
	t, err := c.trakt.GetShow(item.TvdbId)
	if err != nil {
		if errors.Is(err, trakt.ErrItemNotFound) {
			return nil, fmt.Errorf("trakt: get show: show with tvdbId %q: %w", item.TvdbId, ErrItemNotFound)
		}
		return nil, fmt.Errorf("trakt: get show: show with tvdbId %q: %w", item.TvdbId, err)
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

		if oi.Language != "" {
			mi.Languages = append(mi.Languages, strings.Split(oi.Language, ", ")...)
		}

		if oi.Country != "" {
			mi.Country = append(mi.Country, strings.Split(oi.Country, ", ")...)
		}
	}

	return mi, nil
}
