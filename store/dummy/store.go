package dummy

import "context"

type store struct{}

// NewStore creates a new dummy store client
func NewStore() *store {
	return &store{}
}

func (s *store) Validate(ctx context.Context, username, password string) bool {
	// Dummy validation logic
	return username == "admin" && password == "password123"
}
