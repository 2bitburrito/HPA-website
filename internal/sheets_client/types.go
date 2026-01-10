package sheetsclient

import (
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	Service      *sheets.Service
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
type ArticleViewCounts map[string]articleData

type articleData struct {
	Count  int
	rowNum int
}
