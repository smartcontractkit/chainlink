package ccip

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	sui "github.com/block-vision/sui-go-sdk/sui"
	suitx "github.com/block-vision/sui-go-sdk/transaction"
	"github.com/btcsuite/btcutil/bech32"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink-sui/relayer/testutils"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
)

func hexFromSuiBech32PrivKey(bech string) (string, error) {
	hrp, data5, err := bech32.Decode(bech)
	if err != nil {
		return "", err
	}
	if hrp != "suiprivkey" {
		return "", errors.New("unexpected HRP: " + hrp)
	}
	dataBytes, err := bech32.ConvertBits(data5, 5, 8, false)
	if err != nil {
		return "", err
	}
	if len(dataBytes) != 33 {
		return "", fmt.Errorf("decoded privkey wrong length: %d bytes", len(dataBytes))
	}
	seed := dataBytes[1:]
	if len(seed) != 32 {
		return "", fmt.Errorf("unexpected seed length: %d", len(seed))
	}
	return hex.EncodeToString(seed), nil
}

func subscribeSuiTransmitEvents(
	ctx context.Context,
	lggr logger.Logger,
	chainReader pkgtypes.ContractReader,
	onRampAddress string,
	otherChains []uint64,
	srcChainSel uint64,
	loadFinished chan struct{},
	wg *sync.WaitGroup,
	metricPipe chan messageData,
	finalSeqNrCommitChannels map[uint64]chan finalSeqNrReport,
	finalSeqNrExecChannels map[uint64]chan finalSeqNrReport,
) {
	defer wg.Done()
	lggr.Infow("starting sui chain transmit event subscriber for ",
		"srcChain", srcChainSel,
		"otherChains", otherChains,
	)

	boundContracts := []pkgtypes.BoundContract{
		{
			Name:    "OnRamp",
			Address: onRampAddress,
		},
	}
	err := chainReader.Bind(ctx, boundContracts)
	if err != nil {
		lggr.Errorw("failed to bind OnRamp contract", "error", err)
		return
	}

	seqNums := make(map[testhelpers.SourceDestPair]SeqNumRange)
	for _, cs := range otherChains {
		seqNums[testhelpers.SourceDestPair{
			SourceChainSelector: srcChainSel,
			DestChainSelector:   cs,
		}] = SeqNumRange{
			Start: atomic.NewUint64(math.MaxUint64),
			End:   atomic.NewUint64(0),
		}
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			lggr.Errorw("received context cancel signal for transmit watcher",
				"srcChain", srcChainSel)
			return

		case <-loadFinished:
			for _, destChain := range otherChains {
				commitChan := finalSeqNrCommitChannels[destChain]
				execChan := finalSeqNrExecChannels[destChain]

				csPair := testhelpers.SourceDestPair{
					SourceChainSelector: srcChainSel,
					DestChainSelector:   destChain,
				}

				report := finalSeqNrReport{
					sourceChainSelector: srcChainSel,
					expectedSeqNrRange: ccipocr3.SeqNumRange{
						ccipocr3.SeqNum(seqNums[csPair].Start.Load()),
						ccipocr3.SeqNum(seqNums[csPair].End.Load()),
					},
				}

				commitChan <- report
				execChan <- report
			}
			return

		case <-ticker.C:

			type CCIPMessageSentEvent struct {
				DestChainSelector uint64 `json:"destChainSelector"`
				SequenceNumber    uint64 `json:"sequenceNumber"`
				Message           any    `json:"message"`
			}

			boundContract := pkgtypes.BoundContract{
				Name:    "OnRamp",
				Address: onRampAddress,
			}

			filter := query.KeyFilter{
				Key: "CCIPMessageSent",
			}

			limitAndSort := query.LimitAndSort{
				Limit: query.Limit{
					Count:  10000,
					Cursor: "",
				},
			}

			var event CCIPMessageSentEvent
			sequences, err := chainReader.QueryKey(ctx, boundContract, filter, limitAndSort, &event)
			if err != nil {
				lggr.Debugw("error querying transmit events", "error", err)
				continue
			}

			for _, seq := range sequences {
				event := seq.Data.(*CCIPMessageSentEvent)

				destChain := event.DestChainSelector
				seqNum := event.SequenceNumber

				if destChain == 0 || seqNum == 0 {
					lggr.Debugw("skipping invalid event",
						"srcChain", srcChainSel,
						"destChain", destChain,
						"sequenceNumber", seqNum)
					continue
				}

				csPair := testhelpers.SourceDestPair{
					SourceChainSelector: srcChainSel,
					DestChainSelector:   destChain,
				}

				if _, exists := seqNums[csPair]; !exists {
					seqNums[csPair] = SeqNumRange{
						Start: atomic.NewUint64(math.MaxUint64),
						End:   atomic.NewUint64(0),
					}
				}

				isNew := seqNum < seqNums[csPair].Start.Load() || seqNum > seqNums[csPair].End.Load()

				if seqNum < seqNums[csPair].Start.Load() {
					seqNums[csPair].Start.Store(seqNum)
				}

				if seqNum > seqNums[csPair].End.Load() {
					seqNums[csPair].End.Store(seqNum)
				}

				if isNew {
					lggr.Debugw("received sui transmit event for",
						"srcChain", srcChainSel,
						"destChain", destChain,
						"sequenceNumber", seqNum)

					data := messageData{
						eventType: transmitted,
						srcDstSeqNum: srcDstSeqNum{
							src:    srcChainSel,
							dst:    destChain,
							seqNum: seqNum,
						},
						timestamp: uint64(time.Now().Unix()),
					}
					metricPipe <- data
				}
			}
		}
	}
}

