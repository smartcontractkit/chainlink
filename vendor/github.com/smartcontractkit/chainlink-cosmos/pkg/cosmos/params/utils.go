package params

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"golang.org/x/crypto/ripemd160" //nolint: staticcheck
)

// Creates a bech32 address from a hex-encoded secp256k1 public key.
func CreateBech32Address(pubKey, accountPrefix string) (string, error) {
	pubKeyBytes, err := hex.DecodeString(pubKey)
	if err != nil {
		return "", err
	}

	if len(pubKeyBytes) != secp256k1.PubKeySize {
		return "", errors.New("length of pubkey is incorrect")
	}

	sha := sha256.Sum256(pubKeyBytes)
	hasherRIPEMD160 := ripemd160.New()
	hasherRIPEMD160.Write(sha[:]) // does not error
	address := crypto.Address(hasherRIPEMD160.Sum(nil))

	bech32Addr, err := bech32.ConvertAndEncode(accountPrefix, address)
	if err != nil {
		return "", err
	}
	return bech32Addr, nil
}
