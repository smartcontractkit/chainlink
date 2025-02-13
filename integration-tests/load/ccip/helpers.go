package ccip

import (
	"context"
	"errors"
	"fmt"
	"math"
	"slices"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"
	"github.com/smartcontractkit/chainlink/deployment/environment/crib"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	chainselectors "github.com/smartcontractkit/chain-selectors"

	"math/big"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/v1_6_0/offramp"
)

const (
	transmitted = iota
	committed
	executed
	tickerDuration = 30 * time.Second
)

var (
	fundingAmount = new(big.Int).Mul(deployment.UBigInt(10), deployment.UBigInt(1e18)) // 100 eth
)

// todo: Have a different struct for commit/exec?
type LokiMetric struct {
	SequenceNumber uint64 `json:"sequence_number"`
	CommitDuration uint64 `json:"commit_duration"`
	ExecDuration   uint64 `json:"exec_duration"`
}

type finalSeqNrReport struct {
	sourceChainSelector uint64
	expectedSeqNrRange  ccipocr3.SeqNumRange
}

func subscribeCommitEvents(
	ctx context.Context,
	lggr logger.Logger,
	offRamp offramp.OffRampInterface,
	srcChains []uint64,
	startBlock *uint64,
	chainSelector uint64,
	client deployment.OnchainClient,
	finalSeqNrs chan finalSeqNrReport,
	errChan chan error,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()

	lggr.Infow("starting commit event subscriber for ",
		"destChain", chainSelector,
		"startblock", startBlock,
	)
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		// todo: seenMessages should hold a range to avoid hitting memory constraints
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	sink := make(chan *offramp.OffRampCommitReportAccepted)
	// todo: add event.Resubscriber if we move to unreliable rpcs
	subscription, err := offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
		Context: ctx,
		Start:   startBlock,
	}, sink)
	if err != nil {
		errChan <- err
		return
	}
	defer subscription.Unsubscribe()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case subErr := <-subscription.Err():
			errChan <- subErr
			return
		case report := <-sink:
			if len(report.BlessedMerkleRoots)+len(report.UnblessedMerkleRoots) > 0 {
				for _, mr := range append(report.BlessedMerkleRoots, report.UnblessedMerkleRoots...) {
					lggr.Infow("Received commit report ",
						"sourceChain", mr.SourceChainSelector,
						"destChain", chainSelector,
						"minSeqNr", mr.MinSeqNr,
						"maxSeqNr", mr.MaxSeqNr)

					// push metrics to state manager for eventual distribution to loki
					for i := mr.MinSeqNr; i <= mr.MaxSeqNr; i++ {
						blockNum := report.Raw.BlockNumber
						header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
						if err != nil {
							errChan <- err
						}
						data := messageData{
							eventType: committed,
							srcDstSeqNum: srcDstSeqNum{
								src:    mr.SourceChainSelector,
								dst:    chainSelector,
								seqNum: i,
							},
							timestamp: header.Time,
						}
						metricPipe <- data
						seenMessages[mr.SourceChainSelector] = append(seenMessages[mr.SourceChainSelector], i)
					}
				}
			}
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for commit report",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			errChan <- errors.New("timed out waiting for commit report")
			return

		case finalSeqNrUpdate, ok := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else if ok {
				// only add to range if channel is still open
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			lggr.Infow("ticking, checking committed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)
			for srcChain, seqNumRange := range expectedRange {
				// if this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// else, check if all expected sequence numbers have been seen
					// todo: We might need to modify if there are other non-load test txns on network
					if len(seenMessages[srcChain]) >= seqNumRange.Length() && slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						delete(expectedRange, srcChain)
						delete(seenMessages, srcChain)
						lggr.Infow("committed all sequence numbers for ",
							"sourceChain", srcChain,
							"destChain", chainSelector)
					}
				}
			}
			// if all chains have hit expected sequence numbers, return
			// we could instead push complete chains to an incrementer and compare size
			allComplete := true
			for c := range completedSrcChains {
				if !completedSrcChains[c] {
					allComplete = false
					break
				}
			}
			if allComplete {
				lggr.Infof("received commits from expected source chains for all expected sequence numbers to chainSelector %d", chainSelector)
				return
			}
		}
	}
}

