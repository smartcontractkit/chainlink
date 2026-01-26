package ccip

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
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
	_ "github.com/lib/pq" // PostgreSQL driver
	"go.uber.org/atomic"

	selectors "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_3/message_hasher"
	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"
	cldf_sui "github.com/smartcontractkit/chainlink-deployments-framework/chain/sui"
	cldf "github.com/smartcontractkit/chainlink-deployments-framework/deployment"
	sui_common "github.com/smartcontractkit/chainlink-sui/bindings/bind"

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
						"sequenceNumber", seqNum,
						"txHash", seq.Cursor)

					data := messageData{
						eventType: transmitted,
						srcDstSeqNum: srcDstSeqNum{
							src:    srcChainSel,
							dst:    destChain,
							seqNum: seqNum,
						},
						timestamp: uint64(time.Now().Unix()), //nolint:gosec // G115 - Unix time is always positive
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

				allRoots := event.BlessedMerkleRoots
				allRoots = append(allRoots, event.UnblessedMerkleRoots...)

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
								timestamp: uint64(time.Now().Unix()), //nolint:gosec // G115 - Unix time is always positive
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
				MessageID           []byte `json:"messageId"`
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
					"messageId", event.MessageID,
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
						timestamp: uint64(time.Now().Unix()), //nolint:gosec // G115 - Unix time is always positive
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

func GetEVMExtraArgsV2SUI(receiverStateObjectID string) ([]byte, error) {
	// Tag prefix
	SUITag := hexutil.MustDecode("0x21ea4ca9")

	var clockObj [32]byte
	copy(clockObj[:], hexutil.MustDecode(
		"0x0000000000000000000000000000000000000000000000000000000000000006",
	))

	fmt.Printf("Receiver state object id: %s\n", receiverStateObjectID)
	var stateObj [32]byte
	copy(stateObj[:], hexutil.MustDecode(
		receiverStateObjectID,
	))

	receiverObjectIDs := [][32]byte{clockObj, stateObj}

	suiExtraArgsData := message_hasher.ClientSuiExtraArgsV1{
		GasLimit:                 big.NewInt(1000000),
		AllowOutOfOrderExecution: true,
		TokenReceiver:            [32]byte{}, // EOA
		ReceiverObjectIds:        receiverObjectIDs,
	}

	return ccipevm.SerializeExtraArgs(SUITag, "encodeSUIExtraArgsV1", suiExtraArgsData)
}

// SuiCoinObject represents a Sui coin/object with its versioning info
type SuiCoinObject struct {
	ObjectID string
	Version  string
	Digest   string
}

// ConsolidatedCoin represents the result of consolidating multiple coins into one
type ConsolidatedCoin struct {
	ObjectID string
	Version  string
	Digest   string
	Balance  uint64
}

// consolidateSuiCoins fetches all coins of a given type owned by an address and merges them into one.
// This prevents accumulation of small coin objects from previous test runs.
// Returns the consolidated coin info, or nil if no coins exist.
func consolidateSuiCoins(
	ctx context.Context,
	lggr logger.Logger,
	suiClient sui.ISuiAPI,
	deployerAddr string,
	coinType string,
	privateKeyHex string,
	chainSelector uint64,
) (*ConsolidatedCoin, error) {
	lggr.Infow("Starting coin consolidation",
		"chainSelector", chainSelector,
		"coinType", coinType,
		"deployerAddr", deployerAddr)

	// Fetch all coins of this type with pagination
	allCoins := make([]models.CoinData, 0)
	var cursor interface{}
	for {
		coinsResp, err := suiClient.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    deployerAddr,
			CoinType: coinType,
			Cursor:   cursor,
			Limit:    50, // Max limit per page
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query coins: %w", err)
		}

		allCoins = append(allCoins, coinsResp.Data...)

		if !coinsResp.HasNextPage || coinsResp.NextCursor == "" {
			break
		}
		cursor = coinsResp.NextCursor
	}

	lggr.Infow("Found coins to consolidate",
		"chainSelector", chainSelector,
		"coinType", coinType,
		"count", len(allCoins))

	// If 0 or 1 coin, nothing to consolidate
	if len(allCoins) == 0 {
		return nil, nil
	}
	if len(allCoins) == 1 {
		bal, _ := strconv.ParseUint(allCoins[0].Balance, 10, 64)
		return &ConsolidatedCoin{
			ObjectID: allCoins[0].CoinObjectId,
			Version:  allCoins[0].Version,
			Digest:   allCoins[0].Digest,
			Balance:  bal,
		}, nil
	}

	// Find the largest coin to use as merge destination
	var destCoin models.CoinData
	var destBalance uint64
	destIdx := 0
	for i, coin := range allCoins {
		bal, _ := strconv.ParseUint(coin.Balance, 10, 64)
		if bal > destBalance {
			destBalance = bal
			destCoin = coin
			destIdx = i
		}
	}

	// Collect source coins (all except destination)
	sourceCoins := make([]models.CoinData, 0, len(allCoins)-1)
	for i, coin := range allCoins {
		if i != destIdx {
			sourceCoins = append(sourceCoins, coin)
		}
	}

	lggr.Infow("Merging coins into destination",
		"chainSelector", chainSelector,
		"destCoinID", destCoin.CoinObjectId,
		"destBalance", destBalance,
		"sourceCount", len(sourceCoins))

	// Create SDK signer
	seedBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %w", err)
	}
	if len(seedBytes) != 32 {
		return nil, fmt.Errorf("invalid seed length: expected 32 bytes, got %d", len(seedBytes))
	}
	suiSDKSigner := signer.NewSigner(seedBytes)

	// Track destination coin's current version/digest
	currentDestVersion := destCoin.Version
	currentDestDigest := destCoin.Digest
	totalBalance := destBalance

	// Merge in batches to avoid gas limits
	batchSize := 20
	for batch := 0; batch < len(sourceCoins); batch += batchSize {
		endIdx := batch + batchSize
		if endIdx > len(sourceCoins) {
			endIdx = len(sourceCoins)
		}
		batchCoins := sourceCoins[batch:endIdx]

		// Build transaction
		tx := suitx.NewTransaction()
		tx.SetSigner(suiSDKSigner)
		tx.SetSuiClient(suiClient.(*sui.Client))
		tx.SetSender(models.SuiAddress(deployerAddr))

		// Create destination coin reference
		destCoinRef, err := suitx.NewSuiObjectRef(
			models.SuiAddress(destCoin.CoinObjectId),
			currentDestVersion,
			models.ObjectDigest(currentDestDigest),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create dest coin ref: %w", err)
		}

		// For SUI coins, we need to use a separate gas coin or the dest coin as gas
		// Set gas payment to use the destination coin (it will be both gas and merge target)
		tx.SetGasPayment([]suitx.SuiObjectRef{*destCoinRef})
		tx.SetGasBudget(50_000_000) // 0.05 SUI gas budget per batch

		// Use tx.Gas() to reference the gas/destination coin for SUI merges
		// For non-SUI coins, we need to reference the destination coin as an object
		var destArg suitx.Argument
		if coinType == "0x2::sui::SUI" {
			destArg = tx.Gas()
		} else {
			// For non-SUI coins, we need a separate gas coin
			// First, get a SUI gas coin
			gasCoins, err := suiClient.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
				Owner:    deployerAddr,
				CoinType: "0x2::sui::SUI",
				Limit:    1,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to query gas coins: %w", err)
			}
			if len(gasCoins.Data) == 0 {
				return nil, errors.New("no SUI gas coins available")
			}

			gasCoinRef, err := suitx.NewSuiObjectRef(
				models.SuiAddress(gasCoins.Data[0].CoinObjectId),
				gasCoins.Data[0].Version,
				models.ObjectDigest(gasCoins.Data[0].Digest),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create gas coin ref: %w", err)
			}
			tx.SetGasPayment([]suitx.SuiObjectRef{*gasCoinRef})

			// Reference dest coin as owned object
			destCallArg := suitx.CallArg{
				Object: &suitx.ObjectArg{
					ImmOrOwnedObject: destCoinRef,
				},
			}
			destArg = tx.Object(destCallArg)
		}

		// Create source coin arguments
		sourceArgs := make([]suitx.Argument, 0, len(batchCoins))
		for _, coin := range batchCoins {
			coinRef, err := suitx.NewSuiObjectRef(
				models.SuiAddress(coin.CoinObjectId),
				coin.Version,
				models.ObjectDigest(coin.Digest),
			)
			if err != nil {
				return nil, fmt.Errorf("failed to create source coin ref: %w", err)
			}
			coinCallArg := suitx.CallArg{
				Object: &suitx.ObjectArg{
					ImmOrOwnedObject: coinRef,
				},
			}
			sourceArgs = append(sourceArgs, tx.Object(coinCallArg))

			// Add to total balance
			bal, _ := strconv.ParseUint(coin.Balance, 10, 64)
			totalBalance += bal
		}

		// Merge coins
		tx.MergeCoins(destArg, sourceArgs)

		// Execute the transaction
		resp, err := tx.Execute(ctx, models.SuiTransactionBlockOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		}, "WaitForLocalExecution")
		if err != nil {
			return nil, fmt.Errorf("failed to execute merge transaction batch %d: %w", batch/batchSize, err)
		}

		// Update destination coin version for next batch
		for _, change := range resp.ObjectChanges {
			if change.Type == "mutated" && change.ObjectId == destCoin.CoinObjectId {
				currentDestVersion = change.Version
				currentDestDigest = change.Digest
				lggr.Debugw("Updated destination coin reference after merge",
					"newVersion", change.Version,
					"newDigest", change.Digest)
			}
		}

		lggr.Infow("Completed merge batch",
			"batch", batch/batchSize,
			"coinsInBatch", len(batchCoins),
			"totalMergedSoFar", batch+len(batchCoins))
	}

	lggr.Infow("Successfully consolidated coins",
		"chainSelector", chainSelector,
		"coinType", coinType,
		"finalCoinID", destCoin.CoinObjectId,
		"totalBalance", totalBalance,
		"coinsConsolidated", len(allCoins))

	return &ConsolidatedCoin{
		ObjectID: destCoin.CoinObjectId,
		Version:  currentDestVersion,
		Digest:   currentDestDigest,
		Balance:  totalBalance,
	}, nil
}

