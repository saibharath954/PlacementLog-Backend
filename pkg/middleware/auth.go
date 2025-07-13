package middleware

import (
	"net/http"
	"strings"

	"github.com/varnit-ta/PlacementLog/pkg/jwt"
)

/*
AuthMiddleware validates JWT tokens and extracts user information.
This middleware checks for the Authorization header, validates the JWT token,
and adds user information to the request headers for downstream handlers.

The middleware expects:
- Authorization header with format "Bearer <token>"
- Valid JWT token

It adds to request headers:
- X-User-ID: The user ID from the token
- X-User-Role: The role from the token ("user" or "admin")

If validation fails, it returns a 401 Unauthorized response.
*/
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized: missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, role, err := jwt.ValidateJwtToken(token)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		r.Header.Set("X-User-ID", userID)
		r.Header.Set("X-User-Role", role)

		next.ServeHTTP(w, r)
	})
}

/*
UserAuthMiddleware ensures the request is from an authenticated user.
This middleware specifically validates that the JWT token belongs to a user
with the "user" role, not an admin.

The middleware expects:
- Authorization header with format "Bearer <token>"
- Valid JWT token with "user" role

It adds to request headers:
- X-User-ID: The user ID from the token
- X-User-Role: Set to "user"

If validation fails or token is not a user token, it returns a 401 Unauthorized response.
*/
func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized: missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := jwt.ValidateUserToken(token)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		r.Header.Set("X-User-ID", userID)
		r.Header.Set("X-User-Role", "user")

		next.ServeHTTP(w, r)
	})
}

/*
AdminAuthMiddleware ensures the request is from an authenticated admin.
This middleware specifically validates that the JWT token belongs to an admin
with the "admin" role.

The middleware expects:
- Authorization header with format "Bearer <token>"
- Valid JWT token with "admin" role

It adds to request headers:
- X-Admin-ID: The admin ID from the token
- X-User-Role: Set to "admin"

If validation fails or token is not an admin token, it returns a 401 Unauthorized response.
*/
func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "unauthorized: missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		adminID, err := jwt.ValidateAdminToken(token)
		if err != nil {
			http.Error(w, "unauthorized: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add admin info to request context
		r.Header.Set("X-Admin-ID", adminID)
		r.Header.Set("X-User-Role", "admin")

		next.ServeHTTP(w, r)
	})
}