func subscribeSuiCommitEvents(
	ctx context.Context,
	lggr logger.Logger,
	chainReader pkgtypes.ContractReader,
	offRampAddress string,
	srcChains []uint64,
	chainSelector uint64,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting sui commit event subscriber for ",
		"destChain", chainSelector,
	)

	boundContracts := []pkgtypes.BoundContract{
		{
			Name:    "OffRamp",
			Address: offRampAddress,
		},
	}
	err := chainReader.Bind(ctx, boundContracts)
	if err != nil {
		lggr.Errorw("failed to bind OffRamp contract", "error", err)
		return
	}

	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	checkTicker := time.NewTicker(tickerDuration)
	defer checkTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for commit report",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 || finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			type MerkleRoot struct {
				SourceChainSelector uint64 `json:"sourceChainSelector"`
				OnRampAddress       []byte `json:"onRampAddress"`
				MinSeqNr            uint64 `json:"minSeqNr"`
				MaxSeqNr            uint64 `json:"maxSeqNr"`
				MerkleRoot          []byte `json:"merkleRoot"`
			}

			type PriceUpdates struct {
				TokenPriceUpdates []any `json:"tokenPriceUpdates"`
				GasPriceUpdates   []any `json:"gasPriceUpdates"`
			}

			type CommitReportAcceptedEvent struct {
				BlessedMerkleRoots   []MerkleRoot `json:"blessedMerkleRoots"`
				UnblessedMerkleRoots []MerkleRoot `json:"unblessedMerkleRoots"`
				PriceUpdates         PriceUpdates `json:"priceUpdates"`
			}

			boundContract := pkgtypes.BoundContract{
				Name:    "OffRamp",
				Address: offRampAddress,
			}

			filter := query.KeyFilter{
				Key: "CommitReportAccepted",
			}

			limitAndSort := query.LimitAndSort{
				Limit: query.Limit{
					Count:  10000,
					Cursor: "",
				},
			}

			var event CommitReportAcceptedEvent
			sequences, err := chainReader.QueryKey(ctx, boundContract, filter, limitAndSort, &event)
			if err != nil {
				lggr.Debugw("error querying commit events", "error", err)
				continue
			}

			for _, seq := range sequences {
				event := seq.Data.(*CommitReportAcceptedEvent)

				allRoots := append(event.BlessedMerkleRoots, event.UnblessedMerkleRoots...)

				if len(allRoots) == 0 {
					lggr.Debugw("skipping empty commit event", "destChain", chainSelector)
					continue
				}

				for _, mr := range allRoots {
					if mr.SourceChainSelector == 0 || mr.MinSeqNr == 0 {
						lggr.Debugw("skipping invalid merkle root",
							"srcChain", mr.SourceChainSelector,
							"minSeqNr", mr.MinSeqNr,
							"maxSeqNr", mr.MaxSeqNr)
						continue
					}

					lggr.Infow("received sui commit report",
						"srcChain", mr.SourceChainSelector,
						"destChain", chainSelector,
						"minSeqNr", mr.MinSeqNr,
						"maxSeqNr", mr.MaxSeqNr)

					if _, ok := expectedRange[mr.SourceChainSelector]; !ok {
						lggr.Debugw("received sui commit report (expectedRange not yet populated)",
							"srcChain", mr.SourceChainSelector,
							"destChain", chainSelector,
							"minSeqNr", mr.MinSeqNr,
							"maxSeqNr", mr.MaxSeqNr)
					}

					for seqNum := mr.MinSeqNr; seqNum <= mr.MaxSeqNr; seqNum++ {
						if !contains(seenMessages[mr.SourceChainSelector], seqNum) {
							seenMessages[mr.SourceChainSelector] = append(seenMessages[mr.SourceChainSelector], seqNum)

							data := messageData{
								eventType: committed,
								srcDstSeqNum: srcDstSeqNum{
									src:    mr.SourceChainSelector,
									dst:    chainSelector,
									seqNum: seqNum,
								},
								timestamp: uint64(time.Now().Unix()),
							}
							metricPipe <- data
						}
					}
				}
			}

		case <-checkTicker.C:
			allComplete := true
			for srcChain, expectedSeqNrRange := range expectedRange {
				if !completedSrcChains[srcChain] {
					complete := true
					for seqNum := uint64(expectedSeqNrRange.Start()); seqNum <= uint64(expectedSeqNrRange.End()); seqNum++ {
						if !contains(seenMessages[srcChain], seqNum) {
							complete = false
							break
						}
					}
					if complete {
						completedSrcChains[srcChain] = true
						lggr.Infow("all messages committed",
							"srcChain", srcChain,
							"destChain", chainSelector,
							"expectedRange", expectedSeqNrRange)
					} else {
						allComplete = false
					}
				}
			}

			if allComplete && len(expectedRange) > 0 {
				lggr.Infow("all source chains completed commit",
					"destChain", chainSelector)
				return
			}
		}
	}
}

