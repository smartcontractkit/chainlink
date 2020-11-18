package offchainreporting

import (
	"strings"
	"sync"
	"time"

	gethCommon "github.com/ethereum/go-ethereum/common"

	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"

	"github.com/smartcontractkit/libocr/gethwrappers/offchainaggregator"
	"github.com/smartcontractkit/libocr/offchainreporting/confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting/types"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	OCRContractConfigSet = getConfigSetHash()
)

var (
	_ ocrtypes.ContractConfigSubscription = &OCRContractConfigSubscription{}
	_ eth.LogListener                     = &OCRContractConfigSubscription{}
)

const OCRContractConfigSubscriptionHandleLogTimeout = 5 * time.Second

type OCRContractConfigSubscription struct {
	logger            logger.Logger
	contractAddress   gethCommon.Address
	ch                chan ocrtypes.ContractConfig
	chIncoming        chan ocrtypes.ContractConfig
	processLogsWorker utils.SleeperTask
	queue             []ocrtypes.ContractConfig
	queueMu           sync.Mutex
	oc                *OCRContractConfigTracker
	closer            sync.Once
	chStop            chan struct{}
}

func (sub *OCRContractConfigSubscription) start() {
	sub.processLogsWorker = utils.NewSleeperTask(
		utils.SleeperTaskFuncWorker(sub.processLogs),
	)
}

func (sub *OCRContractConfigSubscription) processLogs() {
	sub.queueMu.Lock()
	defer sub.queueMu.Unlock()
	for len(sub.queue) > 0 {
		cc := sub.queue[0]
		sub.queue = sub.queue[1:]

		select {
		// NOTE: This is thread-safe because HandleLog cannot be called concurrently with Unregister due to the design of LogBroadcaster
		// It will never send on closed channel
		case sub.ch <- cc:
		case <-time.After(OCRContractConfigSubscriptionHandleLogTimeout):
			sub.logger.Error("OCRContractConfigSubscription HandleLog timed out waiting on receive channel")
		case <-sub.chStop:
			return
		}
	}
}

// OnConnect complies with LogListener interface
func (sub *OCRContractConfigSubscription) OnConnect() {}

// OnDisconnect complies with LogListener interface
func (sub *OCRContractConfigSubscription) OnDisconnect() {}

// HandleLog complies with LogListener interface
func (sub *OCRContractConfigSubscription) HandleLog(lb eth.LogBroadcast, err error) {
	if err != nil {
		sub.logger.Errorw("OCRContract: error in previous LogListener", "err", err)
		return
	}

	was, err := lb.WasAlreadyConsumed()
	if err != nil {
		sub.logger.Errorw("OCRContract: could not determine if log was already consumed", "error", err)
		return
	} else if was {
		return
	}

	topics := lb.RawLog().Topics
	if len(topics) == 0 {
		return
	}
	switch topics[0] {
	case OCRContractConfigSet:
		raw := lb.RawLog()
		if raw.Address != sub.contractAddress {
			sub.logger.Errorf("log address of 0x%x does not match configured contract address of 0x%x", raw.Address, sub.contractAddress)
			return
		}
		configSet, err2 := sub.oc.contractFilterer.ParseConfigSet(raw)
		if err2 != nil {
			sub.logger.Errorw("could not parse config set", "err", err2)
			return
		}
		configSet.Raw = lb.RawLog()
		cc := confighelper.ContractConfigFromConfigSetEvent(*configSet)

		sub.queueMu.Lock()
		defer sub.queueMu.Unlock()
		sub.queue = append(sub.queue, cc)
		sub.processLogsWorker.WakeUp()
	default:
	}

	err = lb.MarkConsumed()
	if err != nil {
		sub.logger.Errorw("OCRContract: could not mark log consumed", "error", err)
		return
	}
}

// IsV2Job complies with LogListener interface
func (sub *OCRContractConfigSubscription) IsV2Job() bool {
	return true
}

// JobIDV2 complies with LogListener interface
func (sub *OCRContractConfigSubscription) JobIDV2() int32 {
	return sub.oc.jobID
}

// JobID complies with LogListener interface
func (sub *OCRContractConfigSubscription) JobID() *models.ID {
	return &models.ID{}
}

// Configs complies with ContractConfigSubscription interface
func (sub *OCRContractConfigSubscription) Configs() <-chan ocrtypes.ContractConfig {
	return sub.ch
}

// Close complies with ContractConfigSubscription interface
func (sub *OCRContractConfigSubscription) Close() {
	sub.closer.Do(func() {
		close(sub.chStop)
		sub.oc.logBroadcaster.Unregister(sub.oc.contractAddress, sub)
		err := sub.processLogsWorker.Stop()
		if err != nil {
			sub.logger.Error(err)
		}

		close(sub.ch)
	})
}

func getConfigSetHash() gethCommon.Hash {
	abi, err := abi.JSON(strings.NewReader(offchainaggregator.OffchainAggregatorABI))
	if err != nil {
		panic("could not parse OffchainAggregator ABI: " + err.Error())
	}
	return abi.Events["ConfigSet"].ID
}
