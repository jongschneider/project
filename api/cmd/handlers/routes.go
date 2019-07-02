package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	clientSVC "github.com/jongschneider/youtube-project/api/internal/platform/client"
)

// Router creates a new Router with all of our routes attached
func Router(client *clientSVC.Client) http.Handler {
	log := client.Log()

	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
		MaxAge:         1000,
	}).Handler)
	r.Use(middleware.DefaultCompress)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, "Project API")
	})

	r.Get("/health", Health(client))

	// Set up a static file server
	workDir, err := os.Getwd()
	if err != nil {
		log.WithError(err).Fatal()
	}
	filesDir := filepath.Join(workDir, "static")
	FileServer(r, "/static", http.Dir(filesDir))

	r.Route("/auth", func(r chi.Router) {
		// r.Use(mid.Auth)
		r.Mount("/login", Login(client))
	})

	return r
}

// articleRouter is an example of how to create a subrouter used for versioning
func articleRouter() http.Handler {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.Respond(w, r, "hello")
	})
	r.Route("/{articleID}", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			render.Respond(w, r, fmt.Sprintf("article %s", chi.URLParam(r, "articleID")))
		})
	})
	return r
}
