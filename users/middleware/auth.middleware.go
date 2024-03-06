package authmiddleware

import (
	"context"
	"log"
	"net/http"
	"users/model"

	"github.com/dgrijalva/jwt-go"
)

var signingKey = []byte("your-secret-key")

type userContextKey string

const (
	// UserContextKey is the key used to store user details in the request context
	UserContextKey userContextKey = "user"
)

func AdminAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {

			if claims.Type != "admin" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return signingKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {

			log.Printf("Parsed claims: %+v", claims)

			ctx := context.WithValue(r.Context(), UserContextKey, claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	})
}
