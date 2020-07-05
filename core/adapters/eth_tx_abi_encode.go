package adapters

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi"
	gethCommon "github.com/ethereum/go-ethereum/common"
	strpkg "github.com/smartcontractkit/chainlink/core/store"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

const evmWordSize = 32

// EthTxABIEncode holds the Address to send the result to and the FunctionABI
// to use for encoding arguments.
type EthTxABIEncode struct {
	ToAddress                        gethCommon.Address
	FromAddress                      gethCommon.Address
	GasPrice                         *utils.Big
	GasLimit                         uint64
	MinRequiredOutgoingConfirmations uint64

	// ABI of contract function this task calls
	FunctionABI abi.Method
}

// GetToAddress implements the EthTxCommon interface
func (e *EthTxABIEncode) GetToAddress() gethCommon.Address {
	return e.ToAddress
}

// GetFromAddress implements the EthTxCommon interface
func (e *EthTxABIEncode) GetFromAddress() gethCommon.Address {
	return e.FromAddress
}

// GetGasLimit implements the EthTxCommon interface
func (e *EthTxABIEncode) GetGasLimit() uint64 {
	return e.GasLimit
}

// GetGasPrice implements the EthTxCommon interface
func (e *EthTxABIEncode) GetGasPrice() *utils.Big {
	return e.GasPrice
}

// GetMinRequiredOutgoingConfirmations implements the EthTxCommon interface
func (e *EthTxABIEncode) GetMinRequiredOutgoingConfirmations() uint64 {
	return e.MinRequiredOutgoingConfirmations
}

// GetEncodedPayload implements the EthTxCommon interface
func (e *EthTxABIEncode) GetEncodedPayload(input models.RunInput) ([]byte, error) {
	return e.abiEncode(input)
}

// TaskType returns the type of Adapter.
func (e *EthTxABIEncode) TaskType() models.TaskType {
	return TaskTypeEthTxABIEncode
}

// UnmarshalJSON for custom JSON unmarshal that is strict, i.e. doesn't
// accept spurious fields. (In particular, we wan't to ensure that we don't
// get spurious fields in the FunctionABI, so that users don't get any wrong
// ideas about what parts of the ABI we use for encoding data.)
func (e *EthTxABIEncode) UnmarshalJSON(data []byte) error {
	var fields struct {
		Address     gethCommon.Address
		FromAddress gethCommon.Address
		FunctionABI struct {
			Name   string
			Inputs abi.Arguments
		}
		GasPrice *utils.Big
		GasLimit uint64
	}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&fields); err != nil {
		return err
	}

	e.ToAddress = fields.Address
	e.FromAddress = fields.FromAddress
	e.FunctionABI.Name = fields.FunctionABI.Name
	e.FunctionABI.Inputs = fields.FunctionABI.Inputs
	e.GasPrice = fields.GasPrice
	e.GasLimit = fields.GasLimit
	return nil
}

// Perform creates the run result for the transaction if the existing run result
// is not currently pending. Then it confirms the transaction was confirmed on
// the blockchain.
func (e *EthTxABIEncode) Perform(input models.RunInput, store *strpkg.Store) models.RunOutput {
	return findOrInsertEthTx(e, input, store)
}

// abiEncode ABI-encodes the arguments passed in a RunResult's result field
// according to etx.FunctionABI
func (e *EthTxABIEncode) abiEncode(input models.RunInput) ([]byte, error) {
	args, ok := input.Data().Get("result").Value().(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("json result is not an object")
	}
	return abiEncode(&e.FunctionABI, args)
}

