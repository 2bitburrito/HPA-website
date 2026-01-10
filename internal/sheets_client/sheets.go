package sheetsclient

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func CreateSheetsService(spreadsheetID, serviceCredentials string) (*Client, error) {
	scopes := []string{
		"https://www.googleapis.com/auth/spreadsheets.readonly",
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

func (s *Client) getMainTable() ([][]any, error) {
	mainTable, err := s.Service.Spreadsheets.Values.Get(s.creds.spreadsheetID, mainTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve main table data: %w", err)
	}
	return mainTable.Values, nil
}

func (s *Client) getPageData() ([][]any, error) {
	pageData, err := s.Service.Spreadsheets.Values.Get(s.creds.spreadsheetID, pageDataTableBounds).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve page data table data: %w", err)
	}
	return pageData.Values, nil
}
