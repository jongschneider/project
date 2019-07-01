package client

import (
	"time"

	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/pkg/errors"
)

// Client is an object that holds anything that might be necessary in various services.
type Client struct {
	db *database.DB
	tz *time.Location
}

// DB returns the Client's db connection
func (c *Client) DB() *database.DB {
	return c.db
}

// TZ returns the Client's timezone
func (c *Client) TZ() *time.Location {
	return c.tz
}

// New returns a new Client
func New(configBlocks ...ConfigBlock) *Client {
	c := &Client{}

	var err error
	c.tz, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(errors.Wrap(err, "load location"))
	}

	for _, cb := range configBlocks {
		if err := cb.Configure(c); err != nil {
			panic(errors.Wrapf(err, "config block: %s", cb))
		}
	}

	return c
}

// ConfigBlock is an interface used to configure a New Client
// example: adding a DB connection to the client
type ConfigBlock interface {
	Configure(*Client) error
	String() string
}

// DBConfigBlock adds a DB to a Client
type DBConfigBlock struct {
	*database.DB
}

func (d DBConfigBlock) Configure(c *Client) error {
	c.db = d.DB
	return nil
}

func (d DBConfigBlock) String() string {
	return "DBConfigBlock"
}
