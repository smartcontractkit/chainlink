package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

var flagsABI = getABI(flags_wrapper.FlagsABI)

type Flags struct {
	Address common.Address
	flags_wrapper.Flags
}

func NewFlagsContract(address common.Address, backend bind.ContractBackend) (*Flags, error) {
	flags, err := flags_wrapper.NewFlags(address, backend)
	if err != nil {
		return nil, err
	}
	return &Flags{
		Address: address,
		Flags:   *flags,
	}, nil
}

type flagsDecodingLogListener struct {
	contract *Flags
	eth.LogListener
}

var _ eth.LogListener = (*flagsDecodingLogListener)(nil)

func NewFlagsDecodingLogListener(
	contract *Flags,
	innerListener eth.LogListener,
) eth.LogListener {
	return flagsDecodingLogListener{
		contract:    contract,
		LogListener: innerListener,
	}
}

//  TODO - RYAN - test this
func (ll flagsDecodingLogListener) HandleLog(lb eth.LogBroadcast, err error) {
	if err != nil {
		ll.LogListener.HandleLog(lb, err)
		return
	}

	rawLog := lb.RawLog()
	if len(rawLog.Topics) == 0 {
		return
	}
	eventID := rawLog.Topics[0]
	var decodedLog interface{}

	switch eventID {
	case flagsABI.Events["FlagRaised"].ID:
		decodedLog, err = ll.contract.ParseFlagRaised(rawLog)
	case flagsABI.Events["FlagLowered"].ID:
		decodedLog, err = ll.contract.ParseFlagLowered(rawLog)
	default:
		logger.Warnf("Unknown topic for Flags contract: %s", eventID.Hex())
		return // don't pass on unknown/unexpectred events
	}

	lb.SetDecodedLog(decodedLog)
	ll.LogListener.HandleLog(lb, err)
}
