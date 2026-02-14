package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in a JWT token
type JWTClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// JWTManager handles JWT token generation and validation
type JWTManager struct {
	secret []byte
	expiry time.Duration
}

// NewJWTManager creates a new JWT manager
func NewJWTManager(secret string, expiry time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		expiry: expiry,
	}
}

// Generate creates a new JWT token for a username
func (m *JWTManager) Generate(username string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.expiry)

	claims := JWTClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}

// Validate validates a JWT token and returns the claims
func (m *JWTManager) Validate(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Refresh generates a new token for an existing valid token
func (m *JWTManager) Refresh(tokenString string) (string, time.Time, error) {
	claims, err := m.Validate(tokenString)
	if err != nil {
		return "", time.Time{}, err
	}

	return m.Generate(claims.Username)
}
