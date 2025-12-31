package redis

// Option -.
type Option func(*Redis)

// MaxPoolSize -.
func MaxPoolSize(size int) Option {
	return func(c *Redis) {
		c.maxPoolSize = size
	}
}
