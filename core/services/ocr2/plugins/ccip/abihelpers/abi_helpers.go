package abihelpers

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
)

func MustGetEventID(name string, abi2 abi.ABI) common.Hash {
	event, ok := abi2.Events[name]
	if !ok {
		panic(fmt.Sprintf("missing event %s", name))
	}
	return event.ID
}

func MustGetEventInputs(name string, abi2 abi.ABI) abi.Arguments {
	m, ok := abi2.Events[name]
	if !ok {
		panic(fmt.Sprintf("missing event %s", name))
	}
	return m.Inputs
}

func MustGetMethodInputs(name string, abi2 abi.ABI) abi.Arguments {
	m, ok := abi2.Methods[name]
	if !ok {
		panic(fmt.Sprintf("missing method %s", name))
	}
	return m.Inputs
}

func MustParseABI(abiStr string) abi.ABI {
	abiParsed, err := abi.JSON(strings.NewReader(abiStr))
	if err != nil {
		panic(err)
	}
	return abiParsed
}

// ProofFlagsToBits transforms a list of boolean proof flags to a *big.Int
// encoded number.
func ProofFlagsToBits(proofFlags []bool) *big.Int {
	encodedFlags := big.NewInt(0)
	for i := 0; i < len(proofFlags); i++ {
		if proofFlags[i] {
			encodedFlags.SetBit(encodedFlags, i, 1)
		}
	}
	return encodedFlags
}

type AbiDefined interface {
	AbiString() string
}

type AbiDefinedValid interface {
	AbiDefined
	Validate() error
}

func ABIEncode(abiStr string, values ...interface{}) ([]byte, error) {
	inAbi, err := getABI(abiStr, ENCODE)
	if err != nil {
		return nil, err
	}
	res, err := inAbi.Pack("method", values...)
	if err != nil {
		return nil, err
	}
	return res[4:], nil
}

func ABIDecode(abiStr string, data []byte) ([]interface{}, error) {
	inAbi, err := getABI(abiStr, DECODE)
	if err != nil {
		return nil, err
	}
	return inAbi.Unpack("method", data)
}

func EncodeAbiStruct[T AbiDefined](decoded T) ([]byte, error) {
	return ABIEncode(decoded.AbiString(), decoded)
}

func EncodeAddress(address common.Address) ([]byte, error) {
	return ABIEncode(`[{"type":"address"}]`, address)
}

func DecodeAbiStruct[T AbiDefinedValid](encoded []byte) (T, error) {
	var empty T

	decoded, err := ABIDecode(empty.AbiString(), encoded)
	if err != nil {
		return empty, err
	}

	converted := abi.ConvertType(decoded[0], &empty)
	if casted, ok := converted.(*T); ok {
		return *casted, (*casted).Validate()
	}
	return empty, fmt.Errorf("can't cast from %T to %T", converted, empty)
}

func EvmWord(i uint64) common.Hash {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return common.BigToHash(big.NewInt(0).SetBytes(b))
}

func DecodeOCR2Config(encoded []byte) (*ocr2aggregator.OCR2AggregatorConfigSet, error) {
	unpacked := new(ocr2aggregator.OCR2AggregatorConfigSet)
	abiPointer, err := ocr2aggregator.OCR2AggregatorMetaData.GetAbi()
	if err != nil {
		return unpacked, err
	}
	defaultABI := *abiPointer
	err = defaultABI.UnpackIntoInterface(unpacked, "ConfigSet", encoded)
	if err != nil {
		return unpacked, errors.Wrap(err, "failed to unpack log data")
	}
	return unpacked, nil
}

// create const encode and decode
const (
	ENCODE = iota
	DECODE
)

type abiCache struct {
	cache map[string]*abi.ABI
	mu    *sync.RWMutex
}

func newAbiCache() *abiCache {
	return &abiCache{
		cache: make(map[string]*abi.ABI),
		mu:    &sync.RWMutex{},
	}
}

// Global cache for ABIs to avoid parsing the same ABI multiple times
// As the module is already a helper module and not a service, we can keep the cache global
// It's private to the package and can't be accessed from outside
var myAbiCache = newAbiCache()

// This Function is used to get the ABI from the cache or create a new one and cache it for later use
// operationType is used to differentiate between encoding and decoding
// encoding uses a definition with `inputs` and decoding uses a definition with `outputs` (check inDef)
func getABI(abiStr string, operationType uint8) (*abi.ABI, error) {
	var operationStr string
	switch operationType {
	case ENCODE:
		operationStr = "inputs"
	case DECODE:
		operationStr = "outputs"
	default:
		return nil, fmt.Errorf("invalid operation type")
	}

	inDef := fmt.Sprintf(`[{ "name" : "method", "type": "function", "%s": %s}]`, operationStr, abiStr)

	myAbiCache.mu.RLock()
	if cachedAbi, found := myAbiCache.cache[inDef]; found {
		myAbiCache.mu.RUnlock() // unlocking before returning
		return cachedAbi, nil
	}
	myAbiCache.mu.RUnlock()

	res, err := abi.JSON(strings.NewReader(inDef))
	if err != nil {
		return nil, err
	}

	myAbiCache.mu.Lock()
	defer myAbiCache.mu.Unlock()
	myAbiCache.cache[inDef] = &res
	return &res, nil
}
