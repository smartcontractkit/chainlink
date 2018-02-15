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
	"net/http"
	"strconv"
	"strings"
	"time"

	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	uuid "github.com/satori/go.uuid"
)

const (
	HUMAN_TIME_FORMAT = "2006-01-02 15:04:05 MST"
	weiPerEth         = 1e18
)

// ZeroAddress is an empty address, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

func HexToUint64(hex string) (uint64, error) {
	if strings.ToLower(hex[0:2]) == "0x" {
		hex = hex[2:]
	}
	return strconv.ParseUint(hex, 16, 64)
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

func ISO8601UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// BasicAuthPost posts to the HTTP client with the given username and password
// to authenticate at the url with contentType and returns a response.
func BasicAuthPost(username, password, url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	return resp, err
}

// BasicAuthGet uses the given username and password to send a GET request
// at the given URL and returns a response.
func BasicAuthGet(username, password, url string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
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
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

// HexToBytes converts the given array of strings and returns bytes.
func HexToBytes(strs ...string) ([]byte, error) {
	return hex.DecodeString(removeHexPrefix(HexConcat(strs...)))
}

// HexConcat concatenates a given array of strings to return a single
// string.
func HexConcat(strs ...string) string {
	hex := "0x"
	for _, str := range strs {
		hex = hex + removeHexPrefix(str)
	}
	return hex
}

func removeHexPrefix(str string) string {
	if len(str) > 1 && str[0:2] == "0x" {
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

// WeiToEth converts wei amounts to ether.
func WeiToEth(numWei *big.Int) float64 {
	numWeiBigFloat := new(big.Float).SetInt(numWei)
	weiPerEthBigFloat := new(big.Float).SetFloat64(weiPerEth)
	numEthBigFloat := new(big.Float)
	numEthBigFloat.Quo(numWeiBigFloat, weiPerEthBigFloat)
	numEthFloat64, _ := numEthBigFloat.Float64()
	return numEthFloat64
}

// EthToWei converts ether amounts to wei.
func EthToWei(numEth float64) *big.Int {
	numEthBigFloat := new(big.Float).SetFloat64(numEth)
	weiPerEthBigFloat := new(big.Float).SetFloat64(weiPerEth)
	numWeiBigFloat := new(big.Float)
	numWeiBigFloat.Mul(weiPerEthBigFloat, numEthBigFloat)
	numWeiBigInt, _ := numWeiBigFloat.Int(nil)
	return numWeiBigInt
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
	return prependHexPrefix(hex.EncodeToString([]byte(in)))
}

func prependHexPrefix(str string) string {
	if len(str) < 2 || len(str) > 1 && str[0:2] != "0x" {
		str = "0x" + str
	}
	return str
}

// HexToString decodes a hex encoded string.
func HexToString(hex string) (string, error) {
	b, err := HexToBytes(hex)
	return string(b), err
}
