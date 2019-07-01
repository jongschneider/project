package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jongschneider/youtube-project/api/cmd/handlers"
	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
	"github.com/sirupsen/logrus"
)

// Options alter the server in some way
// example: add tls
type Options func(*http.Server)

// New creates a new server
func New(port int, log *logrus.Logger, client *clientSVC.Client, opts ...Options) *http.Server {
	server := http.Server{
		Addr:           fmt.Sprintf(":%d", port),
		Handler:        handlers.Router(log, client),
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	for _, o := range opts {
		o(&server)
	}

	return &server
}
