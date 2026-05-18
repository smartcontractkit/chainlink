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
