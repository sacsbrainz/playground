package main

import (
	"crypto/rand"
	"net/http"
	"os"
	"path"
	_ "playground/app/cache"
	_ "playground/app/db"
	"playground/app/handlers"
	m "playground/app/middleware"
	"playground/app/utils"
	"playground/app/views/auth"
	"playground/app/views/dashboard"
	"playground/app/views/errors"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/gorilla/csrf"
)

func main() {
	router := chi.NewRouter()
	listenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	appEnv := os.Getenv("APP_ENV")

	// A good base middleware stack
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Use(
		cleanPath,
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Compress(5, "text/html", "text/css"),
		middleware.Recoverer,
	)

	// Static files handler for public assets (images, css, js etc.)
	router.Handle("/public/*", utils.StaticFileHandler(appEnv))

	// Public Routes
	router.Group(func(r chi.Router) {
		r.Use(m.NoAuthMiddleware)
		r.Get("/auth", templ.Handler(auth.AuthIndex()).ServeHTTP)
		r.Get("/auth/signup", templ.Handler(auth.SignUpIndex(&utils.GlobalFormState{})).ServeHTTP)
		r.Post("/auth/signup/email", handlers.SignUpEmail)
		r.Post("/auth/signin/email", handlers.LogInEmail)
		r.Get("/auth/login", templ.Handler(auth.LoginIndex(&utils.GlobalFormState{})).ServeHTTP)
		r.Get("/auth/forgot-password", templ.Handler(auth.ForgotPasswordIndex()).ServeHTTP)
	})

	// Private Routes
	// Require Authentication
	router.Group(func(r chi.Router) {
		r.Use(m.AuthMiddleware) // Uncomment and implement your AuthMiddleware
		r.Get("/dashboard", templ.Handler(dashboard.Index()).ServeHTTP)

	})

	// Handle 500
	router.Get("/500", templ.Handler(errors.Error500(), templ.WithStatus(http.StatusInternalServerError)).ServeHTTP)

	// Handle 404
	router.NotFound(templ.Handler(errors.Error404(), templ.WithStatus(http.StatusNotFound)).ServeHTTP)

	csrfMiddleware := csrf.Protect(mustGenerateCSRFKey())
	withCSRFProtection := csrfMiddleware(router)
	http.ListenAndServe(listenAddr, withCSRFProtection)
}

// when this is placed in the middleware folder or accessed from chi directly it crashes (nill pointer)
func cleanPath(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())

		routePath := rctx.RoutePath
		if routePath == "" {
			if r.URL.RawPath != "" {
				routePath = r.URL.RawPath
			} else {
				routePath = r.URL.Path
			}
			rctx.RoutePath = path.Clean(routePath)
		}

		next.ServeHTTP(w, r)
	})
}

func mustGenerateCSRFKey() (key []byte) {
	key = make([]byte, 32)
	n, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	if n != 32 {
		panic("unable to read 32 bytes for CSRF key")
	}
	return
}
