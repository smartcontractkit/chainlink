package eth

import (
	"github.com/smartcontractkit/chainlink/core/eth"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type ConnectedContract interface {
	eth.ContractCodec
	Call(result interface{}, methodName string, args ...interface{}) error
	SubscribeToLogs(listener LogListener) (connected bool, _ UnsubscribeFunc)
}

type connectedContract struct {
	eth.ContractCodec
	address        common.Address
	ethClient      eth.Client
	logBroadcaster LogBroadcaster
}

type UnsubscribeFunc func()

func NewConnectedContract(
	codec eth.ContractCodec,
	address common.Address,
	ethClient eth.Client,
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
	callArgs := eth.CallArgs{To: contract.address, Data: data}
	err = contract.ethClient.Call(&rawResult, "eth_call", callArgs, "latest")
	if err != nil {
		return errors.Wrap(err, "unable to call client")
	}

	err = contract.ABI().Unpack(result, methodName, rawResult)
	return errors.Wrap(err, "unable to unpack values")
}

func (contract *connectedContract) SubscribeToLogs(listener LogListener) (connected bool, _ UnsubscribeFunc) {
	connected = contract.logBroadcaster.Register(contract.address, listener)
	unsub := func() { contract.logBroadcaster.Unregister(contract.address, listener) }
	return connected, unsub
}
