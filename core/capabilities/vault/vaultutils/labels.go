package vaultutils

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/tdh2/go/tdh2/tdh2easy"
)

// WorkflowOwnerToLabel converts a workflow owner string to a 32-byte TDH2 ciphertext
// label using the Ethereum address encoding: 12 zero bytes followed by the 20-byte address.
// This matches the label format used when secrets are encrypted for a vault workflow owner,
// including JWT-derived workflow owners.
func WorkflowOwnerToLabel(owner string) [32]byte {
	var label [32]byte
	addr := common.HexToAddress(owner)
	copy(label[12:], addr.Bytes())
	return label
}

// EncryptSecretWithWorkflowOwner encrypts a secret using a TDH2 public key with a label
// derived from a workflow owner's Ethereum address (left-padded to 32 bytes).
func EncryptSecretWithWorkflowOwner(secret string, masterPublicKey *tdh2easy.PublicKey, owner common.Address) (string, error) {
	var label [32]byte
	copy(label[12:], owner.Bytes())
	return encryptWithLabel(secret, masterPublicKey, label)
}

func encryptWithLabel(secret string, masterPublicKey *tdh2easy.PublicKey, label [32]byte) (string, error) {
	cipher, err := tdh2easy.EncryptWithLabel(masterPublicKey, []byte(secret), label)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt secret: %w", err)
	}
	cipherBytes, err := cipher.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to marshal encrypted secret: %w", err)
	}
	return hex.EncodeToString(cipherBytes), nil
}
