package crypto

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// EmailConfirmationToken provides a signed piece of info that can be
// emailed to a user, and when it "round trips" back to us - we can
// assert that they control the email.
func EmailConfirmationToken(uuid, email, signingToken string) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"email": email,
		"uuid":  uuid,
		"exp":   time.Now().Add(time.Duration(2 * time.Hour)).Unix(),
	}).SignedString([]byte(signingToken))

	if err != nil {
		return "", err
	}
	return token, nil
}

// VerifyEmailConfirmationToken takes a token, and returns the uuid and email if valid.
func VerifyEmailConfirmationToken(tokenString, signingToken string) (email string, uuid string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingToken), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email = claims["email"].(string)
		uuid = claims["uuid"].(string)
		return
	}
	err = fmt.Errorf("unable to verify token")
	return
}

// SessionToken provides a signed piece of info that can be
// stored in a user's browser
func SessionToken(uuid, signingToken string, support bool, dur time.Duration) (string, error) {
	if len(signingToken) == 0 {
		panic("No signing token present")
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"uuid":    uuid,
		"support": support,
		"exp":     time.Now().Add(dur).Unix(),
	}).SignedString([]byte(signingToken))

	if err != nil {
		return "", err
	}
	return token, nil
}

// RetrieveSessionInformation takes a session token, and returns its validated contents
func RetrieveSessionInformation(tokenString, signingToken string) (uuid string, support bool, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingToken), nil
	})

	if err != nil {
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uuid = claims["uuid"].(string)
		support = claims["support"].(bool)
		return
	}
	err = fmt.Errorf("unable to verify token")
	return
}
