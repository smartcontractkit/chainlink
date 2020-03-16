package eth

import (
	"fmt"
	"strings"

	"chainlink/core/logger"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/gobuffalo/packr"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

type ContractCodec interface {
	ABI() *abi.ABI
	GetMethodID(method string) ([]byte, error)
	EncodeMessageCall(method string, args ...interface{}) ([]byte, error)
	UnpackLog(out interface{}, event string, log Log) error
}

// Contract holds the solidity contract's parsed ABI
type contractCodec struct {
	abi abi.ABI
}

func getContractCodec(name string, box packr.Box) (ContractCodec, error) {
	jsonFile, err := box.Find(name + ".json")
	if err != nil {
		return nil, errors.Wrap(err, "unable to read contract JSON")
	}

	abiBytes := gjson.GetBytes(jsonFile, "compilerOutput.abi")
	abiParsed, err := abi.JSON(strings.NewReader(abiBytes.Raw))
	if err != nil {
		return nil, err
	}

	return &contractCodec{abiParsed}, nil
}

// GetContract loads the contract JSON file from ../../evm-contracts/abi/v0.4
// and parses the ABI JSON contents into an abi.ABI object
//
// NB: These contracts can be built by running
//    yarn setup:contracts
// in the base project directory.
func GetContractCodec(name string) (ContractCodec, error) {
	box := packr.NewBox("../../evm-contracts/abi/v0.4")
	return getContractCodec(name, box)
}

// GetV6Contract loads the contract JSON file from ../../evm-contracts/abi/v0.6
// and parses the ABI JSON contents into an abi.ABI object
//
// NB: These contracts can be built by running
//    yarn setup:contracts
// in the base project directory.
func GetV6ContractCodec(name string) (ContractCodec, error) {
	box := packr.NewBox("../../evm-contracts/abi/v0.6")
	return getContractCodec(name, box)
}

func (contract *contractCodec) ABI() *abi.ABI {
	return &contract.abi
}

// EncodeMessageCall encodes method name and arguments into a byte array
// to conform with the contract's ABI
func (cc *contractCodec) EncodeMessageCall(method string, args ...interface{}) ([]byte, error) {
	return cc.abi.Pack(method, args...)
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
func (cc *contractCodec) GetMethodID(method string) ([]byte, error) {
	mabi, found := cc.abi.Methods[method]
	if !found {
		return []byte{}, errors.New("unable to find contract method " + method)
	}
	return mabi.ID(), nil
}

// MustGetV6ContractEventID finds the event for the given contract by searching
// embedded contract assets from evm/, or panics if not found.
func MustGetV6ContractEventID(name, eventName string) common.Hash {
	cc, err := GetV6ContractCodec(name)
	if err != nil {
		logger.Panic(fmt.Errorf("unable to find contract %s", name))
	}

	event, found := cc.ABI().Events[eventName]
	if !found {
		logger.Panic(fmt.Errorf("unable to find event %s for contract %s", eventName, name))
	}
	return event.ID()
}

func (cc *contractCodec) UnpackLog(out interface{}, event string, log Log) error {
	return gethUnpackLog(cc, out, event, log)
}

type ConnectedContract interface {
	ContractCodec
	Call(result interface{}, methodName string, args ...interface{}) error
	SubscribeToLogs(listener LogListener) UnsubscribeFunc
}

type connectedContract struct {
	ContractCodec
	address        common.Address
	ethClient      Client
	logBroadcaster LogBroadcaster
}

type UnsubscribeFunc func()

func NewConnectedContract(
	codec ContractCodec,
	address common.Address,
	ethClient Client,
	logBroadcaster LogBroadcaster,
) ConnectedContract {
	return &connectedContract{codec, address, ethClient, logBroadcaster}
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
	return errors.Wrap(err, "unable to unpack values")
}

func (contract *connectedContract) SubscribeToLogs(listener LogListener) UnsubscribeFunc {
	contract.logBroadcaster.Register(contract.address, listener)
	return func() { contract.logBroadcaster.Unregister(contract.address, listener) }
}
