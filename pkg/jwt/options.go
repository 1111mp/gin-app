package jwt

import "time"

// Option -.
type Option func(*JWTManager)

// Expire -.
func Expire(expire time.Duration) Option {
	return func(j *JWTManager) {
		j.expire = expire
	}
}

// Issuer -.
func Issuer(issuer string) Option {
	return func(j *JWTManager) {
		j.issuer = issuer
	}
}

// Secret -.
func Secret(secret string) Option {
	return func(j *JWTManager) {
		j.secret = []byte(secret)
	}
}
