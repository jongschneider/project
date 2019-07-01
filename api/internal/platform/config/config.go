package config

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/jongschneider/youtube-project/api/internal/platform/env"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"
)

// Base holds the shared config used by each binary in this repo
type Base struct {
	AppConfig env.App
	DBConfig  database.Config
	Port      int `envconfig:"PORT" required:"true" default:"3000"`
}

// configurable is an internal interface to enforce this config as an embedded struct if another program wants to modify it.
type configurable interface {
	env() string
}

func (c Base) env() string {
	return string(c.AppConfig.Env)
}

// Load attempts to gather configuration information from the environment (and .env) and store it in the provided configurable
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
