package sheetsclient

import (
	"context"
	"fmt"
	"time"
)

const flushInterval = time.Minute * 5

func (c *Client) Get(name string) (int, error) {
	if c.ArticleViews == nil {
		return 0, fmt.Errorf("no article views set")
	}
	for _, v := range c.ArticleViews {
		if v.Title == name {
			return v.Count, nil
		}
	}
	return 0, fmt.Errorf("failed to find article in cache")
}

func (c *Client) Increment(name string) {
	for i, v := range c.ArticleViews {
		if v.Title == name {
			c.ArticleViews[i].Count++
			return
		}
	}
}

func (c *Client) SetFlushRoutine(ctx context.Context) {
	ticker := time.NewTicker(flushInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			c.FlushToSheets()
		}
	}
}
