package sonarr

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
	resp, err := rek.Get(util.JoinURL(c.apiURL, "system", "status"), rek.Client(c.http), rek.Headers(c.apiHeaders))
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
	resp, err := rek.Get(util.JoinURL(c.apiURL, "qualityprofile"), rek.Client(c.http), rek.Headers(c.apiHeaders))
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

func (c *Client) getLanguageProfileId(profileName string) (int, error) {
	// send request
	resp, err := rek.Get(util.JoinURL(c.apiURL, "languageprofile"), rek.Client(c.http), rek.Headers(c.apiHeaders))
	if err != nil {
		return 0, fmt.Errorf("request language profiles: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return 0, fmt.Errorf("validate language profiles response: %s", resp.Status())
	}

	// decode response
	b := new([]languageProfile)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return 0, fmt.Errorf("decode language profiles response: %w", err)
	}

	// find quality profile
	for _, profile := range *b {
		if strings.EqualFold(profile.Name, profileName) {
			return profile.Id, nil
		}
	}

	return 0, errors.New("language language not found")
}

func (c *Client) lookupMediaItem(item *media.Item) (*lookupRequest, error) {
	// retrieve and validate media provider data
	mdp, mdi := item.GetProviderData()
	if mdp == "" || mdi == "" {
		return nil, fmt.Errorf("no media provider details found")
	}

	// prepare request
	reqUrl, err := util.URLWithQuery(util.JoinURL(c.apiURL, "series", "lookup"),
		url.Values{"term": []string{fmt.Sprintf("%s:%s", mdp, mdi)}})
	if err != nil {
		return nil, fmt.Errorf("generate series lookup request url: %w", err)
	}

	// send request
	resp, err := rek.Get(reqUrl, rek.Client(c.http), rek.Headers(c.apiHeaders))
	if err != nil {
		return nil, fmt.Errorf("request series lookup: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("validate series lookup response: %s", resp.Status())
	}

	// decode response
	b := new([]lookupRequest)
	if err := json.NewDecoder(resp.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode series lookup response: %w", err)
	}

	// find series
	for _, s := range *b {
		if strconv.Itoa(s.TvdbId) == mdi {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("series lookup %sId: %v: %w", mdp, mdi, ErrItemNotFound)
}

func (c *Client) AddMediaItem(item *media.Item, opts ...nabarr.PvrOption) error {
	// prepare options
	o, err := nabarr.BuildPvrOptions(opts...)
	if err != nil {
		return fmt.Errorf("build options: %v: %w", item.TvdbId, err)
	}

	// prepare request
	tvdbId, err := strconv.Atoi(item.TvdbId)
	if err != nil {
		return fmt.Errorf("converting tvdb id to int: %q", item.TvdbId)
	}

	req := addRequest{
		Title:             item.Title,
		TitleSlug:         item.Slug,
		Year:              item.Year,
		QualityProfileId:  c.qualityProfileId,
		LanguageProfileId: c.languageProfileId,
		Images:            []string{},
		Tags:              []string{},
		Monitored:         o.AddMonitored,
		RootFolderPath:    c.rootFolder,
		AddOptions: addOptions{
			SearchForMissingEpisodes:   o.SearchMissing,
			IgnoreEpisodesWithFiles:    false,
			IgnoreEpisodesWithoutFiles: false,
		},
		Seasons:      []string{},
		SeriesType:   util.StringOrDefault(o.SeriesType, "standard"),
		SeasonFolder: true,
		TvdbId:       tvdbId,
	}

	// send request
	resp, err := rek.Post(util.JoinURL(c.apiURL, "series"), rek.Client(c.http), rek.Headers(c.apiHeaders),
		rek.Json(req))
	if err != nil {
		return fmt.Errorf("request add series: %w", err)
	}
	defer resp.Body().Close()

	// validate response
	if resp.StatusCode() != 200 && resp.StatusCode() != 201 {
		return fmt.Errorf("validate add series response: %s", resp.Status())
	}

	return nil
}
