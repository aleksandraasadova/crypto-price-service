package middleware

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWT(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var raw string

		if authHeader := r.Header.Get("Authorization"); strings.HasPrefix(authHeader, "Bearer ") {
			raw = strings.TrimPrefix(authHeader, "Bearer ")
		} else {
			writeJSONError(w, http.StatusUnauthorized, ErrorResponse{Error: "missing or invalid authorization header"})
			return
		}

		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT")), nil
		})
		if err != nil || !token.Valid {
			writeJSONError(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired token"})
			return
		}

		if claims.Subject == "" {
			writeJSONError(w, http.StatusUnauthorized, ErrorResponse{Error: "invalid or expired token"})
			return
		}

		next(w, r)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
