package trakt

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/lucperkins/rek"
	"net/url"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) GetShow(item *nabarr.FeedItem) (*nabarr.MediaItem, error) {
	// validate item has a valid id
	if item.TvdbId == "" {
		return nil, fmt.Errorf("no tvdbId for: %v", item.Title)
	}

	// prepare request
	reqUrl, err := nabarr.URLWithQuery(nabarr.JoinURL(c.apiURL, fmt.Sprintf("/search/tvdb/%s", item.TvdbId)),
		url.Values{
			"type":     []string{"show"},
			"extended": []string{"full"}})
	if err != nil {
		return nil, fmt.Errorf("generate lookup show request url: %w", err)
	}

	c.log.Trace().
		Str("url", reqUrl).
		Msg("Searching trakt")

	// send request
	c.rl.Take()
	resp, err := rek.Get(reqUrl, rek.Headers(c.getAuthHeaders()), rek.Timeout(c.apiTimeout))
	if err != nil {
		return nil, fmt.Errorf("request show: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate show response: %s", resp.Status())
	}

	// decode response
	b := new([]struct{ Show show })
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode show response: %w", err)
	}

	if len(*b) < 1 {
		return nil, fmt.Errorf("show with tvdbId: %v: %w", item.TvdbId, ErrItemNotFound)
	}

	// translate response
	return (*b)[0].Show.ToMediaItem(item), nil
}

func (c *Client) GetMovie(item *nabarr.FeedItem) (*nabarr.MediaItem, error) {
	// validate item has a valid id
	if item.ImdbId == "" {
		return nil, fmt.Errorf("no imdbId for: %v", item.Title)
	}

	// prepare request
	reqUrl, err := nabarr.URLWithQuery(nabarr.JoinURL(c.apiURL, fmt.Sprintf("/search/imdb/%s", item.ImdbId)),
		url.Values{
			"type":     []string{"movie"},
			"extended": []string{"full"}})
	if err != nil {
		return nil, fmt.Errorf("generate lookup movie request url: %w", err)
	}

	c.log.Trace().
		Str("url", reqUrl).
		Msg("Searching trakt")

	// send request
	c.rl.Take()
	resp, err := rek.Get(reqUrl, rek.Headers(c.getAuthHeaders()), rek.Timeout(c.apiTimeout))
	if err != nil {
		return nil, fmt.Errorf("request movie: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate movie response: %s", resp.Status())
	}

	// decode response
	b := new([]struct{ Movie movie })
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode movie response: %w", err)
	}

	if len(*b) < 1 {
		return nil, fmt.Errorf("movie with imdbId: %v: %w", item.ImdbId, ErrItemNotFound)
	}

	// translate response
	return (*b)[0].Movie.ToMediaItem(item), nil
}
