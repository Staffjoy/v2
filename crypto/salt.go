package crypto

import "crypto/rand"

// NewSalt generates a random salt for use in account
// passwords
func NewSalt() ([]byte, error) {
	key := make([]byte, 60)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