// consolidateSuiGasCoins consolidates all SUI gas coins owned by the deployer into one.
// This should be called before splitting gas coins to clean up leftover coins from previous runs.
func consolidateSuiGasCoins(
	ctx context.Context,
	lggr logger.Logger,
	suiChain cldf_sui.Chain,
	chainSelector uint64,
	privateKeyHex string,
) (*ConsolidatedCoin, error) {
	// Get deployer address from signer
	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get deployer address: %w", err)
	}
	// Add 0x prefix if not present
	if !strings.HasPrefix(deployerAddr, "0x") {
		deployerAddr = "0x" + deployerAddr
	}

	suiClient := sui.NewSuiClient(suiChain.URL)

	return consolidateSuiCoins(
		ctx,
		lggr,
		suiClient,
		deployerAddr,
		"0x2::sui::SUI",
		privateKeyHex,
		chainSelector,
	)
}

// consolidateSuiLinkTokens consolidates all Link tokens owned by the deployer into one.
// This should be called before splitting Link tokens to clean up leftover tokens from previous runs.
func consolidateSuiLinkTokens(
	ctx context.Context,
	lggr logger.Logger,
	suiChain cldf_sui.Chain,
	chainSelector uint64,
	privateKeyHex string,
	linkTokenPkgID string,
) (*ConsolidatedCoin, error) {
	// Get deployer address from signer
	deployerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return nil, fmt.Errorf("failed to get deployer address: %w", err)
	}
	// Add 0x prefix if not present
	if !strings.HasPrefix(deployerAddr, "0x") {
		deployerAddr = "0x" + deployerAddr
	}

	suiClient := sui.NewSuiClient(suiChain.URL)

	coinType := linkTokenPkgID + "::link::LINK"

	return consolidateSuiCoins(
		ctx,
		lggr,
		suiClient,
		deployerAddr,
		coinType,
		privateKeyHex,
		chainSelector,
	)
}

