// Package utils is used for the common functions for dealing with
// conversion to and from hex, bytes, and strings, formatting time.
package utils

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jpillora/backoff"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	null "gopkg.in/guregu/null.v3"
)

const (
	// HumanTimeFormat is the predefined layout for use in Time.Format and time.Parse
	HumanTimeFormat = "2006-01-02 15:04:05 MST"
	// EVMWordByteLen the length of an EVM Word Byte
	EVMWordByteLen = 32
	// EVMWordHexLen the length of an EVM Word Hex
	EVMWordHexLen = EVMWordByteLen * 2
)

// ZeroAddress is an empty address, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

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

// EncodeTxToHex converts the given Ethereum Transaction type and
// returns its hex-value string.
func EncodeTxToHex(tx *types.Transaction) (string, error) {
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return "", err
	}
	return hexutil.Encode(rlp.Bytes()), nil
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
	return t.Sub(time.Now())
}

// FormatJSON applies indent to format a JSON response.
func FormatJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// NewBytes32Length holds the length of bytes needed for Bytes32ID.
const NewBytes32Length = 32

// NewBytes32ID returns a randomly generated UUID that conforms to
// Ethereum bytes32.
func NewBytes32ID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// NewBytes32Secret returns a randomly generated string which conforms to
// Ethereum bytes32.
func NewBytes32Secret() (string, error) {
	b := make([]byte, 24)
	_, err := rand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "new secret bytes32")
	}
	return base64.StdEncoding.EncodeToString(b), nil
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
	return len(str) >= 2 && str[0] == '0' && str[1] == 'x'
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

// ToFilterQueryFor returns a struct that encapsulates desired arguments used to filter
// event logs.
func ToFilterQueryFor(fromBlock *big.Int, addresses []common.Address) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: WithoutZeroAddresses(addresses),
	}
}

// ToFilterArg filters logs with the given FilterQuery
// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L363
func ToFilterArg(q ethereum.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
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

// ConstantSleeper is to assist with reattempts with
// the same sleep duration.
type ConstantSleeper struct {
	interval time.Duration
}

// NewConstantSleeper returns a ConstantSleeper that is configured to
// sleep for a constant duration based on the input.
func NewConstantSleeper(d time.Duration) ConstantSleeper {
	return ConstantSleeper{interval: d}
}

// Reset is a no op since sleep time is constant.
func (cs ConstantSleeper) Reset() {}

// Sleep waits for the given duration before reattempting.
func (cs ConstantSleeper) Sleep() {
	time.Sleep(cs.interval)
}

// After returns the duration.
func (cs ConstantSleeper) After() time.Duration {
	return cs.interval
}

// Duration returns the duration value.
func (cs ConstantSleeper) Duration() time.Duration {
	return cs.interval
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

// MaxUint64 finds the maximum value of a list of uint64s.
func MaxUint64(uints ...uint64) uint64 {
	var max uint64
	for _, n := range uints {
		if n > max {
			max = n
		}
	}
	return max
}

// MaxUint32 finds the maximum value of a list of uint32s.
func MaxUint32(uints ...uint32) uint32 {
	var max uint32
	for _, n := range uints {
		if n > max {
			max = n
		}
	}
	return max
}

// MaxInt finds the maximum value of a list of ints.
func MaxInt(ints ...int) int {
	var max int
	for _, n := range ints {
		if n > max {
			max = n
		}
	}
	return max
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
				return nil, fmt.Errorf("Unable to coerce key %T %v to a string", k, k)
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

// FileExists returns true if a file at the passed string exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
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

// RemoveContents removes everything in a directory but not the directory
// itself.
// https://stackoverflow.com/questions/33450980/how-to-remove-all-contents-of-a-directory-using-golang
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// FilesInDir returns an array of filenames in the directory.
func FilesInDir(dir string) ([]string, error) {
	f, err := os.Open(dir)
	if err != nil {
		return []string{}, err
	}
	defer f.Close()

	r, err := f.Readdirnames(-1)
	if err != nil {
		return []string{}, err
	}

	return r, nil
}

// FileContents returns the contents of a file as a string.
func FileContents(path string) (string, error) {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// JustError takes a tuple and returns the last entry, the error.
func JustError(_ interface{}, err error) error {
	return err
}
