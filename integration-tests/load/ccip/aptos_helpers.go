package ccip

import (
	"context"
	"fmt"
	"maps"
	"math"
	"slices"
	"sync"
	"testing"
	"time"

	"github.com/aptos-labs/aptos-go-sdk"
	"go.uber.org/atomic"

	aptos_ccip_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp"
	module_offramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_offramp/offramp"
	aptos_ccip_onramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp"
	module_onramp "github.com/smartcontractkit/chainlink-aptos/bindings/ccip_onramp/onramp"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared"

	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"

	"github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
)

func fundAdditionalAptosKeys(
	t *testing.T,
	signer aptos.TransactionSigner,
	e cldf.Environment,
	destChains []uint64,
	fundingAmount uint64,
) (map[uint64][]aptos.Account, error) {
	funded := make(map[uint64][]aptos.Account, len(destChains))

	for _, chain := range e.BlockChains.AptosChains() {
		numAccounts := len(destChains)
		funded[chain.ChainSelector()] = make([]aptos.Account, 0, numAccounts)

		for range numAccounts {
			account, err := aptos.NewEd25519Account()
			if err != nil {
				return nil, fmt.Errorf("failed to create new aptos account: %w", err)
			}
			pk, err := account.PrivateKeyString()
			if err != nil {
				return nil, fmt.Errorf("failed to get private key string: %w", err)
			}

			t.Logf(
				"Created new Aptos sender on Chain %d | Address: %s | PrivateKey: %s",
				chain.ChainSelector(),
				account.Address.String(),
				pk,
			)

			memory.FundAptosAccount(t, signer, account.AccountAddress(), fundingAmount, chain.Client)
			funded[chain.ChainSelector()] = append(funded[chain.ChainSelector()], *account)
		}
	}
	return funded, nil
}

