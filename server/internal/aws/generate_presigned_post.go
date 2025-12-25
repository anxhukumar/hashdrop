package aws

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type S3PostResponse struct {
	URL    string            `json:"url"`
	Fields map[string]string `json:"fields"`
}

func GeneratePresignedPOST(
	ctx context.Context,
	cfg aws.Config,
	minDataSize int64,
	maxDataSize int64,
	presignedLinkExpiry time.Duration,
	bucket,
	key,
	region string,
) (*S3PostResponse, error) {

	// Load credentials from config
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	dateShort := now.Format("20060102")
	dateLong := now.Format("20060102T150405Z")
	expiry := now.Add(presignedLinkExpiry).Format("2006-01-02T15:04:05Z")
	credential := fmt.Sprintf("%s/%s/%s/s3/aws4_request", creds.AccessKeyID, dateShort, region)

	// Define the Policy Document with Size Enforcement
	// Conditions are checked by S3 before the upload starts
	policy := map[string]interface{}{
		"expiration": expiry,
		"conditions": []interface{}{
			map[string]string{"bucket": bucket},
			map[string]string{"key": key},
			map[string]string{"x-amz-algorithm": "AWS4-HMAC-SHA256"},
			map[string]string{"x-amz-credential": credential},
			map[string]string{"x-amz-date": dateLong},
			// SOLID ENFORCEMENT: [condition, minBytes, maxBytes]
			[]interface{}{"content-length-range", minDataSize, maxDataSize},
		},
	}

	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return nil, err
	}

	policyBase64 := base64.StdEncoding.EncodeToString(policyBytes)

	// SigV4 Signing (Manual Implementation)
	signingKey := deriveSigningKey(creds.SecretAccessKey, dateShort, region, "s3")
	signature := hmacHex(signingKey, policyBase64)

	return &S3PostResponse{
		URL: fmt.Sprintf("https://%s.s3.%s.amazonaws.com/", bucket, region),
		Fields: map[string]string{
			"key":              key,
			"x-amz-algorithm":  "AWS4-HMAC-SHA256",
			"x-amz-credential": credential,
			"x-amz-date":       dateLong,
			"policy":           policyBase64,
			"x-amz-signature":  signature,
		},
	}, nil
}

// Crypto Helpers
func hmacSha256(key []byte, data string) []byte {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(data))
	return h.Sum(nil)
}

func hmacHex(key []byte, data string) string {
	return hex.EncodeToString(hmacSha256(key, data))
}

func deriveSigningKey(secret, date, region, service string) []byte {
	kDate := hmacSha256([]byte("AWS4"+secret), date)
	kRegion := hmacSha256(kDate, region)
	kService := hmacSha256(kRegion, service)
	return hmacSha256(kService, "aws4_request")
}
