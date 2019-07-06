package cache

import (
	"crypto/tls"
	"fmt"

	"github.com/go-redis/redis"
)

// Config holds all of the configuration for a redes connection
type Config struct {
	Host     string `envconfig:"REDIS_HOST" required:"true" default:"localhost"`
	Port     int    `envconfig:"REDIS_PORT" required:"true" default:"6379"`
	Password string `envconfig:"REDIS_PASSWORD"`
	DB       int    `envconfig:"REDIS_DB"`
	TLS      string `envconfig:"REDIS_TLS"`
}

// New returns a new redis connection
func New(cfg Config) *redis.Client {
	opts := &redis.Options{
		Addr:     getConnectionString(cfg),
		Password: cfg.Password,
		DB:       cfg.DB,
	}

	if cfg.TLS != "" {
		opts.TLSConfig = &tls.Config{
			ServerName: cfg.TLS,
		}
	}

	r := redis.NewClient(opts)

	mustPingRedis(r)

	return r
}

func getConnectionString(cfg Config) string {
	return fmt.Sprintf("%s:%d",
		cfg.Host,
		cfg.Port,
	)
}

func pingRedis(c *redis.Client) error {
	return c.Ping().Err()
}

// mustPingDB is the same a pingDB but it panics on error.
func mustPingRedis(c *redis.Client) {
	err := pingRedis(c)
	if err != nil {
		panic(err)
	}
}
