package main

import (
	"log"

	"github.com/richierichard99/gauther/auth"
	"github.com/richierichard99/gauther/server"
)

func main() {
	authClient, err := auth.NewClientRsa()
	if err != nil {
		log.Fatalf("failed to create auth client: %v", err)
	}

	httpServer := server.NewServer(log.Default(), authClient)
	log.Println("Server running on :8080")
	if err := httpServer.Start(":8080", authClient); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
