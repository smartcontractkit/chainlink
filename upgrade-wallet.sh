#!/usr/bin/env bash
set -euo pipefail

echo "=== Upgrading to Full Sovereign Wallet (PIN + Guardian) ==="

# Add vault package
cat <<'VAULT' > internal/vault/vault.go
package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
)

type EncryptedVault struct {
	EncryptedMnemonic string `json:"encryptedMnemonic"`
	IV                string `json:"iv"`
	Tag               string `json:"tag"`
	Salt              string `json:"salt"`
}

func EncryptMnemonic(mnemonic, pin string) (EncryptedVault, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	key := sha256.Sum256([]byte(pin + hex.EncodeToString(salt)))

	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nil, nonce, []byte(mnemonic), nil)

	return EncryptedVault{
		EncryptedMnemonic: hex.EncodeToString(ciphertext),
		IV:                hex.EncodeToString(nonce),
		Tag:               hex.EncodeToString(ciphertext[len(ciphertext)-gcm.Overhead():]),
		Salt:              hex.EncodeToString(salt),
	}, nil
}

func SaveVault(v EncryptedVault) error {
	os.MkdirAll("keys", 0700)
	data, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile("keys/vault.json", data, 0600)
}
VAULT

# Update main.go with full flow
cat <<'MAIN' > cmd/kinproof/main.go
package main

import (
	"encoding/json"
	"fmt"
	"kinproof/internal/attribution"
	"kinproof/internal/hd"
	"kinproof/internal/identity"
	"kinproof/internal/intent"
	"kinproof/internal/proof"
	"kinproof/internal/vault"
)

func main() {
	fmt.Println("=== KinProof Full Sovereign Wallet ===")

	// === STEP 1: Secure Mnemonic Input (In production use real input) ===
	mnemonic := "test test test test test test test test test test test junk"
	pin := "123456" // Change this in production

	fmt.Println("✓ Mnemonic loaded (encrypted in production)")

	// === STEP 2: Create Sovereign HD Wallet ===
	hdw, err := hd.New(mnemonic)
	if err != nil {
		panic(err)
	}

	receive, _ := hdw.DeriveReceiveOnly()
	fmt.Println("Receive-Only Anchor Address:", receive)

	// === STEP 3: Encrypt + Save Vault ===
	encVault, _ := vault.EncryptMnemonic(mnemonic, pin)
	vault.SaveVault(encVault)
	fmt.Println("✓ Mnemonic encrypted with PIN and saved to keys/vault.json")

	// === STEP 4: Rotate Ephemeral Identity ===
	id, _ := identity.Rotate(hdw, "proofs", 42)

	// === STEP 5: Create Sovereign Intent ===
	env, _ := proof.CreateEnvelope("sovereign-root-hash", id.Address, "sei-transfer", "nonce-001", map[string]string{
		"to":     "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		"amount": "0.1",
	})

	sig, _ := intent.Sign(id.PrivateKey, env)

	receipt := map[string]interface{}{
		"envelope":  env,
		"signature": sig,
		"status":    "intent_recorded",
	}
	receipt = attribution.Attach(receipt, "The Keeper")

	out, _ := json.MarshalIndent(receipt, "", "  ")
	fmt.Println(string(out))
	fmt.Println("\n✓ Full sovereign transaction intent with encrypted vault completed.")
}
MAIN

echo "=== Wallet Upgrade Complete ==="
echo "Run: ./run.sh"
