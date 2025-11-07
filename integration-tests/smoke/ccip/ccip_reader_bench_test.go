package ccip

import (
	"context"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_reader_tester"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/pkg/chainaccessor"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	ccipocr3common "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-evm/pkg/config"

	writer_mocks "github.com/smartcontractkit/chainlink-ccip/mocks/chainlink_common/types"
	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"

	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmchaintypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

// This file contains benchmarks for the CCIPReader methods, we test here only performance of
// CCIPReader methods which are reaching the database to fetch matching logs (e.g. CCIPMessageSent)
// Under the hood, we verify if the CCIPReader/ChainAccessorLayer database queries are efficient enough to
// handle large number of logs.
//
// These tests are not fully e2e, because we don't interact with contracts, but rather
// we insert logs directly into the database, so we can control the amount of logs inserted and their content.
// Also, using contracts would be too slow, because we would have to wait for the
// transactions to be mined and consumed by LogPoller. Therefore, these tests should not be used to verify correctness
// of the CCIPReader methods, but rather their performance under various circumstances
// (different number of logs, different number of source and destination chains, etc.).
//
// For deep dive you can enable logging SQL queries and then testing them manually on
// the database (e.g. with `explain (analyze, buffers, verbose)`)
// Change lggr.SetLogLevel(zapcore.ErrorLevel) to zapcore.DebugLevel

var (
	onrampABI  = evmchaintypes.MustGetABI(onramp.OnRampABI)
	offrampABI = evmchaintypes.MustGetABI(offramp.OffRampABI)
)

// Benchmark_CCIPReader_CCIPMessageSent/LatestMsgSeqNum_5_source_chains_and_5_destination_chains-12            	    2649	    456123 ns/op
// Benchmark_CCIPReader_CCIPMessageSent/LatestMsgSeqNum_70_source_chains_and_10_destination_chains-12          	    2475	    501386 ns/op
// Benchmark_CCIPReader_CCIPMessageSent/MsgsBetweenSeqNums_5_source_chains_and_5_destination_chains-12         	      49	  21842648 ns/op
// Benchmark_CCIPReader_CCIPMessageSent/MsgsBetweenSeqNums_70_source_chains_and_10_destination_chains-12       	     114	  12192915 ns/op
func Benchmark_CCIPReader_CCIPMessageSent(b *testing.B) {
	tests := []struct {
		name                 string
		logsInsertedPerChain int
		startSeqNum          cciptypes.SeqNum
		endSeqNum            cciptypes.SeqNum
		sourceChainsCount    int
		destChainsCount      int

		expectedLogs   int
		expectedLatest cciptypes.SeqNum
	}{
		{
			// Case in which we have 5 chains densely connected generating large volume of logs
			name:                 "5 source chains and 5 destination chains",
			logsInsertedPerChain: 50_000, // 250k logs in total inserted (50k * 5 chains)
			startSeqNum:          5_000,
			endSeqNum:            5_256,
			sourceChainsCount:    5,
			destChainsCount:      5,
			expectedLogs:         257,
			expectedLatest:       9_899, // it's always smaller than latestBlock, because last 500 logs are not finalized
		},
		{
			// Case in which we have multiple a lot of source chains, but only a few destinations are in use
			name:                 "70 source chains and 10 destination chains",
			logsInsertedPerChain: 25_000, // 1.75kk logs in total inserted (25000 * 70 chains)
			startSeqNum:          2_000,
			endSeqNum:            2_300,
			sourceChainsCount:    70,
			destChainsCount:      10,
			expectedLogs:         301,
			expectedLatest:       2_449,
		},
	}

	for _, tt := range tests {
		reader := prepareMessageSentEventsInDb(
			b,
			tt.logsInsertedPerChain,
			tt.sourceChainsCount,
			tt.destChainsCount,
		)

		b.Run("MsgsBetweenSeqNums "+tt.name, func(b *testing.B) {
			for b.Loop() {
				msgs, err := reader.MsgsBetweenSeqNums(
					b.Context(),
					chainS1,
					cciptypes.NewSeqNumRange(tt.startSeqNum, tt.endSeqNum),
				)
				require.NoError(b, err)
				require.Len(b, msgs, tt.expectedLogs)
			}
		})

		b.Run("LatestMsgSeqNum "+tt.name, func(b *testing.B) {
			for b.Loop() {
				latest, err := reader.LatestMsgSeqNum(
					b.Context(),
					chainS1,
				)
				require.NoError(b, err)
				require.Equal(b, tt.expectedLatest, latest)
			}
		})
	}
}

// Benchmark_CCIPReader_CommitReportsGTETimestamp/5_chains,_only_some_logs_are_matched-12         	     304	   4502751 ns/op
// Benchmark_CCIPReader_CommitReportsGTETimestamp/70_chains,_little_logs_matched-12               	    2228	    516133 ns/op
// Benchmark_CCIPReader_CommitReportsGTETimestamp/100_chains,_everything_is_matched-12            	      37	  34224577 ns/op
func Benchmark_CCIPReader_CommitReportsGTETimestamp(b *testing.B) {
	tests := []struct {
		name                 string
		logsInsertedFirst    int
		logsInsertedMatching int
		numberOfChains       int
	}{
		{
			name:                 "5 chains, only some logs are matched",
			logsInsertedFirst:    1_000,
			logsInsertedMatching: 100,
			numberOfChains:       5,
		},
		{
			name:                 "70 chains, little logs matched",
			logsInsertedFirst:    10_000,
			logsInsertedMatching: 1,
			numberOfChains:       70,
		},
		{
			name:                 "100 chains, everything is matched",
			logsInsertedFirst:    1,
			logsInsertedMatching: 1_000,
			numberOfChains:       100,
		},
	}

	for _, tt := range tests {
		reader := prepareCommitReportsEventsInDb(
			b,
			tt.logsInsertedFirst,
			tt.logsInsertedMatching,
			tt.numberOfChains,
		)

		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				reports, err := reader.CommitReportsGTETimestamp(
					b.Context(),
					time.Now().Add(-10*time.Minute),
					primitives.Unconfirmed,
					tt.logsInsertedMatching,
				)
				require.NoError(b, err)
				require.Len(b, reports, tt.logsInsertedMatching)
			}
		})
	}
}