// SuiGasCoinPool manages pre-split SUI gas coins for parallel transaction execution.
// Uses a buffered channel as a FIFO queue to distribute unique gas coins to concurrent workers.
// Each transaction pops a unique coin, preventing object version conflicts.
type SuiGasCoinPool struct {
	coins    chan SuiCoinObject // Buffered channel acts as FIFO queue
	lggr     logger.Logger
	chainSel uint64
	capacity int // Original capacity for logging
}

// NewSuiGasCoinPool creates a new gas coin pool with the given coins
func NewSuiGasCoinPool(lggr logger.Logger, chainSel uint64, coins []SuiCoinObject) *SuiGasCoinPool {
	ch := make(chan SuiCoinObject, len(coins))
	for _, coin := range coins {
		ch <- coin
	}
	return &SuiGasCoinPool{
		coins:    ch,
		lggr:     lggr,
		chainSel: chainSel,
		capacity: len(coins),
	}
}

// Pop retrieves a gas coin from the pool, blocking if empty until context is cancelled.
// Returns error if context is cancelled before a coin becomes available.
func (p *SuiGasCoinPool) Pop(ctx context.Context) (SuiCoinObject, error) {
	select {
	case coin := <-p.coins:
		remaining := len(p.coins)
		if remaining < p.capacity/4 {
			p.lggr.Warnw("Gas coin pool running low",
				"chainSelector", p.chainSel,
				"remaining", remaining,
				"capacity", p.capacity)
		}
		return coin, nil
	case <-ctx.Done():
		return SuiCoinObject{}, fmt.Errorf("context cancelled while waiting for gas coin: %w", ctx.Err())
	}
}