func subscribeSuiExecutionEvents(
	ctx context.Context,
	lggr logger.Logger,
	chainReader pkgtypes.ContractReader,
	offRampAddress string,
	srcChains []uint64,
	chainSelector uint64,
	finalSeqNrs chan finalSeqNrReport,
	wg *sync.WaitGroup,
	metricPipe chan messageData,
) {
	defer wg.Done()
	defer close(finalSeqNrs)

	lggr.Infow("starting sui chain execution event subscriber for ",
		"destChain", chainSelector,
	)

	boundContracts := []pkgtypes.BoundContract{
		{
			Name:    "OffRamp",
			Address: offRampAddress,
		},
	}
	err := chainReader.Bind(ctx, boundContracts)
	if err != nil {
		lggr.Errorw("failed to bind OffRamp contract", "error", err)
		return
	}

	seenMessages := make(map[uint64][]uint64)
	expectedRange := make(map[uint64]ccipocr3.SeqNumRange)
	completedSrcChains := make(map[uint64]bool)
	for _, srcChain := range srcChains {
		seenMessages[srcChain] = make([]uint64, 0)
		completedSrcChains[srcChain] = false
	}

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	checkTicker := time.NewTicker(tickerDuration)
	defer checkTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			lggr.Errorw("timed out waiting for execution",
				"destChain", chainSelector,
				"sourceChains", srcChains,
				"expectedSeqNumbers", expectedRange)
			return

		case finalSeqNrUpdate := <-finalSeqNrs:
			if finalSeqNrUpdate.expectedSeqNrRange.Start() == math.MaxUint64 || finalSeqNrUpdate.expectedSeqNrRange.End() == 0 {
				delete(completedSrcChains, finalSeqNrUpdate.sourceChainSelector)
				delete(seenMessages, finalSeqNrUpdate.sourceChainSelector)
			} else {
				expectedRange[finalSeqNrUpdate.sourceChainSelector] = finalSeqNrUpdate.expectedSeqNrRange
			}

		case <-ticker.C:
			type ExecutionStateChangedEvent struct {
				SourceChainSelector uint64 `json:"sourceChainSelector"`
				SequenceNumber      uint64 `json:"sequenceNumber"`
				MessageId           []byte `json:"messageId"`
				MessageHash         []byte `json:"messageHash"`
				State               byte   `json:"state"`
			}

			boundContract := pkgtypes.BoundContract{
				Name:    "OffRamp",
				Address: offRampAddress,
			}

			filter := query.KeyFilter{
				Key: "ExecutionStateChanged",
			}

			limitAndSort := query.LimitAndSort{
				Limit: query.Limit{
					Count:  10000,
					Cursor: "",
				},
			}

			var event ExecutionStateChangedEvent
			sequences, err := chainReader.QueryKey(ctx, boundContract, filter, limitAndSort, &event)
			if err != nil {
				lggr.Debugw("error querying execution events", "error", err)
				continue
			}

			for _, seq := range sequences {
				event := seq.Data.(*ExecutionStateChangedEvent)

				if event.SourceChainSelector == 0 || event.SequenceNumber == 0 {
					lggr.Debugw("skipping invalid execution event",
						"srcChain", event.SourceChainSelector,
						"seqNum", event.SequenceNumber)
					continue
				}

				lggr.Infow("received execution state changed",
					"srcChain", event.SourceChainSelector,
					"destChain", chainSelector,
					"seqNum", event.SequenceNumber,
					"messageId", event.MessageId,
					"state", event.State)

				if !contains(seenMessages[event.SourceChainSelector], event.SequenceNumber) {
					seenMessages[event.SourceChainSelector] = append(seenMessages[event.SourceChainSelector], event.SequenceNumber)

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
				}
			}

		case <-checkTicker.C:
			allComplete := true
			for srcChain, expectedSeqNrRange := range expectedRange {
				if !completedSrcChains[srcChain] {
					complete := true
					for seqNum := uint64(expectedSeqNrRange.Start()); seqNum <= uint64(expectedSeqNrRange.End()); seqNum++ {
						if !contains(seenMessages[srcChain], seqNum) {
							complete = false
							break
						}
					}
					if complete {
						completedSrcChains[srcChain] = true
						lggr.Infow("all messages executed",
							"srcChain", srcChain,
							"destChain", chainSelector,
							"expectedRange", expectedSeqNrRange)
					} else {
						allComplete = false
					}
				}
			}

			if allComplete && len(expectedRange) > 0 {
				lggr.Infow("all source chains completed execution",
					"destChain", chainSelector)
				return
			}
		}
	}
}

