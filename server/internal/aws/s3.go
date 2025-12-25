package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Initialize AWS client, fetch configurations and set them to server struct
func InitS3(ctx context.Context, region string) (aws.Config, *s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return aws.Config{}, nil, err
	}

	client := s3.NewFromConfig(cfg)

	return cfg, client, nil
}
