package auth

import (
	"net/http"

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
			// Parse the JWT and validate it using the public key.
			_, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
				return &c.key.PublicKey, nil
			})
			if err != nil {
				http.Error(w, "Invalid JWT: "+err.Error(), http.StatusUnauthorized)
				return
			}
			next(w, r)
		}
	}
}