func contains(slice []uint64, val uint64) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func NewChainReaderFromLatestBlock(
	ctx context.Context,
	lgr logger.Logger,
	ptbClient *client.PTBClient,
	chainReaderConfig config.ChainReaderConfig,
	db sqlutil.DataSource,
) (pkgtypes.ContractReader, error) {
	dbStore := database.NewDBStore(db, lgr)

	err := dbStore.EnsureSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure database schema: %w", err)
	}

	lgr.Info("Setting cursors to latest block...")
	err = setAllEventCursorsToLatest(ctx, lgr, ptbClient, dbStore, chainReaderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to set event cursors: %w", err)
	}

	txnIndexer := indexer.NewTransactionsIndexer(
		db,
		lgr,
		ptbClient,
		chainReaderConfig.TransactionsIndexer.PollingInterval,
		chainReaderConfig.TransactionsIndexer.SyncTimeout,
		map[string]*config.ChainReaderEvent{},
	)

	eventsIndexer := indexer.NewEventIndexer(
		db,
		lgr,
		ptbClient,
		[]*client.EventSelector{},
		chainReaderConfig.EventsIndexer.PollingInterval,
		chainReaderConfig.EventsIndexer.SyncTimeout,
	)

	indexerInstance := indexer.NewIndexer(lgr, eventsIndexer, txnIndexer)

	chainReader, err := reader.NewChainReader(ctx, lgr, ptbClient, chainReaderConfig, db, indexerInstance)
	if err != nil {
		return nil, fmt.Errorf("failed to create chain reader: %w", err)
	}

	lgr.Info("Chain reader created with cursors at latest block")
	return chainReader, nil
}

