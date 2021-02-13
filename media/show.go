package media

import (
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr/media/trakt"
	"strconv"
)

func (c *Client) GetShowInfo(item *FeedItem) (*Item, error) {
	// lookup on trakt
	s, err := c.trakt.GetShow(item.TvdbId)
	if err != nil {
		if errors.Is(err, trakt.ErrItemNotFound) {
			return nil, fmt.Errorf("trakt: get show: show with tvdbId %q: %w", item.TvdbId, ErrItemNotFound)
		}
		return nil, fmt.Errorf("trakt: get show: show with tvdbId %q: %w", item.TvdbId, err)
	}

	// transform trakt info to MediaItem
	mi := &Item{
		TvdbId:        strconv.Itoa(s.Ids.Tvdb),
		TmdbId:        strconv.Itoa(s.Ids.Tmdb),
		ImdbId:        s.Ids.Imdb,
		Slug:          s.Ids.Slug,
		Title:         s.Title,
		FeedTitle:     item.Title,
		Summary:       s.Overview,
		Country:       []string{s.Country},
		Network:       s.Network,
		Date:          s.FirstAired,
		Year:          s.FirstAired.Year(),
		Runtime:       s.Runtime,
		Rating:        s.Rating,
		Votes:         s.Votes,
		Status:        s.Status,
		Genres:        s.Genres,
		Languages:     []string{s.Language},
		AiredEpisodes: s.AiredEpisodes,
	}

	// fetch additional info

	return mi, nil
}