// Benchmark_CCIPReader_ExecutedMessages/5_source_chains_and_5_destination_chains-12         	      49	  24335313 ns/op
// Benchmark_CCIPReader_ExecutedMessages/20_dest_chains_and_40_sources_chains-12             	     187	   6190834 ns/op
func Benchmark_CCIPReader_ExecutedMessages(b *testing.B) {
	tests := []struct {
		name                 string
		logsInsertedPerChain int
		sourceChainsCount    int
		destChainsCount      int
		startSeqNum          cciptypes.SeqNum
		endSeqNum            cciptypes.SeqNum

		expectedChains       int
		expectedLogsPerChain int
	}{
		{
			// Case in which we have 5 chains densely connected generating large volume of logs
			name:                 "5 source chains and 5 destination chains",
			logsInsertedPerChain: 50_000, // 250k logs in total inserted (50k * 5 chains)
			startSeqNum:          11,
			endSeqNum:            20,
			sourceChainsCount:    5,
			destChainsCount:      5,
			expectedChains:       4,
			expectedLogsPerChain: 10,
		},
		{
			// Case in which we have multiple a lot of source chains, but only a few destinations are in use
			name:                 "20 dest chains and 40 sources chains",
			logsInsertedPerChain: 70_000, // 1.4kk logs in total inserted
			startSeqNum:          101,
			endSeqNum:            110,
			sourceChainsCount:    40,
			destChainsCount:      20,
			expectedChains:       39,
			expectedLogsPerChain: 10,
		},
	}

	for _, tt := range tests {
		reader := prepareExecutedStateChangesEventsInDb(
			b,
			tt.logsInsertedPerChain,
			tt.sourceChainsCount,
			tt.destChainsCount,
		)

		filters := map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{}
		for i := 0; i < tt.sourceChainsCount; i++ {
			// #nosec G115
			chainSelector := cciptypes.ChainSelector(i + 1)
			if chainSelector == chainD {
				continue
			}
			// This enforces variation in seqNum ranges
			// #nosec G115
			start := cciptypes.SeqNum(i*16) + tt.startSeqNum
			// #nosec G115
			stop := cciptypes.SeqNum(i*16) + tt.endSeqNum

			filters[chainSelector] = append(
				filters[chainSelector],
				cciptypes.NewSeqNumRange(start, stop),
			)
		}

		b.Run(tt.name, func(b *testing.B) {
			for b.Loop() {
				executedRanges, err := reader.ExecutedMessages(
					b.Context(),
					filters,
					primitives.Unconfirmed,
				)

				require.NoError(b, err)
				require.Len(b, executedRanges, tt.expectedChains)
				for _, seqNrs := range executedRanges {
					require.Len(b, seqNrs, tt.expectedLogsPerChain)
				}
			}
		})
	}
}

