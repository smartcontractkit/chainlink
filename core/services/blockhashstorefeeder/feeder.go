package blockhashstorefeeder

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"go.uber.org/multierr"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/blockhash_store"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/null"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/eth"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore"
	"github.com/smartcontractkit/chainlink/core/services/log"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type feeder struct {
	utils.StartStopOnce

	config Config
	lggr   logger.Logger
	txm    bulletprooftxmanager.TxManager
	job    job.Job

	blockhashStore    *blockhash_store.BlockhashStore
	blockhashStoreABI *abi.ABI

	parseLogFunc  log.ParseLogFunc
	requestTopic  common.Hash
	responseTopic common.Hash

	wg       *sync.WaitGroup
	stopChan chan struct{}

	ethClient eth.Client

	ethKeystore keystore.Eth
}

// Close implements job.Service.Close
func (f *feeder) Close() error {
	return f.StopOnce("BlockhashStoreFeeder", func() error {
		close(f.stopChan)
		f.wg.Wait()
		return nil
	})
}

// Start implements job.Service.Start
func (f *feeder) Start() error {
	return f.StartOnce("BlockhashStoreFeeder", func() error {
		spec := f.job.BlockhashStoreFeederSpec

		// Feeder periodically computes a set of logs which can be fulfilled.
		f.wg.Add(1)
		go func() {
			f.runFeeder(spec.PollPeriod, f.wg)
		}()

		return nil
	})
}

func (f *feeder) runFeeder(pollPeriod time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(pollPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-f.stopChan:
			return
		case <-ticker.C:
			f.process()
		}
	}
}

func (f *feeder) process() error {
	logs, latestNumber, err := f.getLogs()
	if err != nil {
		return err
	}

	unfulfilledNumbers, err := f.unfulfilledBlockNumbers(logs, latestNumber)
	if err != nil {
		return errors.Wrap(err, "could not get unfulfilled block numbers")
	}

	return f.storeBlockhashes(unfulfilledNumbers)
}

func (f *feeder) storeBlockhashes(unfulfilledNumbers []uint64) error {
	sendingKeys, err := f.ethKeystore.SendingKeys()
	if err != nil {
		return errors.Wrap(err, "could not get eth sending keys to store blockhashes")
	}
	fromAddress := sendingKeys[0].Address.Address()
	errs := []error{}
	for _, blockNumber := range unfulfilledNumbers {
		f.lggr.Infof("Storing blockhash of block %d", blockNumber)

		// Check if it's already in the store
		_, err := f.blockhashStore.GetBlockhash(nil, big.NewInt(int64(blockNumber)))
		if err == nil {
			// blockhash is in store, nothing to do
			continue
		}

		payload, err := f.blockhashStoreABI.Pack("store", big.NewInt(int64(blockNumber)))
		if err != nil {
			return errors.Wrapf(err, "could not encode BlockhashStore.store call with args: %d", blockNumber)
		}

		// Store it since it's not there with bptxm
		_, err = f.txm.CreateEthTransaction(bulletprooftxmanager.NewTx{
			FromAddress:      fromAddress,
			ToAddress:        f.blockhashStore.Address(),
			EncodedPayload:   payload,
			MinConfirmations: null.Uint32From(uint32(f.config.MinRequiredOutgoingConfirmations())),
			Strategy:         bulletprooftxmanager.NewSendEveryStrategy(false), // TODO: should we simulate?
		})
		if err != nil {
			f.lggr.With("err", err).Error("could not create blockhash storage tx")
			errs = append(errs, err)
			continue
		}
	}
	return multierr.Combine(errs...)
}

func (f *feeder) getLogs() ([]types.Log, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	currentHead, err := f.ethClient.HeadByNumber(ctx, nil)
	if err != nil {
		return nil, -1, errors.Wrap(err, "could not fetch latest head")
	}
	fromBlock := currentHead.Number - int64(f.job.BlockhashStoreFeederSpec.LookbackBlocks)
	logs, err := f.ethClient.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		Addresses: []common.Address{f.job.BlockhashStoreFeederSpec.CoordinatorAddress.Address()},
		Topics: [][]common.Hash{
			{f.requestTopic, f.responseTopic},
		},
	})
	if err != nil {
		return nil, -1, errors.Wrap(err, "could not fetch logs")
	}

	return logs, currentHead.Number, nil
}

