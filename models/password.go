package models

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"

	"github.com/smartcontractkit/chainlink-go/logger"
	"golang.org/x/crypto/pbkdf2"
)

type Password struct {
	Hash []byte `storm:"id"`
	Salt []byte
}

func NewPassword(phrase string) Password {
	salt := generateSalt()
	return Password{
		Hash: generateHash(phrase, salt),
		Salt: salt,
	}
}

func (self Password) Check(phrase string) bool {
	return bytes.Compare(generateHash(phrase, self.Salt), self.Hash) == 0
}

func generateHash(phrase string, salt []byte) []byte {
	return pbkdf2.Key([]byte(phrase), salt, 262144, 32, sha256.New)
}

func generateSalt() []byte {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		logger.Fatal(err)
	}
	return salt
}
