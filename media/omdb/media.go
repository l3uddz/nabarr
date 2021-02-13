package omdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr/util"
	"github.com/lucperkins/rek"
	"net/url"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) GetItem(imdbId string) (*Item, error) {
	// prepare request
	reqUrl, err := util.URLWithQuery(c.apiURL, url.Values{
		"apikey": []string{c.apiKey},
		"i":      []string{imdbId}})
	if err != nil {
		return nil, fmt.Errorf("generate lookup request url: %w", err)
	}

	c.log.Trace().
		Str("url", reqUrl).
		Msg("Searching omdb")

	// send request
	c.rl.Take()
	resp, err := rek.Get(reqUrl, rek.Timeout(c.apiTimeout))
	if err != nil {
		return nil, fmt.Errorf("request lookup: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate lookup response: %s", resp.Status())
	}

	// decode response
	b := new(Item)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode lookup response: %w", err)
	}

	if b.Title == "" {
		return nil, fmt.Errorf("item with imdbId: %v: %w", imdbId, ErrItemNotFound)
	}

	return b, nil
}