// abiEncode ABI-encodes the arguments in args according to fnABI.
func abiEncode(fnABI *abi.Method, args map[string]interface{}) ([]byte, error) {
	if len(fnABI.Inputs) != len(args) {
		return nil, errors.Errorf(
			"json result has wrong length. should have %v entries, one for each argument",
			len(fnABI.Inputs))
	}

	encodedStaticPartSize := 0
	for _, input := range fnABI.Inputs {
		encodedStaticPartSize += staticSize(&input.Type)
	}

	encodedStaticPart := make([]byte, 0, encodedStaticPartSize)
	encodedDynamicPart := make([]byte, 0)
	dynamicOffset := encodedStaticPartSize
	for _, input := range fnABI.Inputs {
		name := input.Name
		jval, ok := args[name]
		if !ok {
			return nil, errors.Errorf("entry for argument %s is missing", name)
		}

		if !isSupportedABIType(&input.Type) {
			return nil, errors.Errorf(
				"argument %s has unsupported ABI type %s",
				name, input.Type)
		}

		static, dynamic, err := enc(&input.Type, jval, name)
		switch {
		case err != nil:
			return nil, err
		case static != nil && dynamic != nil:
			panic("static AND dynamic returned")
		case static != nil:
			assertPadded(static)
			encodedStaticPart = append(encodedStaticPart, static...)
		case dynamic != nil:
			assertPadded(dynamic)
			encodedStaticPart = append(encodedStaticPart, encPositiveInt(dynamicOffset)...)
			dynamicOffset += len(dynamic)
			encodedDynamicPart = append(encodedDynamicPart, dynamic...)
		}
	}

	if len(encodedStaticPart) != encodedStaticPartSize {
		panic("unexpected size of static part")
	}

	result := functionSelector(fnABI)
	result = append(result, encodedStaticPart...)
	result = append(result, encodedDynamicPart...)
	return result, nil
}

// We support every type that solidity contracts as of solc v0.5.11 can decode:
// address, bool, bytes, bytes1, ..., bytes32, int8, ..., int256, string, uint8,
// ..., uint256, as well as fixed size arrays (e.g. int128[6] and bool[3][3])
// and slices (e.g. address[] and bool[3][3][])
//
// We (like solidity) don't support nested dynamic types like address[][] or
// string[].
func isSupportedABIType(typ *abi.Type) bool {
	switch typ.T {
	case abi.StringTy, abi.BytesTy:
		return true
	case abi.SliceTy:
		return isSupportedStaticABIType(typ.Elem)
	default:
		return isSupportedStaticABIType(typ)
	}
}

func isSupportedStaticABIType(typ *abi.Type) bool {
	switch typ.T {
	case abi.AddressTy, abi.BoolTy:
		return true
	case abi.ArrayTy:
		return isSupportedStaticABIType(typ.Elem)
	case abi.IntTy, abi.UintTy:
		return typ.Size%8 == 0 && 0 < typ.Size && typ.Size <= 8*evmWordSize
	case abi.FixedBytesTy:
		return 0 < typ.Size && typ.Size <= evmWordSize
	default:
		return false
	}
}

// enc encodes a JSON value jval of ABI type typ. name is passed for better
// error reporting.
func enc(typ *abi.Type, jval interface{}, name string) (static []byte, dynamic []byte, err error) {
	switch typ.T {
	case abi.BytesTy:
		bytes, err := bytesFromJSON(jval, name)
		if err != nil {
			return nil, nil, err
		}
		return nil, padAndPrefixDynamic(bytes), nil
	case abi.SliceTy:
		s, ok := jval.([]interface{})
		if !ok {
			return nil, nil, errors.Errorf("argument %s is not an array", name)
		}

		result := encPositiveInt(len(s))
		for i, elem := range s {
			encoded, err := encStatic(typ.Elem, elem, fmt.Sprintf("%s[%v]", name, i))
			if err != nil {
				return nil, nil, err
			}
			result = append(result, encoded...)
		}
		return nil, result, nil
	case abi.StringTy:
		s, ok := jval.(string)
		if !ok {
			return nil, nil, errors.Errorf("argument %s is not a string", name)
		}
		bytes := []byte(s)
		return nil, padAndPrefixDynamic(bytes), nil
	default:
		static, err := encStatic(typ, jval, name)
		return static, nil, err
	}
}

// Dynamic types like bytes and string are length-prefixed and padded to a
// multiple of evmWordSize
func padAndPrefixDynamic(bytes []byte) []byte {
	result := encPositiveInt(len(bytes))
	result = append(result, bytes...)
	result = padRight(result, (len(result)+evmWordSize-1)/evmWordSize*evmWordSize)
	return result
}

