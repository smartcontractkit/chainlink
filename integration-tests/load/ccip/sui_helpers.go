package ccip

import (
	"context"
	"fmt"
	"maps"
	"math"
	"os"
	"slices"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/scylladb/go-reflectx"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/pg"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil/sqltest"
	chain_reader_types "github.com/smartcontractkit/chainlink-common/pkg/types"
	sui_query "github.com/smartcontractkit/chainlink-common/pkg/types/query"
	crConfig "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	chainreader "github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"

	"github.com/block-vision/sui-go-sdk/models"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	suitestutils "github.com/smartcontractkit/chainlink-sui/relayer/testutils"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	suiState "github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview/sui"
)

// SUI CCIP Event structures
type SuiCCIPMessageSent struct {
	DestChainSelector uint64
	SequenceNumber    uint64
	Message           SuiRampMessage
}

type SuiRampMessage struct {
	Header         SuiMessageHeader
	Sender         string
	Data           []byte
	Receiver       []byte
	ExtraArgs      []byte
	FeeToken       string
	FeeTokenAmount uint64
	FeeValueJuels  uint64
}

type SuiMessageHeader struct {
	MessageId           []byte
	SourceChainSelector uint64
	DestChainSelector   uint64
	SequenceNumber      uint64
	Nonce               uint64
}

type SuiCommitReport struct {
	SourceChainSelector uint64
	MinSeqNr            uint64
	MaxSeqNr            uint64
	MerkleRoot          []byte
}

type SuiExecutionStateChanged struct {
	SourceChainSelector uint64
	SequenceNumber      uint64
	MessageId           []byte
	State               uint8
}