// TryPop attempts to retrieve a gas coin without blocking.
// Returns the coin and true if available, or empty coin and false if pool is empty.
func (p *SuiGasCoinPool) TryPop() (SuiCoinObject, bool) {
	select {
	case coin := <-p.coins:
		return coin, true
	default:
		return SuiCoinObject{}, false
	}
}

// Return puts a gas coin back into the pool for reuse.
// This should be called after a transaction completes (success or failure)
// since the gas coin object still exists with updated version.
func (p *SuiGasCoinPool) Return(coin SuiCoinObject) {
	select {
	case p.coins <- coin:
		// Successfully returned
	default:
		p.lggr.Warnw("Gas coin pool is full, discarding coin",
			"chainSelector", p.chainSel,
			"objectID", coin.ObjectID)
	}
}

// Size returns the number of available gas coins in the pool
func (p *SuiGasCoinPool) Size() int {
	return len(p.coins)
}

// SuiTokenPool manages a pool of Sui token objects for load testing.
// Uses a buffered channel as a FIFO queue - tokens are consumed (popped) and not reused
// since they are burned/locked during token transfers.
type SuiTokenPool struct {
	tokens   chan string // Buffered channel acts as FIFO queue
	lggr     logger.Logger
	capacity int // Original capacity for logging
}

// NewSuiTokenPool creates a new token pool with the given token object IDs
func NewSuiTokenPool(lggr logger.Logger, tokenObjectIDs []string) *SuiTokenPool {
	ch := make(chan string, len(tokenObjectIDs))
	for _, id := range tokenObjectIDs {
		ch <- id
	}
	return &SuiTokenPool{
		tokens:   ch,
		lggr:     lggr,
		capacity: len(tokenObjectIDs),
	}
}

// Pop retrieves a token object ID from the pool, blocking if empty until context is cancelled.
// Tokens are consumed (not returned) since they are burned/locked during transfers.
func (p *SuiTokenPool) Pop(ctx context.Context) (string, error) {
	select {
	case token := <-p.tokens:
		remaining := len(p.tokens)
		if remaining < p.capacity/4 {
			p.lggr.Warnw("Token pool running low",
				"remaining", remaining,
				"capacity", p.capacity)
		}
		return token, nil
	case <-ctx.Done():
		return "", fmt.Errorf("context cancelled while waiting for token: %w", ctx.Err())
	}
}

// TryPop attempts to retrieve a token without blocking.
// Returns the token and true if available, or empty string and false if pool is empty.
func (p *SuiTokenPool) TryPop() (string, bool) {
	select {
	case token := <-p.tokens:
		return token, true
	default:
		return "", false
	}
}

