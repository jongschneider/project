package auth

import (
	"context"
	"crypto/rsa"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/jongschneider/youtube-project/api/internal/platform/web"
	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type authCtxKey string

var tokenKey authCtxKey = "token"

var (
	// ErrNotUsed is the error returned by a RequestValidator if it does not conclusively prove a valid request, but shouldn't necessarily disqualify the request.
	ErrNotUsed = errors.New("not used")

	// ErrNotAuthorized is the error returned when there are no valid requests.
	ErrNotAuthorized = errors.New("not authorized")

	// ErrGenerateToken is the error returned when we couldn't generate a token.
	ErrGenerateToken = errors.New("couldn't generate token")

	// ErrMissingToken is the error returned when there is no token.
	ErrMissingToken = errors.New("missing token")
)

// RequestValidator is a function that validates a request to see if it's valid for receiving a JWT.
// It should return a nil error if it is valid, a non-nil error if it's not valid, and the error ErrNotUsed if it isn't conclusive.
type RequestValidator func(*http.Request) error

// EnforceFunc is invoked when a request to a secured endpoint shouldn't be considered authorized.
type EnforceFunc func(w http.ResponseWriter, r *http.Request, err error, statusCode int)

// TokenBlockedFunc is invoked when a request for a token is not authorized.
type TokenBlockedFunc func(r *http.Request, err error, statusCode int)

// Service holds all of the setup for an authentication service
type Service struct {
	requestValidators []RequestValidator
	ttl               time.Duration
	enforce           bool
	tz                *time.Location
	cache             *redis.Client
	privateKey        *rsa.PrivateKey
	publicKey         *rsa.PublicKey
	abortRequest      EnforceFunc
	continueRequest   EnforceFunc
	tokenBlocked      TokenBlockedFunc
	keyPrefix         string
	issuer            string
}

// Config holds all of the configuration needed to create an authentication service
type Config struct {
	// Issuer is the name of the issuer of the token
	Issuer string `envconfig:"AUTH_ISSUER" default:"youtube-project"`

	// A valid *time.Location used for logging (all timestamps should be UTC internally)
	TZ *time.Location

	// A PEM-encoded RSA private key
	PrivateKey string

	// Enforce should be true if the auth service will actually reject requests that are invalid/unautenticated.
	// If false, these requests will be logged and passed through.
	Enforce bool

	// RequestValidators are functions that return an error if the request for a JWT isn't a valid request.
	// If any of these functions return an error  that != ErrNotUsed, the request shouldn't be considered valid.
	RequestValidators []RequestValidator

	// AbortRequest is the function that's invoked if the request is unauthorized and therefore about to be aborted.
	// AbortRequest is expected to send a responose on the ResponseWriter.
	AbortRequest EnforceFunc

	// ContinueRequest is the function that's invoked if the request is unauthorized but will ne allowed to continue.
	ContinueRequest EnforceFunc

	// TokenBlocked is the function that is invoked if a token should not be issues.
	TokenBlocked TokenBlockedFunc

	Cache *redis.Client
}

// New returns a new Service with the provided configuration.
// It will validate the private key and panic if it does not pass.
func New(c Config) *Service {
	key := mustParseRSAPrivateKeyFromPEM(c.PrivateKey)

	if c.TZ == nil {
		c.TZ = time.Local
	}

	return &Service{
		requestValidators: c.RequestValidators,
		ttl:               2 * time.Hour,
		enforce:           c.Enforce,
		tz:                c.TZ,
		cache:             c.Cache,
		privateKey:        key,
		publicKey:         &key.PublicKey,
		abortRequest:      c.AbortRequest,
		continueRequest:   c.ContinueRequest,
		tokenBlocked:      c.TokenBlocked,
		keyPrefix:         "auth:",
		issuer:            c.Issuer,
	}
}

func mustParseRSAPrivateKeyFromPEM(privateKey string) *rsa.PrivateKey {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		panic(err)
	}

	return key
}

// NewSignedToken creates a new JWT, persists it to Redis and returns the signed token.
// It can return an error if there is an issue signing the token with th egiven RSA private key or saving to the cache.
func (s *Service) NewSignedToken() (string, error) {
	now := time.Now().UTC()
	expiresAt := now.Add(s.ttl)

	claims := &jwt.StandardClaims{
		ExpiresAt: expiresAt.Unix(),
		Issuer:    s.issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)

	ss, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", errors.Wrap(err, "signed string")
	}

	key := fmt.Sprintf("token:%s", ss)
	val := expiresAt.Unix()
	exp := expiresAt.Sub(now)

	err = s.cache.Set(key, val, exp).Err()
	if err != nil {
		return "", errors.Wrap(err, "error persisting token to cache")
	}

	return ss, nil
}

