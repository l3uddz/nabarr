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
		cacheKey := fmt.Sprintf("tvdb_%s", item.TvdbId)
		if _, err := c.cache.Get(cacheKey); err == nil {
			continue
		}
		_ = c.cache.Set(cacheKey, nil)

		// trakt search item
		show, err := c.t.GetShow(item)
		if err != nil {
			if errors.Is(err, trakt.ErrItemNotFound) {
				c.log.Warn().
					Err(err).
					Str("feed_title", item.Title).
					Str("feed_tvdb_id", item.TvdbId).
					Msg("Item was not found on trakt")
				continue
			}

			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("feed_tvdb_id", item.TvdbId).
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
				Msg("Failed checking show against ignore filters")
			continue
		}

		if ignore {
			c.log.Trace().
				Str("feed_title", item.Title).
				Str("trakt_title", show.Title).
				Str("trakt_tvdb_id", show.TvdbId).
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
					Msg("Item was not found via pvr lookup")
				continue
			}

			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("feed_tvdb_id", item.TvdbId).
				Msg("Failed finding item via pvr lookup")
		}

		if s.Id > 0 {
			// item already existed in pvr
			c.log.Info().
				Str("pvr_title", s.Title).
				Int("pvr_year", s.Year).
				Int("pvr_tvdb_id", s.TvdbId).
				Msg("Item already existed in pvr")
			continue
		}

		// add item to pvr
		slug := s.TitleSlug
		if slug == "" {
			slug = show.Slug
		}

		c.log.Info().
			Str("feed_title", item.Title).
			Str("trakt_title", show.Title).
			Str("trakt_tvdb_id", show.TvdbId).
			Msg("Sending show to PVR")
	}
}
