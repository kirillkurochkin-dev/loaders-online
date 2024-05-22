package middleware

import (
	"context"
	"loaders-online/pkg/util"
	"net/http"
	"strings"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			util.LogHandler("JWTMiddleware", "Missing token", nil)
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		claims, err := util.ValidateJWT(tokenString)
		if err != nil {
			util.LogHandler("JWTMiddleware", "Invalid token", err)
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
		ctx = context.WithValue(ctx, "role", claims.Role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RoleMiddleware(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			role, ok := r.Context().Value("role").(string)
			if !ok {
				util.LogHandler("RoleMiddleware", "Missing role in context", nil)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			for _, allowedRole := range allowedRoles {
				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			util.LogHandler("RoleMiddleware", "Forbidden role: "+role, nil)
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}