// IssueTokenHandler is the http.Handler that can issue JWTs signed with the provided RSA Key
func (s *Service) IssueTokenHandler(w http.ResponseWriter, r *http.Request) {
	if !s.validRequest(r) {
		s.tokenBlocked(r, ErrNotAuthorized, http.StatusUnauthorized)
		if s.enforce {
			web.RespondWithCodedError(w, r, http.StatusUnauthorized, ErrNotAuthorized.Error(), ErrNotAuthorized)
			return
		}
	}

	ss, err := s.NewSignedToken()
	if err != nil {
		s.tokenBlocked(r, ErrGenerateToken, http.StatusInternalServerError)
		if s.enforce {
			web.RespondWithCodedError(w, r, http.StatusInternalServerError, ErrGenerateToken.Error(), ErrGenerateToken)
			return
		}
	}

	response := struct {
		Success bool   `json:"success"`
		Token   string `json:"token"`
	}{
		Success: true,
		Token:   ss,
	}

	web.Respond(w, r, response, http.StatusOK)
}

// validRequest return true if:
// 		- there are 0 RequestValidators
//		-at least one RequestValidators returns nil error and the others return nil error or ErrNotUsed
func (s *Service) validRequest(r *http.Request) bool {
	valid := len(s.requestValidators) == 0

	for _, v := range s.requestValidators {
		err := v(r)
		switch {
		case err == nil:
			valid = true
		case errors.Cause(err) == ErrNotUsed:
			// nothing
		default:
			return false
		}
	}

	return valid
}

// getEnforcer returns an EnforceFunc that wraps the passed in http.Handler.
// If the service is in enforcing mode, it will return something to about the request.
// Otherwise, it'll return something which continues the request normally.
func (s *Service) getEnforcer(next http.Handler) EnforceFunc {
	if s.enforce {
		return s.abortRequest
	}

	return EnforceFunc(func(w http.ResponseWriter, r *http.Request, err error, statusCode int) {
		s.continueRequest(w, r, err, statusCode)
		next.ServeHTTP(w, r)
	})
}

// RequireValidToken is middleware that requires a valid JWT as either header or a querystring parameter.
// If a token is present, it is validated.
// If validated, the token is placed in the request's context for further use down the chain.
func (s *Service) RequireValidToken(next http.Handler) http.Handler {
	enforce := s.getEnforcer(next)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the token out of the request and make sure it's not empty
		rawToken := r.URL.Query().Get("token")
		if rawToken == "" {
			t := r.Header.Get("Authorization")
			if t == "" {
				// There was no token supplied
				enforce(w, r, ErrMissingToken, http.StatusUnauthorized)
				return
			}

			// The request might have multiple tokens passed via the "Authorization" header
			// but we are just interested in the Bearer token, so we split on the ","
			// and look for the Bearer token.
			for _, v := range strings.Split(t, ",") {
				if strings.HasPrefix(strings.TrimSpace(v), "Bearer") {
					rawToken = strings.Replace(v, "Bearer ", "", 1)
				}
			}
		}

		// Parse the token and verify the claims.
		// If the token can't be parsed or the claims aren't valid, the token is unauthorized.

		claims := &jwt.StandardClaims{}
		jwtParser := &jwt.Parser{
			// We want different error messages for a malformed JWT vs on that's expired,
			// so we will validate the claims separately.
			SkipClaimsValidation: true,
		}

		_, err := jwtParser.ParseWithClaims(rawToken, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.Errorf("unexpected signing method: %T", token.Header["alg"])
			}
			return s.publicKey, nil
		})
		if err != nil {
			enforce(w, r, errors.Wrap(err, "parse jwt"), http.StatusUnauthorized)
			return
		}

		err = claims.Valid()
		if err != nil {
			enforce(w, r, errors.Wrap(err, "invalid claims"), http.StatusUnauthorized)
			return
		}

		log.WithField("expiresAt", time.Unix(claims.ExpiresAt, 0).Local()).Info("token authenticated")

		// Put the token in the request context to be used by later middlewares.
		r = r.WithContext(context.WithValue(r.Context(), tokenKey, rawToken))

		next.ServeHTTP(w, r)
	})
}
