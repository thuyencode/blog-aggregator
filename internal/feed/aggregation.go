package feed

import (
	"context"
	"fmt"
	"log"
	"time"
	"uuid"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/thuyencode/blog-aggregator/internal/database"
)

func ScrapeFeeds(queries *database.Queries) error {
	nextFeedToFetch, err := queries.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed to fetch: %w", err)
	}

	feed, err := scrapeFeed(queries, nextFeedToFetch)
	if err != nil {
		return err
	}

	go func() {
		savePosts(queries, feed, nextFeedToFetch.ID)
	}()

	return nil
}

func scrapeFeed(queries *database.Queries, toFetch database.Feed) (*RSSFeed, error) {
	if err := queries.MarkFeedFetched(context.Background(), toFetch.ID); err != nil {
		return nil, fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	feed, err := FetchFeed(context.Background(), toFetch.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}

	return feed, nil
}

func savePosts(queries *database.Queries, feed *RSSFeed, feedID uuid.UUID) {
	for _, item := range feed.Channel.Item {
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("failed to parse pubDate: %s\n", err)
			continue
		}

		_, err = queries.CreatePost(context.Background(), database.CreatePostParams{
			Title:       item.Title,
			Description: item.Description,
			Url:         item.Link,
			PublishedAt: publishedAt,
			FeedID:      feedID,
		})
		if err != nil {
			if pqErr := pq.As(err, pqerror.UniqueViolation); pqErr != nil {
				continue
			}
			log.Printf("failed to save post: %s\n", err)
		}
	}
}
