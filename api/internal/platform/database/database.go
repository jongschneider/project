package database

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // provides the mysql driver for sqlx
	"github.com/jimmysawczuk/try"
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

	// We have to use the try.Try package to connect to the db because when we are setting up the
	// services to run in our dev environment, the api is ready before the db. When the api tries
	// to establish a connection to the db, it fails, because the db isn't running yet.
	// This results in a panic.
	// The try.Try package simply attempts to perform the code inside the function. If it errors,
	// it waits a duration - in this case 5 seconds - and then makes another attempt. It will keep
	// doing this until the timeout has elapsed - in this case 160 seconds.
	// This gives our db time to set up and the api a chance to make a connection.
	var db *sqlx.DB
	var err error
	if terr := try.Try(func() error {
		db, err = sqlx.Open(driverName, connectionString)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = pingDB(ctx, db)
		if err != nil {
			return err
		}
		return nil
	}, 160*time.Second, 5*time.Second); terr != nil {
		panic(terr)
	}

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
