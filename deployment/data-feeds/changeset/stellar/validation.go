package stellar

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/smartcontractkit/chainlink-deployments-framework/datastore"
	"github.com/stellar/go-stellar-sdk/strkey"
)

// validateContract restricts a request's target to the two DF contract types.
func validateContract(t datastore.ContractType) error {
	if t != CacheContract && t != ProxyContract {
		return fmt.Errorf("unsupported contract type %q: must be %q or %q", t, CacheContract, ProxyContract)
	}
	return nil
}

// validateAddress accepts a Stellar account (G...) or contract (C...) strkey.
func validateAddress(s string) error {
	if strkey.IsValidEd25519PublicKey(s) {
		return nil
	}
	if _, err := strkey.Decode(strkey.VersionByteContract, s); err == nil {
		return nil
	}
	return fmt.Errorf("%q is not a valid Stellar account or contract address", s)
}

// dataIDsToBytes converts 0x-prefixed hex feed IDs to [16]byte, left-justified
// with trailing zero padding (data IDs are canonically left-aligned).
func dataIDsToBytes(ids []string) ([][16]byte, error) {
	out := make([][16]byte, 0, len(ids))
	for _, id := range ids {
		v, ok := new(big.Int).SetString(id, 0)
		if !ok {
			return nil, fmt.Errorf("invalid data_id: %q", id)
		}
		if v.BitLen() > 128 {
			return nil, fmt.Errorf("data_id too long: %q (%d bits)", id, v.BitLen())
		}
		var b [16]byte
		copy(b[:], v.Bytes())
		out = append(out, b)
	}
	return out, nil
}

// workflowNameToBytes right-pads an ASCII workflow name into [10]byte.
func workflowNameToBytes(s string) ([10]byte, error) {
	var out [10]byte
	if len(s) > len(out) {
		return out, fmt.Errorf("workflow name %q exceeds %d bytes", s, len(out))
	}
	copy(out[:], s)
	return out, nil
}

// workflowOwnerToBytes decodes a 20-byte 0x-prefixed hex workflow owner.
func workflowOwnerToBytes(hexStr string) ([20]byte, error) {
	var out [20]byte
	b, err := hex.DecodeString(strings.TrimPrefix(hexStr, "0x"))
	if err != nil {
		return out, fmt.Errorf("invalid workflow owner %q: %w", hexStr, err)
	}
	if len(b) != len(out) {
		return out, fmt.Errorf("workflow owner must be %d bytes, got %d", len(out), len(b))
	}
	copy(out[:], b)
	return out, nil
}