func setAllEventCursorsToLatest(
	ctx context.Context,
	lgr logger.Logger,
	ptbClient *client.PTBClient,
	dbStore *database.DBStore,
	chainReaderConfig config.ChainReaderConfig,
) error {
	for _, moduleConfig := range chainReaderConfig.Modules {
		if moduleConfig.Events == nil {
			continue
		}

		for _, eventConfig := range moduleConfig.Events {
			selector := client.EventSelector{
				Package: eventConfig.Package,
				Module:  eventConfig.Module,
				Event:   eventConfig.EventType,
			}

			if selector.Package == "" {
				continue
			}

			err := setEventCursorToLatest(ctx, lgr, ptbClient, dbStore, selector)
			if err != nil {
				lgr.Warnw("Failed to set cursor", "error", err)
			}
		}
	}

	return nil
}

func setEventCursorToLatest(
	ctx context.Context,
	lgr logger.Logger,
	ptbClient *client.PTBClient,
	dbStore *database.DBStore,
	selector client.EventSelector,
) error {
	limit := uint(1)
	queryCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	events, err := ptbClient.QueryEvents(queryCtx, selector, &limit, nil, &client.QuerySortOptions{Descending: true})
	if err != nil {
		return fmt.Errorf("failed to query latest event: %w", err)
	}

	if len(events.Data) == 0 {
		return nil
	}

	latestEvent := events.Data[0]
	eventHandle := fmt.Sprintf("%s::%s::%s", selector.Package, selector.Module, selector.Event)

	block, err := ptbClient.BlockByDigest(queryCtx, latestEvent.Id.TxDigest)
	if err != nil {
		return fmt.Errorf("failed to get block: %w", err)
	}

	markerEvent := database.EventRecord{
		EventAccountAddress: selector.Package,
		EventHandle:         eventHandle,
		EventOffset:         0,
		TxDigest:            latestEvent.Id.TxDigest,
		BlockVersion:        0,
		BlockHeight:         fmt.Sprintf("%d", block.Height),
		BlockHash:           []byte(block.TxDigest),
		BlockTimestamp:      block.Timestamp,
		Data:                map[string]any{"_marker": "latest_cursor"},
	}

	err = dbStore.InsertEvents(queryCtx, []database.EventRecord{markerEvent})
	if err != nil {
		return fmt.Errorf("failed to insert marker: %w", err)
	}

	lgr.Infow("Cursor set to latest", "eventHandle", eventHandle, "checkpoint", block.Height)
	return nil
}

func GetEVMExtraArgsV2SUI(receiverStateObjectId string) ([]byte, error) {
	// Tag prefix
	SUITag := hexutil.MustDecode("0x21ea4ca9")

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	fmt.Printf("Receiver state object id: %s\n", receiverStateObjectId)
	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		receiverStateObjectId,
	))

	recieverObjectIds := [][32]byte{clockObj, stateObj}

	suiExtraArgsData := message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 big.NewInt(1000000),
		AllowOutOfOrderExecution: true,
		TokenReceiver:            [32]byte{}, // EOA
		ReceiverObjectIds:        recieverObjectIds,
	}

	return ccipevm.SerializeExtraArgs(SUITag, "encodeSUIExtraArgsV1", suiExtraArgsData)

}

// SuiTokenPool manages a pool of Sui token objects for load testing
// It cycles through available token objects in a thread-safe manner
type SuiTokenPool struct {
	mu             sync.Mutex
	tokenObjectIds []string
	currentIndex   int
	lggr           logger.Logger
}

// NewSuiTokenPool creates a new token pool
func NewSuiTokenPool(lggr logger.Logger, tokenObjectIds []string) *SuiTokenPool {
	return &SuiTokenPool{
		tokenObjectIds: tokenObjectIds,
		currentIndex:   0,
		lggr:           lggr,
	}
}

