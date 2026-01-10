package sheetsclient

import (
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	Service *sheets.Service
	creds   credentials
}

type credentials struct {
	spreadsheetID      string
	serviceCredentials string
}

const (
	mainTableBounds     = "main_table!A1:Z"
	pageDataTableBounds = "article_data!A1:Z"
)
