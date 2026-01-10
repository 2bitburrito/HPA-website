package sheetsclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"google.golang.org/api/sheets/v4"
)

func TestGetSheetsPages(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping google sheets integration test in short mode")
	}
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatal("Error loading .env file")
	}

	sheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	if sheetID == "" {
		t.Error("skipping test; no google sheets id set")
	}
	serviceCredentials := os.Getenv("GOOGLE_SERVICE_CREDENTIALS")
	if serviceCredentials == "" {
		t.Error("skipping test; no google service credentials set")
	}

	svc, err := CreateSheetsService(sheetID, serviceCredentials)
	if err != nil {
		t.Fatalf("unable to create sheets service: %v", err)
	}
	mainTable, err := svc.getMainTable()
	if err != nil {
		t.Fatalf("unable to get main table: %v", err)
	}
	pageData, err := svc.getPageData()
	if err != nil {
		t.Fatalf("unable to get page data: %v", err)
	}
	fmt.Println(mainTable, pageData)
}

func TestExtractArticlesFromData(t *testing.T) {
	testCases := []struct {
		desc      string
		mtrx      [][]any
		rtnData   ArticleViewCounts
		shouldErr bool
	}{
		{
			desc: "Correctly extracts article titles and view counts",
			mtrx: [][]any{
				{"blog_title", "views"},
				{"Article 1", 10},
				{"Article 2", 20},
				{"Article 3", 30},
			},
			rtnData: ArticleViewCounts{
				"Article 1": articleData{
					Count:  10,
					rowNum: 1,
				},
				"Article 2": articleData{
					Count:  20,
					rowNum: 2,
				},
				"Article 3": articleData{
					Count:  30,
					rowNum: 3,
				},
			},
			shouldErr: false,
		},
		{
			desc: "Fails on incorrect types",
			mtrx: [][]any{
				{"blog_title", "views"},
				{"Article 1", "thisshouldbeanumber"},
				{"", 0},
			},
			rtnData:   nil,
			shouldErr: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			c := &Client{
				Service:      &sheets.Service{},
				creds:        credentials{},
				MainData:     MainData{},
				ArticleViews: ArticleViewCounts{},
			}
			data, err := c.extractArticlesFromData(tC.mtrx)
			if err != nil && !tC.shouldErr {
				t.Fatalf("unexpected error: %v", err)
			}
			for k, v := range tC.rtnData {
				if data[k] != v {
					t.Errorf("expected %v, got %v", v, data[k])
				}
			}
		})
	}
}