func prepareCommitReportsEventsInDb(
	b *testing.B,
	firstInserted int,
	matchingLogs int,
	numberOfChains int,
) ccipreaderpkg.CCIPReader {
	ctx := b.Context()
	s := benchSetup(ctx, b, benchSetupParams{
		ReaderChain:        chainD,
		DestChain:          chainD,
		Cfg:                evmconfig.DestReaderConfig,
		ContractNameToBind: consts.ContractNameOffRamp,
	})

	for j := range numberOfChains {
		// #nosec G115
		orm := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(j+1)), s.dbs, logger.TestLogger(b))
		// #nosec G115
		destChainID := cciptypes.ChainSelector(j + 1)
		populateDatabaseForCommitReportAccepted(
			ctx,
			b,
			s,
			orm,
			destChainID,
			numberOfChains,
			firstInserted,
			0,
		)

		populateDatabaseForCommitReportAccepted(
			ctx,
			b,
			s,
			orm,
			destChainID,
			numberOfChains,
			matchingLogs,
			firstInserted,
		)
	}

	return s.reader
}

func populateDatabaseForCommitReportAccepted(
	ctx context.Context,
	b *testing.B,
	testEnv *benchSetupData,
	orm *logpoller.DSORM,
	destChain cciptypes.ChainSelector,
	numberOfChains int,
	numOfReports int,
	offset int,
) {
	var logs []logpoller.Log
	commitReportEvent, exists := offrampABI.Events[consts.EventNameCommitReportAccepted]
	require.True(b, exists, "Event CommitReportAccepted not found in ABI")

	commitReportEventSig := commitReportEvent.ID
	commitReportAddress := testEnv.contractAddr

	// Calculate timestamp based on whether these are the first logs or matching logs
	var timestamp time.Time
	if offset == 0 {
		// For first set of logs, set timestamp to very old
		timestamp = time.Now().Add(-10 * time.Hour)
	} else {
		// For matching logs, use current time
		timestamp = time.Now()
	}

	for i := range numOfReports {
		// Calculate unique BlockNumber and LogIndex
		blockNumber := int64(offset + i + 1) // Offset ensures unique block numbers
		logIndex := int64(offset + i + 1)    // Offset ensures unique log indices

		// #nosec G115
		sourceChain := cciptypes.ChainSelector(i%numberOfChains + 1)
		if sourceChain == destChain {
			sourceChain++
		}

		// Simulate merkleRoots
		merkleRoots := []offramp.InternalMerkleRoot{
			{
				SourceChainSelector: uint64(sourceChain),
				OnRampAddress:       utils.RandomAddress().Bytes(),
				// #nosec G115
				MinSeqNr: uint64(i*100 + 1),
				// #nosec G115
				MaxSeqNr:   uint64(i*100 + 100),
				MerkleRoot: utils.RandomBytes32(),
			},
		}

		var unblessed []offramp.InternalMerkleRoot

		sourceToken := utils.RandomAddress()

		// Simulate priceUpdates
		priceUpdates := offramp.InternalPriceUpdates{
			TokenPriceUpdates: []offramp.InternalTokenPriceUpdate{
				{SourceToken: sourceToken, UsdPerToken: big.NewInt(8)},
			},
			GasPriceUpdates: []offramp.InternalGasPriceUpdate{
				{DestChainSelector: uint64(1), UsdPerUnitGas: big.NewInt(10)},
			},
		}

		// Combine encoded data
		encodedData, err := commitReportEvent.Inputs.
			Pack(merkleRoots, unblessed, priceUpdates)
		require.NoError(b, err)

		// Topics (first one is the event signature)
		topics := [][]byte{
			commitReportEventSig[:],
		}

		// Create log entry
		logs = append(logs, logpoller.Log{
			EVMChainID:     ubig.New(new(big.Int).SetUint64(uint64(destChain))),
			LogIndex:       logIndex,
			BlockHash:      utils.NewHash(),
			BlockNumber:    blockNumber,
			BlockTimestamp: timestamp,
			EventSig:       commitReportEventSig,
			Topics:         topics,
			Address:        commitReportAddress,
			TxHash:         utils.NewHash(),
			Data:           encodedData,
			CreatedAt:      time.Now(),
		})
	}

	// Insert logs into the database
	require.NoError(b, orm.InsertLogs(ctx, logs))
	require.NoError(b, orm.InsertBlock(ctx, utils.RandomHash(), int64(offset+numOfReports), timestamp, int64(offset+numOfReports), int64(offset+numOfReports)))
}

