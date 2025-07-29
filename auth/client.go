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

func NewClientRsa(pemKey []byte) (*client, error) {
	var key *rsa.PrivateKey
	var err error

	if len(pemKey) == 0 {
		// Load the RSA key from PEM format if provided
		if key, err = jwt.ParseRSAPrivateKeyFromPEM(pemKey); err != nil {
			return nil, fmt.Errorf("failed to parse RSA private key: %w", err)
		}
	} else {
		if key, err = rsa.GenerateKey(rand.Reader, 4096); err != nil {
			return nil, fmt.Errorf("failed to generate RSA key: %w", err)
		}
	}
	return &client{
		key: key,
	}, nil
}

func (c *client) GenerateToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(c.key)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT: %w", err)
	}
	return tokenString, nil

}
