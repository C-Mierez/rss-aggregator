package scrapper

import (
	"context"
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/c-mierez/rss-aggregator/internal/lib/queries"
	"github.com/c-mierez/rss-aggregator/internal/lib/rss"
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

		// Create a wait group to wait for all goroutines to finish
		wg := &sync.WaitGroup{}

		for _, feed := range feeds {
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
		log.Printf("Found Item: %+v", item)
	}
	log.Printf("Found %v items", len(fetchedFeed.Channel.Item))
}
