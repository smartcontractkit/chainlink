package ccip

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"slices"
	"sync"
	"time"

	"go.uber.org/atomic"
	"golang.org/x/sync/errgroup"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"

	"github.com/smartcontractkit/chainlink-testing-framework/seth"

	chainselectors "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/nonce_manager"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/environment/crib"
)

const (
	transmitted = iota
	committed
	executed
	tickerDuration      = 3 * time.Minute
	SubscriptionTimeout = 1 * time.Minute
)

var (
	fundingAmount = new(big.Int).Mul(deployment.UBigInt(100), deployment.UBigInt(1e18)) // 100 eth
)

type finalSeqNrReport struct {
	sourceChainSelector uint64
	expectedSeqNrRange  ccipocr3.SeqNumRange
}

func subscribeTransmitEvents(
	ctx context.Context,
	lggr logger.Logger,
	onRamp onramp.OnRampInterface,
	otherChains []uint64,
	startBlock *uint64,
	srcChainSel uint64,
	loadFinished chan struct{},
	client deployment.OnchainClient,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
	finalSeqNrCommitChannels map[uint64]chan finalSeqNrReport,
	finalSeqNrExecChannels map[uint64]chan finalSeqNrReport,
) {
	defer wg.Done()
	lggr.Infow("starting transmit event subscriber for ",
		"srcChain", srcChainSel,
		"startblock", startBlock,
	)

	seqNums := make(map[testhelpers.SourceDestPair]SeqNumRange)
	for _, cs := range otherChains {
		seqNums[testhelpers.SourceDestPair{
			SourceChainSelector: srcChainSel,
			DestChainSelector:   cs,
		}] = SeqNumRange{
			// we use the maxuint as a sentinel value here to ensure we always get the lowest possible seqnum
			Start: atomic.NewUint64(math.MaxUint64),
			End:   atomic.NewUint64(0),
		}
	}

	sink := make(chan *onramp.OnRampCCIPMessageSent)
	subscription := event.Resubscribe(SubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		return onRamp.WatchCCIPMessageSent(&bind.WatchOpts{
			Context: ctx,
			Start:   startBlock,
		}, sink, nil, nil)
	})
	defer subscription.Unsubscribe()

	for {
		select {
		case <-subscription.Err():
			return
		case event := <-sink:
			lggr.Debugw("received transmit event for",
				"srcChain", srcChainSel,
				"destChain", event.DestChainSelector,
				"sequenceNumber", event.SequenceNumber)

			blockNum := event.Raw.BlockNumber
			header, err := client.HeaderByNumber(ctx, new(big.Int).SetUint64(blockNum))
			if err != nil {
				lggr.Errorw("error getting header by number")
			}
			data := messageData{
				eventType: transmitted,
				srcDstSeqNum: srcDstSeqNum{
					src:    srcChainSel,
					dst:    event.DestChainSelector,
					seqNum: event.SequenceNumber,
				},
			}
			if header != nil {
				data.timestamp = header.Time
			}
			metricPipe <- data
			csPair := testhelpers.SourceDestPair{
				SourceChainSelector: srcChainSel,
				DestChainSelector:   event.DestChainSelector,
			}
			// always store the lowest seen number as the start seq num
			if event.SequenceNumber < seqNums[csPair].Start.Load() {
				seqNums[csPair].Start.Store(event.SequenceNumber)
			}

			// always store the greatest sequence number we have seen as the maximum
			if event.SequenceNumber > seqNums[csPair].End.Load() {
				seqNums[csPair].End.Store(event.SequenceNumber)
			}
		case <-ctx.Done():
			lggr.Errorw("received context cancel signal for transmit watcher",
				"srcChain", srcChainSel)
			return
		case <-loadFinished:
			lggr.Debugw("load finished, closing transmit watchers", "srcChainSel", srcChainSel)
			for csPair, seqNums := range seqNums {
				lggr.Infow("pushing finalized sequence numbers for ",
					"srcChainSelector", srcChainSel,
					"destChainSelector", csPair.DestChainSelector,
					"seqNums", seqNums)
				finalSeqNrCommitChannels[csPair.DestChainSelector] <- finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNums.Start.Load()), ccipocr3.SeqNum(seqNums.End.Load()),
					},
				}

				finalSeqNrExecChannels[csPair.DestChainSelector] <- finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNums.Start.Load()), ccipocr3.SeqNum(seqNums.End.Load()),
					},
				}
			}
			return
		}
	}
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
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

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
	subscription := event.Resubscribe(SubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		return offRamp.WatchCommitReportAccepted(&bind.WatchOpts{
			Context: ctx,
			Start:   startBlock,
		}, sink)
	})
	defer subscription.Unsubscribe()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case <-subscription.Err():
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
							lggr.Errorw("error getting header by number")
						}
						data := messageData{
							eventType: committed,
							srcDstSeqNum: srcDstSeqNum{
								src:    mr.SourceChainSelector,
								dst:    chainSelector,
								seqNum: i,
							},
						}
						if header != nil {
							data.timestamp = header.Time
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
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

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
	subscription := event.Resubscribe(SubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		return offRamp.WatchExecutionStateChanged(&bind.WatchOpts{
			Context: ctx,
			Start:   startBlock,
		}, sink, nil, nil, nil)
	})
	defer subscription.Unsubscribe()
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case subErr := <-subscription.Err():
			lggr.Errorw("error in execution subscription",
				"err", subErr)
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
				lggr.Errorw("error getting header by number")
			}
			data := messageData{
				eventType: executed,
				srcDstSeqNum: srcDstSeqNum{
					src:    event.SourceChainSelector,
					dst:    chainSelector,
					seqNum: event.SequenceNumber,
				},
			}
			if header != nil {
				data.timestamp = header.Time
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

func subscribeAlreadyExecuted(
	ctx context.Context,
	destChain uint64,
	offRamp offramp.OffRampInterface,
	lggr logger.Logger,
) {
	sink := make(chan *offramp.OffRampSkippedAlreadyExecutedMessage)
	subscription := event.Resubscribe(SubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		return offRamp.WatchSkippedAlreadyExecutedMessage(&bind.WatchOpts{
			Context: ctx,
			Start:   nil,
		}, sink)
	})
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case subErr := <-subscription.Err():
			lggr.Errorw("error in alreadyExecuted subscription",
				"destChain", destChain,
				"err", subErr)
			return
		case ev := <-sink:
			lggr.Errorw("received already executed event", "seqNr", ev.SequenceNumber,
				"destChain", destChain,
				"sourceChain", ev.SourceChainSelector)
		}
	}
}

