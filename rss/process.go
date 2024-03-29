package rss

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strings"

	"github.com/lucperkins/rek"

	"github.com/l3uddz/nabarr/media"
	"github.com/l3uddz/nabarr/util"
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

	queuedItems := 0
	for p := range items {
		if j.queueItemWithPvrs(&items[p]) {
			queuedItems++
		}
	}

	if queuedItems == 0 {
		j.log.Debug().Msg("Refreshed, no valid items to queue")
		return nil
	}

	j.log.Info().
		Int("count", len(items)).
		Msg("Queued items")
	return nil
}

func (j *rssJob) queueItemWithPvrs(item *media.FeedItem) bool {
	queued := false
	for _, pvr := range j.pvrs {
		switch {
		case (item.TvdbId != "" || item.TmdbId != "") && util.ContainsTvCategory(item.Categories) && pvr.Type() == "sonarr":
			// tvdbId/tmdbId is present, queue with sonarr
			pvr.QueueFeedItem(item)
			queued = true
		case (item.ImdbId != "" || item.TmdbId != "") && util.ContainsMovieCategory(item.Categories) && pvr.Type() == "radarr":
			// imdbId is present, queue with radarr
			pvr.QueueFeedItem(item)
			queued = true
		}
	}
	return queued
}

func (j *rssJob) getFeed() ([]media.FeedItem, error) {
	// request feed
	res, err := rek.Get(j.url, rek.Client(j.http))
	if err != nil {
		return nil, fmt.Errorf("request feed: %w", err)
	}
	defer res.Body().Close()

	// validate response
	if res.StatusCode() != 200 {
		return nil, fmt.Errorf("validate response: %s", res.Status())
	}

	// decode response
	b := new(media.Rss)
	if err := xml.NewDecoder(res.Body()).Decode(b); err != nil {
		return nil, fmt.Errorf("decode feed: %w", err)
	}

	// prepare result
	items := make([]media.FeedItem, 0)
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

		// item seen before?
		if cacheValue, err := j.cache.Get(j.name, i.GUID); err == nil {
			if string(cacheValue) == j.cacheFiltersHash {
				// item has been seen before and the filter hash has not changed since
				continue
			}
			// item has been seen, however the filter hash has changed, re-process
		}

		// process feed item attributes
		for _, a := range i.Attributes {
			switch strings.ToLower(a.Name) {
			case "category":
				b.Channel.Items[p].Categories = append(b.Channel.Items[p].Categories, a.Value)
			case "language":
				b.Channel.Items[p].Language = a.Value
			case "tvdb", "tvdbid", "thetvdb":
				b.Channel.Items[p].TvdbId = a.Value
			case "imdb", "imdbid":
				if strings.HasPrefix(a.Value, "tt") {
					b.Channel.Items[p].ImdbId = a.Value
				} else {
					b.Channel.Items[p].ImdbId = fmt.Sprintf("tt%s", a.Value)
				}
			case "tmdb", "tmdbid":
				b.Channel.Items[p].TmdbId = a.Value
			}
		}

		// validate item
		switch {
		case b.Channel.Items[p].TvdbId != "" && !util.StringSliceContains([]string{"0", "1"}, b.Channel.Items[p].TvdbId):
			// tvdb id is present, allow processing
			break
		case b.Channel.Items[p].ImdbId != "" && strings.HasPrefix(b.Channel.Items[p].ImdbId, "tt"):
			// imdb id present, allow processing
			break
		case b.Channel.Items[p].TmdbId != "" && !util.StringSliceContains([]string{"0", "1"}, b.Channel.Items[p].TmdbId):
			// tmdb id present, allow processing
			break
		default:
			// skip item as an expected media provider id was not present
			continue
		}

		// add validated item for processing
		b.Channel.Items[p].Feed = j.name
		items = append(items, b.Channel.Items[p])

		// add item to temp cache (to prevent re-processing)
		if err := j.cache.Put(j.name, i.GUID, []byte(j.cacheFiltersHash), j.cacheDuration); err != nil {
			j.log.Error().
				Err(err).
				Str("guid", i.GUID).
				Msg("Failed storing item in temp cache")
		}
	}

	return items, nil
}
