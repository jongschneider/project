package env

import (
	"fmt"
	"strings"
)

type App struct {
	Env Env `envconfig:"APP_ENV" required:"true" default:"local"`
}

// Env represents the environment the project is being run in
type Env string

const (
	Local         Env = "local"
	Development       = "dev"
	Staging           = "staging"
	Preproduction     = "preprod"
	Production        = "prod"
)

func (e *Env) Set(s string) error {
	switch v := Env(strings.ToLower(s)); v {
	case Local, Development, Staging, Preproduction, Production:
		*e = v
	case Env("production"):
		*e = Production
	default:
		return fmt.Errorf("invalid env: %s", s)
	}

	return nil
}