func subscribeSkippedIncorrectNonce(
	ctx context.Context,
	destChain uint64,
	nm nonce_manager.NonceManagerInterface,
	lggr logger.Logger,
) {
	sink := make(chan *nonce_manager.NonceManagerSkippedIncorrectNonce)
	subscription := event.Resubscribe(SubscriptionTimeout, func(_ context.Context) (event.Subscription, error) {
		return nm.WatchSkippedIncorrectNonce(&bind.WatchOpts{
			Context: ctx,
			Start:   nil,
		}, sink)
	})
	defer subscription.Unsubscribe()

	for {
		select {
		case <-ctx.Done():
			return
		case subErr := <-subscription.Err():
			lggr.Errorw("error in skipped incorrect nonce subscription",
				"destChain", destChain,
				"err", subErr)
			return
		case ev := <-sink:
			lggr.Errorw("received an incorrect nonce", "seqNr", ev.Nonce,
				"destChain", destChain,
				"sourceChain", ev.SourceChainSelector)
		}
	}
}

// fundAdditionalKeys will create len(targetChains) new addresses, and send funds to them on every targetChain
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
func reclaimFunds(lggr logger.Logger, e deployment.Environment, addressesByChain map[uint64][]*bind.TransactOpts, returnAddress common.Address) error {
	removeFundsFromAccounts := func(ctx context.Context, lggr logger.Logger, chain deployment.Chain, addresses []*bind.TransactOpts, returnAddress common.Address, sel uint64) error {
		for _, deployer := range addresses {
			balance, err := chain.Client.BalanceAt(ctx, deployer.From, nil)
			if err != nil {
				lggr.Warnw("could not get balance for deployer key",
					"err", err,
					"chain", sel,
					"account", deployer.From)
				continue
			}
			nonce, err := chain.Client.NonceAt(ctx, deployer.From, nil)
			if err != nil {
				lggr.Warnw("could not get latest nonce for deployer key",
					"err", err,
					"chain", sel,
					"account", deployer.From)
				continue
			}

			// Get the current gas price
			gasPrice, err := chain.Client.SuggestGasPrice(ctx)
			if err != nil {
				lggr.Warnw("could not get gas price",
					"err", err,
					"chain", sel)
			}

			tx := gethtypes.NewTx(&gethtypes.LegacyTx{
				Nonce:    nonce,
				To:       &returnAddress,
				Value:    balance.Sub(balance, big.NewInt(1e14)), // leave a little bit for gas
				GasPrice: gasPrice,
				Gas:      21000,
			})

			signedTx, err := deployer.Signer(deployer.From, tx)
			if err != nil {
				lggr.Errorw("could not sign transaction for sending funds to ",
					"chain", sel,
					"account", deployer.From,
					"err", err)
				return err
			}

			lggr.Infow("sending transaction for ", "account", deployer.From.String(), "chain", sel)
			err = chain.Client.SendTransaction(ctx, signedTx)
			if err != nil {
				lggr.Errorw("could not send transaction to address on ",
					"chain", sel,
					"address", deployer.From,
					"err", err)
				return err
			}
		}
		return nil
	}
	g := new(errgroup.Group)
	for sel, addresses := range addressesByChain {
		sel, addresses := sel, addresses
		g.Go(func() error {
			return removeFundsFromAccounts(e.GetContext(), lggr, e.Chains[sel], addresses, returnAddress, sel)
		})
	}

	return g.Wait()
}
