package radarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/lucperkins/rek"
	"net/url"
	"strings"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) getSystemStatus() (*systemStatus, error) {
	// send request
	resp, err := rek.Get(nabarr.JoinURL(c.apiURL, "/system/status"), rek.Headers(c.apiHeaders),
		rek.Timeout(c.apiTimeout))
	if err != nil {
		return nil, fmt.Errorf("request system status: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate system status response: %s", resp.Status())
	}

	// decode response
	b := new(systemStatus)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode system status response: %w", err)
	}

	return b, nil
}

func (c *Client) getQualityProfileId(profileName string) (int, error) {
	// send request
	resp, err := rek.Get(nabarr.JoinURL(c.apiURL, "/profile"), rek.Headers(c.apiHeaders),
		rek.Timeout(c.apiTimeout))
	if err != nil {
		return 0, fmt.Errorf("request quality profiles: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("validate quality profiles response: %s", resp.Status())
	}

	// decode response
	b := new([]qualityProfile)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return 0, fmt.Errorf("decode quality profiles response: %w", err)
	}

	// find quality profile
	for _, profile := range *b {
		if strings.EqualFold(profile.Name, profileName) {
			return profile.Id, nil
		}
	}

	return 0, errors.New("quality profile not found")
}

func (c *Client) lookupMediaItem(item *nabarr.MediaItem) (*lookupRequest, error) {
	// prepare request
	reqUrl, err := nabarr.URLWithQuery(nabarr.JoinURL(c.apiURL, "/movie/lookup"),
		url.Values{"term": []string{fmt.Sprintf("imdb:%s", item.ImdbId)}})
	if err != nil {
		return nil, fmt.Errorf("generate movie lookup request url: %w", err)
	}

	// send request
	resp, err := rek.Get(reqUrl, rek.Headers(c.apiHeaders), rek.Timeout(c.apiTimeout))
	if err != nil {
		return nil, fmt.Errorf("request movie lookup: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate movie lookup response: %s", resp.Status())
	}

	// decode response
	b := new([]lookupRequest)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode movie lookup response: %w", err)
	}

	// find movie
	for _, s := range *b {
		if s.ImdbId == item.ImdbId {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("movie lookup imdbId: %v: %w", item.ImdbId, ErrItemNotFound)
}

func (c *Client) AddMediaItem(item *nabarr.MediaItem) error {
	// prepare request
	req := addRequest{
		Title:               item.Title,
		TitleSlug:           item.Slug,
		Year:                item.Year,
		QualityProfileId:    c.qualityProfileId,
		Images:              []string{},
		Monitored:           true,
		RootFolderPath:      c.rootFolder,
		MinimumAvailability: "released",
		AddOptions: addOptions{
			SearchForMovie:             true,
			IgnoreEpisodesWithFiles:    false,
			IgnoreEpisodesWithoutFiles: false,
		},
		TmdbId: nabarr.Atoi(item.TmdbId, 0),
		ImdbId: item.ImdbId,
	}

	// send request
	resp, err := rek.Post(nabarr.JoinURL(c.apiURL, "/movie"), rek.Headers(c.apiHeaders), rek.Json(req),
		rek.Timeout(c.apiTimeout))
	if err != nil {
		return fmt.Errorf("request add movie: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return fmt.Errorf("validate add movie response: %s", resp.Status())
	}

	return nil
}
