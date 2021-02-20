package radarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/media"
	"github.com/l3uddz/nabarr/util"
	"github.com/lucperkins/rek"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrItemNotFound = errors.New("not found")
)

func (c *Client) getSystemStatus() (*systemStatus, error) {
	// send request
	resp, err := rek.Get(util.JoinURL(c.apiURL, "system", "status"), rek.Headers(c.apiHeaders),
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
	resp, err := rek.Get(util.JoinURL(c.apiURL, "profile"), rek.Headers(c.apiHeaders),
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

func (c *Client) lookupMediaItem(item *media.Item) (*lookupRequest, error) {
	// determine metadata id to use
	mdType := "imdb"
	mdId := item.ImdbId

	if item.TmdbId != "" && item.TmdbId != "0" {
		// radarr prefers tmdb
		mdType = "tmdb"
		mdId = item.TmdbId
	}

	// prepare request
	reqUrl, err := util.URLWithQuery(util.JoinURL(c.apiURL, "movie", "lookup"),
		url.Values{"term": []string{fmt.Sprintf("%s:%s", mdType, mdId)}})
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
		switch mdType {
		case "tmdb":
			if strconv.Itoa(s.TmdbId) == item.TmdbId {
				return &s, nil
			}
		default:
			if s.ImdbId == item.ImdbId {
				return &s, nil
			}
		}
	}

	return nil, fmt.Errorf("movie lookup %sId: %v: %w", mdType, mdId, ErrItemNotFound)
}

func (c *Client) AddMediaItem(item *media.Item, opts ...nabarr.PvrOption) error {
	// prepare options
	_, err := buildOptions(opts...)
	if err != nil {
		return fmt.Errorf("build options: %v: %w", item.TmdbId, err)
	}

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
		TmdbId: util.Atoi(item.TmdbId, 0),
		ImdbId: item.ImdbId,
	}

	// send request
	resp, err := rek.Post(util.JoinURL(c.apiURL, "movie"), rek.Headers(c.apiHeaders), rek.Json(req),
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