// GetNextToken retrieves a token from the pool (non-blocking).
// This is a convenience method that wraps TryPop for backward compatibility.
func (p *SuiTokenPool) GetNextToken() (string, error) {
	token, ok := p.TryPop()
	if !ok {
		return "", errors.New("token pool is empty")
	}
	return token, nil
}

// Size returns the number of available tokens in the pool
func (p *SuiTokenPool) Size() int {
	return len(p.tokens)
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

	// First, consolidate any existing Link token objects from previous runs
	lggr.Infow("Consolidating existing Link tokens before split",
		"chainSelector", chainSelector,
		"linkTokenPkg", linkTokenPkgID)

	consolidatedCoin, err := consolidateSuiLinkTokens(ctx, lggr, suiChain, chainSelector, privateKeyHex, linkTokenPkgID)
	if err != nil {
		return nil, fmt.Errorf("failed to consolidate Link tokens: %w", err)
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

	// Use consolidated coin if available, otherwise query for coins
	var sourceCoinID string
	var sourceBalance *big.Int
	var ownedCoins models.PaginatedCoinsResponse

	if consolidatedCoin != nil {
		// Use the consolidated coin
		sourceCoinID = consolidatedCoin.ObjectID
		sourceBalance = new(big.Int).SetUint64(consolidatedCoin.Balance)

		// Query for the updated coin info (we need the full CoinData for version/digest)
		coinType := linkTokenPkgID + "::link::LINK"
		ownedCoins, err = client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    deployerAddr,
			CoinType: coinType,
			Limit:    10,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query owned Link tokens after consolidation: %w", err)
		}

		// Find the consolidated coin in the list (it should be the only one or the one matching our ID)
		found := false
		for i, coin := range ownedCoins.Data {
			if coin.CoinObjectId == sourceCoinID {
				ownedCoins.Data[0] = coin // Move to front for later use
				if i != 0 {
					ownedCoins.Data[i] = ownedCoins.Data[0]
				}
				found = true
				break
			}
		}
		if !found && len(ownedCoins.Data) > 0 {
			// Use the first available coin if our consolidated coin ID wasn't found
			sourceCoinID = ownedCoins.Data[0].CoinObjectId
			sourceBalance, _ = new(big.Int).SetString(ownedCoins.Data[0].Balance, 10)
		}

		lggr.Infow("Using consolidated Link token for splitting",
			"sourceCoinID", sourceCoinID,
			"sourceBalance", sourceBalance.String())
	} else {
		// No consolidation happened (0 or 1 coin), query normally
		coinType := linkTokenPkgID + "::link::LINK"
		ownedCoins, err = client.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    deployerAddr,
			CoinType: coinType,
			Limit:    10,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query owned Link tokens: %w", err)
		}

		if len(ownedCoins.Data) == 0 {
			return nil, errors.New("deployer account has no Link tokens to split")
		}

		sourceCoinID = ownedCoins.Data[0].CoinObjectId
		sourceBalance, _ = new(big.Int).SetString(ownedCoins.Data[0].Balance, 10)

		lggr.Infow("Found Link token objects owned by deployer",
			"count", len(ownedCoins.Data),
			"firstCoinID", ownedCoins.Data[0].CoinObjectId,
			"firstCoinBalance", ownedCoins.Data[0].Balance)
	}

	// Calculate total amount needed
	//nolint:gosec // G115 - numTokenObjects and amountPerObject are bounded small positive values
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
	tokenObjectIDs := make([]string, 0, numTokenObjects)

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
				switch change.ObjectId {
				case gasCoins.Data[0].CoinObjectId:
					// Update gas coin version and digest for next batch
					gasCoins.Data[0].Version = change.Version
					gasCoins.Data[0].Digest = change.Digest
					lggr.Debugw("Updated gas coin reference",
						"newVersion", change.Version,
						"newDigest", change.Digest)
				case sourceCoinID:
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
						tokenObjectIDs = append(tokenObjectIDs, change.ObjectId)
						lggr.Debugw("Split token object created",
							"objectId", change.ObjectId,
							"batch", batch)
					}
				}
			}
		}

		lggr.Infow("Completed split batch",
			"batch", batch,
			"totalCreated", len(tokenObjectIDs))
	}

	lggr.Infow("Successfully split token objects for load test",
		"chainSelector", chainSelector,
		"numObjects", len(tokenObjectIDs),
		"amountPerObject", amountPerObject)

	return tokenObjectIDs, nil
}