func staticSize(typ *abi.Type) int {
	switch typ.T {
	case abi.AddressTy, abi.BoolTy, abi.BytesTy, abi.FixedBytesTy,
		abi.IntTy, abi.SliceTy, abi.StringTy, abi.UintTy:
		return evmWordSize
	case abi.ArrayTy:
		return typ.Size * staticSize(typ.Elem)
	default:
		panic("Unsupported type")
	}
}

// Encodes JSON value jval according to static ABI type (e.g. int*, uint*, ...)
// typ. name is used for better error messages.
func encStatic(typ *abi.Type, jval interface{}, name string) ([]byte, error) {
	switch typ.T {
	case abi.AddressTy:
		s, ok := jval.(string)
		if !ok {
			return nil, errors.Errorf("argument %s is not an address string", name)
		}
		if !utils.HasHexPrefix(s) {
			return nil, errors.Errorf("argument %s is not a hexstring", name)
		}
		addressBytes, err := hex.DecodeString(s[2:])
		if err != nil {
			return nil, errors.Wrapf(err, "argument %s is not a hexstring", name)
		}
		if len(addressBytes) > 20 {
			return nil, errors.Errorf("argument %s is too long for an address (20 bytes)", name)
		}
		return padLeft(addressBytes, evmWordSize), nil
	case abi.ArrayTy:
		a, ok := jval.([]interface{})
		if !ok {
			return nil, errors.Errorf("argument %s is not an array", name)
		}
		if len(a) != typ.Size {
			return nil, errors.Errorf("argument %s is an array with %v items, but we need %v", name, len(a), typ.Size)
		}
		result := make([]byte, 0, staticSize(typ))
		for i, elem := range a {
			encoded, err := encStatic(typ.Elem, elem, fmt.Sprintf("%s[%v]", name, i))
			if err != nil {
				return nil, err
			}
			result = append(result, encoded...)
		}
		return result, nil
	case abi.BoolTy:
		b, ok := jval.(bool)
		if !ok {
			return nil, errors.Errorf("argument %s is not a boolean", name)
		}
		if b {
			return padLeft([]byte{1}, evmWordSize), nil
		}
		return padLeft([]byte{0}, evmWordSize), nil
	case abi.FixedBytesTy:
		bytes, err := bytesFromJSON(jval, name)
		if err != nil {
			return nil, err
		}
		if len(bytes) != typ.Size {
			return nil, errors.Errorf("argument %s doesn't have exactly %v bytes", name, typ.Size)
		}
		return padRight(bytes, evmWordSize), nil
	case abi.IntTy:
		if _, ok := jval.(float64); ok && typ.Size > 48 {
			return nil, errors.Errorf("argument %s is a json number which isn't suitable for storing integers greater than 2**53", name)
		}
		n, err := bigIntFromJSON(jval, name)
		if err != nil {
			return nil, err
		}
		encoded, err := encSigned(uint(typ.Size/8), n, name)
		if err != nil {
			return nil, err
		}
		return encoded, nil
	case abi.UintTy:
		if _, ok := jval.(float64); ok && typ.Size > 48 {
			return nil, errors.Errorf("argument %s is a json number which isn't suitable for storing integers greater than 2**53", name)
		}
		n, err := bigIntFromJSON(jval, name)
		if err != nil {
			return nil, err
		}
		encoded, err := encUnsigned(uint(typ.Size/8), n, name)
		if err != nil {
			return nil, err
		}
		return encoded, nil
	default:
		panic(fmt.Sprintf("Unsupported type: %s", typ.String()))
	}
}

func encSigned(sizeInBytes uint, n *big.Int, name string) ([]byte, error) {
	min, max := big.NewInt(0), big.NewInt(0)
	min.Lsh(big.NewInt(1), sizeInBytes*8-1)
	max.Sub(min, big.NewInt(1))
	min.Neg(min)
	if n.Cmp(min) < 0 || max.Cmp(n) < 0 {
		return nil, errors.Errorf(
			"argument %s out of valid range for signed integer with %v bytes", name, sizeInBytes)
	}
	if n.Sign() < 0 {
		// Convert to two's complement: make n positive, then subtract from
		// 2**256
		n.Neg(n)
		n.Sub(min.Lsh(big.NewInt(1), evmWordSize*8), n)
	}
	return padLeft(n.Bytes(), evmWordSize), nil
}

