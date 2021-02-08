package sonarr

import (
	"errors"
	"fmt"
	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/trakt"
)

func (c *Client) QueueFeedItem(item *nabarr.FeedItem) {
	c.queue <- item
}

func (c *Client) queueProcessor() {
	for item := range c.queue {
		// stop processing
		if item == nil {
			return
		}

		// validate item has required id(s)
		if item.TvdbId == "" {
			continue
		}

		// check cache / add item to cache
		pvrCacheBucket := fmt.Sprintf("pvr_%s_%s", c.Type(), c.name)
		cacheKey := fmt.Sprintf("tvdb_%s", item.TvdbId)
		if _, err := c.cache.Get(pvrCacheBucket, cacheKey); err == nil {
			// item already exists in the cache (was previously looked up or exists in the pvr)
			continue
		} else {
			// item did not exist in the cache, lets insert a temp cache entry
			if err := c.cache.Put(pvrCacheBucket, cacheKey, nil, c.cacheTempDuration); err != nil {
				c.log.Error().
					Err(err).
					Msg("Failed storing item in temp cache")
			}
		}

		// trakt search item
		show, err := c.t.GetShow(item)
		if err != nil {
			if errors.Is(err, trakt.ErrItemNotFound) {
				c.log.Debug().
					Err(err).
					Str("feed_title", item.Title).
					Str("feed_tvdb_id", item.TvdbId).
					Str("feed_name", item.Feed).
					Msg("Item was not found on trakt")
				continue
			}

			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("feed_tvdb_id", item.TvdbId).
				Str("feed_name", item.Feed).
				Msg("Failed finding item on trakt")
			continue
		}

		// trakt expression check
		ignore, err := c.ShouldIgnore(show)
		if err != nil {
			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("trakt_title", show.Title).
				Str("trakt_tvdb_id", show.TvdbId).
				Str("feed_name", item.Feed).
				Msg("Failed checking show against ignore filters")
			continue
		}

		if ignore {
			c.log.Trace().
				Str("feed_title", item.Title).
				Str("trakt_title", show.Title).
				Str("trakt_tvdb_id", show.TvdbId).
				Str("feed_name", item.Feed).
				Msg("Show matched ignore filters")
			continue
		}

		// lookup item in pvr
		s, err := c.lookupMediaItem(show)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				// the item was not found
				c.log.Warn().
					Err(err).
					Str("feed_title", item.Title).
					Str("feed_tvdb_id", item.TvdbId).
					Str("feed_name", item.Feed).
					Msg("Item was not found via pvr lookup")
				continue
			}

			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("feed_tvdb_id", item.TvdbId).
				Str("feed_name", item.Feed).
				Msg("Failed finding item via pvr lookup")
		}

		if s.Id > 0 {
			// item already existed in pvr
			c.log.Debug().
				Str("pvr_title", s.Title).
				Int("pvr_year", s.Year).
				Int("pvr_tvdb_id", s.TvdbId).
				Str("feed_name", item.Feed).
				Msg("Item already existed in pvr")

			// add item to perm cache (items already in pvr)
			if err := c.cache.Put(pvrCacheBucket, cacheKey, nil, 0); err != nil {
				c.log.Error().
					Err(err).
					Msg("Failed storing item in perm cache")
			}
			continue
		}

		// add item to pvr
		c.log.Debug().
			Str("feed_title", item.Title).
			Str("trakt_title", show.Title).
			Str("trakt_tvdb_id", show.TvdbId).
			Int("trakt_year", show.Year).
			Str("feed_name", item.Feed).
			Msg("Sending show to pvr")

		if s.TitleSlug != "" {
			// use slug from pvr search
			show.Slug = s.TitleSlug
		}

		if err := c.AddMediaItem(show); err != nil {
			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("trakt_title", show.Title).
				Str("trakt_tvdb_id", show.TvdbId).
				Int("trakt_year", show.Year).
				Str("feed_name", item.Feed).
				Msg("Failed adding item to pvr")
		}

		// add item to perm cache (item was added to pvr)
		if err := c.cache.Put(pvrCacheBucket, cacheKey, nil, 0); err != nil {
			c.log.Error().
				Err(err).
				Msg("Failed storing item in perm cache")
		}

		c.log.Info().
			Err(err).
			Str("trakt_title", show.Title).
			Str("trakt_tvdb_id", show.TvdbId).
			Int("trakt_year", show.Year).
			Str("feed_name", item.Feed).
			Msg("Added item")
	}
}
