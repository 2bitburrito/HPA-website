package sheetsclient

import (
	"context"

	"google.golang.org/api/sheets/v4"
)

type Client struct {
	service       *sheets.Service
	creds         credentials
	MainData      MainData
	ArticleViews  ArticleViewCounts
	cancelFlushes context.CancelFunc
}

type credentials struct {
	spreadsheetID      string
	serviceCredentials string
}

type MainData struct {
	HomePageViewCount int
}
type ArticleViewCounts []articleData

type articleData struct {
	Name  string // This is the title of the blog in kebab case (it matches the file name)
	Count int
}
