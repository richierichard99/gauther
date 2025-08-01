package redis

import (
	"context"
	"log"
	"testing"
)

func TestRedisStore(t *testing.T) {
	addr := "localhost:6379"
	store := NewStore(log.Default(), addr, 15)
	ctx := context.Background()

	t.Cleanup(func() {
		if err := store.redis.Del(ctx, "testuser").Err(); err != nil {
			t.Fatalf("Failed to delete test user: %v", err)
		}
	})

	// Test InsertUser
	err := store.InsertUser(ctx, "testuser", "testpassword")
	if err != nil {
		t.Fatalf("InsertUser failed: %v", err)
	}

	// Test Validate
	valid, err := store.Validate(ctx, "testuser", "testpassword")
	if err != nil {
		t.Fatalf("Validate failed: %v", err)
	}
	if !valid {
		t.Error("Validate failed for correct credentials")
	}

	// Test Validate with wrong password
	valid, err = store.Validate(ctx, "testuser", "wrongpassword")
	if err != nil {
		t.Fatalf("Validate failed with wrong password: %v", err)
	}
	if valid {
		t.Error("Validate succeeded for incorrect password")
	}

	// Test Validate with non-existent user
	valid, err = store.Validate(ctx, "nonexistent", "password")
	if err != nil {
		t.Fatalf("Validate failed for non-existent user: %v", err)
	}
	if valid {
		t.Error("Validate succeeded for non-existent user")
	}

	t.Log("Redis store tests passed")
}