// GetNextToken returns the next token object ID from the pool
// It cycles through the available tokens in a round-robin fashion
func (p *SuiTokenPool) GetNextToken() (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.tokenObjectIds) == 0 {
		return "", fmt.Errorf("token pool is empty")
	}

	token := p.tokenObjectIds[p.currentIndex]
	p.currentIndex = (p.currentIndex + 1) % len(p.tokenObjectIds)

	return token, nil
}

// Size returns the number of tokens in the pool
func (p *SuiTokenPool) Size() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.tokenObjectIds)
}

// splitSuiTokens splits existing Link tokens owned by the deployer into multiple small objects
// Returns a list of token object IDs that can be used for load testing
func splitSuiTokens(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	env cldf.Environment,
	state *stateview.CCIPOnChainState,
	chainSelector uint64,
	numTokenObjects int,
	amountPerObject uint64, // amount in smallest unit (e.g., 1e4 for 0.00001 Link)
	privateKeyHex string, // Sui private key in hex format (without 0x prefix)
	tokenToTransferPkdID string,
) ([]string, error) {
	suiChain := env.BlockChains.SuiChains()[chainSelector]

	// Get Link token package ID from state
	linkTokenPkgID := tokenToTransferPkdID
	if linkTokenPkgID == "" {
		return nil, fmt.Errorf("link token not configured for chain %d", chainSelector)
	}

	// Get deployer address from signer
	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get deployer address: %w", err)
	}
	// Add 0x prefix if not present
	if !strings.HasPrefix(deployerAddr, "0x") {
		deployerAddr = "0x" + deployerAddr
	}

	lggr.Infow("Splitting Link tokens for load test",
		"chainSelector", chainSelector,
		"deployerAddr", deployerAddr,
		"numObjects", numTokenObjects,
		"amountPerObject", amountPerObject,
		"linkTokenPkg", linkTokenPkgID)

	// Create Sui client
	client := sui.NewSuiClient(suiChain.URL)

	// Query for owned Link token objects
	coinType := fmt.Sprintf("%s::link::LINK", linkTokenPkgID)
	ownedCoins, err := client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
		Owner:    deployerAddr,
		CoinType: coinType,
		Limit:    10,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query owned Link tokens: %w", err)
	}

	if len(ownedCoins.Data) == 0 {
		return nil, fmt.Errorf("deployer account has no Link tokens to split")
	}

	lggr.Infow("Found Link token objects owned by deployer",
		"count", len(ownedCoins.Data),
		"firstCoinID", ownedCoins.Data[0].CoinObjectId,
		"firstCoinBalance", ownedCoins.Data[0].Balance)

	// Use the first coin object to split from
	sourceCoinID := ownedCoins.Data[0].CoinObjectId
	sourceBalanceStr := ownedCoins.Data[0].Balance
	sourceBalance, ok := new(big.Int).SetString(sourceBalanceStr, 10)
	if !ok {
		return nil, fmt.Errorf("failed to parse source coin balance: %s", sourceBalanceStr)
	}

	// Calculate total amount needed
	totalNeeded := new(big.Int).Mul(big.NewInt(int64(numTokenObjects)), big.NewInt(int64(amountPerObject)))
	if sourceBalance.Cmp(totalNeeded) < 0 {
		return nil, fmt.Errorf("insufficient balance: have %s, need %s", sourceBalance.String(), totalNeeded.String())
	}

	lggr.Infow("Splitting source coin into smaller objects",
		"sourceCoinID", sourceCoinID,
		"sourceBalance", sourceBalance.String(),
		"totalNeeded", totalNeeded.String())

	// Create SDK signer using the provided private key (32-byte seed in hex format)
	// Decode the hex string to get the raw seed bytes
	seedBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %w", err)
	}
	if len(seedBytes) != 32 {
		return nil, fmt.Errorf("invalid seed length: expected 32 bytes, got %d", len(seedBytes))
	}

	suiSDKSigner := signer.NewSigner(seedBytes)

	// Query for SUI gas coins owned by the deployer
	gasCoinType := "0x2::sui::SUI"
	gasCoins, err := client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
		Owner:    deployerAddr,
		CoinType: gasCoinType,
		Limit:    5,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query gas coins: %w", err)
	}
	if len(gasCoins.Data) == 0 {
		return nil, fmt.Errorf("deployer account has no SUI for gas")
	}

	lggr.Infow("Found gas coins for transaction",
		"count", len(gasCoins.Data),
		"firstGasCoinID", gasCoins.Data[0].CoinObjectId,
		"firstGasCoinBalance", gasCoins.Data[0].Balance)

	// Build PTB to split the coin into multiple objects
	tokenObjectIds := make([]string, 0, numTokenObjects)

	// Split in batches to avoid gas limits
	batchSize := 10
	for batch := 0; batch < numTokenObjects; batch += batchSize {
		remaining := numTokenObjects - batch
		if remaining > batchSize {
			remaining = batchSize
		}

		// Build transaction
		tx := suitx.NewTransaction()
		tx.SetSigner(suiSDKSigner)
		tx.SetSuiClient(client.(*sui.Client))
		tx.SetSender(models.SuiAddress(deployerAddr))

		// Set gas payment - use the first gas coin
		gasObjRef, err := suitx.NewSuiObjectRef(
			models.SuiAddress(gasCoins.Data[0].CoinObjectId),
			gasCoins.Data[0].Version,
			models.ObjectDigest(gasCoins.Data[0].Digest),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create gas object ref: %w", err)
		}
		tx.SetGasPayment([]suitx.SuiObjectRef{*gasObjRef})
		tx.SetGasBudget(100000000) // 0.1 SUI gas budget

		// Reference the source coin as an owned object
		// Need to create proper object reference with version and digest
		sourceCoinRef, err := suitx.NewSuiObjectRef(
			models.SuiAddress(sourceCoinID),
			ownedCoins.Data[0].Version,
			models.ObjectDigest(ownedCoins.Data[0].Digest),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create source coin ref: %w", err)
		}

		// Create CallArg for the owned coin
		coinCallArg := suitx.CallArg{
			Object: &suitx.ObjectArg{
				ImmOrOwnedObject: sourceCoinRef,
			},
		}
		coinArg := tx.Object(coinCallArg)

		// Split the coins one at a time to get individual coin objects
		// This ensures each split creates a distinct coin object we can track
		splitCoins := make([]suitx.Argument, 0, remaining)
		for i := 0; i < remaining; i++ {
			// Split one coin at a time
			splitCoin := tx.SplitCoins(coinArg, []suitx.Argument{tx.Pure(amountPerObject)})
			splitCoins = append(splitCoins, splitCoin)
		}

		// Transfer all split coins to the deployer (finalizes their creation as distinct objects)
		if len(splitCoins) > 0 {
			tx.TransferObjects(splitCoins, tx.Pure(deployerAddr))
		}

		// Execute the transaction
		resp, err := tx.Execute(ctx, models.SuiTransactionBlockOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		}, "WaitForLocalExecution")
		if err != nil {
			return nil, fmt.Errorf("failed to execute split transaction batch %d: %w", batch, err)
		}

		// After successful execution, update object references for next batch
		// Both gas coin and source coin versions change after each transaction
		for _, change := range resp.ObjectChanges {
			if change.Type == "mutated" {
				if change.ObjectId == gasCoins.Data[0].CoinObjectId {
					// Update gas coin version and digest for next batch
					gasCoins.Data[0].Version = change.Version
					gasCoins.Data[0].Digest = change.Digest
					lggr.Debugw("Updated gas coin reference",
						"newVersion", change.Version,
						"newDigest", change.Digest)
				} else if change.ObjectId == sourceCoinID {
					// Update source coin version and digest for next batch
					ownedCoins.Data[0].Version = change.Version
					ownedCoins.Data[0].Digest = change.Digest
					lggr.Debugw("Updated source coin reference",
						"newVersion", change.Version,
						"newDigest", change.Digest)
				}
			}
		}

		// Extract new coin object IDs from object changes
		if len(resp.ObjectChanges) > 0 {
			for _, change := range resp.ObjectChanges {
				if change.Type == "created" && change.ObjectType != "" {
					// Check if it's a Link coin
					if strings.Contains(change.ObjectType, "::link::LINK") {
						tokenObjectIds = append(tokenObjectIds, change.ObjectId)
						lggr.Debugw("Split token object created",
							"objectId", change.ObjectId,
							"batch", batch)
					}
				}
			}
		}

		lggr.Infow("Completed split batch",
			"batch", batch,
			"totalCreated", len(tokenObjectIds))
	}

	lggr.Infow("Successfully split token objects for load test",
		"chainSelector", chainSelector,
		"numObjects", len(tokenObjectIds),
		"amountPerObject", amountPerObject)

	return tokenObjectIds, nil
}

