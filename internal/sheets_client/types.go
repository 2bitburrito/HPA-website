package sheetsclient

import (
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	service      *sheets.Service
	creds        credentials
	MainData     MainData
	ArticleViews ArticleViewCounts
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
	Title string
	Count int
}
