package media

import (
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr/media/trakt"
	"strconv"
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
		c.log.Trace().
			Err(err).
			Str("imdb_id", item.ImdbId).
			Msg("Failed finding item on omdb")
	} else if oi != nil {
		mi.Omdb = *oi
	}

	return mi, nil
}
