package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/jongschneider/youtube-project/api/internal/platform/auth"
	"github.com/jongschneider/youtube-project/api/internal/platform/cache"
	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/jongschneider/youtube-project/api/internal/platform/env"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

// Base holds the shared config used by each binary in this repo
type Base struct {
	AppConfig   env.App
	DBConfig    database.Config
	CacheConfig cache.Config
	AuthConfig  auth.Config
	Port        int  `envconfig:"PORT" required:"true" default:"3000"`
	Debug       bool `envconfig:"DEBUG" default:"false"`
}

// configurable is an internal interface to enforce this config as an embedded struct if another program wants to modify it.
type configurable interface {
	env() string
}

func (c Base) env() string {
	return string(c.AppConfig.Env)
}

// LogFields outputs all of the configuration values in debug mode, but less when not in debug mode.
func (c Base) LogFields() logrus.Fields {
	fields := logrus.Fields{
		"env":   c.AppConfig.Env,
		"debug": c.Debug,
		"port":  c.Port,
	}

	if c.Debug {
		fields["db_user"] = c.DBConfig.User
		fields["db_pass"] = c.DBConfig.Password
		fields["db_host"] = c.DBConfig.Host
		fields["db_port"] = c.DBConfig.Port
		fields["db_name"] = c.DBConfig.DBName
		fields["db_tls"] = c.DBConfig.TLS
		fields["db_multistatements"] = c.DBConfig.MultiStatements

		fields["redis_host"] = c.CacheConfig.Host
		fields["redis_port"] = c.CacheConfig.Port
		fields["redis_pass"] = c.CacheConfig.Password
		fields["redis_db"] = c.CacheConfig.DB
		fields["redis_tls"] = c.CacheConfig.TLS
	}

	return fields
}

// Load attempts to gather configuration information from the environment (and .env) and store it in the provided configurable.
func Load(c configurable) (err error) {
	err = godotenv.Load()
	if err != nil {
		logrus.Info(errors.Wrap(err, "godotenv"))
	}

	err = envconfig.Process("", c)
	if err != nil {
		logrus.Error(errors.Wrap(err, "envconfig: process"))
		return errors.Wrap(err, "envconfig: process")
	}

	return nil
}

// SetLogrusFormatter sets the formatter for the logrus logger.
func SetLogrusFormatter(l *logrus.Logger) {
	var formatter logrus.Formatter
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		formatter = &logrus.TextFormatter{}
	} else {
		formatter = &logrus.JSONFormatter{}
	}

	if l != nil {
		l.Formatter = formatter
	}
	logrus.SetFormatter(formatter)
}
