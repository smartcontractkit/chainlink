// Package utils is used for the common functions for dealing with
// conversion to and from hex, bytes, and strings, formatting time,
// and basic HTTP authentication.
package utils

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jpillora/backoff"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
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

var weiPerEth = big.NewInt(1e18)

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
	return common.ToHex(rlp.Bytes()), nil
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

// BasicAuthPost sends a POST request to the HTTP client with the given username
// and password to authenticate at the url with contentType and returns a response.
func BasicAuthPost(username, password, url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	return resp, err
}

// BasicAuthGet uses the given username and password to send a GET request
// at the given URL and returns a response.
func BasicAuthGet(username, password, url string, headers ...map[string]string) (*http.Response, error) {
	var h map[string]string
	if len(headers) > 0 {
		h = headers[0]
	} else {
		h = map[string]string{}
	}

	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(username, password)
	for key, value := range h {
		request.Header.Add(key, value)
	}
	resp, err := client.Do(request)
	return resp, err
}

// BasicAuthPatch sends a PATCH request to the HTTP client with the given username
// and password to authenticate at the url with contentType and returns a response.
func BasicAuthPatch(username, password, url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	return resp, err
}

// BasicAuthDelete sends a DELETE request to the HTTP client with the given username
// and password to authenticate at the url with contentType and returns a response.
func BasicAuthDelete(username, password, url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("DELETE", url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	return resp, err
}

// FormatJSON applies indent to format a JSON response.
func FormatJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// GetStringKeys returns an array of strings from the keys of
// Unmarshalled JSON given as input.
// For example, if `j` were our JSON:
//  var value map[string]interface{}
//  err = json.Unmarshal(j, &value)
//  keys := utils.GetStringKeys(value)
func GetStringKeys(v map[string]interface{}) []string {
	keys := make([]string, len(v))

	i := 0
	for k := range v {
		keys[i] = k
		i++
	}

	return keys
}

// NewBytes32ID returns a randomly generated UUID that conforms to
// Ethereum bytes32.
func NewBytes32ID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}

// HexToBytes converts the given array of strings and returns bytes.
func HexToBytes(strs ...string) ([]byte, error) {
	return hex.DecodeString(RemoveHexPrefix(HexConcat(strs...)))
}

// HexConcat concatenates a given array of strings to return a single
// string.
func HexConcat(strs ...string) string {
	hex := "0x"
	for _, str := range strs {
		hex = hex + RemoveHexPrefix(str)
	}
	return hex
}

// RemoveHexPrefix removes the prefix (0x) of a given hex string.
func RemoveHexPrefix(str string) string {
	if len(str) > 1 && strings.ToLower(str[0:2]) == "0x" {
		return str[2:]
	}
	return str
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

// HexToString decodes a hex encoded string.
func HexToString(hex string) (string, error) {
	b, err := HexToBytes(hex)
	return string(b), err
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
	Duration() time.Duration
}

// BackoffSleeper is a counter to assist with reattempts.
type BackoffSleeper struct {
	*backoff.Backoff
}

// NewBackoffSleeper returns a BackoffSleeper that is configured to
// sleep for 1 second minimum, and 10 seconds maximum.
func NewBackoffSleeper() BackoffSleeper {
	return BackoffSleeper{&backoff.Backoff{
		Min: 1 * time.Second,
		Max: 10 * time.Second,
	}}
}

// Sleep waits for the given duration before reattempting.
func (bs BackoffSleeper) Sleep() {
	time.Sleep(bs.Backoff.Duration())
}

// Duration returns the current duration value.
func (bs BackoffSleeper) Duration() time.Duration {
	return bs.ForAttempt(bs.Attempt())
}

// ConstantSleeper is to assist with reattempts with
// the same sleep duration.
type ConstantSleeper struct {
	Sleeper
	interval time.Duration
}

// NewConstantSleeper returns a ConstantSleeper that is configured to
// sleep for a constant duration based on the input.
func NewConstantSleeper(d time.Duration) ConstantSleeper {
	return ConstantSleeper{interval: d}
}

// Sleep waits for the given duration before reattempting.
func (cs ConstantSleeper) Sleep() {
	time.Sleep(cs.interval)
}

// Duration returns the duration value.
func (cs ConstantSleeper) Duration() time.Duration {
	return cs.interval
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

// EVMSignedHexNumber formats a number as a 32 byte hex string
// Twos compliment representation if a minus number
func EVMSignedHexNumber(val *big.Int) (string, error) {
	var sh string
	if val.Sign() == -1 {
		evmUint256Max, ok := (&big.Int{}).SetString(strings.Repeat("f", 64), 16)
		if !ok {
			return sh, fmt.Errorf("could not parse evmUint256 max")
		}
		sh = EVMHexNumber((&big.Int{}).Add(evmUint256Max, val.Add(val, big.NewInt(1))))
	} else {
		sh = EVMHexNumber(val)
	}
	return sh, nil
}

// EVMHexNumber formats a number as a 32 byte hex string.
func EVMHexNumber(val interface{}) string {
	return fmt.Sprintf("0x%064x", val)
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

// ParseUintHex parses an unsigned integer out of a hex string.
func ParseUintHex(hex string) (*big.Int, error) {
	amount, ok := new(big.Int).SetString(hex, 0)
	if !ok {
		return amount, fmt.Errorf("unable to decode hex to integer: %s", hex)
	}
	return amount, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
