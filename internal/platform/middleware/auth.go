package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key int

const (
	userIDKey key = 0
	tenantIDKey key = 1
)

type Claims struct {
	UserID   int64 `json:"user_id"`
	TenantID int64 `json:"tenant_id"`
	jwt.RegisteredClaims
}

func JWTAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondWithError(w, http.StatusUnauthorized, "Authorization header required")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				respondWithError(w, http.StatusUnauthorized, "Invalid authorization header format")
				return
			}

			tokenString := parts[1]
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(secretKey), nil
			})

			if err != nil {
				respondWithError(w, http.StatusUnauthorized, "Invalid token")
				return
			}

			claims, ok := token.Claims.(*Claims)
			if !ok {
				mapClaims, ok := token.Claims.(jwt.MapClaims)
				if !ok || !token.Valid {
					respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
					return
				}
				userIDFloat, ok := mapClaims["user_id"].(float64)
				if !ok {
					respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
					return
				}
				tenantIDFloat, ok := mapClaims["tenant_id"].(float64)
				if !ok {
					respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
					return
				}
				ctx := context.WithValue(r.Context(), userIDKey, int64(userIDFloat))
				ctx = context.WithValue(ctx, tenantIDKey, int64(tenantIDFloat))
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if !token.Valid {
				respondWithError(w, http.StatusUnauthorized, "Invalid token claims")
				return
			}

			ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, tenantIDKey, claims.TenantID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserIDFromContext(ctx context.Context) int64 {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0
	}
	return userID
}

func SetUserIDInContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetTenantIDFromContext(ctx context.Context) int64 {
	tenantID, ok := ctx.Value(tenantIDKey).(int64)
	if !ok {
		return 0
	}
	return tenantID
}

func SetTenantIDInContext(ctx context.Context, tenantID int64) context.Context {
	return context.WithValue(ctx, tenantIDKey, tenantID)
}

func respondWithError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write([]byte(`{"error":"` + message + `"}`))
}