// splitSuiGasCoins splits SUI gas coins owned by the deployer into multiple small objects
// for parallel transaction execution. Returns a SuiGasCoinPool that can be used to
// distribute unique gas coins to concurrent transaction workers.
func splitSuiGasCoins(
	ctx context.Context,
	lggr logger.Logger,
	suiChain cldf_sui.Chain,
	chainSelector uint64,
	numCoins int,
	amountPerCoin uint64, // amount in MIST (1 SUI = 1e9 MIST), e.g., 500_000_000 for 0.5 SUI
	privateKeyHex string, // Sui private key in hex format (without 0x prefix)
) (*SuiGasCoinPool, error) {
	// First, consolidate any existing SUI gas coin objects from previous runs
	lggr.Infow("Consolidating existing SUI gas coins before split",
		"chainSelector", chainSelector)

	consolidatedCoin, err := consolidateSuiGasCoins(ctx, lggr, suiChain, chainSelector, privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to consolidate SUI gas coins: %w", err)
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

	lggr.Infow("Splitting SUI gas coins for parallel load test",
		"chainSelector", chainSelector,
		"deployerAddr", deployerAddr,
		"numCoins", numCoins,
		"amountPerCoin", amountPerCoin)

	// Create Sui client
	suiClient := sui.NewSuiClient(suiChain.URL)

	// Use consolidated coin if available
	var sourceCoin models.CoinData
	var largestBalance uint64

	if consolidatedCoin != nil {
		// Query for the updated coin info after consolidation
		gasCoinType := "0x2::sui::SUI"
		ownedCoins, err := suiClient.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    deployerAddr,
			CoinType: gasCoinType,
			Limit:    50,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query owned SUI coins after consolidation: %w", err)
		}

		if len(ownedCoins.Data) == 0 {
			return nil, fmt.Errorf("deployer account has no SUI coins after consolidation")
		}

		// Find the consolidated coin (should be the only one or the largest)
		for _, coin := range ownedCoins.Data {
			if coin.CoinObjectId == consolidatedCoin.ObjectID {
				sourceCoin = coin
				largestBalance = consolidatedCoin.Balance
				break
			}
		}

		// If we couldn't find it by ID, use the largest available
		if sourceCoin.CoinObjectId == "" {
			for _, coin := range ownedCoins.Data {
				bal, _ := strconv.ParseUint(coin.Balance, 10, 64)
				if bal > largestBalance {
					largestBalance = bal
					sourceCoin = coin
				}
			}
		}

		lggr.Infow("Using consolidated SUI coin for splitting",
			"sourceCoinID", sourceCoin.CoinObjectId,
			"sourceBalance", largestBalance)
	} else {
		// No consolidation happened, query for coins
		gasCoinType := "0x2::sui::SUI"
		ownedCoins, err := suiClient.SuiXGetCoins(ctx, models.SuiXGetCoinsRequest{
			Owner:    deployerAddr,
			CoinType: gasCoinType,
			Limit:    50, // Get more coins to find one large enough
		})
		if err != nil {
			return nil, fmt.Errorf("failed to query owned SUI coins: %w", err)
		}

		if len(ownedCoins.Data) == 0 {
			return nil, fmt.Errorf("deployer account has no SUI coins")
		}

		// Find the largest coin to split from
		for _, coin := range ownedCoins.Data {
			bal, _ := strconv.ParseUint(coin.Balance, 10, 64)
			if bal > largestBalance {
				largestBalance = bal
				sourceCoin = coin
			}
		}

		lggr.Infow("Found SUI coins owned by deployer",
			"count", len(ownedCoins.Data),
			"largestCoinID", sourceCoin.CoinObjectId,
			"largestBalance", sourceCoin.Balance)
	}

	// Calculate total amount needed (plus gas for the split transactions)
	totalNeeded := uint64(numCoins) * amountPerCoin //nolint:gosec // G115 - numCoins is a bounded small positive value
	gasBuffer := uint64(100_000_000)                // 0.1 SUI for gas during splits
	if largestBalance < totalNeeded+gasBuffer {
		return nil, fmt.Errorf("insufficient SUI balance: have %d, need %d (including gas buffer)", largestBalance, totalNeeded+gasBuffer)
	}

	lggr.Infow("Splitting source SUI coin into smaller objects",
		"sourceCoinID", sourceCoin.CoinObjectId,
		"sourceBalance", sourceCoin.Balance,
		"totalNeeded", totalNeeded)

	// Create SDK signer using the provided private key (32-byte seed in hex format)
	seedBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key hex: %w", err)
	}
	if len(seedBytes) != 32 {
		return nil, fmt.Errorf("invalid seed length: expected 32 bytes, got %d", len(seedBytes))
	}

	suiSDKSigner := signer.NewSigner(seedBytes)

	// Collect split gas coin objects
	gasCoinObjects := make([]SuiCoinObject, 0, numCoins)

	// Track the source coin's current version/digest (it changes after each transaction)
	currentSourceVersion := sourceCoin.Version
	currentSourceDigest := sourceCoin.Digest

	// Split in batches to avoid gas limits
	batchSize := 20 // Larger batch for SUI splits since they're simpler
	for batch := 0; batch < numCoins; batch += batchSize {
		remaining := numCoins - batch
		if remaining > batchSize {
			remaining = batchSize
		}

		// Build transaction
		tx := suitx.NewTransaction()
		tx.SetSigner(suiSDKSigner)
		tx.SetSuiClient(suiClient.(*sui.Client))
		tx.SetSender(models.SuiAddress(deployerAddr))

		// Create object reference for source coin with current version
		sourceCoinRef, err := suitx.NewSuiObjectRef(
			models.SuiAddress(sourceCoin.CoinObjectId),
			currentSourceVersion,
			models.ObjectDigest(currentSourceDigest),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create source coin ref: %w", err)
		}

		// Set gas payment to use the source coin
		tx.SetGasPayment([]suitx.SuiObjectRef{*sourceCoinRef})
		tx.SetGasBudget(50_000_000) // 0.05 SUI gas budget per batch

		// Use tx.Gas() to reference the gas coin - this avoids the "mutable object appears twice" error
		// When splitting SUI coins, we split from the gas coin itself
		gasArg := tx.Gas()

		// Split coins one at a time from the gas coin
		splitCoins := make([]suitx.Argument, 0, remaining)
		for i := 0; i < remaining; i++ {
			splitCoin := tx.SplitCoins(gasArg, []suitx.Argument{tx.Pure(amountPerCoin)})
			splitCoins = append(splitCoins, splitCoin)
		}

		// Transfer all split coins to the deployer
		if len(splitCoins) > 0 {
			tx.TransferObjects(splitCoins, tx.Pure(deployerAddr))
		}

		// Execute the transaction
		resp, err := tx.Execute(ctx, models.SuiTransactionBlockOptions{
			ShowEffects:       true,
			ShowObjectChanges: true,
		}, "WaitForLocalExecution")
		if err != nil {
			return nil, fmt.Errorf("failed to execute gas split transaction batch %d: %w", batch, err)
		}

		// Update source coin version for next batch and collect created coins
		for _, change := range resp.ObjectChanges {
			if change.Type == "mutated" && change.ObjectId == sourceCoin.CoinObjectId {
				currentSourceVersion = change.Version
				currentSourceDigest = change.Digest
				lggr.Debugw("Updated source coin reference",
					"newVersion", change.Version,
					"newDigest", change.Digest)
			}
			if change.Type == "created" && change.ObjectType == "0x2::coin::Coin<0x2::sui::SUI>" {
				gasCoinObjects = append(gasCoinObjects, SuiCoinObject{
					ObjectID: change.ObjectId,
					Version:  change.Version,
					Digest:   change.Digest,
				})
				lggr.Debugw("Split gas coin created",
					"objectId", change.ObjectId,
					"batch", batch)
			}
		}

		lggr.Infow("Completed gas coin split batch",
			"batch", batch,
			"totalCreated", len(gasCoinObjects))
	}

	lggr.Infow("Successfully split SUI gas coins for parallel load test",
		"chainSelector", chainSelector,
		"numCoins", len(gasCoinObjects),
		"amountPerCoin", amountPerCoin)

	return NewSuiGasCoinPool(lggr, chainSelector, gasCoinObjects), nil
}

