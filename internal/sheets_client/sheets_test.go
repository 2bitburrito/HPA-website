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
	pageData, err := svc.getArticleData()
	if err != nil {
		t.Fatalf("unable to get page data: %v", err)
	}
	fmt.Println("main table: ", mainTable)
	fmt.Println("page data: ", pageData)

	allData, err := svc.batchGetAllSheetData()
	if err != nil {
		t.Fatalf("unable to get all data")
	}
	for i, v := range allData {
		fmt.Println("batchdata #", i, v.Values)
	}
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
				articleData{
					Title: "Article 1",
					Count: 10,
				},
				articleData{
					Count: 20,
					Title: "Article 2",
				},
				articleData{
					Title: "Article 3",
					Count: 30,
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
				service:      &sheets.Service{},
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
