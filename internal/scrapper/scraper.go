package scrapper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/rss"
	"github.com/google/uuid"
)

const (
	expiredFeedThreshold = 2 * time.Minute
)

func StartScraping(
	db *queries.Queries,
	concurrencyAmount int,
	requestInterval time.Duration,
) {
	log.Printf("Scraping on %v goroutines every %v", concurrencyAmount, requestInterval)

	// Create a ticker that will fire every requestInterval
	ticker := time.NewTicker(requestInterval)

	// Create a loop that will run every time the ticker fires
	for ; ; <-ticker.C {
		// Fetch all expired feeds
		feeds, err := db.GetNextFeedsToFetch(context.Background(), int32(concurrencyAmount))
		if err != nil {
			log.Printf("Error fetching feeds: %s", err.Error())
			continue
		}

		// If there are no feeds to fetch, continue
		if len(feeds) == 0 {
			continue
		}

		// Filter all feeds that are within the expiredFeedThreshold
		var expiredFeeds []queries.Feed
		for _, feed := range feeds {
			if feed.LastFetchedAt.Valid && time.Since(feed.LastFetchedAt.Time) < expiredFeedThreshold {
				continue
			}
			expiredFeeds = append(expiredFeeds, feed)
		}

		// Create a wait group to wait for all goroutines to finish
		wg := &sync.WaitGroup{}

		for _, feed := range expiredFeeds {
			wg.Add(1)

			go scrapeFeed(wg, db, feed)
		}

		wg.Wait()
	}
}

func scrapeFeed(wg *sync.WaitGroup, db *queries.Queries, feed queries.Feed) {
	defer wg.Done()

	// Mark feed as fetching
	_, err := db.UpdateFeedFetchTime(context.Background(), queries.UpdateFeedFetchTimeParams{
		ID:            feed.ID,
		LastFetchedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	if err != nil {
		log.Printf("Error updating feed fetch time: %s", err.Error())
		return
	}

	// Fetch feed
	fetchedFeed, err := rss.URLToFeed(feed.Url)
	if err != nil {
		log.Printf("Error fetching feed: %s", err.Error())
		return
	}

	for _, item := range fetchedFeed.Channel.Item {
		log.Printf("Found Item: %+v", item.Title)
		// Parse Publication Date
		// 	TODO: This should be made more robust to support all types of date formats
		parsedPubDate, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing publication date: %s", err.Error())
			continue
		}
		// Insert item
		_, err = db.CreatePost(context.Background(), queries.CreatePostParams{
			ID:          uuid.New(),
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: parsedPubDate,
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error inserting post: %s", err.Error())
			continue
		}

	}
	log.Printf("Found %v items", len(fetchedFeed.Channel.Item))
}
