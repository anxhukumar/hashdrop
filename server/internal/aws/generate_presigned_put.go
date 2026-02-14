package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3PutResponse struct {
	URL string `json:"url"`
}

func GeneratePresignedPUT(
	ctx context.Context,
	client *s3.Client,
	presignedLinkExpiry time.Duration,
	bucket string,
	key string,
) (*S3PutResponse, error) {

	if client == nil {
		return nil, fmt.Errorf("s3 client is nil")
	}

	presigner := s3.NewPresignClient(client)

	req, err := presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(presignedLinkExpiry))

	if err != nil {
		return nil, err
	}

	return &S3PutResponse{
		URL: req.URL,
	}, nil
}