func fundAptosLoadAccountsWithBnM(
	lggr logger.Logger,
	e cldf.Environment,
	aptosChainSelector uint64,
	loadAccounts []*aptos.Account,
) error {
	lggr.Infow("Funding Aptos load test accounts with BnM tokens",
		"chain", aptosChainSelector,
		"numAccounts", len(loadAccounts))

	addresses, err := e.ExistingAddresses.AddressesForChain(aptosChainSelector)
	if err != nil {
		return fmt.Errorf("failed to get addresses for chain %d: %w", aptosChainSelector, err)
	}

	var bnmTokenAddr aptos.AccountAddress
	var found bool

	for addrStr, typeAndVersion := range addresses {
		if typeAndVersion.Type == cldf.ContractType(shared.CCIPBnMSymbol) {
			err := bnmTokenAddr.ParseStringRelaxed(addrStr)
			if err != nil {
				return fmt.Errorf("failed to parse BnM address %s: %w", addrStr, err)
			}
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("CCIP-BnM FA metadata address not found for chain %d", aptosChainSelector)
	}

	deployer := e.BlockChains.AptosChains()[aptosChainSelector].DeployerSigner
	client := e.BlockChains.AptosChains()[aptosChainSelector].Client

	deployerAddr := deployer.AccountAddress()
	lggr.Infow("Starting BnM token transfers for Aptos accounts",
		"tokenAddress", bnmTokenAddr.String(),
		"deployerAddress", deployerAddr.String())

	// Transfer small amount since token is 8 decimals
	transferAmount := uint64(1e7)

	aptosClient, ok := client.(*aptos.Client)
	if !ok {
		return fmt.Errorf("client is not of type *aptos.Client")
	}

	for i, loadAccount := range loadAccounts {
		lggr.Infow("Transferring BnM tokens to Aptos load account",
			"accountIndex", i,
			"recipient", loadAccount.Address.String(),
			"amount", transferAmount)

		payload, err := aptos.FungibleAssetPrimaryStoreTransferPayload(
			&bnmTokenAddr,       // BnM token metadata address
			loadAccount.Address, // recipient address
			transferAmount,      // amount to transfer
		)
		if err != nil {
			return fmt.Errorf("failed to create transfer payload for account %s: %w", loadAccount.Address.String(), err)
		}

		tx, err := aptosClient.BuildSignAndSubmitTransaction(deployer, aptos.TransactionPayload{
			Payload: payload,
		})
		if err != nil {
			return fmt.Errorf("failed to build/sign/submit transfer transaction for account %s: %w", loadAccount.Address.String(), err)
		}

		res, err := aptosClient.WaitForTransaction(tx.Hash)
		if err != nil {
			return fmt.Errorf("failed to wait for BnM transfer transaction for account %s: %w", loadAccount.Address.String(), err)
		}

		if !res.Success {
			return fmt.Errorf("BnM transfer transaction failed for account %s: %s", loadAccount.Address.String(), res.VmStatus)
		}

		lggr.Infow("Successfully transferred BnM tokens to Aptos account",
			"recipient", loadAccount.Address.String(),
			"amount", transferAmount,
			"txHash", tx.Hash)
	}

	lggr.Infow("Successfully funded all Aptos load test accounts with BnM tokens",
		"chain", aptosChainSelector,
		"numAccounts", len(loadAccounts),
		"amountPerAccount", transferAmount)

	return nil
}

func subscribeAptosTransmitEvents(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	onrampAddress aptos.AccountAddress,
	otherChains []uint64,
	startVersion *uint64,
	srcChainSel uint64,
	loadFinished chan struct{},
	client aptos.AptosRpcClient,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
	finalSeqNrCommitChannels map[uint64]chan finalSeqNrReport,
	finalSeqNrExecChannels map[uint64]chan finalSeqNrReport,
) {
	defer wg.Done()
	lggr.Infow("starting aptos chain transmit event subscriber for ",
		"srcChain", srcChainSel,
		"otherChains", otherChains,
		"startVersion", startVersion,
	)

	seqNums := make(map[testhelpers.SourceDestPair]SeqNumRange)
	for _, cs := range otherChains {
		seqNums[testhelpers.SourceDestPair{
			SourceChainSelector: srcChainSel,
			DestChainSelector:   cs,
		}] = SeqNumRange{
			// Use maxuint as a sentinel value to ensure we get the lowest possible seqnum
			Start: atomic.NewUint64(math.MaxUint64),
			End:   atomic.NewUint64(0),
		}
	}

	done := make(chan any)
	boundOnRamp := aptos_ccip_onramp.Bind(onrampAddress, client)

	onrampStateAddress, err := boundOnRamp.Onramp().GetStateAddress(nil)
	if err != nil {
		lggr.Errorw("Error getting onramp state address",
			"err", err,
			"onrampAddress", onrampAddress.String())
		return
	}

	sink, errCh := testhelpers.AptosEventEmitter[module_onramp.CCIPMessageSent](
		t,
		client,
		onrampStateAddress,
		onrampAddress.StringLong()+"::onramp::OnRampState",
		"ccip_message_sent_events",
		startVersion,
		done,
	)
	defer close(done)

	for {
		select {
		case err := <-errCh:
			lggr.Errorw("error in aptos event emitter for subscribing transmit events",
				"srcChain", srcChainSel,
				"err", err)
			return

		case eventWithVersion := <-sink:
			event := eventWithVersion.Event
			lggr.Debugw("received aptos transmit event for",
				"srcChain", srcChainSel,
				"destChain", event.DestChainSelector,
				"sequenceNumber", event.SequenceNumber,
				"version", eventWithVersion.Version)

			// Push metrics to state manager
			data := messageData{
				eventType: transmitted,
				srcDstSeqNum: srcDstSeqNum{
					src:    srcChainSel,
					dst:    event.DestChainSelector,
					seqNum: event.SequenceNumber,
				},
				timestamp: uint64(time.Now().Unix()), // todo: do we require a real timestamp here?
			}

			metricPipe <- data
			csPair := testhelpers.SourceDestPair{
				SourceChainSelector: srcChainSel,
				DestChainSelector:   event.DestChainSelector,
			}

			// Initialize the sequence number range if it doesn't exist
			if seqNums[csPair].Start == nil {
				lggr.Infow("Initializing sequence number range for new chain pair", "csPair", csPair)
				seqNums[csPair] = SeqNumRange{
					Start: atomic.NewUint64(math.MaxUint64),
					End:   atomic.NewUint64(0),
				}
			}

			if event.SequenceNumber < seqNums[csPair].Start.Load() {
				seqNums[csPair].Start.Store(event.SequenceNumber)
			}
			if event.SequenceNumber > seqNums[csPair].End.Load() {
				seqNums[csPair].End.Store(event.SequenceNumber)
			}

		case <-ctx.Done():
			lggr.Errorw("received context cancel signal for transmit watcher",
				"srcChain", srcChainSel)
			return

		case <-loadFinished:
			// When load is finished, notify commit and execution subscribers about sequence numbers
			for csPair, seqNumRange := range maps.Clone(seqNums) {
				lggr.Infow("pushing finalized sequence numbers for ",
					"csPair", csPair,
					"seqNumRange", seqNumRange)

				report := finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNumRange.Start.Load()),
						ccipocr3.SeqNum(seqNumRange.End.Load()),
					},
				}

				finalSeqNrCommitChannels[csPair.DestChainSelector] <- report
				finalSeqNrExecChannels[csPair.DestChainSelector] <- report
			}
			return
		}
	}
}

