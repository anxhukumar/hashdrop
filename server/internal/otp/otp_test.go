package otp

import "testing"

func TestOtpGeneration(t *testing.T) {

	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Generate 6 digit otp",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			otp, err := GenerateOTP()

			// Test if we get errors when expected
			if (err != nil) != tt.wantErr {
				t.Fatalf("GenerateOTP() | error = %v, wantErr = %v", err, tt.wantErr)
			}

			// Test if we get 6 digit otp
			if len(otp) != 6 {
				t.Fatal("GenerateOTP() returned less than 6 digit otp")
			}

			// Check if it only returns digits
			for _, r := range otp {
				if r < '0' || r > '9' {
					t.Fatalf("GenerateOTP() returned non-digit character: %q", r)
				}
			}

		})
	}
}

func TestHashAndVerifyOTP(t *testing.T) {
	otp := "123456"
	secret := "super-secret-key"

	hash := HashOTP(otp, secret)

	tests := []struct {
		name       string
		inputOTP   string
		secret     string
		storedHash string
		wantOK     bool
	}{
		{
			name:       "Correct OTP and secret",
			inputOTP:   otp,
			secret:     secret,
			storedHash: hash,
			wantOK:     true,
		},
		{
			name:       "Wrong OTP",
			inputOTP:   "654321",
			secret:     secret,
			storedHash: hash,
			wantOK:     false,
		},
		{
			name:       "Wrong secret",
			inputOTP:   otp,
			secret:     "wrong-secret",
			storedHash: hash,
			wantOK:     false,
		},
		{
			name:       "Tampered hash",
			inputOTP:   otp,
			secret:     secret,
			storedHash: hash + "00",
			wantOK:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ok := VerifyOTP(tt.inputOTP, tt.storedHash, tt.secret)
			if ok != tt.wantOK {
				t.Fatalf("VerifyOTP() = %v, want %v", ok, tt.wantOK)
			}
		})
	}

	// Also test determinism of HashOTP
	hash2 := HashOTP(otp, secret)
	if hash != hash2 {
		t.Fatalf("HashOTP() is not deterministic: %s != %s", hash, hash2)
	}
}
