package fluxmonitor

import (
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/flux_aggregator_wrapper"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/log"
)

type fluxAggregatorDecodingLogListener struct {
	wrapper flux_aggregator_wrapper.FluxAggregator
	log.Listener
}

var _ log.Listener = (*fluxAggregatorDecodingLogListener)(nil)

func newFluxAggregatorDecodingLogListener(
	address common.Address,
	backend bind.ContractBackend,
	innerListener log.Listener,
) (log.Listener, error) {
	wrapper, err := flux_aggregator_wrapper.NewFluxAggregator(address, backend)
	if err != nil {
		return nil, err
	}
	return fluxAggregatorDecodingLogListener{
		wrapper:  *wrapper,
		Listener: innerListener,
	}, nil
}

func (ll fluxAggregatorDecodingLogListener) HandleLog(lb log.Broadcast, err error) {
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
	case fluxAggregatorABI.Events["NewRound"].ID:
		decodedLog, err = ll.wrapper.ParseNewRound(rawLog)
	case fluxAggregatorABI.Events["AnswerUpdated"].ID:
		decodedLog, err = ll.wrapper.ParseAnswerUpdated(rawLog)
	default:
		logger.Warnf("Unknown topic for FluxAggregator contract: %s", eventID.Hex())
		return // don't pass on unknown/unexpected events
	}

	lb.SetDecodedLog(decodedLog)
	ll.Listener.HandleLog(lb, err)
}
