package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (c *client) VerifyJwt() func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Extract the JWT from the request header.
			jwtToken := r.Header.Get("Authorization")
			if jwtToken == "" {
				http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
				return
			}
			if len(jwtToken) < 7 && jwtToken[:7] != "Bearer " {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			parsedToken := jwtToken[7:] // Remove "Bearer " prefix

			claims := jwt.RegisteredClaims{}
			if _, err := jwt.ParseWithClaims(parsedToken, &claims, func(token *jwt.Token) (interface{}, error) {
				return &c.key.PublicKey, nil
			}); err != nil {
				http.Error(w, "Invalid JWT: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// Only check expiry (exp)
			if claims.ExpiresAt != nil && !claims.ExpiresAt.After(time.Now()) {
				http.Error(w, "JWT expired", http.StatusUnauthorized)
				return
			}

			next(w, r)
		}
	}
}
