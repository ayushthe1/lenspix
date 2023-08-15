package rand

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// https://go.dev/play/p/MShqzz907vQ

// Bytes generates and returns a slice of random bytes of length 'n'
func Bytes(n int) ([]byte, error) {
	b := make([]byte, n)
	nread, err := rand.Read(b) // generate random bytes and fill the slice b with them
	if err != nil {
		return nil, fmt.Errorf("bytes: %w", err)
	}
	if nread < n {
		return nil, fmt.Errorf("bytes: didn't read enough random bytes")
	}

	return b, nil
}

// string returns a random string using crypto/rand
// n is the number of bytes being used to generate random string
func String(n int) (string, error) {
	b, err := Bytes(n) // b is like [23,4,5,12,99,..nth]
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil // return something like 'Obfmp8X6dQ=='
}
