package client

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jongschneider/youtube-project/api/internal/platform/config"
	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/pkg/errors"
)

// Client is an object that holds anything that might be necessary in various services.
type Client struct {
	tz  *time.Location
	key string
	db  *database.DB

	log *logrus.Logger
}

// DB returns the Client's db connection
func (c *Client) DB() *database.DB {
	return c.db
}

// TZ returns the Client's timezone
func (c *Client) TZ() *time.Location {
	return c.tz
}

// Key returns the Client's key
func (c *Client) Key() string {
	return c.key
}

// Log returns the Client's logger
func (c *Client) Log() *logrus.Logger {
	if c.log == nil {
		c.log = logrus.New()
		config.SetLogrusFormatter(c.log)
	}

	return c.log
}

// Config configures a new *Client
type Config struct {
	DB  *database.DB
	Log *logrus.Logger
	Key string
}

// New returns a new Client
func New(cfg Config) *Client {
	c := Client{
		db:  cfg.DB,
		log: cfg.Log,
		key: cfg.Key,
	}

	var err error
	c.tz, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(errors.Wrap(err, "load location"))
	}

	return &c
}
