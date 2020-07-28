// Package utils is used for the common functions for dealing with
// conversion to and from hex, bytes, and strings, formatting time.
package utils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/smartcontractkit/chainlink/core/logger"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	null "gopkg.in/guregu/null.v3"
)

const (
	// DefaultSecretSize is the entroy in bytes to generate a base64 string of 64 characters.
	DefaultSecretSize = 48
	// EVMWordByteLen the length of an EVM Word Byte
	EVMWordByteLen = 32
	// EVMWordHexLen the length of an EVM Word Hex
	EVMWordHexLen = EVMWordByteLen * 2
)

// ZeroAddress is an address of all zeroes, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

// EmptyHash is a hash of all zeroes, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000000000000000000000000000
var EmptyHash = common.Hash{}

// WithoutZeroAddresses returns a list of addresses excluding the zero address.
func WithoutZeroAddresses(addresses []common.Address) []common.Address {
	var withoutZeros []common.Address
	for _, address := range addresses {
		if address != ZeroAddress {
			withoutZeros = append(withoutZeros, address)
		}
	}
	return withoutZeros
}

// HexToUint64 converts a given hex string to 64-bit unsigned integer.
func HexToUint64(hex string) (uint64, error) {
	return strconv.ParseUint(RemoveHexPrefix(hex), 16, 64)
}

// Uint64ToHex converts the given uint64 value to a hex-value string.
func Uint64ToHex(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

var maxUint256 = common.HexToHash("0x" + strings.Repeat("f", 64)).Big()

// Uint256ToBytes(x) is x represented as the bytes of a uint256
func Uint256ToBytes(x *big.Int) (uint256 []byte, err error) {
	if x.Cmp(maxUint256) > 0 {
		return nil, fmt.Errorf("too large to convert to uint256")
	}
	uint256 = common.LeftPadBytes(x.Bytes(), EVMWordByteLen)
	if x.Cmp(big.NewInt(0).SetBytes(uint256)) != 0 {
		panic("failed to round-trip uint256 back to source big.Int")
	}
	return uint256, err
}

// ISO8601UTC formats given time to ISO8601.
func ISO8601UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// NullISO8601UTC returns formatted time if valid, empty string otherwise.
func NullISO8601UTC(t null.Time) string {
	if t.Valid {
		return ISO8601UTC(t.Time)
	}
	return ""
}

// DurationFromNow returns the amount of time since the Time
// field was last updated.
func DurationFromNow(t time.Time) time.Duration {
	return time.Until(t)
}

// FormatJSON applies indent to format a JSON response.
func FormatJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// NewBytes32ID returns a randomly generated UUID that conforms to
// Ethereum bytes32.
func NewBytes32ID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// NewSecret returns a new securely random sequence of n bytes of entropy.  The
// result is a base64 encoded string.
//
// Panics on failed attempts to read from system's PRNG.
func NewSecret(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(errors.Wrap(err, "generating secret failed"))
	}
	return base64.StdEncoding.EncodeToString(b)
}

// RemoveHexPrefix removes the prefix (0x) of a given hex string.
func RemoveHexPrefix(str string) string {
	if HasHexPrefix(str) {
		return str[2:]
	}
	return str
}