// CalculateRequiredGasCoins calculates the number of gas coins needed for a load test
func CalculateRequiredGasCoins(loadDurationSec int, requestFreqSec int, numDestinations int) int {
	if requestFreqSec <= 0 {
		requestFreqSec = 1
	}
	txPerDest := loadDurationSec / requestFreqSec
	// Total = txPerDest * numDestinations * 2 (need 2 coins per tx: one for gas, one for fee token) * 1.2 (20% buffer)
	total := int(float64(txPerDest*numDestinations*2) * 1.2)
	if total < 20 {
		total = 20 // Minimum 20 coins (10 transactions worth)
	}
	return total
}

// cleanupSuiDatabase cleans up Sui events and transactions tables before test runs
func cleanupSuiDatabase(dbURL string, lggr logger.Logger) {
	db, err := sqlx.Open("postgres", dbURL)
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
		30*time.Second,
		suiKeystore,
		2000,
		"WaitForEffectsCert",
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

func SplitCoin(ctx context.Context, env cldf.Environment) (string, error) {
	// signerAddr, _ := suiSigner.GetAddress()
	suiSelectors := env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(selectors.FamilySui))
	if len(suiSelectors) == 0 {
		return "", errors.New("no Sui chains available")
	}
	suiSelector := suiSelectors[0]

	suiChain := env.BlockChains.SuiChains()[suiSelector]
	client := sui.NewSuiClient(suiChain.URL)

	signerAddr, err := suiChain.Signer.GetAddress()
	if err != nil {
		return "", err
	}

	coinsResp, err := client.SuiXGetAllCoins(ctx, models.SuiXGetAllCoinsRequest{
		Owner: signerAddr,
		Limit: 10,
	})
	if err != nil {
		return "", err
	}
	coins := coinsResp.Data

	fmt.Println("Available Coins:", coins)
	if coins == nil {
		return "", errors.New("not enough balance")
	}

	var largestCoin models.CoinData
	var largestBalance uint64

	for _, c := range coins {
		if c.CoinType != "0x2::sui::SUI" {
			continue
		}

		// Convert string balance to number if needed
		var bal uint64
		switch v := any(c.Balance).(type) {
		case string:
			bal, _ = strconv.ParseUint(v, 10, 64)
		case uint64:
			bal = v
		default:
			return "", errors.New("unexpected balance")
		}

		if bal > largestBalance {
			largestBalance = bal
			largestCoin = c
		}
	}

	tokenToSplit := largestCoin.CoinObjectId
	fmt.Printf("Selected base coin (largest): %s with balance %v\n", tokenToSplit, largestCoin.Balance)

	req := models.PayRequest{
		Signer:      signerAddr,             // your wallet address
		SuiObjectId: []string{tokenToSplit}, // coin(s) to draw from (can include gas coin)
		Recipient:   []string{signerAddr},   // send to yourself to create a new coin
		Amount:      []string{"500000000"},  // in MIST (0.5 SUI)
		Gas:         nil,                    // let the node pick gas (or use a different one)
		GasBudget:   "10000000",             // budget for this tx
	}

	resp, err := client.Pay(ctx, req)
	if err != nil {
		return "", err
	}

	decodedSplitTx, err := base64.StdEncoding.DecodeString(resp.TxBytes)
	if err != nil {
		return "", err
	}

	splitTx, err := sui_common.SignAndSendTx(ctx, suiChain.Signer, client, decodedSplitTx, true)
	if err != nil {
		return "", err
	}

	fmt.Printf("Split coins")
	return splitTx.Effects.Created[0].Reference.ObjectId, nil
}
