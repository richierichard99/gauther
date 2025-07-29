package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type store struct {
	redis *redis.Client
}

// NewStore creates a new Redis client
func NewStore(addr string, db int) *store {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &store{redis: rdb}
}

// Validate checks the username and password
func (s *store) Validate(ctx context.Context, username, password string) bool {
	storedHash, err := s.redis.Get(ctx, username).Result()
	if err != nil {
		return false
	}
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) != nil {
		return false
	}
	return true
}

// InsertUser inserts a username and hashed password into Redis
func (s *store) InsertUser(ctx context.Context, username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	if err := s.redis.Set(ctx, username, hashedPassword, 0).Err(); err != nil {
		return fmt.Errorf("failed to set user in Redis: %w", err)
	}
	return nil
}
