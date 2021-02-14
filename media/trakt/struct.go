package trakt

import (
	"time"
)

type ShowIds struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	Tvdb  int    `json:"tvdb"`
	Imdb  string `json:"imdb"`
	Tmdb  int    `json:"tmdb"`
}

type MovieIds struct {
	Trakt int    `json:"trakt"`
	Slug  string `json:"slug"`
	Imdb  string `json:"imdb"`
	Tmdb  int    `json:"tmdb"`
}

type Show struct {
	Type                  string    `json:"type"`
	Title                 string    `json:"title"`
	Year                  int       `json:"year"`
	Ids                   ShowIds   `json:"ids"`
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

type Movie struct {
	Type                  string   `json:"type"`
	Title                 string   `json:"title"`
	Year                  int      `json:"year"`
	Ids                   MovieIds `json:"ids"`
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
