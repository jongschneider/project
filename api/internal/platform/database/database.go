package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // provides the mysql driver for sqlx
	"github.com/jmoiron/sqlx"
)

var driverName = "mysql"

// Config holds all of the configuration for a database connection
type Config struct {
	User            string `envconfig:"MYSQL_USER" required:"true" default:"root"`
	Password        string `envconfig:"MYSQL_PASSWORD" default:""`
	Host            string `envconfig:"MYSQL_HOST" required:"true" default:"localhost"`
	Port            int    `envconfig:"MYSQL_PORT" required:"true" default:"3306"`
	DBName          string `envconfig:"MYSQL_DBNAME" required:"true" default:"example"`
	TLS             bool   `envconfig:"MYSQL_TLS" required:"true" default:"false"`
	MultiStatements bool   `envconfig:"MYSQL_MULTISTATEMENTS" required:"true" default:"true"`
}

// DB represents a db connection
type DB struct {
	*sqlx.DB
}

// New returns a new db connection
func New(cfg Config) *DB {
	connectionString := getConnectionString(cfg)

	db := sqlx.MustOpen(driverName, connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	mustPingDB(ctx, db)

	return &DB{db}
}

func getConnectionString(cfg Config) string {
	if cfg.Port == 0 {
		cfg.Port = 3306
	}

	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=%t&multiStatements=%t",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.TLS,
		cfg.MultiStatements,
	)
}

func pingDB(ctx context.Context, db *sqlx.DB) error {
	return db.PingContext(ctx)
}

// mustPingDB is the same a pingDB but it panics on error.
func mustPingDB(ctx context.Context, db *sqlx.DB) {
	err := pingDB(ctx, db)
	if err != nil {
		panic(err)
	}
}
