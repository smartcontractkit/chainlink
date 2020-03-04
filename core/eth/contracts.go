package eth

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"

	"chainlink/core/logger"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type Contract interface {
	ABI() abi.ABI
	GetMethodID(method string) ([]byte, error)
	EncodeMessageCall(method string, args ...interface{}) ([]byte, error)
	UnpackLog(out interface{}, event string, log Log) error
}

// Contract holds the solidity contract's parsed ABI
type contract struct {
	abi abi.ABI
}

func getContract(name string, box packr.Box) (Contract, error) {
	jsonFile, err := box.Find(name + ".json")
	if err != nil {
		return nil, errors.Wrap(err, "unable to read contract JSON")
	}

	abiBytes := gjson.GetBytes(jsonFile, "compilerOutput.abi")
	abiParsed, err := abi.JSON(strings.NewReader(abiBytes.Raw))
	if err != nil {
		return nil, err
	}

	return &contract{abiParsed}, nil
}

// GetContract loads the contract JSON file from ../../evm-contracts/abi/v0.4
// and parses the ABI JSON contents into an abi.ABI object
//
// NB: These contracts can be built by running
//    yarn setup:contracts
// in the base project directory.
func GetContract(name string) (Contract, error) {
	box := packr.NewBox("../../evm-contracts/abi/v0.4")
	return getContract(name, box)
}

// GetV6Contract loads the contract JSON file from ../../evm-contracts/abi/v0.6
// and parses the ABI JSON contents into an abi.ABI object
//
// NB: These contracts can be built by running
//    yarn setup:contracts
// in the base project directory.
func GetV6Contract(name string) (Contract, error) {
	box := packr.NewBox("../../evm-contracts/abi/v0.6")
	return getContract(name, box)
}

func (contract *contract) ABI() abi.ABI {
	return contract.abi
}

// EncodeMessageCall encodes method name and arguments into a byte array
// to conform with the contract's ABI
func (contract *contract) EncodeMessageCall(method string, args ...interface{}) ([]byte, error) {
	return contract.abi.Pack(method, args...)
}

// GetMethodID returns the first 4 bytes of the keccak256 hash of the method
// signature. The passed method is simply the method name, not the parameters,
// as defined by go-ethereum ABI Methods
//
// e.g.
// There are two functions have same name:
// * foo(int,int)
// * foo(uint,uint)
// The method name of the first one will be resolved as foo while the second one
// will be resolved as foo0.
func (contract *contract) GetMethodID(method string) ([]byte, error) {
	mabi, found := contract.abi.Methods[method]
	if !found {
		return []byte{}, errors.New("unable to find contract method " + method)
	}
	return mabi.ID(), nil
}

// MustGetV6ContractEventID finds the event for the given contract by searching
// embedded contract assets from evm/, or panics if not found.
func MustGetV6ContractEventID(name, eventName string) common.Hash {
	contract, err := GetV6Contract(name)
	if err != nil {
		logger.Panic(fmt.Errorf("unable to find contract %s", name))
	}

	event, found := contract.ABI().Events[eventName]
	if !found {
		logger.Panic(fmt.Errorf("unable to find event %s for contract %s", eventName, name))
	}
	return event.ID()
}

func (contract *contract) UnpackLog(out interface{}, event string, log Log) error {
	if len(log.Data) > 0 {
		if err := contract.abi.Unpack(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range contract.abi.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return parseTopics(out, indexed, log.Topics[1:])
}

var (
	reflectHash    = reflect.TypeOf(common.Hash{})
	reflectAddress = reflect.TypeOf(common.Address{})
	reflectBigInt  = reflect.TypeOf(new(big.Int))
)

func parseTopics(out interface{}, fields abi.Arguments, topics []common.Hash) error {
	// Sanity check that the fields and topics match up
	if len(fields) != len(topics) {
		return errors.New("topic/field count mismatch")
	}
	// Iterate over all the fields and reconstruct them from topics
	for _, arg := range fields {
		if !arg.Indexed {
			return errors.New("non-indexed field in topic reconstruction")
		}
		field := reflect.ValueOf(out).Elem().FieldByName(capitalise(arg.Name))

		// Try to parse the topic back into the fields based on primitive types
		switch field.Kind() {
		case reflect.Bool:
			if topics[0][common.HashLength-1] == 1 {
				field.Set(reflect.ValueOf(true))
			}
		case reflect.Int8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int8(num.Int64())))

		case reflect.Int16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int16(num.Int64())))

		case reflect.Int32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int32(num.Int64())))

		case reflect.Int64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Int64()))

		case reflect.Uint8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint8(num.Uint64())))

		case reflect.Uint16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint16(num.Uint64())))

		case reflect.Uint32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint32(num.Uint64())))

		case reflect.Uint64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Uint64()))

		default:
			// Ran out of plain primitive types, try custom types

			switch field.Type() {
			case reflectHash: // Also covers all dynamic types
				field.Set(reflect.ValueOf(topics[0]))

			case reflectAddress:
				var addr common.Address
				copy(addr[:], topics[0][common.HashLength-common.AddressLength:])
				field.Set(reflect.ValueOf(addr))

			case reflectBigInt:
				num := new(big.Int).SetBytes(topics[0][:])
				if arg.Type.T == abi.IntTy {
					if num.Cmp(abi.MaxInt256) > 0 {
						num.Add(abi.MaxUint256, big.NewInt(0).Neg(num))
						num.Add(num, big.NewInt(1))
						num.Neg(num)
					}
				}
				field.Set(reflect.ValueOf(num))

			default:
				// Ran out of custom types, try the crazies
				switch {
				// static byte array
				case arg.Type.T == abi.FixedBytesTy:
					reflect.Copy(field, reflect.ValueOf(topics[0][:arg.Type.Size]))
				default:
					return fmt.Errorf("unsupported indexed type: %v", arg.Type)
				}
			}
		}
		topics = topics[1:]
	}
	return nil
}

func capitalise(input string) string {
	return abi.ToCamelCase(input)
}

type ConnectedContract interface {
	Contract
	Call(result interface{}, methodName string, args ...interface{}) error
}

type connectedContract struct {
	Contract
	ethClient Client
	address   common.Address
}

func NewConnectedContract(contract Contract, ethClient Client, address common.Address) ConnectedContract {
	return &connectedContract{contract, ethClient, address}
}

func (contract *connectedContract) Call(result interface{}, methodName string, args ...interface{}) error {
	data, err := contract.EncodeMessageCall(methodName, args...)
	if err != nil {
		return errors.Wrap(err, "unable to encode message call")
	}

	var rawResult hexutil.Bytes
	callArgs := CallArgs{To: contract.address, Data: data}
	err = contract.ethClient.Call(&rawResult, "eth_call", callArgs, "latest")
	if err != nil {
		return errors.Wrap(err, "unable to call client")
	}

	err = contract.ABI().Unpack(result, methodName, rawResult)
	if err != nil {
		return errors.Wrap(err, "unable to unpack values")
	}
	return nil
}
