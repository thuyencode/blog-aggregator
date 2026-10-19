package feed

import (
	"context"
	"fmt"
	"log"

	"github.com/thuyencode/blog-aggregator/internal/database"
)

func ScrapeFeeds(queries *database.Queries) error {
	nextFeedToFetch, err := queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	if err = queries.MarkFeedFetched(context.Background(), nextFeedToFetch.ID); err != nil {
		return fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	feed, err := FetchFeed(context.Background(), nextFeedToFetch.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch feed: %w", err)
	}

	log.Printf("feed %q fetched, posts list:\n", nextFeedToFetch.Name)
	for _, item := range feed.Channel.Item {
		fmt.Printf("* Title: %q\n", item.Title)
	}

	return nil
}
