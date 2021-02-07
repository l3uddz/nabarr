package trakt

import (
	"github.com/l3uddz/nabarr"
	"strconv"
	"time"
)

type showIds struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	Tvdb  int    `json:"tvdb"`
	Imdb  string `json:"imdb"`
	Tmdb  int    `json:"tmdb"`
}

type movieIds struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	Imdb  string `json:"imdb"`
	Tmdb  int    `json:"tmdb"`
}

type show struct {
	Type                  string    `json:"type"`
	Title                 string    `json:"title"`
	Year                  int       `json:"year"`
	Ids                   showIds   `json:"ids"`
	Overview              string    `json:"overview"`
	FirstAired            time.Time `json:"first_aired"`
	Runtime               int       `json:"runtime"`
	Certification         string    `json:"certification"`
	Network               string    `json:"network"`
	Country               string    `json:"country"`
	Trailer               string    `json:"trailer"`
	Homepage              string    `json:"homepage"`
	Status                string    `json:"status"`
	Rating                float64   `json:"rating"`
	Votes                 int       `json:"votes"`
	CommentCount          int       `json:"comment_count"`
	Language              string    `json:"language"`
	AvailableTranslations []string  `json:"available_translations"`
	Genres                []string  `json:"genres"`
	AiredEpisodes         int       `json:"aired_episodes"`
	Character             string    `json:"character"`
}

type movie struct {
	Type                  string   `json:"type"`
	Title                 string   `json:"title"`
	Year                  int      `json:"year"`
	Ids                   movieIds `json:"ids"`
	Tagline               string   `json:"tagline"`
	Overview              string   `json:"overview"`
	Released              string   `json:"released"`
	Runtime               int      `json:"runtime"`
	Country               string   `json:"country"`
	Trailer               string   `json:"trailer"`
	Homepage              string   `json:"homepage"`
	Status                string   `json:"status"`
	Rating                float64  `json:"rating"`
	Votes                 int      `json:"votes"`
	CommentCount          int      `json:"comment_count"`
	Language              string   `json:"language"`
	AvailableTranslations []string `json:"available_translations"`
	Genres                []string `json:"genres"`
	Certification         string   `json:"certification"`
	Character             string   `json:"character"`
}

func (s *show) ToMediaItem(item *nabarr.FeedItem) *nabarr.MediaItem {
	return &nabarr.MediaItem{
		TvdbId:    strconv.Itoa(s.Ids.Tvdb),
		TmdbId:    strconv.Itoa(s.Ids.Tmdb),
		ImdbId:    s.Ids.Imdb,
		Slug:      s.Ids.Slug,
		Title:     s.Title,
		FeedTitle: item.Title,
		Summary:   s.Overview,
		Country:   []string{s.Country},
		Network:   s.Network,
		Date:      s.FirstAired,
		Year:      s.FirstAired.Year(),
		Runtime:   s.Runtime,
		Rating:    s.Rating,
		Status:    s.Status,
		Genres:    s.Genres,
		Languages: []string{s.Language},
	}
}

func (m *movie) ToMediaItem(item *nabarr.FeedItem) *nabarr.MediaItem {
	date, err := time.Parse("2006-01-02", m.Released)
	if err != nil {
		date = time.Time{}
	}

	return &nabarr.MediaItem{
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
		Status:    m.Status,
		Genres:    m.Genres,
		Languages: []string{m.Language},
	}
}
