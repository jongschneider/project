package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-redis/redis"
	"github.com/jongschneider/youtube-project/api/cmd/handler"
	"github.com/jongschneider/youtube-project/api/internal/platform/auth"
	"github.com/jongschneider/youtube-project/api/internal/platform/cache"
	"github.com/jongschneider/youtube-project/api/internal/platform/config"
	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var cfg config.Base
var log *logrus.Logger

func init() {
	// Set up the global logger for the api
	log = logrus.New()
	config.SetLogrusFormatter(log)

	// Load in the configuration via a .env file
	err := config.Load(&cfg)
	if err != nil {
		log.WithError(err).Fatal("config: load")
	}

	if !cfg.Debug {
		log.SetFormatter(&logrus.JSONFormatter{})
	}

	log.WithFields(cfg.LogFields()).Info("Startup Config")

}

func main() {

	db := database.New(cfg.DBConfig)
	cacheSVC := cache.New(cfg.CacheConfig)
	authSVC := getAuthClient(cacheSVC)
	// Create a handler
	h := handler.New(
		handler.Config{
			DB:    db,
			Cache: cacheSVC,
			Auth:  authSVC,
			Log:   log,
		})

	// Create a new server with all of the routes attached to the server's handler
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Port),
		Handler:        h.Handler,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
		TLSConfig: &tls.Config{
			Certificates:       []tls.Certificate{mustLoadCert()},
			InsecureSkipVerify: cfg.Debug, // when in debug mode (working locally) skip verification of the certs
		},
	}

	// Gracefully handle shutdowns
	go shutdown(srv, time.Second*30)

	// Start the server listening for requests.
	log.Printf("listening on port%s", srv.Addr)
	err := srv.ListenAndServeTLS("", "")
	if err != nil && err != http.ErrServerClosed {
		log.Fatalln(errors.Wrap(err, "start server"))
	}

}

// shutdown handles graceful shutdowns
func shutdown(srv *http.Server, timeout time.Duration) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info(errors.Wrap(err, "shutdown server"))
	}
}

func getAuthClient(c *redis.Client) *auth.Service {
	return auth.New(auth.Config{
		Issuer:            cfg.AuthConfig.Issuer,
		PrivateKey:        mustLoadAuthKey(),
		Enforce:           cfg.AuthConfig.Enforce,
		RequestValidators: []auth.RequestValidator{},
		AbortRequest: func(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
			log.WithError(err).WithFields(logrus.Fields{
				"statuscode": statusCode,
				"enforcing":  true,
			}).Info("auth: abort")

			web.RespondWithCodedError(w, r, statusCode, http.StatusText(statusCode), err)
		},
		ContinueRequest: func(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
			log.WithError(err).WithFields(logrus.Fields{
				"statuscode": statusCode,
				"enforcing":  false,
			}).Info("auth: continue")

		},
		TokenBlocked: func(r *http.Request, err error, statusCode int) {
			log.WithError(err).WithField("statuscode", statusCode).Info("token not granted")
		},
		Cache: c,
	})
}
