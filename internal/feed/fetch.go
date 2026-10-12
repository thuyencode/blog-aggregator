package feed

import (
	"context"
	"encoding/xml"
	"html"
	"net/http"
)

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

func FetchFeed(ctx context.Context, feedURL string) (feed *RSSFeed, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", "gator/0.0.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer func() {
		err = resp.Body.Close()
	}()

	decoder := xml.NewDecoder(resp.Body)
	if err = decoder.Decode(&feed); err != nil {
		return
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i, item := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(item.Title)
		feed.Channel.Item[i].Description = html.UnescapeString(item.Description)
	}

	return
}
