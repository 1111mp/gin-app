package rediskey

import "strings"

const (
	__defaultPrefix = "prefix"
)

type RedisKey interface {
	Key(parts ...string) string
}

type redisKeyImpl struct {
	env    string
	prefix string
}

func New(prefix, env string) RedisKey {
	if prefix == "" {
		prefix = __defaultPrefix
	}

	return &redisKeyImpl{
		env:    env,
		prefix: prefix,
	}
}

func (k *redisKeyImpl) Key(parts ...string) string {
	all := []string{k.env, k.prefix}
	all = append(all, parts...)
	return strings.Join(all, ":")
}
