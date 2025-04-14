package auth

import (
	"errors"
	"fmt"
	"ops_service/internal/model"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID int            `json:"user_id"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

type TokenService struct {
	signingKey []byte
	tokenTTL   time.Duration
}

func NewTokenService(signingKey string, tokenTTL time.Duration) *TokenService {
	return &TokenService{
		signingKey: []byte(signingKey),
		tokenTTL:   tokenTTL,
	}
}

func (s *TokenService) GenerateToken(userID int, role model.UserRole) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.signingKey)
}

func (s *TokenService) GenerateDummyToken(role model.UserRole) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: 0, // Dummy user ID
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.signingKey)
}

func (s *TokenService) ParseToken(accessToken string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		accessToken,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method")
			}
			return s.signingKey, nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, errors.New("token claims are not of expected type")
	}

	return claims, nil
}

func (s *TokenService) GetUserFromToken(accessToken string) (int, model.UserRole, error) {
	claims, err := s.ParseToken(accessToken)
	if err != nil {
		return 0, "", fmt.Errorf("failed to parse token: %w", err)
	}

	return claims.UserID, claims.Role, nil
}
