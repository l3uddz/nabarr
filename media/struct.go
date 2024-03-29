package media

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/l3uddz/nabarr/media/omdb"
	"github.com/l3uddz/nabarr/media/tvdb"
	"github.com/l3uddz/nabarr/util"
)

/* Media Item(s) */

type Item struct {
	TvdbId        string    `json:"TvdbId,omitempty"`
	TmdbId        string    `json:"TmdbId,omitempty"`
	ImdbId        string    `json:"ImdbId,omitempty"`
	Slug          string    `json:"Slug,omitempty"`
	FeedTitle     string    `json:"FeedTitle,omitempty"`
	Title         string    `json:"Title,omitempty"`
	Summary       string    `json:"Summary,omitempty"`
	Country       []string  `json:"Country,omitempty"`
	Network       string    `json:"Network,omitempty"`
	Date          time.Time `json:"Date"`
	Year          int       `json:"Year,omitempty"`
	Runtime       int       `json:"Runtime,omitempty"`
	Rating        float64   `json:"Rating,omitempty"`
	Votes         int       `json:"Votes,omitempty"`
	Status        string    `json:"Status,omitempty"`
	Genres        []string  `json:"Genres,omitempty"`
	Languages     []string  `json:"Languages,omitempty"`
	AiredEpisodes int       `json:"AiredEpisodes,omitempty"`

	// additional media provider data
	Omdb omdb.Item `json:"Omdb,omitempty"`
	Tvdb tvdb.Item `json:"Tvdb,omitempty"`
}

func (i *Item) GetProviderData() (string, string) {
	switch {
	case i.TvdbId != "" && i.TvdbId != "0":
		return "tvdb", i.TvdbId
	case i.TmdbId != "" && i.TmdbId != "0":
		return "tmdb", i.TmdbId
	case i.ImdbId != "":
		return "imdb", i.ImdbId
	}
	return "", ""
}

/* Rss Item(s) */

type Rss struct {
	Channel struct {
		Items []FeedItem `xml:"item"`
	} `xml:"channel"`
}

type FeedItem struct {
	Title      string `xml:"title,omitempty"`
	Category   string `xml:"category,omitempty"`
	Categories []string
	GUID       string `xml:"guid,omitempty"`
	PubDate    Time   `xml:"pubDate,omitempty"`

	// set by processor
	Feed string

	// attributes
	Language string
	TvdbId   string `xml:"tvdb,omitempty"`
	TvMazeId string
	ImdbId   string `xml:"imdb,omitempty"`
	TmdbId   string `xml:"tmdb,omitempty"`

	Attributes []struct {
		XMLName xml.Name
		Name    string `xml:"name,attr"`
		Value   string `xml:"value,attr"`
	} `xml:"attr"`
}

func (f *FeedItem) GetProviderData() (string, string) {
	switch {
	case f.TvdbId != "" && !util.StringSliceContains([]string{"0", "1"}, f.TvdbId):
		return "tvdb", f.TvdbId
	case f.TmdbId != "" && !util.StringSliceContains([]string{"0", "1"}, f.TmdbId):
		return "tmdb", f.TmdbId
	case f.ImdbId != "" && strings.HasPrefix(f.ImdbId, "tt"):
		return "imdb", f.ImdbId
	}
	return "", ""
}

// Time credits: https://github.com/mrobinsn/go-newznab/blob/cd89d9c56447859fa1298dc9a0053c92c45ac7ef/newznab/structs.go#L150
type Time struct {
	time.Time
}

func (t *Time) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if err := e.EncodeToken(start); err != nil {
		return errors.Wrap(err, "failed to encode xml token")
	}
	if err := e.EncodeToken(xml.CharData([]byte(t.UTC().Format(time.RFC1123Z)))); err != nil {
		return errors.Wrap(err, "failed to encode xml token")
	}
	if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
		return errors.Wrap(err, "failed to encode xml token")
	}
	return nil
}

func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var raw string

	err := d.DecodeElement(&raw, &start)
	if err != nil {
		return err
	}
	date, err := time.Parse(time.RFC1123Z, raw)

	if err != nil {
		return err
	}

	*t = Time{date}
	return nil
}
