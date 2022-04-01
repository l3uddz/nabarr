package tvdb

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/l3uddz/nabarr/util"
	"github.com/lucperkins/rek"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) GetItem(tvdbId string) (*Item, error) {
	// empty item when appropriate
	if c.apiKey == "" || tvdbId == "" {
		return nil, nil
	}

	// prepare request
	reqUrl := util.JoinURL(c.apiURL, "series", tvdbId)

	// send request
	resp, err := rek.Get(reqUrl, rek.Client(c.http), rek.Headers(c.apiHeaders))
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

	if b.Data.SeriesName == "" {
		return nil, fmt.Errorf("item with tvdbId: %v: %w", tvdbId, ErrItemNotFound)
	}

	return &Item{
		Runtime:         util.Atoi(b.Data.Runtime, 0),
		Language:        b.Data.Language,
		Network:         b.Data.Network,
		Genre:           b.Data.Genre,
		AirsDayOfWeek:   b.Data.AirsDayOfWeek,
		SiteRating:      b.Data.SiteRating,
		SiteRatingCount: b.Data.SiteRatingCount,
	}, nil
}
