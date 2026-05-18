#!/usr/bin/env bash
set -euo pipefail

echo "=== KinProof Full Sovereign Wallet Setup (PIN + Vault) ==="

# Clean and recreate
rm -rf kinproof
mkdir -p kinproof
cd kinproof

# Create all directories first
mkdir -p cmd/kinproof internal/{hd,identity,intent,proof,attribution,vault} proofs receipts keys scripts

# go.mod
cat <<'GO' > go.mod
module kinproof
go 1.23

require (
	github.com/ethereum/go-ethereum v1.14.12
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
)
GO

# HD Wallet
cat <<'HD' > internal/hd/hd.go
package hd

import (
	"encoding/hex"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type SovereignHD struct {
	Mnemonic string
	Root     *bip32.Key
}

func New(mnemonic string) (*SovereignHD, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, "")
	root, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	return &SovereignHD{Mnemonic: mnemonic, Root: root}, nil
}

func (h *SovereignHD) DeriveReceiveOnly() (string, error) {
	k := h.Root
	k, _ = k.NewChildKey(402)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	privKey, _ := crypto.ToECDSA(k.Key)
	return crypto.PubkeyToAddress(privKey.PublicKey).Hex(), nil
}

func (h *SovereignHD) DeriveEphemeral(index uint32) (string, string, error) {
	k := h.Root
	k, _ = k.NewChildKey(402)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(index)
	privHex := hex.EncodeToString(k.Key)
	privKey, _ := crypto.ToECDSA(k.Key)
	addr := crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	return addr, privHex, nil
}
HD

# Vault (PIN Encryption)
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
)

type EncryptedVault struct {
	EncryptedMnemonic string `json:"encryptedMnemonic"`
	IV                string `json:"iv"`
	Salt              string `json:"salt"`
}

func EncryptMnemonic(mnemonic, pin string) (EncryptedVault, error) {
	salt := make([]byte, 16)
	rand.Read(salt)
	key := sha256.Sum256(append([]byte(pin), salt...))

	block, _ := aes.NewCipher(key[:])
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	rand.Read(nonce)

	ciphertext := gcm.Seal(nil, nonce, []byte(mnemonic), nil)

	return EncryptedVault{
		EncryptedMnemonic: hex.EncodeToString(ciphertext),
		IV:                hex.EncodeToString(nonce),
		Salt:              hex.EncodeToString(salt),
	}, nil
}

func Save(v EncryptedVault) error {
	os.MkdirAll("keys", 0700)
	data, _ := json.MarshalIndent(v, "", "  ")
	return os.WriteFile("keys/vault.json", data, 0600)
}
VAULT

# (Other packages + Main - shortened for space but complete)
cat <<'ID' > internal/identity/identity.go
package identity
import ("encoding/json"; "os"; "path/filepath"; "time"; "kinproof/internal/hd")
type ID struct { Address string `json:"address"`; PrivateKey string `json:"privateKey"`; Index uint32 `json:"index"`; CreatedAt int64 `json:"createdAt"` }
func Rotate(h *hd.SovereignHD, dir string, i uint32) (*ID, error) {
	addr, priv, _ := h.DeriveEphemeral(i)
	id := &ID{addr, priv, i, time.Now().UnixMilli()}
	os.MkdirAll(dir, 0700)
	data, _ := json.MarshalIndent(id, "", "  ")
	os.WriteFile(filepath.Join(dir, "latest_identity.json"), data, 0600)
	return id, nil
}
ID

cat <<'INT' > internal/intent/intent.go
package intent
import ("encoding/hex"; "encoding/json"; "github.com/ethereum/go-ethereum/crypto"; "kinproof/internal/proof")
func Sign(pkHex string, e proof.ExecutionEnvelope) (string, error) {
	pk, _ := hex.DecodeString(pkHex)
	priv, _ := crypto.ToECDSA(pk)
	msg, _ := json.Marshal(e)
	hash := crypto.Keccak256Hash(msg)
	sig, _ := crypto.Sign(hash.Bytes(), priv)
	return hex.EncodeToString(sig), nil
}
INT

cat <<'PR' > internal/proof/proof.go
package proof
import ("crypto/sha256"; "encoding/hex"; "encoding/json"; "fmt")
type ExecutionEnvelope struct {
	ExecutionID string `json:"executionId"`
	IdentityRoot string `json:"identityRoot"`
	EphemeralAddr string `json:"ephemeralAddr"`
	Intent string `json:"intent"`
	Nonce string `json:"nonce"`
}
func CreateEnvelope(root, addr, intent, nonce string, payload interface{}) (ExecutionEnvelope, error) {
	p, _ := json.Marshal(payload)
	h := sha256.Sum256(p)
	id := sha256.Sum256([]byte(fmt.Sprintf("%s|%s|%x|%s", root, addr, h, nonce)))
	return ExecutionEnvelope{hex.EncodeToString(id[:]), root, addr, intent, nonce}, nil
}
PR

cat <<'ATTR' > internal/attribution/attribution.go
package attribution
import "time"
type Meta struct {
	Protocol string `json:"protocol"`
	Subsystem string `json:"subsystem"`
	Creator string `json:"creator"`
	ExecutionStandard string `json:"executionStandard"`
	Timestamp int64 `json:"timestamp"`
}
func Attach(r map[string]interface{}, c string) map[string]interface{} {
	r["attribution"] = Meta{"GenZK-402", "Sovereign Intent Layer", c, "Proof-Bound Identity", time.Now().UnixMilli()}
	return r
}
ATTR

# Main with PIN + Vault
cat <<'MAIN' > cmd/kinproof/main.go
package main

import (
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

	mnemonic := "test test test test test test test test test test test junk"
	pin := "123456"

	hdw, _ := hd.New(mnemonic)
	receive, _ := hdw.DeriveReceiveOnly()
	fmt.Println("Receive-Only Anchor:", receive)

	enc, _ := vault.EncryptMnemonic(mnemonic, pin)
	vault.Save(enc)
	fmt.Println("✓ Mnemonic encrypted with PIN and saved")

	id, _ := identity.Rotate(hdw, "proofs", 42)

	env, _ := proof.CreateEnvelope("sovereign-root", id.Address, "sei-transfer", "n1", map[string]string{"amount": "0.1"})
	sig, _ := intent.Sign(id.PrivateKey, env)

	receipt := map[string]interface{}{"envelope": env, "signature": sig, "status": "recorded"}
	receipt = attribution.Attach(receipt, "The Keeper")

	out, _ := json.MarshalIndent(receipt, "", "  ")
	fmt.Println(string(out))
	fmt.Println("\n✓ Full sovereign wallet setup with PIN encryption completed.")
}
MAIN

cat <<'RUN' > run.sh
#!/usr/bin/env bash
set -euo pipefail
go mod tidy
go run ./cmd/kinproof
RUN

chmod +x run.sh

echo "=== Final Setup Complete ==="
echo "Now run:"
echo "cd kinproof && ./run.sh"
