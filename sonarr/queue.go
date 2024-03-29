package sonarr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lefelys/state"

	"github.com/l3uddz/nabarr"
	"github.com/l3uddz/nabarr/media"
	"github.com/l3uddz/nabarr/util"
)

func (c *Client) QueueFeedItem(item *media.FeedItem) {
	c.queue <- item
}

func (c *Client) Start() state.State {
	st, tail := state.WithShutdown()
	go c.queueProcessor(tail)
	return st
}

func (c *Client) queueProcessor(tail state.ShutdownTail) {
	for {
		select {
		case <-tail.End():
			// shutdown
			tail.Done()
			return
		case feedItem := <-c.queue:
			// stop processing
			if feedItem == nil {
				tail.Done()
				return
			}

			// retrieve and validate media provider data
			mdp, mdi := feedItem.GetProviderData()
			if mdp == "" || mdi == "" {
				continue
			}

			// check cache / add item to cache
			cacheBucket := fmt.Sprintf("pvr_%s_%s", c.Type(), c.name)
			cacheKey := fmt.Sprintf("%s_%s", mdp, mdi)
			if !c.testMode {
				// not running in test mode, so use cache
				if cacheValue, err := c.cache.Get(cacheBucket, cacheKey); err == nil {
					// item already exists in the cache
					switch string(cacheValue) {
					case c.name:
						// perm cache entry, item exists in the pvr
						continue
					case c.cacheFiltersHash:
						// temp cache entry, item recently checked with the same filters
						continue
					}
				}

				// insert temp cache entry
				if err := c.cache.Put(cacheBucket, cacheKey, []byte(c.cacheFiltersHash), c.cacheTempDuration); err != nil {
					c.log.Error().
						Err(err).
						Msg("Failed storing item in temp cache")
				}
			}

			// get media info
			mediaItem, err := c.m.GetShowInfo(feedItem)
			if err != nil {
				if errors.Is(err, media.ErrItemNotFound) {
					c.log.Debug().
						Err(err).
						Str("feed_title", feedItem.Title).
						Str(fmt.Sprintf("feed_%s_id", mdp), mdi).
						Str("feed_name", feedItem.Feed).
						Msg("Item was not found on trakt")
					continue
				}

				c.log.Error().
					Err(err).
					Str("feed_title", feedItem.Title).
					Str(fmt.Sprintf("feed_%s_id", mdp), mdi).
					Str("feed_name", feedItem.Feed).
					Msg("Failed finding item on trakt")
				continue
			}

			if c.testMode {
				c.log.Debug().
					Interface("trakt_item", mediaItem).
					Msg("Item found on trakt")
			}

			// trakt expression check
			ignore, filter, err := c.ShouldIgnore(mediaItem)
			if err != nil {
				c.log.Error().
					Err(err).
					Str("feed_title", mediaItem.FeedTitle).
					Str("trakt_title", mediaItem.Title).
					Str("trakt_tvdb_id", mediaItem.TvdbId).
					Str("feed_name", feedItem.Feed).
					Str("ignore_filter", filter).
					Msg("Failed checking item against ignore filters")
				continue
			}

			if ignore {
				c.log.Debug().
					Str("feed_title", mediaItem.FeedTitle).
					Str("trakt_title", mediaItem.Title).
					Str("trakt_tvdb_id", mediaItem.TvdbId).
					Str("feed_name", feedItem.Feed).
					Str("ignore_filter", filter).
					Msg("Item matched ignore filters")
				continue
			}

			// lookup item in pvr
			s, err := c.lookupMediaItem(mediaItem)
			if err != nil {
				if errors.Is(err, ErrItemNotFound) {
					// the item was not found
					c.log.Warn().
						Err(err).
						Str("feed_title", mediaItem.FeedTitle).
						Str(fmt.Sprintf("feed_%s_id", mdp), mdi).
						Str("feed_name", feedItem.Feed).
						Msg("Item was not found via pvr lookup")
					continue
				}

				c.log.Error().
					Err(err).
					Str("feed_title", mediaItem.FeedTitle).
					Str(fmt.Sprintf("feed_%s_id", mdp), mdi).
					Str("feed_name", feedItem.Feed).
					Msg("Failed finding item via pvr lookup")
				continue
			}

			if s.Id > 0 {
				// item already existed in pvr
				c.log.Debug().
					Str("pvr_title", s.Title).
					Int("pvr_year", s.Year).
					Int("pvr_tvdb_id", s.TvdbId).
					Str("feed_name", feedItem.Feed).
					Msg("Item already existed in pvr")

				// add item to perm cache (items already in pvr)
				if !c.testMode {
					if err := c.cache.Put(cacheBucket, cacheKey, []byte(c.name), 0); err != nil {
						c.log.Error().
							Err(err).
							Msg("Failed storing item in perm cache")
					}
				}
				continue
			}

			// lookup pvr exclusions
			exclusions, err := c.getExclusions()
			if err != nil {
				c.log.Error().
					Err(err).
					Msg("Failed retrieving pvr exclusions")
				continue
			}

			if _, ok := exclusions[s.TvdbId]; ok {
				// item was in pvr exclusions
				c.log.Debug().
					Str("pvr_title", s.Title).
					Int("pvr_year", s.Year).
					Int("pvr_tvdb_id", s.TvdbId).
					Str("feed_name", feedItem.Feed).
					Msg("Item found in pvr exclusions")
				continue
			}

			// set appropriate series type
			switch {
			case util.StringSliceContains(mediaItem.Genres, "anime"), util.StringSliceContains(mediaItem.Tvdb.Genre, "anime"):
				s.SeriesType = "anime"
			}

			// check if item should be skipped (skip options)
			if c.skipAnime && strings.EqualFold(s.SeriesType, "anime") {
				c.log.Debug().
					Str("trakt_title", mediaItem.Title).
					Str("trakt_tvdb_id", mediaItem.TvdbId).
					Int("trakt_year", mediaItem.Year).
					Str("feed_name", feedItem.Feed).
					Msg("Skipping item (skip_anime enabled)")
				continue
			}

			// add item to pvr
			c.log.Debug().
				Str("feed_title", mediaItem.FeedTitle).
				Str("trakt_title", mediaItem.Title).
				Str("trakt_tvdb_id", mediaItem.TvdbId).
				Int("trakt_year", mediaItem.Year).
				Str("feed_name", feedItem.Feed).
				Msg("Sending item to pvr")

			if s.TitleSlug != "" {
				// use slug from pvr search
				mediaItem.Slug = s.TitleSlug
			}

			if c.testMode && !c.testModeAdd {
				c.log.Info().
					Str("trakt_title", mediaItem.Title).
					Str("trakt_tvdb_id", mediaItem.TvdbId).
					Int("trakt_year", mediaItem.Year).
					Str("feed_name", feedItem.Feed).
					Msg("Added item (test mode)")
				continue
			}

			opts := []nabarr.PvrOption{
				nabarr.WithSeriesType(s.SeriesType),
				nabarr.WithAddMonitored(c.addMonitored),
				nabarr.WithSearchMissing(c.searchMissing),
			}

			if err := c.AddMediaItem(mediaItem, opts...); err != nil {
				c.log.Error().
					Err(err).
					Str("feed_title", mediaItem.FeedTitle).
					Str("trakt_title", mediaItem.Title).
					Str("trakt_tvdb_id", mediaItem.TvdbId).
					Int("trakt_year", mediaItem.Year).
					Str("feed_name", feedItem.Feed).
					Msg("Failed adding item to pvr")
				continue
			}

			// add item to perm cache (item was added to pvr)
			if !c.testMode {
				if err := c.cache.Put(cacheBucket, cacheKey, []byte(c.name), 0); err != nil {
					c.log.Error().
						Err(err).
						Msg("Failed storing item in perm cache")
				}
			}

			c.log.Info().
				Str("trakt_title", mediaItem.Title).
				Str("trakt_tvdb_id", mediaItem.TvdbId).
				Int("trakt_year", mediaItem.Year).
				Str("feed_name", feedItem.Feed).
				Msg("Added item")
		}
	}
}
