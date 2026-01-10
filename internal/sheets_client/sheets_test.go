package sheetsclient

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetSheetsPages(t *testing.T) {
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
