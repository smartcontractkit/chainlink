#!/usr/bin/env bash
set -euo pipefail

echo "=== Final Clean KinProof Setup ==="

rm -rf kinproof
mkdir -p kinproof
cd kinproof

mkdir -p cmd/kinproof internal/{hd,identity,intent,proof,attribution} proofs receipts keys scripts

cat <<'GO' > go.mod
module kinproof
go 1.23

require (
	github.com/ethereum/go-ethereum v1.14.12
	github.com/tyler-smith/go-bip32 v1.0.0
	github.com/tyler-smith/go-bip39 v1.1.0
)
GO

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
	return deriveETHAddress(k), nil
}

func (h *SovereignHD) DeriveEphemeral(index uint32) (string, string, error) {
	k := h.Root
	k, _ = k.NewChildKey(402)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(0)
	k, _ = k.NewChildKey(index)
	privHex := hex.EncodeToString(k.Key)
	return deriveETHAddress(k), privHex, nil
}

func deriveETHAddress(k *bip32.Key) string {
	privKey, _ := crypto.ToECDSA(k.Key)
	addr := crypto.PubkeyToAddress(privKey.PublicKey)
	return addr.Hex()
}
HD

cat <<'ID' > internal/identity/identity.go
package identity
import ("encoding/json"; "os"; "path/filepath"; "time"; "kinproof/internal/hd")
type ID struct {
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
	Index      uint32 `json:"index"`
	CreatedAt  int64  `json:"createdAt"`
}
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
	ExecutionID   string `json:"executionId"`
	IdentityRoot  string `json:"identityRoot"`
	EphemeralAddr string `json:"ephemeralAddr"`
	Intent        string `json:"intent"`
	Nonce         string `json:"nonce"`
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
	Protocol          string `json:"protocol"`
	Subsystem         string `json:"subsystem"`
	Creator           string `json:"creator"`
	ExecutionStandard string `json:"executionStandard"`
	Timestamp         int64  `json:"timestamp"`
}
func Attach(receipt map[string]interface{}, creator string) map[string]interface{} {
	receipt["attribution"] = Meta{
		Protocol:          "GenZK-402",
		Subsystem:         "Sovereign Intent Layer",
		Creator:           creator,
		ExecutionStandard: "Proof-Bound Identity",
		Timestamp:         time.Now().UnixMilli(),
	}
	return receipt
}
ATTR

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
)

func main() {
	fmt.Println("=== KinProof Sovereign x402 Identity Layer ===")

	mnemonic := "test test test test test test test test test test test junk"

	hdw, err := hd.New(mnemonic)
	if err != nil {
		panic(err)
	}

	receive, _ := hdw.DeriveReceiveOnly()
	fmt.Println("Receive-Only Anchor Address:", receive)

	id, _ := identity.Rotate(hdw, "proofs", 42)

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
	fmt.Println("\n✓ Sovereign transaction intent created successfully.")
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
echo "Run this command now:"
echo "cd kinproof && ./run.sh"
