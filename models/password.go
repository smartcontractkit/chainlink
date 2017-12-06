package models

import (
	"crypto/rand"
	"crypto/sha256"

	"github.com/smartcontractkit/chainlink-go/logger"
	"golang.org/x/crypto/pbkdf2"
)

type Password struct {
	Hash  []byte `storm:"id"`
	Salt []byte
}

func NewPassword(phrase string) Password {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		logger.Fatal(err)
	}

	hash := pbkdf2.Key([]byte(phrase), salt, 262144, 32, sha256.New)
	return Password{hash, salt}
}
