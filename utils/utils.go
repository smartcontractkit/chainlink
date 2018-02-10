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

var ZeroAddress = common.Address{}

func SenderFromTxHex(value string, chainID uint64) (common.Address, error) {
	tx, err := DecodeTxFromHex(value, chainID)
	if err != nil {
		return common.Address{}, err
	}
	signer := types.NewEIP155Signer(big.NewInt(int64(chainID)))
	return types.Sender(signer, &tx)
}

func DecodeTxFromHex(value string, chainID uint64) (types.Transaction, error) {
	buffer := bytes.NewBuffer(common.FromHex(value))
	rlpStream := rlp.NewStream(buffer, 0)
	tx := types.Transaction{}
	err := tx.DecodeRLP(rlpStream)
	return tx, err
}

func HexToUint64(hex string) (uint64, error) {
	if strings.ToLower(hex[0:2]) == "0x" {
		hex = hex[2:]
	}
	return strconv.ParseUint(hex, 16, 64)
}

func Uint64ToHex(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

func EncodeTxToHex(tx *types.Transaction) (string, error) {
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return "", err
	}
	return common.ToHex(rlp.Bytes()), nil
}

func StringToHash(str string) (common.Hash, error) {
	b, err := hexutil.Decode(str)
	if err != nil {
		return common.Hash{}, err
	}
	return common.BytesToHash(b), nil
}

func StringToAddress(str string) (common.Address, error) {
	b, err := hexutil.Decode(str)
	if err != nil {
		return common.Address{}, err
	}
	return common.BytesToAddress(b), nil
}

func StringToBytes(str string) (hexutil.Bytes, error) {
	var b hexutil.Bytes
	err := b.UnmarshalText([]byte(addHexPrefix(str)))
	return b, err
}

func TimeParse(s string) time.Time {
	t, err := dateparse.ParseAny(s)
	if err != nil {
		log.Fatal(err)
	}
	return t
}

func ISO8601UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

func BasicAuthPost(username, password, url string, contentType string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", url, body)
	request.Header.Set("Content-Type", contentType)
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	return resp, err
}

func BasicAuthGet(username, password, url string) (*http.Response, error) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", url, nil)
	request.SetBasicAuth(username, password)
	resp, err := client.Do(request)
	return resp, err
}

func FormatJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

func ParseISO8601(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		panic(err)
	}
	return t
}

func NullableTime(t time.Time) null.Time {
	return null.Time{Time: t, Valid: true}
}

func ParseNullableTime(s string) null.Time {
	return NullableTime(ParseISO8601(s))
}

func GetStringKeys(v map[string]interface{}) []string {
	keys := make([]string, len(v))

	i := 0
	for k := range v {
		keys[i] = k
		i++
	}

	return keys
}

func NewBytes32ID() string {
	return strings.Replace(uuid.NewV4().String(), "-", "", -1)
}

func HexToBytes(strs ...string) ([]byte, error) {
	return hex.DecodeString(removeHexPrefix(HexConcat(strs...)))
}

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

func BytesToHex(bytes ...[]byte) string {
	str := "0x"
	for _, b := range bytes {
		str = str + hex.EncodeToString(b)
	}
	return str
}

func DecodeEthereumTx(hex string) (types.Transaction, error) {
	var tx types.Transaction
	b, err := StringToBytes(hex)
	if err != nil {
		return tx, err
	}
	return tx, rlp.DecodeBytes(b, &tx)
}
