package contracts

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
)

var flagsABI = MustGetABI(flags_wrapper.FlagsABI)

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
	log.Listener
}

var _ log.Listener = (*flagsDecodingLogListener)(nil)

func NewFlagsDecodingLogListener(
	contract *Flags,
	innerListener log.Listener,
) log.Listener {
	return flagsDecodingLogListener{
		contract: contract,
		Listener: innerListener,
	}
}

func (ll flagsDecodingLogListener) HandleLog(lb log.Broadcast, err error) {
	if err != nil {
		ll.Listener.HandleLog(lb, err)
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
		return // don't pass on unknown/unexpected events
	}

	lb.SetDecodedLog(decodedLog)
	ll.Listener.HandleLog(lb, err)
}
