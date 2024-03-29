package trakt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/lucperkins/rek"

	"github.com/l3uddz/nabarr/util"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) GetShow(providerType string, providerId string) (*Show, error) {
	// prepare request
	reqUrl, err := util.URLWithQuery(util.JoinURL(c.apiURL, "search", providerType, providerId),
		url.Values{
			"type":     []string{"show"},
			"extended": []string{"full"}})
	if err != nil {
		return nil, fmt.Errorf("generate lookup show request url: %w", err)
	}

	// send request
	resp, err := rek.Get(reqUrl, rek.Client(c.http), rek.Headers(c.getAuthHeaders()))
	if err != nil {
		return nil, fmt.Errorf("request show: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate show response: %s", resp.Status())
	}

	// decode response
	b := new([]struct{ Show Show })
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode show response: %w", err)
	}

	if len(*b) < 1 {
		return nil, ErrItemNotFound
	}

	// translate response
	show := &(*b)[0].Show
	show.Ids.Imdb = util.StripNonAlphaNumeric(show.Ids.Imdb)

	return show, nil
}

func (c *Client) GetMovie(providerType string, providerId string) (*Movie, error) {
	// prepare request
	reqUrl, err := util.URLWithQuery(util.JoinURL(c.apiURL, "search", providerType, providerId),
		url.Values{
			"type":     []string{"movie"},
			"extended": []string{"full"}})
	if err != nil {
		return nil, fmt.Errorf("generate lookup movie request url: %w", err)
	}

	// send request
	resp, err := rek.Get(reqUrl, rek.Client(c.http), rek.Headers(c.getAuthHeaders()))
	if err != nil {
		return nil, fmt.Errorf("request movie: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate movie response: %s", resp.Status())
	}

	// decode response
	b := new([]struct{ Movie Movie })
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode movie response: %w", err)
	}

	if len(*b) < 1 {
		return nil, ErrItemNotFound
	}

	// translate response
	movie := &(*b)[0].Movie
	movie.Ids.Imdb = util.StripNonAlphaNumeric(movie.Ids.Imdb)

	return movie, nil
}
