package evm

import (
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
)

var (
	// We need only ABI part of the contract that describes the event structure
	// https://github.com/circlefin/evm-cctp-contracts/blob/377c9bd813fb86a42d900ae4003599d82aef635a/src/MessageTransmitter.sol#L41
	MessageTransmitterABI = `[
	  {
		"anonymous": false,
		"inputs": [
		  {
			"indexed": false,
			"internalType": "bytes",
			"name": "",
			"type": "bytes"
		  }
		],
		"name": "MessageSent",
		"type": "event"
	  }
	]`
	_ = evmtypes.MustGetABI(MessageTransmitterABI)
)
