package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	userID := uuid.New()
	validTokenSecret := "super-secret-token"

	tests := []struct {
		name             string
		userID           uuid.UUID
		signingSecret    string
		validationSecret string
		expiresIn        time.Duration
		matchUserID      bool
		wantErr          bool
	}{
		{
			name:             "Valid token",
			userID:           userID,
			signingSecret:    validTokenSecret,
			validationSecret: validTokenSecret,
			expiresIn:        time.Second * 5,
			wantErr:          false,
			matchUserID:      true,
		},
		{
			name:             "Expired token",
			userID:           userID,
			signingSecret:    validTokenSecret,
			validationSecret: validTokenSecret,
			expiresIn:        time.Millisecond * 5,
			wantErr:          true,
			matchUserID:      false,
		},
		{
			name:             "Invalid secret",
			userID:           userID,
			signingSecret:    validTokenSecret,
			validationSecret: "invalid-secret",
			expiresIn:        time.Second * 5,
			wantErr:          true,
			matchUserID:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := MakeJWT(tt.userID, tt.signingSecret, tt.expiresIn)
			if err != nil {
				t.Fatalf("Failed to generate JWT: %s", err)
			}

			userID, err := ValidateJWT(token, tt.validationSecret)
			match := tt.userID == userID

			if !tt.wantErr && !match {
				t.Errorf("Failed to validate JWT: %s", err)
			}
		})
	}
}
