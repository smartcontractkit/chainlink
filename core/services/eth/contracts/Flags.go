package contracts

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flags_wrapper"
	"github.com/smartcontractkit/chainlink/core/services/eth"
)

var FlagsABI = getFlagsABI()

func getFlagsABI() abi.ABI {
	abi, err := abi.JSON(strings.NewReader(flags_wrapper.FlagsABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	return abi
}

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
	case FlagsABI.Events["FlagRaised"].ID:
		decodedLog, err = ll.contract.ParseFlagRaised(rawLog)
	case FlagsABI.Events["FlagLowered"].ID:
		decodedLog, err = ll.contract.ParseFlagLowered(rawLog)
	default:
		err = errors.Errorf("Unknown topic for Flags contract: %s", eventID.Hex())
	}

	lb.SetDecodedLog(decodedLog)
	ll.LogListener.HandleLog(lb, err)
}