// NOTE: hopefully will be better when generics rolls around.
func (f *feeder) unfulfilledBlockNumbers(logs []types.Log, currentBlockNumber int64) (blockNumbers []uint64, err error) {
	if f.job.BlockhashStoreFeederSpec.VRFVersion == 1 {
		return f.unfulfilledBlockNumbersV1(logs, currentBlockNumber)
	}
	return f.unfulfilledBlockNumbersV2(logs, currentBlockNumber)
}

func (f *feeder) oldEnough(reqBlockNumber, currentBlockNumber uint64) bool {
	return reqBlockNumber+uint64(f.job.BlockhashStoreFeederSpec.BlocksToWaitBeforeStoring) <= currentBlockNumber
}

func (f *feeder) unfulfilledBlockNumbersV1(logs []types.Log, currentBlockNumber int64) (blockNumbers []uint64, err error) {
	requests := map[[32]byte]*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}

	// extract requested events
	for _, log := range logs {
		abiLog, err := f.parseLogFunc(log)
		if err != nil {
			continue
		}
		if abiLog.Topic() == f.requestTopic {
			req, ok := abiLog.(*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest)
			if !ok {
				// programming error
				return nil, fmt.Errorf("could not cast abi log to VRFCoordinatorRandomnessRequest: %+v", abiLog)
			}
			requests[req.RequestID] = req
		}
	}

	// delete any requests that were fulfilled
	for _, log := range logs {
		abiLog, err := f.parseLogFunc(log)
		if err != nil {
			continue
		}
		if abiLog.Topic() == f.responseTopic {
			resp, ok := abiLog.(*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled)
			if !ok {
				// programming error
				return nil, fmt.Errorf("could not cast abi log to VRFCoordinatorRandomnessRequestFulfilled: %+v", abiLog)
			}
			delete(requests, resp.RequestId)
		}
	}

	// delete any requests that are less than configured wait block
	unfulfilled := map[[32]byte]*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest{}
	for requestID, request := range requests {
		if f.oldEnough(request.Raw.BlockNumber, uint64(currentBlockNumber)) {
			unfulfilled[requestID] = request
		}
	}

	return f.getBlockNumbersV1(unfulfilled), nil
}

func (f *feeder) getBlockNumbersV1(unfulfilledRequests map[[32]byte]*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequest) (blockNumbers []uint64) {
	for _, req := range unfulfilledRequests {
		blockNumbers = append(blockNumbers, req.Raw.BlockNumber)
	}
	return blockNumbers
}

func (f *feeder) unfulfilledBlockNumbersV2(logs []types.Log, currentBlockNumber int64) (blockNumbers []uint64, err error) {
	requests := map[*big.Int]*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}

	// extract requested events
	for _, log := range logs {
		abiLog, err := f.parseLogFunc(log)
		if err != nil {
			continue
		}
		if abiLog.Topic() == f.requestTopic {
			req, ok := abiLog.(*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested)
			if !ok {
				// programming error
				return nil, fmt.Errorf("could not cast abi log to VRFCoordinatorV2RandomWordsRequested: %+v", abiLog)
			}
			requests[req.RequestId] = req
		}
	}

	// delete any requests that were fulfilled
	for _, log := range logs {
		abiLog, err := f.parseLogFunc(log)
		if err != nil {
			continue
		}
		if abiLog.Topic() == f.responseTopic {
			resp, ok := abiLog.(*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsFulfilled)
			if !ok {
				// programming error
				return nil, fmt.Errorf("could not cast abi log to VRFCoordinatorV2RandomWordsFulfilled: %+v", abiLog)
			}
			delete(requests, resp.RequestId)
		}
	}

	// delete any requests that are less than configured wait block
	unfulfilled := map[*big.Int]*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested{}
	for requestID, request := range requests {
		if f.oldEnough(request.Raw.BlockNumber, uint64(currentBlockNumber)) {
			unfulfilled[requestID] = request
		}
	}

	return f.getBlockNumbersV2(unfulfilled), nil
}

func (f *feeder) getBlockNumbersV2(unfulfilledRequests map[*big.Int]*vrf_coordinator_v2.VRFCoordinatorV2RandomWordsRequested) (blockNumbers []uint64) {
	for _, request := range unfulfilledRequests {
		blockNumbers = append(blockNumbers, request.Raw.BlockNumber)
	}
	return blockNumbers
}
