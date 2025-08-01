package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type authClient interface {
	GenerateToken(claims jwt.Claims) (string, error)
	VerifyJwt() func(http.HandlerFunc) http.HandlerFunc
}

type userStore interface {
	Validate(context context.Context, username, password string) (bool, error)
}

// Server struct holds the auth client
type httpServer struct {
	authClient authClient
	userStore  userStore
	logger     *log.Logger
}

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewServer creates a new Server instance
func NewServer(log *log.Logger, authClient authClient, userStore userStore) *httpServer {
	return &httpServer{authClient: authClient, userStore: userStore, logger: log}
}

// middleware is a definition of  what a middleware is,
// take in one handlerfunc and wrap it within another handlerfunc
type middleware func(http.HandlerFunc) http.HandlerFunc

// buildChain builds the middlware chain recursively, functions are first class
func buildChain(f http.HandlerFunc, m ...middleware) http.HandlerFunc {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](buildChain(f, m[1:cap(m)]...))
}

func (s *httpServer) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/login", s.loginHandler)
	// Use JWT middleware for /validate
	mux.HandleFunc("/validate", buildChain(s.validateHandler, s.authClient.VerifyJwt()))
}

// loginHandler handles login and returns a JWT if credentials are valid
func (s *httpServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	valid, err := s.userStore.Validate(r.Context(), c.Username, c.Password)
	if err != nil {
		s.logger.Printf("failed to validate user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !valid {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwt.MapClaims{
		"exp": jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		"iat": jwt.NewNumericDate(time.Now()),
		"sub": c.Username,
	}
	token, err := s.authClient.GenerateToken(claims)
	if err != nil {
		s.logger.Printf("failed to generate token: %v", err)
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// validateHandler is protected by JWT middleware
func (s *httpServer) validateHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// Start runs the HTTP server
func (s *httpServer) Start(addr string) error {
	mux := http.NewServeMux()
	s.RegisterRoutes(mux)
	return http.ListenAndServe(addr, mux)
}
