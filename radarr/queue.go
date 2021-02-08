package radarr

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
		if item.ImdbId == "" {
			continue
		}

		// check cache / add item to cache
		pvrCacheBucket := fmt.Sprintf("pvr_%s_%s", c.Type(), c.name)
		cacheKey := fmt.Sprintf("imdb_%s", item.ImdbId)
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
		movie, err := c.t.GetMovie(item)
		if err != nil {
			if errors.Is(err, trakt.ErrItemNotFound) {
				c.log.Debug().
					Err(err).
					Str("feed_title", item.Title).
					Str("feed_imdb_id", item.ImdbId).
					Str("feed_name", item.Feed).
					Msg("Item was not found on trakt")
				continue
			}

			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("feed_imdb_id", item.ImdbId).
				Str("feed_name", item.Feed).
				Msg("Failed finding item on trakt")
			continue
		}

		// validate tmdbId was found
		if movie.TmdbId == "" || movie.TmdbId == "0" {
			c.log.Warn().
				Str("feed_title", item.Title).
				Str("feed_imdb_id", item.ImdbId).
				Str("feed_name", item.Feed).
				Msg("Item had no tmdbId on trakt")
			continue
		}

		// trakt expression check
		ignore, err := c.ShouldIgnore(movie)
		if err != nil {
			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("trakt_title", movie.Title).
				Str("trakt_imdb_id", movie.ImdbId).
				Str("feed_name", item.Feed).
				Msg("Failed checking movie against ignore filters")
			continue
		}

		if ignore {
			c.log.Trace().
				Str("feed_title", item.Title).
				Str("trakt_title", movie.Title).
				Str("trakt_imdb_id", movie.ImdbId).
				Str("feed_name", item.Feed).
				Msg("Show matched ignore filters")
			continue
		}

		// lookup item in pvr
		s, err := c.lookupMediaItem(movie)
		if err != nil {
			if errors.Is(err, ErrItemNotFound) {
				// the item was not found
				c.log.Warn().
					Err(err).
					Str("feed_title", item.Title).
					Str("feed_imdb_id", item.ImdbId).
					Str("feed_name", item.Feed).
					Msg("Item was not found via pvr lookup")
				continue
			}

			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("feed_imdb_id", item.ImdbId).
				Str("feed_name", item.Feed).
				Msg("Failed finding item via pvr lookup")
		}

		if s.Id > 0 {
			// item already existed in pvr
			c.log.Debug().
				Str("pvr_title", s.Title).
				Int("pvr_year", s.Year).
				Str("pvr_imdb_id", s.ImdbId).
				Int("pvr_tmdb_id", s.TmdbId).
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
			Str("trakt_title", movie.Title).
			Str("trakt_imdb_id", movie.ImdbId).
			Str("trakt_tmdb_id", movie.TmdbId).
			Int("trakt_year", movie.Year).
			Str("feed_name", item.Feed).
			Msg("Sending movie to pvr")

		if s.TitleSlug != "" {
			// use slug from pvr search
			movie.Slug = s.TitleSlug
		}

		if err := c.AddMediaItem(movie); err != nil {
			c.log.Error().
				Err(err).
				Str("feed_title", item.Title).
				Str("trakt_title", movie.Title).
				Str("trakt_imdb_id", movie.ImdbId).
				Str("trakt_tmdb_id", movie.TmdbId).
				Int("trakt_year", movie.Year).
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
			Str("trakt_title", movie.Title).
			Str("trakt_imdb_id", movie.ImdbId).
			Str("trakt_tmdb_id", movie.TmdbId).
			Int("trakt_year", movie.Year).
			Str("feed_name", item.Feed).
			Msg("Added item")
	}
}
