package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jongschneider/youtube-project/api/internal/platform/database"

	"github.com/jongschneider/youtube-project/api/cmd/server"
	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
	"github.com/jongschneider/youtube-project/api/internal/platform/config"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var cfg struct {
	config.Base
}

var log *logrus.Logger

func init() {
	// Set up the global logger for the api
	log = logrus.New()
	config.SetLogrusFormatter(log)
}

func main() {
	// Load in the configuration via a .env file
	err := config.Load(&cfg)
	if err != nil {
		log.WithError(err).Fatal("config: load")
	}

	// Connect to DB
	db := database.New(cfg.DBConfig)

	// Create a client
	client := clientSVC.New(
		clientSVC.LogConfigBlock{Logger: log},
		clientSVC.DBConfigBlock{DB: db},
	)

	// Create a new server with all of the routes attached to the server's handler
	srv := server.New(cfg.Port, client)

	// Gracefully handle shutdowns
	go shutdown(srv, time.Second*30)

	// Start the server listening for requests.
	log.Printf("listening on port%s", srv.Addr)
	err = srv.ListenAndServe()
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