// cleanupSuiDatabase cleans up Sui events and transactions tables before test runs
func cleanupSuiDatabase(dbUrl string, lggr logger.Logger) {
	db, err := sqlx.Open("postgres", dbUrl)
	if err != nil {
		lggr.Warnw("Failed to open database for cleanup", "error", err)
		return
	}
	defer db.Close()

	lggr.Infow("Cleaning up Sui events database for fresh test run")

	// Check if schema exists, create if needed
	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS sui")
	if err != nil {
		lggr.Warnw("Failed to create sui schema", "error", err)
	}

	// Truncate events table (table might not exist on first run)
	_, err = db.Exec("TRUNCATE TABLE sui.events CASCADE")
	if err != nil {
		lggr.Debugw("sui.events table doesn't exist or failed to truncate (this is ok on first run)", "error", err)
	} else {
		lggr.Infow("Truncated sui.events table")
	}

	// Truncate transactions table (table might not exist on first run)
	_, err = db.Exec("TRUNCATE TABLE sui.transactions CASCADE")
	if err != nil {
		lggr.Debugw("sui.transactions table doesn't exist or failed to truncate (this is ok on first run)", "error", err)
	} else {
		lggr.Infow("Truncated sui.transactions table")
	}

	lggr.Infow("Database cleanup complete")
}

