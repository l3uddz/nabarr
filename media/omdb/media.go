package omdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/lucperkins/rek"

	"github.com/l3uddz/nabarr/util"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) GetItem(imdbId string) (*Item, error) {
	// empty item when appropriate
	if c.apiKey == "" || imdbId == "" {
		return nil, nil
	}

	// prepare request
	reqUrl, err := util.URLWithQuery(c.apiURL, url.Values{
		"apikey": []string{c.apiKey},
		"i":      []string{imdbId}})
	if err != nil {
		return nil, fmt.Errorf("generate lookup request url: %w", err)
	}

	// send request
	resp, err := rek.Get(reqUrl, rek.Client(c.http))
	if err != nil {
		return nil, fmt.Errorf("request lookup: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate lookup response: %s", resp.Status())
	}

	// decode response
	b := new(lookupResponse)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode lookup response: %w", err)
	}

	if b.Title == "" {
		return nil, fmt.Errorf("item with imdbId: %v: %w", imdbId, ErrItemNotFound)
	}

	// transform response
	rt := 0
	for _, rating := range b.Ratings {
		if strings.EqualFold(rating.Source, "Rotten Tomatoes") {
			rt = util.Atoi(strings.TrimSuffix(rating.Value, "%"), 0)
			break
		}
	}

	return &Item{
		Metascore:      util.Atoi(b.Metascore, 0),
		RottenTomatoes: rt,
		ImdbRating:     util.Atof64(b.ImdbRating, 0.0),
		ImdbVotes:      util.Atoi(util.StripNonNumeric(b.ImdbVotes), 0),
		Language:       b.Language,
		Country:        b.Country,
		Rated:          b.Rated,
	}, nil
}
