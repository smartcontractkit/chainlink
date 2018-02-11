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
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	uuid "github.com/satori/go.uuid"
	null "gopkg.in/guregu/null.v3"
)

const HUMAN_TIME_FORMAT = "2006-01-02 15:04:05 MST"

// ZeroAddress is an empty address, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

// SenderFromTxHex returns the sender's address from a given transaction.
func SenderFromTxHex(value string, chainID uint64) (common.Address, error) {
	tx, err := DecodeTxFromHex(value, chainID)
	if err != nil {
		return common.Address{}, err
	}
	signer := types.NewEIP155Signer(big.NewInt(int64(chainID)))
	return types.Sender(signer, &tx)
}

// DecodeTxFromHex returns an Ethereum transaction type from the given
// transaction Transaction ID.
func DecodeTxFromHex(value string, chainID uint64) (types.Transaction, error) {
	buffer := bytes.NewBuffer(common.FromHex(value))
	rlpStream := rlp.NewStream(buffer, 0)
	tx := types.Transaction{}
	err := tx.DecodeRLP(rlpStream)
	return tx, err
}

// HexToUint64 converts the given hex string to uint64.
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

// TimeParse returns the given string as a Time type.
func TimeParse(s string) time.Time {
	t, err := dateparse.ParseAny(s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

// ISO8601UTC returns time formatted as "2018-02-11T14:16:47Z"
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

// ParseISO8601 parses the given string as RFC3339Nanoand returns an
// instance of Time.
func ParseISO8601(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return t
}

// NullableTime allows for the given time to be null. Marshals
// to null for JSON serialization if null.
func NullableTime(t time.Time) null.Time {
	return null.Time{Time: t, Valid: true}
}

// ParseNullableTime parses the given string and will allow
// for time to be null.
func ParseNullableTime(s string) null.Time {
	return NullableTime(ParseISO8601(s))
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

func addHexPrefix(str string) string {
	if len(str) > 1 && str[0:2] != "0x" {
		return "0x" + str
	}
	return str
}

// BytesToHex converts and returns the given bytes to its hex
// value as a string.
func BytesToHex(bytes ...[]byte) string {
	str := "0x"
	for _, b := range bytes {
		str = str + hex.EncodeToString(b)
	}
	return str
}

// DecodeEthereumTx parses RLP-encoded data into an Ethereum transaction.
func DecodeEthereumTx(hex string) (types.Transaction, error) {
	var tx types.Transaction
	b, err := hexutil.Decode(hex)
	if err != nil {
		return tx, err
	}
	return tx, rlp.DecodeBytes(b, &tx)
}
