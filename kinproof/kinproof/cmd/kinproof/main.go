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
