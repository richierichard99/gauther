package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type client struct {
	key *rsa.PrivateKey
}

func NewClientRsa() (*client, error) {
	rsaKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key: %w", err)
	}
	return &client{
		key: rsaKey,
	}, nil
}

func (c *client) GenerateToken() (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(c.key)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenString, nil

}