func prepareMessageSentEventsInDb(b *testing.B, logsInserted int, sourceChainsCount, destChainsCount int) ccipreaderpkg.CCIPReader {
	ctx := b.Context()

	s := benchSetup(ctx, b, benchSetupParams{
		ReaderChain:        chainS1,
		DestChain:          chainD,
		Cfg:                evmconfig.SourceReaderConfig,
		ContractNameToBind: consts.ContractNameOnRamp,
	})

	if logsInserted > 0 {
		for j := range sourceChainsCount {
			// #nosec G115
			orm := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(j+1)), s.dbs, logger.TestLogger(b))

			// #nosec G115
			populateDatabaseForMessageSent(ctx, b, s, orm, cciptypes.ChainSelector(j+1), destChainsCount, logsInserted, 0)
		}
	}

	return s.reader
}

func populateDatabaseForMessageSent(
	ctx context.Context,
	b *testing.B,
	testEnv *benchSetupData,
	orm *logpoller.DSORM,
	sourceChain cciptypes.ChainSelector,
	destChainCount int,
	numOfEvents int,
	offset int,
) {
	var logs []logpoller.Log
	messageSentEvent, exists := onrampABI.Events[consts.EventNameCCIPMessageSent]
	require.True(b, exists, "Event CCIPMessageSent not found in ABI")

	messageSentEventSig := messageSentEvent.ID
	messageSentEventAddress := testEnv.contractAddr

	largePayload := make([]byte, 8*1024)
	_, err := rand.Read(largePayload)
	require.NoError(b, err)

	for i := range numOfEvents {
		// Calculate unique BlockNumber and LogIndex
		blockNumber := int64(offset + i + 1) // Offset ensures unique block numbers
		logIndex := int64(offset + i + 1)    // Offset ensures unique log indices

		// Every event targets a different destination chain
		// #nosec G115
		destChainSelector := uint64(i%destChainCount + 1)
		// Every destination chain has its own sequence number
		// #nosec G115
		sequenceNumber := uint64(i / destChainCount)

		// Create InternalRampMessageHeader struct
		header := onramp.InternalRampMessageHeader{
			MessageId:           utils.NewHash(),
			SourceChainSelector: uint64(sourceChain),
			DestChainSelector:   destChainSelector,
			SequenceNumber:      sequenceNumber,
			// #nosec G115
			Nonce: uint64(i),
		}

		// Create InternalEVM2AnyTokenTransfer slice
		tokenTransfers := []onramp.InternalEVM2AnyTokenTransfer{
			{
				SourcePoolAddress: utils.RandomAddress(),
				DestTokenAddress:  []byte{0x01, 0x02},
				ExtraData:         []byte{0x03},
				// #nosec G115
				Amount:       big.NewInt(1000 + int64(i)),
				DestExecData: []byte{},
			},
		}

		// Make it large every 1000th event to simulate large payloads
		// especially to verify lack of errors related to index sizes
		// e.g. index row requires 9560 bytes, maximum size is 8191 (SQLSTATE 54000)
		data := []byte{0x04, 0x05}
		if i%1000 == 0 {
			data = largePayload
		}

		// Create InternalEVM2AnyRampMessage struct
		message := onramp.InternalEVM2AnyRampMessage{
			Header:    header,
			Sender:    utils.RandomAddress(),
			Data:      data,
			Receiver:  []byte{0x06, 0x07},
			ExtraArgs: []byte{0x08},
			FeeToken:  utils.RandomAddress(),
			// #nosec G115
			FeeTokenAmount: big.NewInt(2000 + int64(i)),
			// #nosec G115

			FeeValueJuels: big.NewInt(3000 + int64(i)),
			TokenAmounts:  tokenTransfers,
		}

		// Encode the non-indexed event data
		encodedData, err := messageSentEvent.Inputs.NonIndexed().Pack(
			message,
		)
		require.NoError(b, err)

		// Topics (event signature and indexed fields)
		topics := [][]byte{
			messageSentEventSig[:],                       // Event signature
			logpoller.EvmWord(destChainSelector).Bytes(), // Indexed destChainSelector
			logpoller.EvmWord(sequenceNumber).Bytes(),    // Indexed sequenceNumber
		}

		// Create log entry
		logs = append(logs, logpoller.Log{
			EVMChainID:     ubig.New(big.NewInt(0).SetUint64(uint64(sourceChain))),
			LogIndex:       logIndex,
			BlockHash:      utils.NewHash(),
			BlockNumber:    blockNumber,
			BlockTimestamp: time.Now(),
			EventSig:       messageSentEventSig,
			Topics:         topics,
			Address:        messageSentEventAddress,
			TxHash:         utils.NewHash(),
			Data:           encodedData,
			CreatedAt:      time.Now(),
		})
	}

	// Insert logs into the database
	require.NoError(b, orm.InsertLogs(ctx, logs))
	latestBlock := int64(numOfEvents)
	finalityDepth := int64(500)
	require.NoError(
		b,
		orm.InsertBlock(
			ctx,
			utils.RandomHash(),
			latestBlock,
			time.Now(),
			latestBlock-finalityDepth,
			latestBlock-finalityDepth,
		))
}

