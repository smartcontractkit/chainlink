package utils

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jpillora/backoff"
	"golang.org/x/crypto/sha3"

	"github.com/smartcontractkit/chainlink-common/pkg/utils/hex"
)

// EVMWordByteLen the length of an EVM Word Byte
const EVMWordByteLen = 32

// DefaultQueryTimeout is the default timeout for database queries
const DefaultQueryTimeout = 10 * time.Second

// ZeroAddress is an address of all zeroes, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

// EmptyHash is a hash of all zeroes, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000000000000000000000000000
var EmptyHash = common.Hash{}

func RandomAddress() common.Address {
	b := make([]byte, 20)
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return common.BytesToAddress(b)
}

func RandomHash() common.Hash {
	b := make([]byte, 32)
	_, _ = rand.Read(b) // Assignment for errcheck. Only used in tests so we can ignore.
	return common.BytesToHash(b)
}

// IsEmptyAddress checks that the address is empty, synonymous with the zero
// account/address. No logs can come from this address, as there is no contract
// present there.
//
// See https://stackoverflow.com/questions/48219716/what-is-address0-in-solidity
// for the more info on the zero address.
func IsEmptyAddress(addr common.Address) bool {
	return addr == ZeroAddress
}

func RandomBytes32() (r [32]byte) {
	b := make([]byte, 32)
	_, _ = rand.Read(b[:]) // Assignment for errcheck. Only used in tests so we can ignore.
	copy(r[:], b)
	return
}

func Bytes32ToSlice(a [32]byte) (r []byte) {
	r = append(r, a[:]...)
	return
}

// Uint256ToBytes is x represented as the bytes of a uint256
func Uint256ToBytes(x *big.Int) (uint256 []byte, err error) {
	if x.Cmp(MaxUint256) > 0 {
		return nil, fmt.Errorf("too large to convert to uint256")
	}
	uint256 = common.LeftPadBytes(x.Bytes(), EVMWordByteLen)
	if x.Cmp(big.NewInt(0).SetBytes(uint256)) != 0 {
		panic("failed to round-trip uint256 back to source big.Int")
	}
	return uint256, err
}

// NewHash return random Keccak256
func NewHash() common.Hash {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(b)
}

// PadByteToHash returns a hash with zeros padded on the left of the given byte.
func PadByteToHash(b byte) common.Hash {
	var h [32]byte
	h[31] = b
	return h
}

// Uint256ToBytes32 returns the bytes32 encoding of the big int provided
func Uint256ToBytes32(n *big.Int) []byte {
	if n.BitLen() > 256 {
		panic("vrf.uint256ToBytes32: too big to marshal to uint256")
	}
	return common.LeftPadBytes(n.Bytes(), 32)
}

// MustHash returns the keccak256 hash, or panics on failure.
func MustHash(in string) common.Hash {
	out, err := Keccak256([]byte(in))
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(out)
}

// HexToUint256 returns the uint256 represented by s, or an error if it doesn't
// represent one.
func HexToUint256(s string) (*big.Int, error) {
	rawNum, err := hexutil.Decode(s)
	if err != nil {
		return nil, fmt.Errorf("error while parsing %s as hex: %w", s, err)
	}
	rv := big.NewInt(0).SetBytes(rawNum) // can't be negative number
	if err := CheckUint256(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

var zero = big.NewInt(0)

// CheckUint256 returns an error if n is out of bounds for a uint256
func CheckUint256(n *big.Int) error {
	if n.Cmp(zero) < 0 || n.Cmp(MaxUint256) >= 0 {
		return fmt.Errorf("number out of range for uint256")
	}
	return nil
}

// Keccak256 is a simplified interface for the legacy SHA3 implementation that
// Ethereum uses.
func Keccak256(in []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(in)
	return hash.Sum(nil), err
}

func Keccak256Fixed(in []byte) [32]byte {
	hash := sha3.NewLegacyKeccak256()
	// Note this Keccak256 cannot error https://github.com/golang/crypto/blob/master/sha3/sha3.go#L126
	// if we start supporting hashing algos which do, we can change this API to include an error.
	hash.Write(in)
	var h [32]byte
	copy(h[:], hash.Sum(nil))
	return h
}

// EIP55CapitalizedAddress returns true iff possibleAddressString has the correct
// capitalization for an Ethereum address, per EIP 55
func EIP55CapitalizedAddress(possibleAddressString string) bool {
	possibleAddressString = hex.EnsurePrefix(possibleAddressString)
	EIP55Capitalized := common.HexToAddress(possibleAddressString).Hex()
	return possibleAddressString == EIP55Capitalized
}

// ParseEthereumAddress returns addressString as a go-ethereum Address, or an
// error if it's invalid, e.g. if EIP 55 capitalization check fails
func ParseEthereumAddress(addressString string) (common.Address, error) {
	if !common.IsHexAddress(addressString) {
		return common.Address{}, fmt.Errorf(
			"not a valid Ethereum address: %s", addressString)
	}
	address := common.HexToAddress(addressString)
	if !EIP55CapitalizedAddress(addressString) {
		return common.Address{}, fmt.Errorf(
			"%s treated as Ethereum address, but it has an invalid capitalization! "+
				"The correctly-capitalized address would be %s, but "+
				"check carefully before copying and pasting! ",
			addressString, address.Hex())
	}
	return address, nil
}

// NewRedialBackoff is a standard backoff to use for redialling or reconnecting to
// unreachable network endpoints
func NewRedialBackoff() backoff.Backoff {
	return backoff.Backoff{
		Min:    1 * time.Second,
		Max:    15 * time.Second,
		Jitter: true,
	}
}

// RetryWithBackoff retries the sleeper and backs off if not Done
func RetryWithBackoff(ctx context.Context, fn func() (retry bool)) {
	sleeper := NewBackoffSleeper()
	sleeper.Reset()
	for {
		retry := fn()
		if !retry {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-time.After(sleeper.After()):
			continue
		}
	}
}

// NewBackoffSleeper returns a BackoffSleeper that is configured to
// sleep for 0 seconds initially, then backs off from 1 second minimum
// to 10 seconds maximum.
func NewBackoffSleeper() *BackoffSleeper {
	return &BackoffSleeper{
		Backoff: backoff.Backoff{
			Min: 1 * time.Second,
			Max: 10 * time.Second,
		},
	}
}

// BackoffSleeper is a sleeper that backs off on subsequent attempts.
type BackoffSleeper struct {
	backoff.Backoff
	beenRun atomic.Bool
}

// Sleep waits for the given duration, incrementing the back off.
func (bs *BackoffSleeper) Sleep() {
	if bs.beenRun.CompareAndSwap(false, true) {
		return
	}
	time.Sleep(bs.Backoff.Duration())
}

// After returns the duration for the next stop, and increments the backoff.
func (bs *BackoffSleeper) After() time.Duration {
	if bs.beenRun.CompareAndSwap(false, true) {
		return 0
	}
	return bs.Backoff.Duration()
}

// Duration returns the current duration value.
func (bs *BackoffSleeper) Duration() time.Duration {
	if !bs.beenRun.Load() {
		return 0
	}
	return bs.ForAttempt(bs.Attempt())
}

// Reset resets the backoff intervals.
func (bs *BackoffSleeper) Reset() {
	bs.beenRun.Store(false)
	bs.Backoff.Reset()
}

// RandUint256 generates a random bigNum up to 2 ** 256 - 1
func RandUint256() *big.Int {
	n, err := rand.Int(rand.Reader, MaxUint256)
	if err != nil {
		panic(err)
	}
	return n
}
