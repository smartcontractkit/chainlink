package contracts

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/log"
)

type ConnectedContract interface {
	eth.ContractCodec
	Call(result interface{}, methodName string, args ...interface{}) error
	SubscribeToLogs(listener log.Listener) (connected bool, _ UnsubscribeFunc)
}

type connectedContract struct {
	eth.ContractCodec
	address        common.Address
	ethClient      eth.Client
	logBroadcaster log.Broadcaster
}

type UnsubscribeFunc func()

func NewConnectedContract(
	codec eth.ContractCodec,
	address common.Address,
	ethClient eth.Client,
	logBroadcaster log.Broadcaster,
) ConnectedContract {
	return &connectedContract{codec, address, ethClient, logBroadcaster}
}

func (contract *connectedContract) Call(result interface{}, methodName string, args ...interface{}) (err error) {
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
	err = contract.ABI().UnpackIntoInterface(result, methodName, rawResult)
	return errors.Wrap(err, "unable to unpack values")
}

func (contract *connectedContract) SubscribeToLogs(listener log.Listener) (connected bool, _ UnsubscribeFunc) {
	connected = contract.logBroadcaster.Register(contract.address, listener)
	unsub := func() { contract.logBroadcaster.Unregister(contract.address, listener) }
	return connected, unsub
}