func prepareExecutedStateChangesEventsInDb(
	b *testing.B,
	logsInsertedPerChain int,
	sourceChainsCount int,
	destChainsCount int,
) ccipreaderpkg.CCIPReader {
	ctx := b.Context()
	s := benchSetup(ctx, b, benchSetupParams{
		ReaderChain:        chainD,
		DestChain:          chainD,
		Cfg:                evmconfig.DestReaderConfig,
		ContractNameToBind: consts.ContractNameOffRamp,
	})

	if logsInsertedPerChain > 0 {
		for j := range destChainsCount {
			// #nosec G115
			orm := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(j+1)), s.dbs, logger.TestLogger(b))

			// #nosec G115
			populateDatabaseForExecutionStateChanged(
				ctx,
				b,
				s,
				orm,
				cciptypes.ChainSelector(j+1),
				sourceChainsCount,
				logsInsertedPerChain,
				0,
			)
		}
	}

	return s.reader
}

func populateDatabaseForExecutionStateChanged(
	ctx context.Context,
	b *testing.B,
	testEnv *benchSetupData,
	orm *logpoller.DSORM,
	destChain cciptypes.ChainSelector,
	sourceChainCount int,
	numOfEvents int,
	offset int,
) {
	var logs []logpoller.Log
	executionStateEvent, exists := offrampABI.Events[consts.EventNameExecutionStateChanged]
	require.True(b, exists, "Event ExecutionStateChanged not found in ABI")

	executionStateEventSig := executionStateEvent.ID
	executionStateEventAddress := testEnv.contractAddr

	for i := range numOfEvents {
		// Calculate unique BlockNumber and LogIndex
		blockNumber := int64(offset + i + 1) // Offset ensures unique block numbers
		logIndex := int64(offset + i + 1)    // Offset ensures unique log indices

		// Every source chain will have its own message
		// #nosec G115
		sourceChainSelector := uint64(i%sourceChainCount + 1)
		// #nosec G115
		sequenceNumber := uint64(i / sourceChainCount)
		messageID := utils.NewHash()
		messageHash := utils.NewHash()
		state := uint8(2)
		returnData := []byte{0x01, 0x02}
		gasUsed := big.NewInt(int64(10000 + i))

		// Encode the non indexed event data
		encodedData, err := executionStateEvent.Inputs.NonIndexed().Pack(
			messageHash,
			state,
			returnData,
			gasUsed,
		)
		require.NoError(b, err)

		// Topics (event signature and indexed fields)
		topics := [][]byte{
			executionStateEventSig[:],                      // Event signature
			logpoller.EvmWord(sourceChainSelector).Bytes(), // Indexed sourceChainSelector
			logpoller.EvmWord(sequenceNumber).Bytes(),      // Indexed sequenceNumber
			messageID[:], // Indexed messageId
		}

		// Create log entry
		logs = append(logs, logpoller.Log{
			EVMChainID:     ubig.New(big.NewInt(0).SetUint64(uint64(destChain))),
			LogIndex:       logIndex,
			BlockHash:      utils.NewHash(),
			BlockNumber:    blockNumber,
			BlockTimestamp: time.Now(),
			EventSig:       executionStateEventSig,
			Topics:         topics,
			Address:        executionStateEventAddress,
			TxHash:         utils.NewHash(),
			Data:           encodedData,
			CreatedAt:      time.Now(),
		})
	}

	// Insert logs into the database
	require.NoError(b, orm.InsertLogs(ctx, logs))
	require.NoError(b, orm.InsertBlock(ctx, utils.RandomHash(), int64(offset+numOfEvents), time.Now(), int64(offset+numOfEvents), int64(offset+numOfEvents)))
}

