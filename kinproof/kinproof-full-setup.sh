#!/usr/bin/env bash
set -euo pipefail

echo "=== KinProof Sovereign Wallet Full Setup (Username + PIN) ==="

rm -rf kinproof
mkdir -p kinproof
cd kinproof

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

# Vault with PIN
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
VAULT

# Main with interactive setup
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
	fmt.Println("=== KinProof Sovereign Wallet Setup ===")

	// === Username Setup ===
	var username string
	fmt.Print("Enter your username (e.g. @totalwine2339): ")
	fmt.Scanln(&username)
	if username == "" {
		username = "TheKeeper"
	}
	fmt.Println("Username set:", username)

	// === PIN Setup ===
	var pin string
	fmt.Print("Enter your PIN (6+ characters): ")
	fmt.Scanln(&pin)
	if len(pin) < 4 {
		pin = "123456"
		fmt.Println("Using default PIN for demo.")
	}

	// === Mnemonic (Change this in production) ===
	mnemonic := "test test test test test test test test test test test junk"

	// === Create HD Wallet ===
	hdw, _ := hd.New(mnemonic)
	receive, _ := hdw.DeriveReceiveOnly()
	fmt.Println("Receive-Only Anchor Address:", receive)

	// === Encrypt Vault ===
	err := vault.EncryptAndSave(mnemonic, pin, username)
	if err == nil {
		fmt.Println("✓ Mnemonic encrypted with PIN and saved to keys/vault.json")
	}

	// === Create Execution Intent ===
	id, _ := identity.Rotate(hdw, "proofs", 42)

	env, _ := proof.CreateEnvelope("sovereign-root", id.Address, "sei-transfer", "n1", map[string]string{
		"to":     "0x742d35Cc6634C0532925a3b844Bc454e4438f44e",
		"amount": "0.1",
	})

	sig, _ := intent.Sign(id.PrivateKey, env)

	receipt := map[string]interface{}{
		"username":  username,
		"envelope":  env,
		"signature": sig,
		"status":    "intent_recorded",
	}
	receipt = attribution.Attach(receipt, "The Keeper")

	out, _ := json.MarshalIndent(receipt, "", "  ")
	fmt.Println(string(out))
	fmt.Println("\n✓ Full sovereign wallet setup completed with username + PIN.")
}
MAIN

# Create remaining minimal files
cat <<'MIN' > internal/identity/identity.go
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
MIN

cat <<'MIN' > internal/intent/intent.go
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
MIN

cat <<'MIN' > internal/proof/proof.go
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
MIN

cat <<'MIN' > internal/attribution/attribution.go
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
MIN

cat <<'RUN' > run.sh
#!/usr/bin/env bash
set -euo pipefail
go mod tidy
go run ./cmd/kinproof
RUN

chmod +x run.sh

echo "=== Full Setup Complete ==="
echo "Run the following command to start setup:"
echo "cd kinproof && ./run.sh"
