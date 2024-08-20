package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"playground/app/cache"
	"playground/app/utils"
	"strings"
	"time"
)

type (
	RequestKey         struct{}
	ResponseHeadersKey struct{}
	UserContextKey     string
)

var UserKey UserContextKey = "user" // because of warning: should not use built-in type string as key for value; define your own type to avoid collisions (SA1029)

func WithRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestKey{}, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NoAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		session, err := cache.GetCache().Str().Get("auth:" + cookie.Value)

		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		splitValue := strings.Split(session.String(), "::")
		if len(splitValue) < 2 {
			next.ServeHTTP(w, r)
			return
		}

		utils.Redirect(w, r, http.StatusSeeOther, "/dashboard")
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")

		if err != nil {
			slog.Error("get token", err)
			utils.Redirect(w, r, http.StatusSeeOther, "/auth/login")
			return
		}

		session, err := cache.GetCache().Str().Get("auth:" + cookie.Value)

		if err != nil {
			slog.Error("check session id in cache", err)
			http.SetCookie(w, &http.Cookie{
				Name:    "token",
				MaxAge:  -1,
				Expires: time.Now().Add(-100 * time.Hour),
				Path:    "/",
			})
			utils.Redirect(w, r, http.StatusSeeOther, "/auth/login")
			return
		}

		splitValue := strings.Split(session.String(), "::")
		if len(splitValue) < 2 {
			utils.Redirect(w, r, http.StatusSeeOther, "/auth/login")
			return
		}

		userID := splitValue[0]
		// ip := splitValue[1]
		// ua := splitValue[2]
		// createdAt := splitValue[3]
		// expiresAt := splitValue[4]
		ctx := context.WithValue(r.Context(), UserKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
