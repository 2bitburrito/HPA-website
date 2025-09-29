// Package email sends an email
package email

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type SendEmailParams struct {
	SendingAddress string
	SenderName     string
	Message        string
}

type EmailSender interface {
	SendEmail(params SendEmailParams) error
}

func SendContactEmail(awsCfg aws.Config, params SendEmailParams) error {
	if params.SendingAddress == "" {
		return errors.New("missing required email in send contact email")
	}

	client := sesv2.NewFromConfig(awsCfg)
	fromAddress := "mail@hughpalmer.com.au"

	sesEmailParams := sesv2.SendEmailInput{
		FromEmailAddress: &fromAddress,
		ReplyToAddresses: []string{params.SendingAddress},
		Destination: &types.Destination{
			ToAddresses: []string{"hughpalmerproduction@gmail.com"},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(fmt.Sprintf("New Email From: %v", params.SenderName)),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data: &params.Message,
					},
				},
			},
		},
	}
	// returns a MessageId and result metadata
	_, err := client.SendEmail(context.Background(), &sesEmailParams)
	if err != nil {
		return err
	}
	return nil
}
