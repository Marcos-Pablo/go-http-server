package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	})

	signedStr, err := token.SignedString([]byte(tokenSecret))

	if err != nil {
		return "", fmt.Errorf("Couldn't create jwt token: %w", err)
	}

	return signedStr, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Couldn't validate JWT: %w", err)
	}

	subject, err := token.Claims.GetSubject()

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Couldn't validate JWT: %w", err)
	}

	userID, err := uuid.Parse(subject)

	if err != nil {
		return uuid.UUID{}, fmt.Errorf("Couldn't validate JWT: %w", err)
	}
	return userID, nil
}
