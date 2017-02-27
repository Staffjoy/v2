package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword generates a bcrypt hash of the password using work factor 14.
func HashPassword(salt, password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(append(salt, password...), 12)
}

// CheckPasswordHash securely compares a bcrypt hashed password with its possible
// plaintext equivalent.  Returns nil on success, or an error on failure.
func CheckPasswordHash(hash, salt, password []byte) error {
	return bcrypt.CompareHashAndPassword(hash, append(salt, password...))
}

// ComputeHmac256 function is used by Intercom to sign a user id
func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