func subscribeSuiTransmitEvents(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	suiChain cldf_sui.Chain,
	chainState suiState.CCIPChainState,
	otherChains []uint64,
	startSlot uint64,
	srcChainSel uint64,
	loadFinished chan struct{},
	wg *sync.WaitGroup,
	metricPipe chan messageData,
	finalSeqNrCommitChannels map[uint64]chan finalSeqNrReport,
	finalSeqNrExecChannels map[uint64]chan finalSeqNrReport,
) {
	defer wg.Done()
	lggr.Infow("starting SUI transmit event subscriber for ",
		"srcChain", srcChainSel,
		"otherChains", otherChains,
		"startSlot", startSlot,
	)

	seqNums := make(map[testhelpers.SourceDestPair]SeqNumRange)
	processedEvents := make(map[string]bool) // Track processed events to avoid duplicates
	lastCheckpoint := startSlot - 1          // Start from the checkpoint before startSlot to ensure we don't miss any

	for _, cs := range otherChains {
		csPair := testhelpers.SourceDestPair{
			SourceChainSelector: srcChainSel,
			DestChainSelector:   cs,
		}

		// Initialize the sequence number range if it doesn't exist
		if seqNums[csPair].Start == nil {
			lggr.Infow("Initializing sequence number range for new chain pair", "csPair", csPair)
			seqNums[csPair] = SeqNumRange{
				Start: atomic.NewUint64(math.MaxUint64),
				End:   atomic.NewUint64(0),
			}
		}
	}

	// Create ChainReader for SUI events
	chainReader, indexerInstance, err := createSuiChainReader(ctx, t, lggr, suiChain, chainState)
	if err != nil {
		lggr.Errorw("Failed to create SUI chain reader", "error", err)
		return
	}
	defer chainReader.Close()
	defer indexerInstance.Close()

	// Start the chain reader and indexer
	err = chainReader.Start(ctx)
	if err != nil {
		lggr.Errorw("Failed to start SUI chain reader", "error", err)
		return
	}

	err = indexerInstance.Start(ctx)
	if err != nil {
		lggr.Errorw("Failed to start SUI indexer", "error", err)
		return
	}

	// Bind to the onramp contract
	err = chainReader.Bind(context.Background(), []chain_reader_types.BoundContract{{
		Name:    "onramp",
		Address: chainState.OnRampAddress,
	}})
	if err != nil {
		lggr.Errorw("Failed to bind onramp contract", "error", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second) // Poll every 1 second for responsive metrics
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			lggr.Errorw("received context cancel signal for SUI transmit watcher",
				"srcChain", srcChainSel)
			return
		case <-loadFinished:
			for csPair, seqNumRange := range maps.All(seqNums) {
				lggr.Infow("pushing finalized sequence numbers for ",
					"csPair", csPair,
					"seqNumRange", seqNumRange)
				finalSeqNrCommitChannels[csPair.DestChainSelector] <- finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNumRange.Start.Load()), ccipocr3.SeqNum(seqNumRange.End.Load()),
					},
				}

				finalSeqNrExecChannels[csPair.DestChainSelector] <- finalSeqNrReport{
					sourceChainSelector: csPair.SourceChainSelector,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNumRange.Start.Load()), ccipocr3.SeqNum(seqNumRange.End.Load()),
					},
				}
			}
			return
		case <-ticker.C:
			// Get the current latest checkpoint
			currentCheckpoint, err := suiChain.Client.SuiGetLatestCheckpointSequenceNumber(ctx)
			if err != nil {
				lggr.Errorw("Failed to get current SUI checkpoint", "error", err)
				continue
			}

			// Only query if there are new checkpoints since our last query
			if currentCheckpoint <= lastCheckpoint {
				continue
			}

			// Query for new CCIPMessageSent events from the last processed checkpoint
			var ccipSendEvent SuiCCIPMessageSent
			sequences, err := chainReader.QueryKey(
				ctx,
				chain_reader_types.BoundContract{
					Name:    "onramp",
					Address: chainState.OnRampAddress,
				},
				sui_query.KeyFilter{
					Key: "CCIPMessageSent",
					// TODO: Add checkpoint filtering when supported by the QueryKey API
					// For now, we'll filter in code below
				},
				sui_query.LimitAndSort{
					Limit: sui_query.Limit{
						Count:  100,
						Cursor: "",
					},
				},
				&ccipSendEvent,
			)
			if err != nil {
				lggr.Errorw("Failed to query SUI transmit events", "error", err)
				continue
			}

			newEventsFound := false
			for _, sequence := range sequences {
				event, ok := sequence.Data.(*SuiCCIPMessageSent)
				if !ok {
					lggr.Errorw("Failed to cast SUI event data")
					continue
				}

				// Create a unique key for this event to track if it's been processed
				eventKey := fmt.Sprintf("%d_%d_%d", srcChainSel, event.DestChainSelector, event.SequenceNumber)

				// Skip if we've already processed this event
				if processedEvents[eventKey] {
					continue
				}

				// Mark this event as processed
				processedEvents[eventKey] = true
				newEventsFound = true

				lggr.Debugw("received SUI transmit event for",
					"srcChain", srcChainSel,
					"destChain", event.DestChainSelector,
					"sequenceNumber", event.SequenceNumber)

				data := messageData{
					eventType: transmitted,
					srcDstSeqNum: srcDstSeqNum{
						src:    srcChainSel,
						dst:    event.DestChainSelector,
						seqNum: event.SequenceNumber,
					},
					timestamp: uint64(time.Now().Unix()),
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

				// Update sequence number ranges
				if event.SequenceNumber < seqNums[csPair].Start.Load() {
					seqNums[csPair].Start.Store(event.SequenceNumber)
				}
				if event.SequenceNumber > seqNums[csPair].End.Load() {
					seqNums[csPair].End.Store(event.SequenceNumber)
				}
			}

			// Update last processed checkpoint
			lastCheckpoint = currentCheckpoint

			if newEventsFound {
				lggr.Debugw("Processed new SUI transmit events",
					"srcChain", srcChainSel,
					"currentCheckpoint", currentCheckpoint,
					"startSlot", startSlot)
			}
		}
	}
}

