package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

// Error represents a custom error
type Error struct {
	Message string `json:"message"`
}

// Error implements the error interface
func (e Error) Error() string {
	return e.Message
}

func newError(msg string) Error {
	return Error{
		Message: msg,
	}
}

type middleware []filter
type filter func(http.Handler) http.Handler

func newMiddleware(logOut io.Writer, secretKey string, auth bool) middleware {
	log.SetOutput(logOut)

	m := middleware{
		RequestLogger,
		SecureRequest,
	}

	if auth {
		m = append(m, IsAuthorized(secretKey))
	}

	return m
}

// Wrap applies all middlewares to a base handler
func (m middleware) Wrap(base http.Handler) http.Handler {
	return doWrap(base, m...)
}

func doWrap(base http.Handler, middleware ...filter) http.Handler {
	if len(middleware) == 0 {
		return base
	}
	first := middleware[0]
	remaining := doWrap(base, middleware[1:]...)
	return first(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remaining.ServeHTTP(w, r)
	}))
}

// RequestLogger intercepts and logs data about incoming requests
func RequestLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestLog := fmt.Sprintf("httpMethod=%s, httpRequestURI=%s, header=%+v", r.Method, r.URL.RequestURI(), r.Header)
		if len(r.URL.Query()) > 0 {
			requestLog = fmt.Sprintf("httpMethod=%s, httpRequestURI=%s, queryParams=%+v, header=%+v", r.Method, r.URL.RequestURI(), r.URL.Query(), r.Header)
		}

		// prevent auth token logging
		if r.Header.Get("Authorization") != "" {
			requestLog = strings.ReplaceAll(requestLog, r.Header.Get("Authorization"), "<sensitive>")
		}
		log.Println("Request received: ", requestLog)
		h.ServeHTTP(w, r)
	})
}

func IsAuthorized(secretKey string) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.Header.Get("Authorization")
			if tokenStr == "" {
				serveResponse(newError("Missing authorization header"), http.StatusUnauthorized, w)
				return
			}

			tokenStr = tokenStr[len("Bearer "):]

			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return secretKey, nil
			})

			if err != nil {
				log.Println("error", err.Error())
				werr := newError(err.Error())
				serveResponse(werr, http.StatusInternalServerError, w)
				return
			}

			if !token.Valid {
				serveResponse(newError("Invalid token"), http.StatusUnauthorized, w)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Println("Welcome to the the protected area")
			h.ServeHTTP(w, r)
		})
	}
}

func SecureRequest(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubdomains; preload")
		w.Header().Set("Cache-Control", "no-store")
		w.Header().Set("Pragma", "no-cache")
		h.ServeHTTP(w, r)
	})
}

func serveResponse(res interface{}, status int, w http.ResponseWriter) {
	bb, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(bb)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(bb)
}
