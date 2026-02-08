package otp

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const fromAddr = "noreply@hashdrop.dev"

func SendOTP(ctx context.Context, toEmail, otp string, sesClient *sesv2.Client) error {
	subject := "Your Hashdrop OTP"
	bodyText := fmt.Sprintf(
		"Hi,\n\n"+
			"Your one-time password (OTP) for verifying your Hashdrop account is:\n\n"+
			"%s\n\n"+
			"This code is valid for 10 minutes and can be used only once.\n\n"+
			"You’re receiving this email because someone (hopefully you) just tried to create or verify a Hashdrop account using this email address. "+
			"We don’t send marketing emails or newsletters. You’ll only hear from us for essential account actions like verification and security alerts.\n\n"+
			"If you didn’t request this, you can safely ignore this email—no account will be created without the correct OTP.\n\n"+
			"Thanks,\n"+
			"Hashdrop Team\n",
		otp,
	)

	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(fromAddr),
		Destination: &types.Destination{
			ToAddresses: []string{toEmail},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
				Body: &types.Body{
					Text: &types.Content{
						Data:    aws.String(bodyText),
						Charset: aws.String("UTF-8"),
					},
				},
			},
		},
	}

	_, err := sesClient.SendEmail(ctx, input)
	return err
}