// createSuiChainReader creates a chain reader for a Sui chain with configured OnRamp and OffRamp events
func createSuiChainReader(
	ctx context.Context,
	t *testing.T,
	lggr logger.Logger,
	env cldf.Environment,
	db *sqlx.DB,
	chainSelector uint64,
	onRampAddress string,
	offRampAddress string,
) (pkgtypes.ContractReader, error) {
	suiKeystore := testutils.NewTestKeystore(t)
	ptbClient, err := client.NewPTBClient(
		lggr,
		env.BlockChains.SuiChains()[chainSelector].URL,
		nil,
		10*time.Second,
		suiKeystore,
		5,
		"WaitForLocalExecution",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create PTB client: %w", err)
	}

	chainReaderConfig := config.ChainReaderConfig{
		Modules: map[string]*config.ChainReaderModule{
			"OnRamp": {
				Name: "onramp",
				Events: map[string]*config.ChainReaderEvent{
					"CCIPMessageSent": {
						Name:      "onramp",
						EventType: "CCIPMessageSent",
						EventSelector: client.EventSelector{
							Package: onRampAddress,
							Module:  "onramp",
							Event:   "CCIPMessageSent",
						},
					},
				},
			},
			"OffRamp": {
				Name: "offramp",
				Events: map[string]*config.ChainReaderEvent{
					"CommitReportAccepted": {
						Name:      "offramp",
						EventType: "CommitReportAccepted",
						EventSelector: client.EventSelector{
							Package: offRampAddress,
							Module:  "offramp",
							Event:   "CommitReportAccepted",
						},
					},
					"ExecutionStateChanged": {
						Name:      "offramp",
						EventType: "ExecutionStateChanged",
						EventSelector: client.EventSelector{
							Package: offRampAddress,
							Module:  "offramp",
							Event:   "ExecutionStateChanged",
						},
					},
				},
			},
		},
	}

	return NewChainReaderFromLatestBlock(ctx, lggr, ptbClient, chainReaderConfig, db)
}
