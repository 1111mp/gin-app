package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	__defaultExpire = 24 * time.Hour
	__deaultIssuer  = "gin-app"
)

// JWTManagerInterface -.
type JWTManagerInterface interface {
	GenerateToken(userId int) (string, error)
	ParseToken(t string) (*Claims, error)
}

// JWTManager -.
type JWTManager struct {
	expire time.Duration
	issuer string
	secret []byte
}

// NewJWTManager -.
func NewJWTManager(opts ...Option) *JWTManager {
	m := &JWTManager{
		expire: __defaultExpire,
		issuer: __deaultIssuer,
		secret: nil,
	}

	// Custom options
	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Claims -.
type Claims struct {
	UserId int `json:"userId"`
	jwt.RegisteredClaims
}

// GenerateToken -.
func (m *JWTManager) GenerateToken(userId int) (string, error) {
	claims := &Claims{
		userId,
		jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expire)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.secret)
}

// ParseToken -.
func (m *JWTManager) ParseToken(t string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(t, &Claims{}, func(t *jwt.Token) (any, error) {
		return m.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}