func benchSetup(
	ctx context.Context,
	t testing.TB,
	params benchSetupParams,
) *benchSetupData {
	sb, auth := setupSimulatedBackendAndAuth(t)

	address, _, _, err := ccip_reader_tester.DeployCCIPReaderTester(auth, sb.Client())
	require.NoError(t, err)
	sb.Commit()

	// Setup contract client
	contract, err := ccip_reader_tester.NewCCIPReaderTester(address, sb.Client())
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	// Change that to DEBUG if you want to see SQL queries generated by ChainReader
	lggr.SetLogLevel(zapcore.ErrorLevel)

	var dbs sqlutil.DataSource
	{
		// Heavyweight database for benchmarks
		_, db := heavyweight.FullTestDBV2(t, nil)
		dbs = sqlutil.WrapDataSource(db, lggr, sqlutil.MonitorHook(func() bool {
			return true
		}))
	}

	lpOpts := logpoller.Opts{
		PollPeriod:               10 * time.Minute,
		FinalityDepth:            params.FinalityDepth,
		BackfillBatchSize:        10,
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, sb, big.NewInt(0).SetUint64(uint64(params.ReaderChain)))
	headTracker := headstest.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	orm := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(params.ReaderChain)), dbs, lggr)
	lp := logpoller.NewLogPoller(
		orm,
		cl,
		lggr,
		headTracker,
		lpOpts,
	)
	require.NoError(t, lp.Start(ctx))

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, params.Cfg)
	require.NoError(t, err)

	extendedCr := contractreader.NewExtendedContractReader(cr)

	err = extendedCr.Bind(ctx, []types.BoundContract{
		{
			Address: address.String(),
			Name:    params.ContractNameToBind,
		},
	})
	require.NoError(t, err)

	err = cr.Start(ctx)
	require.NoError(t, err)

	contractReaders := map[cciptypes.ChainSelector]contractreader.Extended{params.ReaderChain: extendedCr}
	contractWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)
	mockAddrCodec := newMockAddressCodec(t)
	mockContractWriter := writer_mocks.NewMockContractWriter(t)
	readerChainAccessor, err := chainaccessor.NewDefaultAccessor(
		lggr,
		params.ReaderChain,
		extendedCr,
		mockContractWriter,
		mockAddrCodec,
	)
	require.NoError(t, err)
	chainAccessors := map[ccipocr3common.ChainSelector]ccipocr3common.ChainAccessor{params.ReaderChain: readerChainAccessor}

	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(
		ctx,
		lggr,
		chainAccessors,
		contractReaders,
		contractWriters,
		params.DestChain,
		nil,
		mockAddrCodec,
	)

	t.Cleanup(func() {
		require.NoError(t, cr.Close())
		require.NoError(t, lp.Close())
	})

	return &benchSetupData{
		contractAddr: address,
		contract:     contract,
		reader:       reader,
		extendedCR:   extendedCr,
		dbs:          dbs,
	}
}

type benchSetupParams struct {
	ReaderChain        cciptypes.ChainSelector
	DestChain          cciptypes.ChainSelector
	Cfg                config.ChainReaderConfig
	ContractNameToBind string
	FinalityDepth      int64
}

type benchSetupData struct {
	contractAddr common.Address
	contract     *ccip_reader_tester.CCIPReaderTester
	reader       ccipreaderpkg.CCIPReader
	extendedCR   contractreader.Extended
	dbs          sqlutil.DataSource
}