func subscribeAptosCommitEvents(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	offrampAddress aptos.AccountAddress,
	srcChains []uint64,
	startVersion *uint64,
	chainSelector uint64,
	client aptos.AptosRpcClient,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting aptos commit event subscriber for ",
		"destChain", chainSelector,
		"startVersion", startVersion,
	)

	// Track seen messages and expected ranges
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)

	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	done := make(chan any)
	boundOffRamp := aptos_ccip_offramp.Bind(offrampAddress, client)
	offRampStateAddress, err := boundOffRamp.Offramp().GetStateAddress(nil)
	if err != nil {
		lggr.Errorw("Error getting offramp state address",
			"err", err,
			"offrampAddress", offrampAddress.String())
		return
	}

	sink, errCh := testhelpers.AptosEventEmitter[module_offramp.CommitReportAccepted](
		t,
		client,
		offRampStateAddress,
		offrampAddress.StringLong()+"::offramp::OffRampState",
		"commit_report_accepted_events",
		startVersion,
		done,
	)
	defer close(done)

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			lggr.Errorw("error in aptos event emitter for subscribing commit events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains,
				"err", err)
			return

		case eventWithVersion := <-sink:
			report := eventWithVersion.Event

			// Process both blessed and unblessed merkle roots
			allRoots := append(report.BlessedMerkleRoots, report.UnblessedMerkleRoots...)
			for _, mr := range allRoots {
				// Push metrics for each sequence number in the range
				for i := mr.MinSeqNr; i <= mr.MaxSeqNr; i++ {
					data := messageData{
						eventType: committed,
						srcDstSeqNum: srcDstSeqNum{
							src:    mr.SourceChainSelector,
							dst:    chainSelector,
							seqNum: i,
						},
						timestamp: uint64(time.Now().Unix()), // todo: do we require a real timestamp here?
					}
					metricPipe <- data
					seenMessages[mr.SourceChainSelector] = append(seenMessages[mr.SourceChainSelector], i)
				}
			}

		case <-ctx.Done():
			lggr.Errorw("timed out waiting for commit report",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 ||
				finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			lggr.Infow("ticking, checking committed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)

			for srcChain, seqNumRange := range expectedRange {
				// If this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// Check if all expected sequence numbers have been seen
					if len(seenMessages[srcChain]) >= seqNumRange.Length() &&
						slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						delete(expectedRange, srcChain)
						delete(seenMessages, srcChain)
						lggr.Infow("committed all sequence numbers for ",
							"sourceChain", srcChain,
							"destChain", chainSelector)
					}
				}
			}

			// Check if all chains have hit expected sequence numbers
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

func subscribeAptosExecutionEvents(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	offrampAddress aptos.AccountAddress,
	srcChains []uint64,
	startVersion *uint64,
	chainSelector uint64,
	client aptos.AptosRpcClient,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting aptos execution event subscriber for ",
		"destChain", chainSelector,
		"startVersion", startVersion,
	)

	// Track seen messages and expected ranges
	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)

	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	done := make(chan any)
	boundOffRamp := aptos_ccip_offramp.Bind(offrampAddress, client)
	offRampStateAddress, err := boundOffRamp.Offramp().GetStateAddress(nil)
	if err != nil {
		lggr.Errorw("Error getting offramp state address",
			"err", err,
			"offrampAddress", offrampAddress.String())
		return
	}

	sink, errCh := testhelpers.AptosEventEmitter[module_offramp.ExecutionStateChanged](
		t,
		client,
		offRampStateAddress,
		offrampAddress.StringLong()+"::offramp::OffRampState",
		"execution_state_changed_events",
		startVersion,
		done,
	)
	defer close(done)

	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	for {
		select {
		case err := <-errCh:
			lggr.Errorw("error in aptos event emitter for execution events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains,
				"err", err)
			return

		case eventWithVersion := <-sink:
			event := eventWithVersion.Event

			// Skip events that are not in the Success state
			if event.State != testhelpers.EXECUTION_STATE_SUCCESS {
				continue
			}

			lggr.Infow("received aptos execution event for",
				"destChain", chainSelector,
				"sourceChain", event.SourceChainSelector,
				"sequenceNumber", event.SequenceNumber,
				"version", eventWithVersion.Version,
				"executionState", event.State)

			// Push metrics
			data := messageData{
				eventType: executed,
				srcDstSeqNum: srcDstSeqNum{
					src:    event.SourceChainSelector,
					dst:    chainSelector,
					seqNum: event.SequenceNumber,
				},
				timestamp: uint64(time.Now().Unix()), // todo: do we require a real timestamp here?
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
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 ||
				finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
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
				// If this chain has already been marked as completed, skip
				if !completedSrcChains[srcChain] {
					// Check if all expected sequence numbers have been seen
					if len(seenMessages[srcChain]) >= seqNumRange.Length() &&
						slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						lggr.Infow("executed all sequence numbers for ",
							"destChain", chainSelector,
							"sourceChain", srcChain,
							"seqNumRange", seqNumRange)
					}
				}
			}

			// Check if all chains have hit expected sequence numbers
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
