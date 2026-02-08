package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
)

// Initialize AWS and return generalAwsConfig, s3Client, sesClient
func InitAWS(ctx context.Context, region string) (aws.Config, *s3.Client, *sesv2.Client, error) {
	cfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(region))
	if err != nil {
		return aws.Config{}, nil, nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	sesClient := sesv2.NewFromConfig(cfg)

	return cfg, s3Client, sesClient, nil
}
