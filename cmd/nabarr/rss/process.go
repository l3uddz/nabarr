package rss

import (
	"encoding/xml"
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/lucperkins/rek"
	"sort"
	"strings"
	"time"
)

func (j *rssJob) process() error {
	// retrieve feed items
	j.log.Debug().Msg("Refreshing")
	items, err := j.getFeed()
	if err != nil {
		return fmt.Errorf("get feed: %w", err)
	}

	// add feed items to pvrs
	if len(items) == 0 {
		j.log.Debug().Msg("Refreshed, no items to queue")
		return nil
	}

	for p, _ := range items {
		j.queueItemWithPvrs(&items[p])
	}

	j.log.Info().
		Int("count", len(items)).
		Msg("Queued items")
	return nil
}

func (j *rssJob) queueItemWithPvrs(item *nabarr.FeedItem) {
	for _, pvr := range j.pvrs {
		switch {
		case item.TvdbId != "" && pvr.Type() == "sonarr":
			// tvdbId is present, queue with sonarr
			pvr.QueueFeedItem(item)
		case item.ImdbId != "" && pvr.Type() == "radarr":
			// imdbId is present, queue with radarr
			pvr.QueueFeedItem(item)
		}
	}
}

func (j *rssJob) getFeed() ([]nabarr.FeedItem, error) {
	// request feed
	res, err := rek.Get(j.url, rek.Timeout(30*time.Minute))
	if err != nil {
		return nil, fmt.Errorf("request feed: %w", err)
	}
	defer res.Body().Close()

	// validate response
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("validate response: %s", res.Status())
	}

	// decode response
	b := new(nabarr.Rss)
	if err := xml.NewDecoder(res.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode feed: %w", err)
	}

	// prepare result
	items := make([]nabarr.FeedItem, 0)
	if len(b.Channel.Items) < 1 {
		return items, nil
	}

	// sort response items
	sort.SliceStable(b.Channel.Items, func(i, j int) bool {
		return b.Channel.Items[i].PubDate.After(b.Channel.Items[j].PubDate.Time)
	})

	// process feed items
	for p, i := range b.Channel.Items {
		// ignore items
		if i.GUID == "" {
			// items must always have a guid
			continue
		}

		// guid seen before?
		cacheKey := fmt.Sprintf("%s_%s", j.name, i.GUID)
		if cacheValue, err := j.cache.Get(j.name, cacheKey); err == nil {
			if string(cacheValue) == j.cacheFiltersHash {
				// item has been seen before and the filters have not changed
				continue
			}
			// item has been seen, however the filters have changed since it was last seen, re-process
		}

		// process feed item attributes
		for _, a := range i.Attributes {
			switch strings.ToLower(a.Name) {
			case "language":
				b.Channel.Items[p].Language = a.Value
			case "tvdb", "tvdbid":
				b.Channel.Items[p].TvdbId = a.Value
			case "imdb", "imdbid":
				if strings.HasPrefix(a.Value, "tt") {
					b.Channel.Items[p].ImdbId = a.Value
				} else {
					b.Channel.Items[p].ImdbId = fmt.Sprintf("tt%s", a.Value)
				}
			}
		}

		// validate item
		if b.Channel.Items[p].TvdbId == "" && b.Channel.Items[p].ImdbId == "" {
			continue
		}

		// add validated item for processing
		b.Channel.Items[p].Feed = j.name
		items = append(items, b.Channel.Items[p])

		// add item to temp cache (to prevent re-processing)
		if err := j.cache.Put(j.name, cacheKey, []byte(j.cacheFiltersHash), j.cacheDuration); err != nil {
			j.log.Error().
				Err(err).
				Str("guid", i.GUID).
				Msg("Failed storing item in temp cache")
		}
	}

	return items, nil
}
