package aws

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func newTestS3Client(t *testing.T) *s3.Client {
	t.Helper()

	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("test", "test", ""),
		),
	)
	if err != nil {
		t.Fatalf("failed to load test aws config: %v", err)
	}

	return s3.NewFromConfig(cfg)
}

func TestGeneratePresignedPUT(t *testing.T) {

	// Get dummy s3 client for testing
	s3Client := newTestS3Client(t)

	tests := []struct {
		name                string
		ctx                 context.Context
		client              *s3.Client
		presignedLinkExpiry time.Duration
		bucket              string
		key                 string
		wantErr             bool
	}{
		{
			name:                "Empty s3Client",
			ctx:                 context.Background(),
			client:              nil,
			presignedLinkExpiry: time.Minute * 15,
			bucket:              "test_bucket",
			key:                 "usrh-kds434/dfs321",
			wantErr:             true,
		},
		{
			name:                "Valid s3Client",
			ctx:                 context.Background(),
			client:              s3Client,
			presignedLinkExpiry: time.Minute * 15,
			bucket:              "test_bucket",
			key:                 "usrh-kds434/dfs321",
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s3Response, err := GeneratePresignedPUT(
				tt.ctx,
				tt.client,
				tt.presignedLinkExpiry,
				tt.bucket,
				tt.key,
			)

			// Test if we get errors when expected
			if (err != nil) != tt.wantErr {
				t.Fatalf("GeneratePresignedPUT() | error = %v, wantErr = %v", err, tt.wantErr)
			}

			// Test if url is non-empty
			if !tt.wantErr && len(s3Response.URL) < 1 {
				t.Error("GeneratePresignedPUT() returned empty URL")
			}

		})
	}
}
