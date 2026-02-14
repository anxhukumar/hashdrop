package cloudfrontguard

import (
	"testing"
)

func TestGenerateSignedCloudfrontURL(t *testing.T) {
	validKeyPath := writeTempPrivateKey(t)

	tests := []struct {
		name           string
		urlPrefix      string
		objectPath     string
		keyPairID      string
		privateKeyPath string
		wantErr        bool
	}{
		{
			name:           "Valid inputs",
			urlPrefix:      "https://example.cloudfront.net/",
			objectPath:     "file/test.txt",
			keyPairID:      "K1234567890",
			privateKeyPath: validKeyPath,
			wantErr:        false,
		},
		{
			name:           "Invalid private key path",
			urlPrefix:      "https://example.cloudfront.net/",
			objectPath:     "file/test.txt",
			keyPairID:      "K1234567890",
			privateKeyPath: "/no/such/key.pem",
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := GenerateSignedCloudfrontURL(
				tt.urlPrefix,
				tt.objectPath,
				tt.keyPairID,
				tt.privateKeyPath,
			)

			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateSignedCloudfrontURL() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if !tt.wantErr && url == "" {
				t.Fatalf("expected non-empty signed url")
			}
		})
	}
}
