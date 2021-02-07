package nabarr

import (
	"encoding/xml"
	"github.com/pkg/errors"
	"time"
)

type Rss struct {
	Channel struct {
		Items []FeedItem `xml:"item"`
	} `xml:"channel"`
}

type FeedItem struct {
	Title    string `xml:"title,omitempty"`
	Category string `xml:"category,omitempty"`
	GUID     string `xml:"guid,omitempty"`
	PubDate  Time   `xml:"pubDate,omitempty"`

	// attributes
	Language string
	TvdbId   string `xml:"tvdb,omitempty"`
	TvMazeId string
	ImdbId   string `xml:"imdb,omitempty"`

	Attributes []struct {
		XMLName xml.Name
		Name    string `xml:"name,attr"`
		Value   string `xml:"value,attr"`
	} `xml:"attr"`
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

type MediaItem struct {
	TvdbId        string
	TmdbId        string
	ImdbId        string
	Slug          string
	FeedTitle     string
	Title         string
	Summary       string
	Country       []string
	Network       string
	Date          time.Time
	Year          int
	Runtime       int
	Rating        float64
	Status        string
	Genres        []string
	Languages     []string
	AiredEpisodes int
}

type ExprEnv struct {
	MediaItem
	Now func() time.Time
}

func GetExprEnv(media *MediaItem) *ExprEnv {
	return &ExprEnv{
		MediaItem: *media,
		Now:       func() time.Time { return time.Now().UTC() },
	}
}

func StringOrDefault(currentValue *string, defaultValue string) string {
	if currentValue == nil {
		return defaultValue
	}

	return *currentValue
}
func Uint64OrDefault(currentValue *uint64, defaultValue uint64) uint64 {
	if currentValue == nil {
		return defaultValue
	}

	return *currentValue
}
