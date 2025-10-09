package ccip

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/btcsuite/btcutil/bech32"
	_ "github.com/lib/pq"
	"go.uber.org/atomic"

	"github.com/smartcontractkit/chainlink-common/pkg/logger"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	pkgtypes "github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query"

	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/config"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/database"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/indexer"
	"github.com/smartcontractkit/chainlink-sui/relayer/chainreader/reader"
	"github.com/smartcontractkit/chainlink-sui/relayer/client"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
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
			type Sui2AnyRampMessage struct {
				Header struct {
					MessageID           []byte `json:"message_id"`
					SourceChainSelector uint64 `json:"source_chain_selector"`
					DestChainSelector   uint64 `json:"dest_chain_selector"`
					SequenceNumber      uint64 `json:"sequence_number"`
					Nonce               uint64 `json:"nonce"`
				} `json:"header"`
				Sender         string `json:"sender"`
				Data           []byte `json:"data"`
				Receiver       []byte `json:"receiver"`
				ExtraArgs      []byte `json:"extra_args"`
				FeeToken       string `json:"fee_token"`
				FeeTokenAmount uint64 `json:"fee_token_amount"`
				FeeValueJuels  string `json:"fee_value_juels"`
				TokenAmounts   []any  `json:"token_amounts"`
			}

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
					Count:  100,
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
					lggr.Debugw("skipping marker/invalid event",
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
					Count:  100,
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
					lggr.Debugw("skipping empty commit event (likely marker)", "destChain", chainSelector)
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
					Count:  100,
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
					lggr.Debugw("skipping marker/invalid execution event",
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

				if _, ok := expectedRange[event.SourceChainSelector]; !ok {
					lggr.Debugw("received sui execution event (expectedRange not yet populated)",
						"srcChain", event.SourceChainSelector,
						"destChain", chainSelector,
						"seqNum", event.SequenceNumber)
				}

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
