package middleware

import (
	"context"
	"net/http"
	"strings"

	"mariadiezmaback/internal/domain"
	"mariadiezmaback/internal/service"
	"mariadiezmaback/pkg/response"
)

type contextKey string

const userClaimsKey contextKey = "user_claims"

func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Unauthorized(w, "invalid authorization format, expected 'Bearer <token>'")
				return
			}

			claims, err := authService.ValidateToken(parts[1])
			if err != nil {
				response.Unauthorized(w, "invalid or expired token")
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRoles(roles ...domain.UserRole) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims := GetClaims(r.Context())
			if claims == nil {
				response.Unauthorized(w, "unauthenticated")
				return
			}

			hasRole := false
			for _, r := range roles {
				if claims.Role == r {
					hasRole = true
					break
				}
			}

			if !hasRole {
				response.Forbidden(w, "you do not have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetClaims(ctx context.Context) *service.JWTClaims {
	claims, ok := ctx.Value(userClaimsKey).(*service.JWTClaims)
	if !ok {
		return nil
	}
	return claims
}
