package main

import (
	"encoding/base64"
	"log"
	"os"

	"github.com/richierichard99/gauther/auth"
	"github.com/richierichard99/gauther/server"
)

func main() {
	b64 := os.Getenv("GUATH_PRIVATE_KEY")
	pemKey, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		log.Fatalf("failed to decode private key: %v", err)
	}
	authClient, err := auth.NewClientRsa(pemKey)
	if err != nil {
		log.Fatalf("failed to create auth client: %v", err)
	}

	httpServer := server.NewServer(log.Default(), authClient)
	log.Println("Server running on :8080")
	if err := httpServer.Start(":8080", authClient); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
