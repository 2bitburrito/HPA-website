// Package setup is the dependency struct
package setup

import (
	"context"
	"fmt"
	"os"

	"github.com/2bitburrito/hpa-website/internal/blog"
	sheetsclient "github.com/2bitburrito/hpa-website/internal/sheets_client"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Dependencies struct {
	Aws           aws.Config
	Blogs         blog.Blogs
	SheetsService *sheetsclient.Client
}

func Setup() (Dependencies, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-southeast-2"),
	)
	if err != nil {
		return Dependencies{}, err
	}

	blogs, err := blog.ReadBlogData()
	if err != nil {
		return Dependencies{}, err
	}

	spreadsheetID := os.Getenv("GOOGLE_SPREADSHEET_ID")
	serviceCredentials := os.Getenv("GOOGLE_SERVICE_CREDENTIALS")
	if spreadsheetID == "" {
		return Dependencies{}, fmt.Errorf("google sheets spreadsheetID not set")
	}
	if serviceCredentials == "" {
		return Dependencies{}, fmt.Errorf("google sheets credentials not set")
	}
	sheetsSvc, err := sheetsclient.CreateSheetsService(spreadsheetID, serviceCredentials)
	if err != nil {
		return Dependencies{}, fmt.Errorf("unable to create sheets service: %w", err)
	}
	err = sheetsSvc.GetAllData(blogs)
	if err != nil {
		return Dependencies{}, fmt.Errorf("unable to get all data from sheets: %w", err)
	}

	return Dependencies{
		Aws:           awsCfg,
		Blogs:         blogs,
		SheetsService: sheetsSvc,
	}, nil
}
