package util

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labovector/vecsys-api/infrastructure/config"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func NewResetPasswordClaims(userId string) *Claims {
	return &Claims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "vecsys",
			Subject:   "user",
			Audience:  jwt.ClaimStrings{"user"},
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 5)),
		},
	}
}

type Maker struct {
	SecretKey string
}

// NewJWTMaker creates a new Maker instance with the provided configuration.
// It extracts the secret key from the provided JWTConfig and initializes
// the Maker with it, which is used for signing JWT tokens.

func NewJWTMaker(config *config.JWTConfig) *Maker {
	secretKey := config.SecretKey
	return &Maker{
		SecretKey: secretKey,
	}
}

// GenerateResetPasswordToken generates a JWT token for user to reset password.
//
// The token contains the user's ID and expires in 5 minutes.
//
// The function returns an error if it fails to generate the token.
func (maker *Maker) GenerateResetPasswordToken(userId string) (string, error) {
	claims := NewResetPasswordClaims(userId)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(maker.SecretKey))
}

// VerifyResetPasswordToken verifies the provided JWT token for password reset.
// It checks the token's signature and ensures it has not expired.
// Returns true if the token is valid, otherwise returns false.

func (maker *Maker) VerifyResetPasswordToken(token string) bool {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (any, error) {
		return []byte(maker.SecretKey), nil
	})
	if err != nil {
		return false
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return false
	}
	return true
}
