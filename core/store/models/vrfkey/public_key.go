package vrfkey

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"go.dedis.ch/kyber/v3"

	"chainlink/core/services/signatures/secp256k1"
	"chainlink/core/utils"
)

const UncompressedPublicKeyLength = 64

// Set sets k to the public key represented by l
func (k *PublicKey) Set(l PublicKey) {
	if copy(k[:], l[:]) != UncompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy entire public key %x to %x", l, k))
	}
}

// Point returns the secp256k1 point corresponding to k
func (k *PublicKey) Point() (kyber.Point, error) {
	return secp256k1.LongUnmarshal(k[:])
}

// NewPublicKey returns the PublicKey corresponding to rawKey
func NewPublicKey(rawKey [64]byte) *PublicKey {
	rv := PublicKey(rawKey)
	return &rv
}

// NewPublicKeyFromHex returns the PublicKey encoded by 0x-hex string hex, or errors
func NewPublicKeyFromHex(hex string) (*PublicKey, error) {
	rawKey, err := hexutil.Decode(hex)
	if err != nil {
		return nil, err
	}
	if l := len(rawKey); l != UncompressedPublicKeyLength {
		return nil, fmt.Errorf("wrong length for public key: %s of length %d", rawKey, l)
	}
	k := &PublicKey{}
	if c := copy(k[:], rawKey[:]); c != UncompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy entire key to return value"))
	}
	return k, err
}

// SetFromHex sets k to the public key represented by hex, which must represent
// the uncompressed binary format
func (k *PublicKey) SetFromHex(hex string) error {
	nk, err := NewPublicKeyFromHex(hex)
	if err != nil {
		return err
	}
	k.Set(*nk)
	return nil
}

// String returns k's binary uncompressed representation, as 0x-hex
func (k *PublicKey) String() string {
	return "0x" + hex.EncodeToString(k[:])
}

// Hash returns the solidity Keccak256 hash of k. Corresponds to hashOfKey on
// VRFCoordinator.
func (k *PublicKey) Hash() common.Hash {
	rv, err := utils.Keccak256(k[:])
	if err != nil {
		panic(errors.Wrapf(err, "while computing hash of public key %s", k))
	}
	return common.BytesToHash(rv)
}

// Address returns the Ethereum address of k
func (k *PublicKey) Address() common.Address {
	return common.BytesToAddress(k.Hash().Bytes()[12:])
}
