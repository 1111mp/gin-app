package jwt

import "time"

// Option -.
type Option func(*jwtImpl)

// Expire -.
func Expire(expire time.Duration) Option {
	return func(j *jwtImpl) {
		j.expire = expire
	}
}

// Issuer -.
func Issuer(issuer string) Option {
	return func(j *jwtImpl) {
		j.issuer = issuer
	}
}

// Secret -.
func Secret(secret string) Option {
	return func(j *jwtImpl) {
		j.secret = []byte(secret)
	}
}
