package handler

import (
	"net/http"
	"os"
	"path/filepath"

	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/go-redis/redis"
	"github.com/jongschneider/youtube-project/api/internal/platform/auth"
	"github.com/jongschneider/youtube-project/api/internal/platform/database"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Handler is an object that holds anything that might be necessary in various services.
type Handler struct {
	tz    *time.Location
	db    *database.DB
	cache *redis.Client
	log   *logrus.Logger
	auth  *auth.Service
	http.Handler
}

// Config configures a new *Handler
type Config struct {
	DB    *database.DB
	Cache *redis.Client
	Auth  *auth.Service
	Log   *logrus.Logger
	Key   string
}

// New returns a new Handler
func New(cfg Config) *Handler {
	h := Handler{
		db:    cfg.DB,
		cache: cfg.Cache,
		auth:  cfg.Auth,
		log:   cfg.Log,
	}

	var err error
	h.tz, err = time.LoadLocation("America/New_York")
	if err != nil {
		panic(errors.Wrap(err, "load location"))
	}

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

	r.Get("/health", h.Health)

	// Set up a static file server
	workDir, err := os.Getwd()
	if err != nil {
		h.log.WithError(err).Fatal()
	}
	filesDir := filepath.Join(workDir, "static")
	h.FileServer(r, "/static", http.Dir(filesDir))

	r.Get("/token", h.auth.IssueTokenHandler)

	r.Route("/auth", func(r chi.Router) {
		r.Use(h.auth.RequireValidToken)
		r.Mount("/user", h.userRouter())
	})

	h.Handler = r

	return &h
}

// userRouter is an example of how to create a subrouter used for versioning
func (h *Handler) userRouter() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.Create)
	r.Get("/", h.GetAllUsers)
	r.Get("/{ID}", h.GetUser)
	r.Delete("/{ID}", h.Delete)
	r.Put("/{ID}", h.Update)
	r.Post("/login", h.Login)
	return r
}
