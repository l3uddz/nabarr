package media

import (
	"fmt"
	"strconv"
	"time"
)

func (c *Client) GetMovieInfo(item *FeedItem) (*Item, error) {
	// lookup on trakt
	m, err := c.trakt.GetMovie(item.ImdbId)
	if err != nil {
		return nil, fmt.Errorf("trakt: get movie: %w", err)
	}

	// transform trakt info to MediaItem
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

	// fetch additional info

	return mi, nil
}
