package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type authClient interface {
	GenerateToken() (string, error)
	VerifyJwt() func(http.HandlerFunc) http.HandlerFunc
}

// Server struct holds the auth client
type httpServer struct {
	authClient authClient
	logger     *log.Logger
}

type creds struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewServer creates a new Server instance
func NewServer(log *log.Logger, authClient authClient) *httpServer {
	return &httpServer{authClient: authClient, logger: log}
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

// loginHandler handles login and returns a JWT if credentials are valid
func (s *httpServer) loginHandler(w http.ResponseWriter, r *http.Request) {
	var c creds
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// Hardcoded credentials
	if c.Username != "admin" || c.Password != "password123" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token, err := s.authClient.GenerateToken()
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
func (s *httpServer) Start(addr string, authClient authClient) error {
	http.HandleFunc("/login", s.loginHandler)
	// Use JWT middleware for /validate
	http.HandleFunc("/validate", buildChain(s.validateHandler, authClient.VerifyJwt()))
	return http.ListenAndServe(addr, nil)
}
