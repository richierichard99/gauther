package redis

import (
	"context"
	"testing"
)

func TestRedisStore(t *testing.T) {
	addr := "localhost:6379"
	store := NewStore(addr, 15)
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
	valid := store.Validate(ctx, "testuser", "testpassword")
	if !valid {
		t.Error("Validate failed for correct credentials")
	}

	// Test Validate with wrong password
	valid = store.Validate(ctx, "testuser", "wrongpassword")
	if valid {
		t.Error("Validate succeeded for incorrect password")
	}

	// Test Validate with non-existent user
	valid = store.Validate(ctx, "nonexistent", "password")
	if valid {
		t.Error("Validate succeeded for non-existent user")
	}
	// Clean up by deleting the test user

	t.Log("Redis store tests passed")
}