func subscribeSuiCommitEvents(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	suiChain cldf_sui.Chain,
	chainState suiState.CCIPChainState,
	srcChains []uint64,
	startSlot uint64,
	chainSelector uint64,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting SUI commit event subscriber for ",
		"destChain", chainSelector,
		"startSlot", startSlot,
	)

	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	processedCommits := make(map[string]bool) // Track processed commit reports to avoid duplicates
	lastCheckpoint := startSlot - 1           // Start from the checkpoint before startSlot

	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	// Create ChainReader for SUI offramp events
	chainReader, indexerInstance, err := createSuiChainReader(ctx, t, lggr, suiChain, chainState)
	if err != nil {
		lggr.Errorw("Failed to create SUI chain reader for commit events", "error", err)
		return
	}
	defer chainReader.Close()
	defer indexerInstance.Close()

	// Start the chain reader and indexer
	err = chainReader.Start(ctx)
	if err != nil {
		lggr.Errorw("Failed to start SUI chain reader for commit events", "error", err)
		return
	}

	err = indexerInstance.Start(ctx)
	if err != nil {
		lggr.Errorw("Failed to start SUI indexer for commit events", "error", err)
		return
	}

	// Bind to the offramp contract
	err = chainReader.Bind(context.Background(), []chain_reader_types.BoundContract{{
		Name:    "offramp",
		Address: chainState.OffRampAddress,
	}})
	if err != nil {
		lggr.Errorw("Failed to bind offramp contract", "error", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second) // Use 1 second for consistency with transmit events
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for SUI commit report",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			return

		case finalSeqNrUpdate, ok := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 || finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else if ok {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			// Get the current latest checkpoint
			currentCheckpoint, err := suiChain.Client.SuiGetLatestCheckpointSequenceNumber(ctx)
			if err != nil {
				lggr.Errorw("Failed to get current SUI checkpoint for commit events", "error", err)
				continue
			}

			// Only query if there are new checkpoints since our last query
			if currentCheckpoint <= lastCheckpoint {
				continue
			}

			// Query for new commit events
			var commitReport SuiCommitReport
			sequences, err := chainReader.QueryKey(
				ctx,
				chain_reader_types.BoundContract{
					Name:    "offramp",
					Address: chainState.OffRampAddress,
				},
				sui_query.KeyFilter{
					Key: "CommitReportAccepted",
				},
				sui_query.LimitAndSort{
					Limit: sui_query.Limit{
						Count:  100,
						Cursor: "",
					},
				},
				&commitReport,
			)
			if err != nil {
				lggr.Errorw("Failed to query SUI commit events", "error", err)
				continue
			}

			newCommitsFound := false
			for _, sequence := range sequences {
				report, ok := sequence.Data.(*SuiCommitReport)
				if !ok {
					lggr.Errorw("Failed to cast SUI commit report data")
					continue
				}

				// Create a unique key for this commit report
				commitKey := fmt.Sprintf("%d_%d_%d_%d", report.SourceChainSelector, chainSelector, report.MinSeqNr, report.MaxSeqNr)

				// Skip if we've already processed this commit
				if processedCommits[commitKey] {
					continue
				}

				// Mark this commit as processed
				processedCommits[commitKey] = true
				newCommitsFound = true

				lggr.Infow("Received SUI commit report ",
					"sourceChain", report.SourceChainSelector,
					"destChain", chainSelector,
					"minSeqNr", report.MinSeqNr,
					"maxSeqNr", report.MaxSeqNr)

				// Push metrics for each sequence number in the range
				for i := report.MinSeqNr; i <= report.MaxSeqNr; i++ {
					data := messageData{
						eventType: committed,
						srcDstSeqNum: srcDstSeqNum{
							src:    report.SourceChainSelector,
							dst:    chainSelector,
							seqNum: i,
						},
						timestamp: uint64(time.Now().Unix()),
					}
					metricPipe <- data
					seenMessages[report.SourceChainSelector] = append(seenMessages[report.SourceChainSelector], i)
				}
			}

			// Update last processed checkpoint
			lastCheckpoint = currentCheckpoint

			if newCommitsFound {
				lggr.Debugw("Processed new SUI commit events",
					"destChain", chainSelector,
					"currentCheckpoint", currentCheckpoint,
					"startSlot", startSlot)
			}

			lggr.Infow("ticking, checking SUI committed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)

			for srcChain, seqNumRange := range expectedRange {
				if !completedSrcChains[srcChain] {
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

			// Check if all chains are complete
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

func subscribeSuiExecutionEvents(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	suiChain cldf_sui.Chain,
	chainState suiState.CCIPChainState,
	srcChains []uint64,
	startSlot uint64,
	chainSelector uint64,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting SUI execution event subscriber for ",
		"destChain", chainSelector,
		"startSlot", startSlot,
	)

	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	processedExecutions := make(map[string]bool) // Track processed execution events to avoid duplicates
	lastCheckpoint := startSlot - 1              // Start from the checkpoint before startSlot

	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	// Create ChainReader for SUI offramp execution events
	chainReader, indexerInstance, err := createSuiChainReader(ctx, t, lggr, suiChain, chainState)
	if err != nil {
		lggr.Errorw("Failed to create SUI chain reader for execution events", "error", err)
		return
	}
	defer chainReader.Close()
	defer indexerInstance.Close()

	// Start the chain reader and indexer
	err = chainReader.Start(ctx)
	if err != nil {
		lggr.Errorw("Failed to start SUI chain reader for execution events", "error", err)
		return
	}

	err = indexerInstance.Start(ctx)
	if err != nil {
		lggr.Errorw("Failed to start SUI indexer for execution events", "error", err)
		return
	}

	// Bind to the offramp contract
	err = chainReader.Bind(context.Background(), []chain_reader_types.BoundContract{{
		Name:    "offramp",
		Address: chainState.OffRampAddress,
	}})
	if err != nil {
		lggr.Errorw("Failed to bind offramp contract for execution events", "error", err)
		return
	}

	ticker := time.NewTicker(1 * time.Second) // Use 1 second for consistency with transmit events
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for SUI execution event",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange,
				"seenMessages", seenMessages,
				"completedSrcChains", completedSrcChains)
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 || finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			// Get the current latest checkpoint
			currentCheckpoint, err := suiChain.Client.SuiGetLatestCheckpointSequenceNumber(ctx)
			if err != nil {
				lggr.Errorw("Failed to get current SUI checkpoint for execution events", "error", err)
				continue
			}

			// Only query if there are new checkpoints since our last query
			if currentCheckpoint <= lastCheckpoint {
				continue
			}

			// Query for new execution events
			var executionEvent SuiExecutionStateChanged
			sequences, err := chainReader.QueryKey(
				ctx,
				chain_reader_types.BoundContract{
					Name:    "offramp",
					Address: chainState.OffRampAddress,
				},
				sui_query.KeyFilter{
					Key: "ExecutionStateChanged",
				},
				sui_query.LimitAndSort{
					Limit: sui_query.Limit{
						Count:  100,
						Cursor: "",
					},
				},
				&executionEvent,
			)
			if err != nil {
				lggr.Errorw("Failed to query SUI execution events", "error", err)
				continue
			}

			newExecutionsFound := false
			for _, sequence := range sequences {
				event, ok := sequence.Data.(*SuiExecutionStateChanged)
				if !ok {
					lggr.Errorw("Failed to cast SUI execution event data")
					continue
				}

				// Create a unique key for this execution event
				executionKey := fmt.Sprintf("%d_%d_%d", event.SourceChainSelector, chainSelector, event.SequenceNumber)

				// Skip if we've already processed this execution
				if processedExecutions[executionKey] {
					continue
				}

				// Mark this execution as processed
				processedExecutions[executionKey] = true
				newExecutionsFound = true

				lggr.Debugw("received SUI execution event for",
					"sourceChain", event.SourceChainSelector,
					"destChain", chainSelector,
					"sequenceNumber", event.SequenceNumber)

				// Push metrics
				data := messageData{
					eventType: executed,
					srcDstSeqNum: srcDstSeqNum{
						src:    event.SourceChainSelector,
						dst:    chainSelector,
						seqNum: event.SequenceNumber,
					},
					timestamp: uint64(time.Now().Unix()),
				}
				metricPipe <- data
				seenMessages[event.SourceChainSelector] = append(seenMessages[event.SourceChainSelector], event.SequenceNumber)
			}

			// Update last processed checkpoint
			lastCheckpoint = currentCheckpoint

			if newExecutionsFound {
				lggr.Debugw("Processed new SUI execution events",
					"destChain", chainSelector,
					"currentCheckpoint", currentCheckpoint,
					"startSlot", startSlot)
			}

			lggr.Infow("ticking, checking SUI executed events",
				"destChain", chainSelector,
				"seenMessages", seenMessages,
				"expectedRange", expectedRange,
				"completedSrcChains", completedSrcChains)

			for srcChain, seqNumRange := range expectedRange {
				if !completedSrcChains[srcChain] {
					if len(seenMessages[srcChain]) >= seqNumRange.Length() && slices.Contains(seenMessages[srcChain], uint64(seqNumRange.End())) {
						completedSrcChains[srcChain] = true
						lggr.Infow("executed all sequence numbers for ",
							"destChain", chainSelector,
							"sourceChain", srcChain,
							"seqNumRange", seqNumRange)
					}
				}
			}

			// Check if all chains are complete
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

// createSuiChainReader creates a SUI ChainReader instance for event subscriptions
func createSuiChainReader(ctx context.Context, t *testing.T, lggr logger.Logger, suiChain cldf_sui.Chain, chainState suiState.CCIPChainState) (chain_reader_types.ContractReader, *indexer.Indexer, error) {
	chainReaderConfig := crConfig.ChainReaderConfig{
		IsLoopPlugin: false,
		Modules: map[string]*crConfig.ChainReaderModule{
			"onramp": {
				Name: "onramp",
				Events: map[string]*crConfig.ChainReaderEvent{
					"CCIPMessageSent": {
						Name:      "CCIPMessageSent",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventSelector{
							Package: chainState.OnRampAddress,
							Module:  "onramp",
							Event:   "CCIPMessageSent",
						},
					},
				},
			},
			"offramp": {
				Name: "offramp",
				Events: map[string]*crConfig.ChainReaderEvent{
					"CommitReportAccepted": {
						Name:      "CommitReportAccepted",
						EventType: "CommitReportAccepted",
						EventSelector: client.EventSelector{
							Package: chainState.OffRampAddress,
							Module:  "offramp",
							Event:   "CommitReportAccepted",
						},
					},
					"ExecutionStateChanged": {
						Name:      "ExecutionStateChanged",
						EventType: "ExecutionStateChanged",
						EventSelector: client.EventSelector{
							Package: chainState.OffRampAddress,
							Module:  "offramp",
							Event:   "ExecutionStateChanged",
						},
					},
				},
			},
		},
	}

	dbURL := os.Getenv("CL_DATABASE_URL")
	if dbURL == "" {
		return nil, nil, fmt.Errorf("CL_DATABASE_URL environment variable is required")
	}

	err := sqltest.RegisterTxDB(dbURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to register database: %w", err)
	}

	db, err := sqlx.Open(pg.DriverTxWrappedPostgres, uuid.New().String())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.MapperFunc(reflectx.CamelToSnakeASCII)

	_, err = db.Connx(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	keystoreInstance := suitestutils.NewTestKeystore(t)
	priv, err := cldf_sui.PrivateKey(suiChain.Signer)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get private key: %w", err)
	}
	keystoreInstance.AddKey(priv)

	relayerClient, err := client.NewPTBClient(lggr, suiChain.URL, nil, 30*time.Second, keystoreInstance, 5, "WaitForEffectsCert")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create PTB client: %w", err)
	}

	// Create the indexers
	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		lggr,
		relayerClient,
		10*time.Second,
		10*time.Second,
		map[string]*crConfig.ChainReaderEvent{},
	)
	evIndexer := indexer.NewEventIndexer(
		db,
		lggr,
		relayerClient,
		[]*client.EventSelector{},
		10*time.Second,
		10*time.Second,
	)
	indexerInstance := indexer.NewIndexer(
		lggr,
		evIndexer,
		txnIndexer,
	)

	chainReader, err := chainreader.NewChainReader(ctx, lggr, relayerClient, chainReaderConfig, db, indexerInstance)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create chain reader: %w", err)
	}

	return chainReader, indexerInstance, nil
}

// FundSuiAccount transfers SUI tokens from a signer account to a destination account
func FundSuiAccount(t *testing.T, suiChain cldf_sui.Chain, toAddress string, amount uint64) error {
	ctx := testhelpers.Context(t)

	signerAddress, err := suiChain.Signer.GetAddress()
	if err != nil {
		return fmt.Errorf("failed to get signer address: %w", err)
	}

	coinObjects, err := suiChain.Client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
		Owner:    signerAddress,
		CoinType: "0x2::sui::SUI",
		Limit:    1,
	})
	if err != nil {
		return fmt.Errorf("failed to get coins for signer: %w", err)
	}

	if len(coinObjects.Data) == 0 {
		return fmt.Errorf("no SUI coins available for transfer from signer %s", signerAddress)
	}

	coinObjectId := coinObjects.Data[0].CoinObjectId

	transferReq := models.TransferSuiRequest{
		Signer:      signerAddress,
		SuiObjectId: coinObjectId,
		GasBudget:   "10000000", // 0.01 SUI for gas budget
		Recipient:   toAddress,
		Amount:      fmt.Sprintf("%d", amount),
	}

	txnMetadata, err := suiChain.Client.TransferSui(ctx, transferReq)
	if err != nil {
		return fmt.Errorf("failed to create SUI transfer transaction: %w", err)
	}

	signAndExecuteReq := models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: txnMetadata,
		Options: models.SuiTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForEffectsCert", // Wait for transaction to be finalized
	}

	response, err := suiChain.Client.SignAndExecuteTransactionBlock(ctx, signAndExecuteReq)
	if err != nil {
		return fmt.Errorf("failed to sign and execute SUI transfer transaction: %w", err)
	}

	if response.Effects.Status.Status != "success" {
		return fmt.Errorf("SUI transfer transaction failed with status: %s", response.Effects.Status.Status)
	}

	t.Logf("Funded SUI account %s from %s with %d MIST (SUI smallest unit)",
		toAddress, signerAddress, amount)

	return nil
}