func encUnsigned(sizeInBytes uint, n *big.Int, name string) ([]byte, error) {
	min, max := big.NewInt(0), big.NewInt(0)
	max.Lsh(big.NewInt(1), sizeInBytes*8)
	max.Sub(max, big.NewInt(1))
	if n.Cmp(min) < 0 || max.Cmp(n) < 0 {
		return nil, errors.Errorf(
			"argument %s out of valid range for unsigned integer with %v bytes", name, sizeInBytes)
	}

	return padLeft(n.Bytes(), evmWordSize), nil
}

func encPositiveInt(i int) []byte {
	if i < 0 {
		panic("Negative int not allowed")
	}
	encoded, err := encUnsigned(evmWordSize, big.NewInt(int64(i)), "")
	if err != nil {
		panic("Unexpected error")
	}
	return encoded
}

var hexDigitsRegexp = regexp.MustCompile("^[0-9a-fA-F]*$")

func bigIntFromJSON(jval interface{}, name string) (*big.Int, error) {
	switch val := jval.(type) {
	case string:
		var n = big.NewInt(0)
		valid := false
		if utils.HasHexPrefix(val) {
			if !hexDigitsRegexp.MatchString(val[2:]) {
				return nil, errors.Errorf("argument %s starts with '0x', but doesn't encode valid hex number", name)
			}
			n, valid = big.NewInt(0).SetString(val[2:], 16)
			if !valid {
				return nil, errors.Errorf("argument %s starts with '0x', but doesn't encode valid hex number", name)
			}
		} else {
			n, valid = big.NewInt(0).SetString(val, 10)
			if !valid {
				return nil, errors.Errorf("argument %s doesn't encode valid decimal number", name)
			}
		}
		return n, nil
	case float64:
		n, accuracy := big.NewFloat(val).Int(big.NewInt(0))
		if accuracy != big.Exact {
			return nil, errors.Errorf("argument %s isn't a whole number", name)
		}
		return n, nil
	default:
		return nil, errors.Errorf("argument %s isn't a JSON string or number", name)
	}
}

func bytesFromJSON(jval interface{}, name string) ([]byte, error) {
	switch val := jval.(type) {
	case string:
		if !utils.HasHexPrefix(val) {
			return nil, errors.Errorf("argument %s is not a hexstring", name)
		}
		bytes, err := hex.DecodeString(val[2:])
		if err != nil {
			return nil, errors.Wrapf(err, "argument %s is not a hexstring", name)
		}
		return bytes, nil
	case []interface{}:
		result := make([]byte, 0)
		for i, elem := range val {
			n, ok := elem.(float64)
			if !ok {
				return nil, errors.Errorf("argument %s[%v] is not a number", name, i)
			}
			if float64(byte(n)) != n {
				return nil, errors.Errorf("argument %s[%v] is not a byte", name, i)
			}
			result = append(result, byte(n))
		}
		return result, nil
	default:
		return nil, errors.Errorf("argument %s is not a hexstring or array", name)
	}
}

func padRight(b []byte, n int) []byte {
	if n < len(b) {
		panic("input to pad longer than desired output length")
	}
	for len(b) < n {
		b = append(b, 0)
	}
	return b
}

func padLeft(b []byte, n int) []byte {
	if n < len(b) {
		panic("input to pad longer than desired output length")
	}
	result := make([]byte, 0, n)
	for len(result)+len(b) < n {
		result = append(result, 0)
	}
	result = append(result, b...)
	return result
}

func assertPadded(b []byte) {
	if len(b)%evmWordSize != 0 {
		panic("ABI encoded data isn't padded properly")
	}
}

func functionSelector(fnABI *abi.Method) []byte {
	types := []string{}
	for _, input := range fnABI.Inputs {
		types = append(types, input.Type.String())
	}
	signature := fmt.Sprintf("%v(%v)", fnABI.Name, strings.Join(types, ","))
	return utils.MustHash(signature).Bytes()[:4]
}
