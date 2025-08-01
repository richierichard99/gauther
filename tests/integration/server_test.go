package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/richierichard99/gauther/auth"
	"github.com/richierichard99/gauther/server"
	redisStore "github.com/richierichard99/gauther/store/redis"
)

type tokenResponse struct {
	Token string `json:"token"`
}

func TestServer(t *testing.T) {
	t.Cleanup(cleanup)
	// Create a new HTTP server with the dummy store and auth client
	server := newTestServer(t)
	defer server.Close()

	res := testLoginRequest(t, server, "testuser", "testpassword")
	// Ensure the response is valid
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", res.Status)
	}
	// Read the response body
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}
	var tokenResp tokenResponse
	if err := json.Unmarshal(responseBody, &tokenResp); err != nil {
		t.Fatalf("Failed to unmarshal response body: %v", err)
	}
	if tokenResp.Token == "" {
		t.Fatal("Expected a token in the response, got empty string")
	}

	// Now test the /validate endpoint with the token
	validateRes := testValidateRequest(t, server, tokenResp.Token)
	// Ensure the response is valid
	defer validateRes.Body.Close()
	if validateRes.StatusCode != http.StatusOK {
		t.Fatalf("Expected status OK, got %v", validateRes.Status)
	}

}

func newTestServer(t *testing.T) *httptest.Server {
	authClient, err := auth.NewClientRsa(nil)
	if err != nil {
		t.Fatalf("Failed to create auth client: %v", err)
	}

	userStore := redisStore.NewStore(log.Default(), "localhost:6379", 15)
	if err := userStore.InsertUser(context.Background(), "testuser", "testpassword"); err != nil {
		t.Fatalf("Failed to insert test user: %v", err)
	}

	httpServer := server.NewServer(log.Default(), authClient, userStore)
	mux := http.NewServeMux()
	httpServer.RegisterRoutes(mux)
	return httptest.NewServer(mux)
}

func testLoginRequest(t *testing.T, server *httptest.Server, username, password string) *http.Response {
	requestBody := map[string]string{
		"username": username,
		"password": password,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatalf("Failed to marshal request body: %v", err)
	}

	res, err := http.Post(fmt.Sprintf("%s/login", server.URL), "application/json", bytes.NewReader(requestBodyBytes))
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	return res
}

func testValidateRequest(t *testing.T, server *httptest.Server, token string) *http.Response {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/validate", server.URL), nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	return res
}

func cleanup() {
	// Cleanup logic if needed, e.g., stopping the server or deleting test users
	ctx := context.Background()
	store := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15,
	})
	defer store.Close()
	if err := store.Del(ctx, "testuser").Err(); err != nil {
		log.Printf("Failed to delete test user: %v", err)
	}
}
