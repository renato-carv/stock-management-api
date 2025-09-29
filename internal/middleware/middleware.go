package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "user_id"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, response *http.Request) {
		authHeader := response.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(writer, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]
		secret := os.Getenv("JWT_SECRET")
		if secret == "" {
			http.Error(writer, "Missing JWT_SECRET environment variable", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(writer, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || claims["user_id"] == nil {
			http.Error(writer, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		userId, ok := claims["user_id"].(string)
		if !ok {
			http.Error(writer, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(response.Context(), UserIDKey, userId)
		next.ServeHTTP(writer, response.WithContext(ctx))
	}
}
