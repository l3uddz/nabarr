package media

import (
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr/media/trakt"
	"strconv"
	"time"
)

func (c *Client) GetMovieInfo(item *FeedItem) (*Item, error) {
	// lookup on trakt
	m, err := c.trakt.GetMovie(item.ImdbId)
	if err != nil {
		if errors.Is(err, trakt.ErrItemNotFound) {
			return nil, fmt.Errorf("trakt: get movie: movie with imdbId %q: %w", item.ImdbId, ErrItemNotFound)
		}
		return nil, fmt.Errorf("trakt: get movie: movie with imdbId %q: %w", item.ImdbId, err)
	}

	// transform trakt info
	date, err := time.Parse("2006-01-02", m.Released)
	if err != nil {
		date = time.Time{}
	}

	mi := &Item{
		TvdbId:    "",
		TmdbId:    strconv.Itoa(m.Ids.Tmdb),
		ImdbId:    m.Ids.Imdb,
		Slug:      m.Ids.Slug,
		Title:     m.Title,
		FeedTitle: item.Title,
		Summary:   m.Overview,
		Country:   []string{m.Country},
		Network:   "",
		Date:      date,
		Year:      date.Year(),
		Runtime:   m.Runtime,
		Rating:    m.Rating,
		Votes:     m.Votes,
		Status:    m.Status,
		Genres:    m.Genres,
		Languages: []string{m.Language},
	}

	// omdb
	if oi, err := c.omdb.GetItem(item.ImdbId); err != nil {
		c.log.Debug().
			Err(err).
			Str("imdb_id", item.ImdbId).
			Msg("Failed finding item on omdb")
	} else if oi != nil {
		mi.Omdb = *oi
	}

	return mi, nil
}
