package sheetsclient

import (
	"context"
	"fmt"
	"log"
	"time"
)

const flushInterval = time.Second * 60

// GetViews returns the number of views for the given article from the cache
// If the article is not found in the cache, it returns an error
func (c *Client) GetViews(name string) (int, error) {
	if c.ArticleViews == nil {
		return 0, fmt.Errorf("no article views set")
	}
	for _, v := range c.ArticleViews {
		if v.Name == name {
			return v.Count, nil
		}
	}
	return 0, fmt.Errorf("failed to find article in cache")
}

func (c *Client) IncrementMain() {
	c.MainData.HomePageViewCount++
}

func (c *Client) Increment(name string) error {
	for i, v := range c.ArticleViews {
		if v.Name == name {
			c.ArticleViews[i].Count++
			return nil
		}
	}
	return fmt.Errorf("failed to find article in cache")
}

func (c *Client) SetFlushRoutine(ctx context.Context) {
	ticker := time.NewTicker(flushInterval)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return
		case <-ticker.C:
			err := c.FlushToSheets()
			if err != nil {
				log.Printf("failed to flush to sheets: %v", err)
			}
		}
	}
}
