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
	t, err := c.trakt.GetMovie(item.ImdbId)
	if err != nil {
		if errors.Is(err, trakt.ErrItemNotFound) {
			return nil, fmt.Errorf("trakt: get movie: movie with imdbId %q: %w", item.ImdbId, ErrItemNotFound)
		}
		return nil, fmt.Errorf("trakt: get movie: movie with imdbId %q: %w", item.ImdbId, err)
	}

	// transform trakt info
	date, err := time.Parse("2006-01-02", t.Released)
	if err != nil {
		date = time.Time{}
	}

	mi := &Item{
		TvdbId:    "",
		TmdbId:    strconv.Itoa(t.Ids.Tmdb),
		ImdbId:    t.Ids.Imdb,
		Slug:      t.Ids.Slug,
		Title:     t.Title,
		FeedTitle: item.Title,
		Summary:   t.Overview,
		Country:   []string{t.Country},
		Network:   "",
		Date:      date,
		Year:      date.Year(),
		Runtime:   t.Runtime,
		Rating:    t.Rating,
		Votes:     t.Votes,
		Status:    t.Status,
		Genres:    t.Genres,
		Languages: []string{t.Language},
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

	return mi, nil
}
