package secp256k1

import (
	"database/sql/driver"
	"fmt"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/chainlink-integrations/evm/utils"
)

// PublicKey is a secp256k1 point in compressed format
type PublicKey [CompressedPublicKeyLength]byte

// CompressedPublicKeyLength is the length of a secp256k1 public key's x
// coordinate as a uint256, concatenated with 00 if y is even, 01 if odd.
const CompressedPublicKeyLength = 33

func init() {
	if CompressedPublicKeyLength != (&Secp256k1{}).Point().MarshalSize() {
		panic("disparity in expected public key lengths")
	}
}

// Set sets k to the public key represented by l
func (k *PublicKey) Set(l PublicKey) {
	if copy(k[:], l[:]) != CompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy entire public key %x to %x", l, k))
	}
}

// Point returns the secp256k1 point corresponding to k
func (k *PublicKey) Point() (kyber.Point, error) {
	p := (&Secp256k1{}).Point()
	return p, p.UnmarshalBinary(k[:])
}

// NewPublicKeyFromHex returns the PublicKey encoded by 0x-hex string hex, or errors
func NewPublicKeyFromHex(hex string) (PublicKey, error) {
	rawKey, err := hexutil.Decode(hex)
	if err != nil {
		return PublicKey{}, err
	}
	return NewPublicKeyFromBytes(rawKey)
}

// NewPublicKeyFromBytes returns the PublicKey built from the given bytes, or errors
func NewPublicKeyFromBytes(rawKey []byte) (PublicKey, error) {
	if l := len(rawKey); l != CompressedPublicKeyLength {
		return PublicKey{}, fmt.Errorf(
			"wrong length for public key: %s of length %d", rawKey, l)
	}
	var k PublicKey
	if c := copy(k[:], rawKey); c != CompressedPublicKeyLength {
		panic(fmt.Errorf("failed to copy entire key to return value"))
	}
	return k, nil
}

// SetFromHex sets k to the public key represented by hex, which must represent
// the compressed binary format
func (k *PublicKey) SetFromHex(hex string) error {
	nk, err := NewPublicKeyFromHex(hex)
	if err != nil {
		return err
	}
	k.Set(nk)
	return nil
}

// String returns k's binary compressed representation, as 0x-hex
func (k PublicKey) String() string {
	return hexutil.Encode(k[:])
}

// StringUncompressed returns k's binary uncompressed representation, as 0x-hex
func (k *PublicKey) StringUncompressed() (string, error) {
	p, err := k.Point()
	if err != nil {
		return "", err
	}
	return hexutil.Encode(LongMarshal(p)), nil
}

// Hash returns the solidity Keccak256 hash of k. Corresponds to hashOfKey on
// VRFCoordinator.
func (k *PublicKey) Hash() (common.Hash, error) {
	p, err := k.Point()
	if err != nil {
		return common.Hash{}, err
	}
	return utils.MustHash(string(LongMarshal(p))), nil
}

// MustHash is like Hash, but panics on error. Useful for testing.
func (k *PublicKey) MustHash() common.Hash {
	hash, err := k.Hash()
	if err != nil {
		panic(fmt.Sprintf("Failed to compute hash of public vrf key %v", k))
	}
	return hash
}

// Address returns the Ethereum address of k or 0 if the key is invalid
func (k *PublicKey) Address() common.Address {
	hash, err := k.Hash()
	if err != nil {
		return common.Address{}
	}
	return common.BytesToAddress(hash.Bytes()[12:])
}

// IsZero returns true iff k is the zero value for PublicKey
func (k *PublicKey) IsZero() bool {
	return *k == PublicKey{}
}

// MarshalText renders k as a text string
func (k PublicKey) MarshalText() ([]byte, error) {
	return []byte(k.String()), nil
}

// UnmarshalText reads a PublicKey into k from text, or errors
func (k *PublicKey) UnmarshalText(text []byte) error {
	if err := k.SetFromHex(string(text)); err != nil {
		return errors.Wrapf(err, "while parsing %s as public key", text)
	}
	return nil
}

// Value marshals PublicKey to be saved in the DB
func (k PublicKey) Value() (driver.Value, error) {
	return k.String(), nil
}

// Scan reconstructs a PublicKey from a DB record of it.
func (k *PublicKey) Scan(value interface{}) error {
	rawKey, ok := value.(string)
	if !ok {
		return errors.Wrap(fmt.Errorf("unable to convert %+v of type %T to PublicKey", value, value), "scan failure")
	}
	if err := k.SetFromHex(rawKey); err != nil {
		return errors.Wrapf(err, "while scanning %s as PublicKey", rawKey)
	}
	return nil
}
