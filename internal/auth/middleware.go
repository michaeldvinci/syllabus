package auth

import (
	"context"
	"net/http"
)

type contextKey string

const UserContextKey contextKey = "user"

// Middleware provides authentication middleware functionality
type Middleware struct {
	store *Store
}

// NewMiddleware creates a new authentication middleware
func NewMiddleware(store *Store) *Middleware {
	return &Middleware{store: store}
}

// RequireAuth is middleware that requires authentication
func (m *Middleware) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session cookie
		cookie, err := r.Cookie("session_token")
		if err != nil {
			// No session cookie, redirect to login
			if r.Header.Get("Accept") == "application/json" {
				http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Validate session
		user, err := m.store.GetUserBySession(cookie.Value)
		if err != nil {
			// Invalid session, clear cookie and redirect
			http.SetCookie(w, &http.Cookie{
				Name:     "session_token",
				Value:    "",
				Path:     "/",
				MaxAge:   -1,
				HttpOnly: true,
				Secure:   r.TLS != nil,
				SameSite: http.SameSiteStrictMode,
			})
			
			if r.Header.Get("Accept") == "application/json" {
				http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth is middleware that adds user context if authenticated, but doesn't require it
func (m *Middleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for session cookie
		cookie, err := r.Cookie("session_token")
		if err == nil {
			// Try to get user from session
			user, err := m.store.GetUserBySession(cookie.Value)
			if err == nil {
				// Add user to request context
				ctx := context.WithValue(r.Context(), UserContextKey, user)
				r = r.WithContext(ctx)
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequireAdmin is middleware that requires admin role
func (m *Middleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r)
		if !ok {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		
		if !user.IsAdmin() {
			if r.Header.Get("Accept") == "application/json" {
				http.Error(w, `{"error":"admin access required"}`, http.StatusForbidden)
				return
			}
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext extracts the user from request context
func GetUserFromContext(r *http.Request) (*User, bool) {
	user, ok := r.Context().Value(UserContextKey).(*User)
	return user, ok
}