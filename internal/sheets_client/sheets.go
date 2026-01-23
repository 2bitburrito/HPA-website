// Package sheetsclient provides a client for interacting with the Google Sheets API
// It acts as both a cache and a fetching service to pull from the API
//
// There are some bad seperation of concerns here, but it isn't ever scaling so it is what it is
package sheetsclient

import (
	"context"
	"fmt"
	"time"

	"github.com/2bitburrito/hpa-website/internal/blog"
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
		service: srv,
		creds: credentials{
			spreadsheetID:      spreadsheetID,
			serviceCredentials: serviceCredentials,
		},
	}, nil
}

// GetAllData retrieves all data from the sheets API and stores it in the client
// It is the main call after creating the client to get fresh data from the sheets API
func (c *Client) GetAllData(blgs blog.Blogs) error {
	allRanges, err := c.batchGetAllSheetData()
	if err != nil {
		return fmt.Errorf("unable to retrieve all sheet data: %w", err)
	}

	// The first slice will be the main table
	c.MainData, err = c.extractMainTableFromData(allRanges[0].Values)
	if err != nil {
		return fmt.Errorf("unable to extract main table data: %w", err)
	}

	// The second slice will be the article data
	c.ArticleViews, err = c.extractArticlesFromData(allRanges[1].Values)
	if err != nil {
		return fmt.Errorf("unable to extract article data: %w", err)
	}

	c.EnsureAllBlogsExist(blgs)

	return nil
}

// GetMainData returns the main table data for the site
// It ignores any caching and always retrieves the data from the sheets API
func (c *Client) GetMainData() (MainData, error) {
	t, err := c.getMainTable()
	if err != nil {
		return MainData{}, err
	}

	return c.extractMainTableFromData(t)
}

func (c *Client) GetAllArticleViews() (ArticleViewCounts, error) {
	pd, err := c.getArticleData()
	if err != nil {
		return ArticleViewCounts{}, err
	}
	return c.extractArticlesFromData(pd)
}

func (c *Client) FlushToSheets() error {
	md := c.restructureMainTable()
	pd := c.restructureArticleData()

	br := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW",
		Data: []*sheets.ValueRange{
			{
				Range:  mainTableBounds,
				Values: md,
			},
			{
				Range:  pageDataTableBounds,
				Values: pd,
			},
		},
	}
	_, err := c.service.Spreadsheets.Values.BatchUpdate(c.creds.spreadsheetID, br).Do()
	return err
}

// extractArticlesFromData takes raw matrix data from page_data and extracts the article titles,
// view counts and row numbers
func (c *Client) extractArticlesFromData(pageData [][]any) (ArticleViewCounts, error) {
	var articleViews []articleData
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
		articleViews = append(articleViews, articleData{
			Title: title,
			Count: n,
		})
	}
	c.ArticleViews = articleViews
	return c.ArticleViews, nil
}

// extractMainTableFromData takes raw matrix data from main_table and extracts the view count
//
// It is currently just extracting the home page view count as that is all we are storing in the main table
func (c *Client) extractMainTableFromData(mainTable [][]any) (MainData, error) {
	views, ok := mainTable[1][1].(int)
	if !ok {
		return MainData{}, fmt.Errorf("unable to convert main table view count to int")
	}
	c.MainData = MainData{
		HomePageViewCount: views,
	}
	return c.MainData, nil
}

func (c *Client) getMainTable() ([][]any, error) {
	mainTable, err := c.service.Spreadsheets.Values.Get(c.creds.spreadsheetID, mainTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve main table data: %w", err)
	}
	return mainTable.Values, nil
}

func (c *Client) getArticleData() ([][]any, error) {
	pageData, err := c.service.Spreadsheets.Values.Get(c.creds.spreadsheetID, pageDataTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve page data table data: %w", err)
	}
	return pageData.Values, nil
}

func (c *Client) batchGetAllSheetData() ([]*sheets.ValueRange, error) {
	dat, err := c.service.Spreadsheets.Values.BatchGet(c.creds.spreadsheetID).Ranges(mainTableBounds, pageDataTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve all sheet data: %w", err)
	}
	return dat.ValueRanges, nil
}

// restructureMainTable restructures the main table data to be in the format expected by the sheets API
func (c *Client) restructureMainTable() [][]any {
	n := c.MainData.HomePageViewCount
	t := [][]any{
		{"homepage_views"},
		{n},
	}
	return t
}

func (c *Client) restructureArticleData() [][]any {
	d := [][]any{{"blog_title", "views"}}

	for _, v := range c.ArticleViews {
		d = append(d, []any{v.Title, v.Count})
	}

	return d
}

func (c *Client) EnsureAllBlogsExist(blgs blog.Blogs) error {
	missing := findMissingBlogs(blgs, c.ArticleViews)

	err := c.addNewArticleRows(missing)
	if err != nil {
		return err
	}
	return nil
}

func findMissingBlogs(nuBlgs blog.Blogs, exstBlgs ArticleViewCounts) [][]any {
	var doesntExist [][]any

	// Loop through the blog data and check agains all existing rows returned from the sheets API
	// If it doesn't exist in the sheet already then return it to be appended
	for _, nuBlg := range nuBlgs {
		exists := false
		for _, exstBlg := range exstBlgs {
			if nuBlg.Title == exstBlg.Title {
				exists = true
				break
			}
		}
		if !exists {
			doesntExist = append(doesntExist, []any{nuBlg.Title, 0})
		}
	}

	return doesntExist
}

func (c *Client) addNewArticleRows(missing [][]any) error {
	if len(missing) == 0 {
		return nil
	}
	_, err := c.service.Spreadsheets.Values.Append(
		c.creds.spreadsheetID,
		pageDataTableBounds,
		&sheets.ValueRange{
			Values: missing,
		}).ValueInputOption("RAW").Do()
	return err
}