// HasHexPrefix returns true if the string starts with 0x.
func HasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// DecodeEthereumTx takes an RLP hex encoded Ethereum transaction and
// returns a Transaction struct with all the fields accessible.
func DecodeEthereumTx(hex string) (types.Transaction, error) {
	var tx types.Transaction
	b, err := hexutil.Decode(hex)
	if err != nil {
		return tx, err
	}
	return tx, rlp.DecodeBytes(b, &tx)
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

// StringToHex converts a standard string to a hex encoded string.
func StringToHex(in string) string {
	return AddHexPrefix(hex.EncodeToString([]byte(in)))
}

// AddHexPrefix adds the previx (0x) to a given hex string.
func AddHexPrefix(str string) string {
	if len(str) < 2 || len(str) > 1 && strings.ToLower(str[0:2]) != "0x" {
		str = "0x" + str
	}
	return str
}

// Sleeper interface is used for tasks that need to be done on some
// interval, excluding Cron, like reconnecting.
type Sleeper interface {
	Reset()
	Sleep()
	After() time.Duration
	Duration() time.Duration
}

// BackoffSleeper is a sleeper that backs off on subsequent attempts.
type BackoffSleeper struct {
	backoff.Backoff
	beenRun bool
}

// NewBackoffSleeper returns a BackoffSleeper that is configured to
// sleep for 0 seconds initially, then backs off from 1 second minimum
// to 10 seconds maximum.
func NewBackoffSleeper() *BackoffSleeper {
	return &BackoffSleeper{Backoff: backoff.Backoff{
		Min: 1 * time.Second,
		Max: 10 * time.Second,
	}}
}

// Sleep waits for the given duration, incrementing the back off.
func (bs *BackoffSleeper) Sleep() {
	if !bs.beenRun {
		time.Sleep(0)
		bs.beenRun = true
		return
	}
	time.Sleep(bs.Backoff.Duration())
}

// After returns the duration for the next stop, and increments the backoff.
func (bs *BackoffSleeper) After() time.Duration {
	if !bs.beenRun {
		bs.beenRun = true
		return 0
	}
	return bs.Backoff.Duration()
}

// Duration returns the current duration value.
func (bs *BackoffSleeper) Duration() time.Duration {
	if !bs.beenRun {
		return 0
	}
	return bs.ForAttempt(bs.Attempt())
}

// Reset resets the backoff intervals.
func (bs *BackoffSleeper) Reset() {
	bs.beenRun = false
	bs.Backoff.Reset()
}

func RetryWithBackoff(chCancel <-chan struct{}, errPrefix string, fn func() error) (aborted bool) {
	sleeper := NewBackoffSleeper()
	sleeper.Reset()
	for {
		err := fn()
		if err == nil {
			return false
		}

		logger.Errorf("%v: %v", errPrefix, err)
		select {
		case <-chCancel:
			return true
		case <-time.After(sleeper.After()):
			continue
		}
	}
}

// MinBigs finds the minimum value of a list of big.Ints.
func MinBigs(first *big.Int, bigs ...*big.Int) *big.Int {
	min := first
	for _, n := range bigs {
		if min.Cmp(n) > 0 {
			min = n
		}
	}
	return min
}

// MaxBigs finds the maximum value of a list of big.Ints.
func MaxBigs(first *big.Int, bigs ...*big.Int) *big.Int {
	max := first
	for _, n := range bigs {
		if max.Cmp(n) < 0 {
			max = n
		}
	}
	return max
}

// MaxUint32 finds the maximum value of a list of uint32s.
func MaxUint32(first uint32, uints ...uint32) uint32 {
	max := first
	for _, n := range uints {
		if n > max {
			max = n
		}
	}
	return max
}

// MaxInt finds the maximum value of a list of ints.
func MaxInt(first int, ints ...int) int {
	max := first
	for _, n := range ints {
		if n > max {
			max = n
		}
	}
	return max
}

// MinUint finds the minimum value of a list of uints.
func MinUint(first uint, vals ...uint) uint {
	min := first
	for _, n := range vals {
		if n < min {
			min = n
		}
	}
	return min
}

// CoerceInterfaceMapToStringMap converts map[interface{}]interface{} (interface maps) to
// map[string]interface{} (string maps) and []interface{} with interface maps to string maps.
// Relevant when serializing between CBOR and JSON.
func CoerceInterfaceMapToStringMap(in interface{}) (interface{}, error) {
	switch typed := in.(type) {
	case map[string]interface{}:
		for k, v := range typed {
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			typed[k] = coerced
		}
		return typed, nil
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range typed {
			coercedKey, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("unable to coerce key %T %v to a string", k, k)
			}
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			m[coercedKey] = coerced
		}
		return m, nil
	case []interface{}:
		r := make([]interface{}, len(typed))
		for i, v := range typed {
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			r[i] = coerced
		}
		return r, nil
	default:
		return in, nil
	}
}

// HashPassword wraps around bcrypt.GenerateFromPassword for a friendlier API.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash wraps around bcrypt.CompareHashAndPassword for a friendlier API.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// Keccak256 is a simplified interface for the legacy SHA3 implementation that
// Ethereum uses.
func Keccak256(in []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(in)
	return hash.Sum(nil), err
}