func subscribeExecutionEvents(
	ctx context.Context,
	lggr logger.Logger,
	offRamp offramp.OffRampInterface,
	srcChains []uint64,
	startBlock *uint64,
	chainSelector uint64,
	client deployment.OnchainClient,
	finalSeqNrs chan finalSeqNrReport,
	errChan chan error,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()

	lggr.Infow("starting execution event subscriber for ",
		"destChain", chainSelector,
		"startblock", startBlock,
	)
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	sink := make(chan *offramp.OffRampExecutionStateChanged)
	// todo: add event.Resubscriber if we move to unreliable rpcs
	subscription, err := offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
		Context: ctx,
		Start:   startBlock,
	}, sink, nil, nil, nil)
	if err != nil {
		errChan <- err
		return
	}
	defer subscription.Unsubscribe()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case subErr := <-subscription.Err():
			lggr.Errorw("error in execution subscription",
				"err", subErr)
			errChan <- subErr
			return
		case event := <-sink:
			lggr.Debugw("received execution event for",
				"destChain", chainSelector,
				"sourceChain", event.SourceChainSelector,
				"sequenceNumber", event.SequenceNumber,
				"blockNumber", event.Raw.BlockNumber)
			// push metrics to loki here
			blockNum := event.Raw.BlockNumber
			header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
			if err != nil {
				errChan <- err
			}
			data := messageData{
				eventType: executed,
				srcDstSeqNum: srcDstSeqNum{
					src:    event.SourceChainSelector,
					dst:    chainSelector,
					seqNum: event.SequenceNumber,
				},
				timestamp: header.Time,
			}
			metricPipe <- data
			seenMessages[event.SourceChainSelector] = append(seenMessages[event.SourceChainSelector], event.SequenceNumber)

		case <-ctx.Done():
			lggr.Errorw("timed out waiting for execution event",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange,
				"seenMessages", seenMessages,
				"completedSrcChains", completedSrcChains)
			errChan <- errors.New("timed out waiting for execution event")
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			lggr.Infow("ticking, checking executed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)

			for srcChain, seqNumRange := range expectedRange {
				// if this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// else, check if all expected sequence numbers have been seen
					if len(seenMessages[srcChain]) >= seqNumRange.Length() && slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						lggr.Infow("executed all sequence numbers for ",
							"destChain", chainSelector,
							"sourceChain", srcChain,
							"seqNumRange", seqNumRange)
					}
				}
			}
			// if all chains have hit expected sequence numbers, return
			allComplete := true
			for c := range completedSrcChains {
				if !completedSrcChains[c] {
					allComplete = false
					break
				}
			}
			if allComplete {
				lggr.Infow("all messages have been executed for all expected sequence numbers",
					"destChain", chainSelector)
				return
			}
		}
	}
}

// this function will create len(targetChains) new addresses, and send funds to them on every targetChain
func fundAdditionalKeys(lggr logger.Logger, e deployment.Environment, destChains []uint64) (map[uint64][]*bind.TransactOpts, error) {
	deployerMap := make(map[uint64][]*bind.TransactOpts)
	addressMap := make(map[uint64][]common.Address)
	numAccounts := len(destChains)
	for chain := range e.Chains {
		deployerMap[chain] = make([]*bind.TransactOpts, 0, numAccounts)
		addressMap[chain] = make([]common.Address, 0, numAccounts)
		for range numAccounts {
			addr, pk, err := seth.NewAddress()
			if err != nil {
				return nil, fmt.Errorf("failed to create new address: %w", err)
			}
			pvtKey, err := crypto.HexToECDSA(pk)
			if err != nil {
				return nil, fmt.Errorf("failed to convert private key to ECDSA: %w", err)
			}
			chainID, err := chainselectors.ChainIdFromSelector(chain)
			if err != nil {
				return nil, fmt.Errorf("could not get chain id from selector: %w", err)
			}

			deployer, err := bind.NewKeyedTransactorWithChainID(pvtKey, new(big.Int).SetUint64(chainID))
			if err != nil {
				return nil, fmt.Errorf("failed to create transactor: %w", err)
			}
			deployerMap[chain] = append(deployerMap[chain], deployer)
			addressMap[chain] = append(addressMap[chain], common.HexToAddress(addr))
		}
	}

	g := new(errgroup.Group)
	for sel, addresses := range addressMap {
		sel, addresses := sel, addresses
		g.Go(func() error {
			return crib.SendFundsToAccounts(e.GetContext(), lggr, e.Chains[sel], addresses, fundingAmount, sel)
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}
	return deployerMap, nil
}
