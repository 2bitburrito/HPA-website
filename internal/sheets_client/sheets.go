package sheetsclient

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	mainTableBounds     = "main_table!A1:Z"
	pageDataTableBounds = "article_data!A1:Z"
)

func CreateSheetsService(spreadsheetID, serviceCredentials string) (*Client, error) {
	scopes := []string{
		"https://www.googleapis.com/auth/spreadsheets", // This is full read/write access
	}
	config, err := google.JWTConfigFromJSON([]byte(serviceCredentials), scopes...)
	if err != nil {
		return nil, fmt.Errorf("unable to create JWT configuration: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve sheets service: %w", err)
	}
	return &Client{
		Service: srv,
		creds: credentials{
			spreadsheetID:      spreadsheetID,
			serviceCredentials: serviceCredentials,
		},
	}, nil
}

// GetMainData returns the main table data for the site
// It ignores any caching and always retrieves the data from the sheets API
func (c *Client) GetMainData() (MainData, error) {
	t, err := c.getMainTable()
	if err != nil {
		return MainData{}, err
	}
	views, ok := t[0][1].(int)
	if !ok {
		return MainData{}, fmt.Errorf("unable to convert main table view count to int")
	}
	c.MainData = MainData{
		HomePageViewCount: views,
	}

	return c.MainData, nil
}

func (c *Client) GetAllArticleViews() (ArticleViewCounts, error) {
	pd, err := c.getPageData()
	if err != nil {
		return ArticleViewCounts{}, err
	}
	return c.extractArticlesFromData(pd)
}

// extractArticlesFromData takes raw matrix data from page_data and extracts the article titles,
// view counts and row numbers
func (c *Client) extractArticlesFromData(pageData [][]any) (ArticleViewCounts, error) {
	for i, row := range pageData {
		if i == 0 {
			continue
		}
		title, ok := row[0].(string)
		if !ok {
			return nil, fmt.Errorf("unable to cast %v to string", row[0])
		}
		n, ok := row[1].(int)
		if !ok {
			return nil, fmt.Errorf("couldn't cast %v to int", row[0])
		}
		c.ArticleViews[title] = articleData{
			Count:  n,
			rowNum: i,
		}
	}
	return c.ArticleViews, nil
}

func (c *Client) getMainTable() ([][]any, error) {
	mainTable, err := c.Service.Spreadsheets.Values.Get(c.creds.spreadsheetID, mainTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve main table data: %w", err)
	}
	return mainTable.Values, nil
}

func (c *Client) getPageData() ([][]any, error) {
	pageData, err := c.Service.Spreadsheets.Values.Get(c.creds.spreadsheetID, pageDataTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve page data table data: %w", err)
	}
	return pageData.Values, nil
}
