package core

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

type rss struct {
	articles []Item
	list     string
}

func (r *rss) fetchRSSData() {
	xmlParser := gofeed.NewParser()
	feed, err := xmlParser.ParseURL(RssURL)
	if err != nil {
		logger.Sugar().Errorf("error parsing rss: %s\n", err.Error())
		return
	}
	feedLength := len(feed.Items)
	for i := 0; i < feedLength-1; i++ {
		feedItem := feed.Items[i]
		r.articles = append(r.articles, Item{
			Title:     feedItem.Title,
			Permalink: feedItem.Link,
			Updated:   feedItem.Updated,
		})
		r.list += fmt.Sprintf(ListItem, feedItem.Link, feedItem.Title, PublishedAt, feedItem.PublishedParsed.String()[:10])
	}
}
