package crypto

import (
	"fmt"
	"strings"

	"github.com/smartcontractkit/chainlink-common/keystore"
	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/stellarkey"
)

// StellarKey holds a generated Stellar (ed25519) keystore key for a Local CRE
// node. Account is the canonical StrKey "G..." address; the same key is used as
// the TXM transmitter account. Mirrors AptosKey (see aptos.go); Stellar
// addresses are StrKey, not hex, so no re-encoding normalization is required.
type StellarKey struct {
	EncryptedJSON []byte
	PublicKey     string
	Account       string
	Password      string
}

func NewStellarKey(password string) (*StellarKey, error) {
	key, err := stellarkey.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create stellar key: %w", err)
	}

	enc, err := key.ToEncryptedJSON(password, keystore.DefaultScryptParams)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt stellar key: %w", err)
	}

	account, err := NormalizeStellarAccount(key.Account())
	if err != nil {
		return nil, fmt.Errorf("failed to normalize stellar account: %w", err)
	}

	return &StellarKey{
		EncryptedJSON: enc,
		PublicKey:     key.PublicKeyStr(),
		Account:       account,
		Password:      password,
	}, nil
}

// NormalizeStellarAccount validates and canonicalizes a StrKey account address.
// stellarkey.Account() already returns a canonical "G..." StrKey, so this is a
// light guard for externally-supplied addresses. TODO(write-path): use
// github.com/stellar/go/strkey for full Base32+CRC validation once that
// dependency is pulled into system-tests/lib for the forwarder feature.
func NormalizeStellarAccount(raw string) (string, error) {
	addr := strings.TrimSpace(raw)
	if len(addr) != 56 || !strings.HasPrefix(addr, "G") {
		return "", fmt.Errorf("invalid stellar account address %q: expected 56-char StrKey starting with 'G'", raw)
	}
	return addr, nil
}
