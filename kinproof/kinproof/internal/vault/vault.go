package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
)

type EncryptedVault struct {
	EncryptedMnemonic string `json:"encryptedMnemonic"`
	IV                string `json:"iv"`
	Salt              string `json:"salt"`
	Username          string `json:"username"`
}

func EncryptAndSave(mnemonic, pin, username string) error {
	salt := make([]byte, 16)
	rand.Read(salt)
	key := sha256.Sum256(append([]byte(pin), salt...))

	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nil, nonce, []byte(mnemonic), nil)

	v := EncryptedVault{
		EncryptedMnemonic: hex.EncodeToString(ciphertext),
		IV:                hex.EncodeToString(nonce),
		Salt:              hex.EncodeToString(salt),
		Username:          username,
	}

	os.MkdirAll("keys", 0700)
	data, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile("keys/vault.json", data, 0600)
}
