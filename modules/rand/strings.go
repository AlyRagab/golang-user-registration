package rand

import (
	"crypto/rand"
	"encoding/base64"
)

// RememberTokenBytes is used as a default value of creating byte size of token
const RememberTokenBytes = 32

// Bytes will help us to generate random bytes
// Returns error if there is one
// It uses the crypto/rand to be safe to use the remember tokens
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// CryptoString will generate a byte slice of size nByte
// then return a string that is the bae64 URL encoded of that byte slice
func CryptoString(nByte int) (string, error) {
	b, err := Bytes(nByte)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// RememberToken is a helper func to generate remember tokens
// with predefined byte size "32"
func RememberToken() (string, error) {
	return CryptoString(RememberTokenBytes)
}
