package log

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/utils"
)

func (lb *broadcaster) ExportedAppendLogChannel(ch1, ch2 <-chan types.Log) chan types.Log {
	return lb.subscriber.appendLogChannel(ch1, ch2)
}

func ExportedNewSubscriber(orm ORM, ethClient eth.Client, config Config, relayer *relayer, dependentAwaiter utils.DependentAwaiter) *subscriber {
	return newSubscriber(orm, ethClient, config, relayer, dependentAwaiter)
}

func ExportedNewRelayer(orm ORM, config Config, dependentAwaiter utils.DependentAwaiter) *relayer {
	return newRelayer(orm, config, dependentAwaiter)
}

func (s *subscriber) ExportedContracts() map[common.Address]uint64 {
	return s.contracts
}