// Sha256 returns a hexadecimal encoded string of a hashed input
func Sha256(in string) (string, error) {
	hasher := sha3.New256()
	_, err := hasher.Write([]byte(in))
	if err != nil {
		return "", errors.Wrap(err, "sha256 write error")
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// StripBearer removes the 'Bearer: ' prefix from the HTTP Authorization header.
func StripBearer(authorizationStr string) string {
	return strings.TrimPrefix(strings.TrimSpace(authorizationStr), "Bearer ")
}

// IsQuoted checks if the first and last characters are either " or '.
func IsQuoted(input []byte) bool {
	return len(input) >= 2 &&
		((input[0] == '"' && input[len(input)-1] == '"') ||
			(input[0] == '\'' && input[len(input)-1] == '\''))
}

// RemoveQuotes removes the first and last character if they are both either
// " or ', otherwise it is a noop.
func RemoveQuotes(input []byte) []byte {
	if IsQuoted(input) {
		return input[1 : len(input)-1]
	}
	return input
}

// EIP55CapitalizedAddress returns true iff possibleAddressString has the correct
// capitalization for an Ethereum address, per EIP 55
func EIP55CapitalizedAddress(possibleAddressString string) bool {
	if !HasHexPrefix(possibleAddressString) {
		possibleAddressString = "0x" + possibleAddressString
	}
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

// MustHash returns the keccak256 hash, or panics on failure.
func MustHash(in string) common.Hash {
	out, err := Keccak256([]byte(in))
	if err != nil {
		panic(err)
	}
	return common.BytesToHash(out)
}

// LogListeningAddress returns the LogListeningAddress
func LogListeningAddress(address common.Address) string {
	if address == ZeroAddress {
		return "[all]"
	}
	return address.String()
}

// JustError takes a tuple and returns the last entry, the error.
func JustError(_ interface{}, err error) error {
	return err
}

var zero = big.NewInt(0)

// CheckUint256(n) returns an error if n is out of bounds for a uint256
func CheckUint256(n *big.Int) error {
	if n.Cmp(zero) < 0 || n.Cmp(maxUint256) >= 0 {
		return fmt.Errorf("number out of range for uint256")
	}
	return nil
}

// HexToUint256 returns the uint256 represented by s, or an error if it doesn't
// represent one.
func HexToUint256(s string) (*big.Int, error) {
	rawNum, err := hexutil.Decode(s)
	if err != nil {
		return nil, errors.Wrapf(err, "while parsing %s as hex: ", s)
	}
	rv := big.NewInt(0).SetBytes(rawNum) // can't be negative number
	if err := CheckUint256(rv); err != nil {
		return nil, err
	}
	return rv, nil
}

// Uint256ToHex returns the hex representation of n, or error if out of bounds
func Uint256ToHex(n *big.Int) (string, error) {
	if err := CheckUint256(n); err != nil {
		return "", err
	}
	return common.BigToHash(n).Hex(), nil
}

func DecimalFromBigInt(i *big.Int, precision int32) decimal.Decimal {
	return decimal.NewFromBigInt(i, -precision)
}

func WaitGroupChan(wg *sync.WaitGroup) <-chan struct{} {
	chAwait := make(chan struct{})
	go func() {
		defer close(chAwait)
		wg.Wait()
	}()
	return chAwait
}

type DependentAwaiter interface {
	AwaitDependents() <-chan struct{}
	AddDependents(n int)
	DependentReady()
}

type dependentAwaiter struct {
	wg *sync.WaitGroup
	ch <-chan struct{}
}

func NewDependentAwaiter() DependentAwaiter {
	return &dependentAwaiter{
		wg: &sync.WaitGroup{},
	}
}

func (da *dependentAwaiter) AwaitDependents() <-chan struct{} {
	if da.ch == nil {
		da.ch = WaitGroupChan(da.wg)
	}
	return da.ch
}

func (da *dependentAwaiter) AddDependents(n int) {
	da.wg.Add(n)
}

func (da *dependentAwaiter) DependentReady() {
	da.wg.Done()
}

// FIFO queue that discards older items when it reaches its capacity.
type BoundedQueue struct {
	capacity uint
	items    []interface{}
	mu       *sync.RWMutex
}

func NewBoundedQueue(capacity uint) *BoundedQueue {
	return &BoundedQueue{
		capacity: capacity,
		mu:       &sync.RWMutex{},
	}
}

func (q *BoundedQueue) Add(x interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, x)
	if uint(len(q.items)) > q.capacity {
		excess := uint(len(q.items)) - q.capacity
		q.items = q.items[excess:]
	}
}

func (q *BoundedQueue) Take() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	x := q.items[0]
	q.items = q.items[1:]
	return x
}

func (q *BoundedQueue) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return len(q.items) == 0
}

func (q *BoundedQueue) Full() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return uint(len(q.items)) >= q.capacity
}

type BoundedPriorityQueue struct {
	queues     map[uint]*BoundedQueue
	priorities []uint
	capacities map[uint]uint
	mu         *sync.RWMutex
}

func NewBoundedPriorityQueue(capacities map[uint]uint) *BoundedPriorityQueue {
	queues := make(map[uint]*BoundedQueue)
	var priorities []uint
	for priority, capacity := range capacities {
		priorities = append(priorities, priority)
		queues[priority] = NewBoundedQueue(capacity)
	}
	sort.Slice(priorities, func(i, j int) bool { return priorities[i] < priorities[j] })
	return &BoundedPriorityQueue{
		queues:     queues,
		priorities: priorities,
		capacities: capacities,
		mu:         &sync.RWMutex{},
	}
}

func (q *BoundedPriorityQueue) Add(priority uint, x interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	subqueue, exists := q.queues[priority]
	if !exists {
		panic(fmt.Sprintf("nonexistent priority: %v", priority))
	}

	subqueue.Add(x)
}

func (q *BoundedPriorityQueue) Take() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()

	for _, priority := range q.priorities {
		queue := q.queues[priority]
		if queue.Empty() {
			continue
		}
		return queue.Take()
	}
	return nil
}

func (q *BoundedPriorityQueue) Empty() bool {
	q.mu.RLock()
	defer q.mu.RUnlock()

	for _, priority := range q.priorities {
		queue := q.queues[priority]
		if !queue.Empty() {
			return false
		}
	}
	return true
}

// WrapIfError decorates an error with the given message.  It is intended to
// be used with `defer` statements, like so:
//
// func SomeFunction() (err error) {
//     defer WrapIfError(&err, "error in SomeFunction:")
//
//     ...
// }
func WrapIfError(err *error, msg string) {
	if *err != nil {
		errors.Wrap(*err, msg)
	}
}
