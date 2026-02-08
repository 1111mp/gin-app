package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

//go:generate mockgen -source=jwt.go -destination=../../internal/modules/user/mocks_jwt_test.go -package=user_test

const (
	__defaultExpire = 24 * time.Hour
	__deaultIssuer  = "gin-app"
)

// JWT -.
type JWT interface {
	GenerateToken(userId int) (string, error)
	ParseToken(t string) (*Claims, error)
}

// jwtImpl -.
type jwtImpl struct {
	expire time.Duration
	issuer string
	secret []byte
}

// NewJWT -.
func NewJWT(opts ...Option) JWT {
	m := &jwtImpl{
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
func (m *jwtImpl) GenerateToken(userId int) (string, error) {
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
func (m *jwtImpl) ParseToken(t string) (*Claims, error) {
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
