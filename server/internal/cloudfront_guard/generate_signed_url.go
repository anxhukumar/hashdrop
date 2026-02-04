package cloudfrontguard

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
)

const signedURLduration = 2 // minutes

func GenerateSignedCloudfrontURL(cloudfrontURLPrefix, objectPath, cloudfrontKeyPairID, privateKeyPath string) (string, error) {
	privateKey, err := loadPrivateKey(privateKeyPath)
	if err != nil {
		return "", fmt.Errorf("error while loading private key: %w", err)
	}

	signer := sign.NewURLSigner(cloudfrontKeyPairID, privateKey)

	signedURL, err := signer.Sign(
		cloudfrontURLPrefix+objectPath,
		time.Now().Add(signedURLduration*time.Minute),
	)
	if err != nil {
		return "", fmt.Errorf("error generating signedURL: %w", err)
	}

	return signedURL, nil
}
