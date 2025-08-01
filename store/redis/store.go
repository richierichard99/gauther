package redis

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type store struct {
	logger *log.Logger
	redis  *redis.Client
}

// NewStore creates a new Redis client
func NewStore(log *log.Logger, addr string, db int) *store {
	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   db,
	})

	return &store{redis: rdb, logger: log}
}

// Validate checks the username and password
func (s *store) Validate(ctx context.Context, username, password string) (bool, error) {
	storedHash, err := s.redis.Get(ctx, username).Result()
	if err == redis.Nil {
		s.logger.Printf("User %s not found", username)
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("failed to get user from Redis: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(password)) != nil {
		return false, nil
	}
	return true, nil
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
