package pkg

import (
	"book-lending-api/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims defines the custom claims carried inside tokens issued by
// this application.  It embeds jwt.RegisteredClaims to get standard
// fields like expiry and issue time.
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTUtil encapsulates the secret used to sign and validate tokens.
type JWTUtil struct {
	secret string
}

// NewJWTUtil returns a new JWT utility with the given secret.
func NewJWTUtil(secret string) *JWTUtil {
	return &JWTUtil{secret: secret}
}

// GenerateToken creates a signed JWT for the provided user.  Tokens
// expire after 24 hours.
func (j *JWTUtil) GenerateToken(user *domain.User) (string, error) {
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

// ValidateToken parses and validates the given token string.  It
// returns the claims if valid or an error otherwise.
func (j *JWTUtil) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
