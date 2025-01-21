package ccip

import (
	"context"
	"fmt"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/ccipevm"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"

	"github.com/smartcontractkit/chainlink-ccip/plugintypes"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/pgtest"

	readermocks "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/contractreader"

	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"

	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_reader_tester"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/fee_quoter"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"

	evmchaintypes "github.com/smartcontractkit/chainlink/v2/core/chains/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/offramp"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/onramp"
)

const (
	chainS1 = cciptypes.ChainSelector(1)
	chainS2 = cciptypes.ChainSelector(2)
	chainS3 = cciptypes.ChainSelector(3)
	chainD  = cciptypes.ChainSelector(4)
)

var (
	defaultGasPrice = assets.GWei(10)
)

var (
	onrampABI  = evmchaintypes.MustGetABI(onramp.OnRampABI)
	offrampABI = evmchaintypes.MustGetABI(offramp.OffRampABI)
)

func setupGetCommitGTETimestampTest(ctx context.Context, t testing.TB, finalityDepth int64, useHeavyDB bool) (*testSetupData, int64, common.Address) {
	sb, auth := setupSimulatedBackendAndAuth(t)
	onRampAddress := utils.RandomAddress()
	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:    chainD,
		DestChain:      chainD,
		OnChainSeqNums: nil,
		Cfg:            evmconfig.DestReaderConfig,
		ToMockBindings: map[cciptypes.ChainSelector][]types.BoundContract{
			chainS1: {
				{
					Address: onRampAddress.Hex(),
					Name:    consts.ContractNameOnRamp,
				},
			},
		},
		BindTester:         true,
		ContractNameToBind: consts.ContractNameOffRamp,
		SimulatedBackend:   sb,
		Auth:               auth,
		FinalityDepth:      finalityDepth,
		UseHeavyDB:         useHeavyDB,
	})

	return s, finalityDepth, onRampAddress
}

func setupExecutedMessageRangesTest(ctx context.Context, t testing.TB, useHeavyDB bool) *testSetupData {
	sb, auth := setupSimulatedBackendAndAuth(t)
	return testSetup(ctx, t, testSetupParams{
		ReaderChain:    chainD,
		DestChain:      chainD,
		OnChainSeqNums: nil,
		Cfg:            evmconfig.DestReaderConfig,
		// Cfg:              cfg,
		ToBindContracts:    nil,
		ToMockBindings:     nil,
		BindTester:         true,
		ContractNameToBind: consts.ContractNameOffRamp,
		SimulatedBackend:   sb,
		Auth:               auth,
		UseHeavyDB:         useHeavyDB,
	})
}

func setupMsgsBetweenSeqNumsTest(ctx context.Context, t testing.TB, useHeavyDB bool) *testSetupData {
	sb, auth := setupSimulatedBackendAndAuth(t)
	return testSetup(ctx, t, testSetupParams{
		ReaderChain:        chainS1,
		DestChain:          chainD,
		OnChainSeqNums:     nil,
		Cfg:                evmconfig.SourceReaderConfig,
		ToBindContracts:    nil,
		ToMockBindings:     nil,
		BindTester:         true,
		ContractNameToBind: consts.ContractNameOnRamp,
		SimulatedBackend:   sb,
		Auth:               auth,
		UseHeavyDB:         useHeavyDB,
	})
}

func emitCommitReports(ctx context.Context, t *testing.T, s *testSetupData, numReports int, tokenA common.Address, onRampAddress common.Address) uint64 {
	var firstReportTs uint64
	for i := uint8(0); int(i) < numReports; i++ {
		_, err := s.contract.EmitCommitReportAccepted(s.auth, ccip_reader_tester.OffRampCommitReport{
			PriceUpdates: ccip_reader_tester.InternalPriceUpdates{
				TokenPriceUpdates: []ccip_reader_tester.InternalTokenPriceUpdate{
					{
						SourceToken: tokenA,
						UsdPerToken: big.NewInt(1000),
					},
				},
				GasPriceUpdates: []ccip_reader_tester.InternalGasPriceUpdate{
					{
						DestChainSelector: uint64(chainD),
						UsdPerUnitGas:     big.NewInt(90),
					},
				},
			},
			MerkleRoots: []ccip_reader_tester.InternalMerkleRoot{
				{
					SourceChainSelector: uint64(chainS1),
					MinSeqNr:            10,
					MaxSeqNr:            20,
					MerkleRoot:          [32]byte{i + 1},
					OnRampAddress:       common.LeftPadBytes(onRampAddress.Bytes(), 32),
				},
			},
			RmnSignatures: []ccip_reader_tester.IRMNRemoteSignature{
				{
					R: [32]byte{1},
					S: [32]byte{2},
				},
				{
					R: [32]byte{3},
					S: [32]byte{4},
				},
			},
		})
		require.NoError(t, err)
		bh := s.sb.Commit()
		b, err := s.sb.Client().BlockByHash(ctx, bh)
		require.NoError(t, err)
		if firstReportTs == 0 {
			firstReportTs = b.Time()
		}
	}
	return firstReportTs
}

func TestCCIPReader_GetOffRampConfigDigest(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	sb, auth := setupSimulatedBackendAndAuth(t)

	addr, _, _, err := offramp.DeployOffRamp(auth, sb.Client(), offramp.OffRampStaticConfig{
		ChainSelector:        uint64(chainD),
		GasForCallExactCheck: 5_000,
		RmnRemote:            utils.RandomAddress(),
		TokenAdminRegistry:   utils.RandomAddress(),
		NonceManager:         utils.RandomAddress(),
	}, offramp.OffRampDynamicConfig{
		FeeQuoter:                               utils.RandomAddress(),
		PermissionLessExecutionThresholdSeconds: 1,
		IsRMNVerificationDisabled:               true,
		MessageInterceptor:                      utils.RandomAddress(),
	}, []offramp.OffRampSourceChainConfigArgs{})
	require.NoError(t, err)
	sb.Commit()

	offRamp, err := offramp.NewOffRamp(addr, sb.Client())
	require.NoError(t, err)

	commitConfigDigest := utils.RandomBytes32()
	execConfigDigest := utils.RandomBytes32()

	_, err = offRamp.SetOCR3Configs(auth, []offramp.MultiOCR3BaseOCRConfigArgs{
		{
			ConfigDigest:                   commitConfigDigest,
			OcrPluginType:                  consts.PluginTypeCommit,
			F:                              1,
			IsSignatureVerificationEnabled: true,
			Signers:                        []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()},
			Transmitters:                   []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()},
		},
		{
			ConfigDigest:                   execConfigDigest,
			OcrPluginType:                  consts.PluginTypeExecute,
			F:                              1,
			IsSignatureVerificationEnabled: false,
			Signers:                        []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()},
			Transmitters:                   []common.Address{utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress(), utils.RandomAddress()},
		},
	})
	require.NoError(t, err)
	sb.Commit()

	commitConfigDetails, err := offRamp.LatestConfigDetails(&bind.CallOpts{
		Context: ctx,
	}, consts.PluginTypeCommit)
	require.NoError(t, err)
	require.Equal(t, commitConfigDigest, commitConfigDetails.ConfigInfo.ConfigDigest)

	execConfigDetails, err := offRamp.LatestConfigDetails(&bind.CallOpts{
		Context: ctx,
	}, consts.PluginTypeExecute)
	require.NoError(t, err)
	require.Equal(t, execConfigDigest, execConfigDetails.ConfigInfo.ConfigDigest)

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, sb, big.NewInt(1337))
	headTracker := headtracker.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	orm := logpoller.NewORM(big.NewInt(1337), db, lggr)
	lp := logpoller.NewLogPoller(
		orm,
		cl,
		lggr,
		headTracker,
		lpOpts,
	)
	require.NoError(t, lp.Start(ctx))
	t.Cleanup(func() { require.NoError(t, lp.Close()) })

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, evmconfig.DestReaderConfig)
	require.NoError(t, err)
	err = cr.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, cr.Close()) })

	extendedCr := contractreader.NewExtendedContractReader(cr)
	err = extendedCr.Bind(ctx, []types.BoundContract{
		{
			Address: addr.Hex(),
			Name:    consts.ContractNameOffRamp,
		},
	})
	require.NoError(t, err)

	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(
		ctx,
		lggr,
		map[cciptypes.ChainSelector]contractreader.Extended{
			chainD: extendedCr,
		},
		nil,
		chainD,
		addr.Bytes(),
		ccipevm.NewExtraArgsCodec(),
	)

	ccipReaderCommitDigest, err := reader.GetOffRampConfigDigest(ctx, consts.PluginTypeCommit)
	require.NoError(t, err)
	require.Equal(t, commitConfigDigest, ccipReaderCommitDigest)

	ccipReaderExecDigest, err := reader.GetOffRampConfigDigest(ctx, consts.PluginTypeExecute)
	require.NoError(t, err)
	require.Equal(t, execConfigDigest, ccipReaderExecDigest)
}

func TestCCIPReader_CommitReportsGTETimestamp(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	s, _, onRampAddress := setupGetCommitGTETimestampTest(ctx, t, 0, false)

	tokenA := common.HexToAddress("123")
	const numReports = 5

	firstReportTs := emitCommitReports(ctx, t, s, numReports, tokenA, onRampAddress)

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var reports []plugintypes.CommitPluginReportWithMeta
	var err error
	require.Eventually(t, func() bool {
		reports, err = s.reader.CommitReportsGTETimestamp(
			ctx,
			chainD,
			// Skips first report
			//nolint:gosec // this won't overflow
			time.Unix(int64(firstReportTs)+1, 0),
			10,
		)
		require.NoError(t, err)
		return len(reports) == numReports-1
	}, 30*time.Second, 50*time.Millisecond)

	assert.Len(t, reports, numReports-1)
	assert.Len(t, reports[0].Report.MerkleRoots, 1)
	assert.Equal(t, chainS1, reports[0].Report.MerkleRoots[0].ChainSel)
	assert.Equal(t, onRampAddress.Bytes(), []byte(reports[0].Report.MerkleRoots[0].OnRampAddress))
	assert.Equal(t, cciptypes.SeqNum(10), reports[0].Report.MerkleRoots[0].SeqNumsRange.Start())
	assert.Equal(t, cciptypes.SeqNum(20), reports[0].Report.MerkleRoots[0].SeqNumsRange.End())
	assert.Equal(t, "0x0200000000000000000000000000000000000000000000000000000000000000",
		reports[0].Report.MerkleRoots[0].MerkleRoot.String())
	assert.Equal(t, tokenA.String(), string(reports[0].Report.PriceUpdates.TokenPriceUpdates[0].TokenID))
	assert.Equal(t, uint64(1000), reports[0].Report.PriceUpdates.TokenPriceUpdates[0].Price.Uint64())
	assert.Equal(t, chainD, reports[0].Report.PriceUpdates.GasPriceUpdates[0].ChainSel)
	assert.Equal(t, uint64(90), reports[0].Report.PriceUpdates.GasPriceUpdates[0].GasPrice.Uint64())
}

func TestCCIPReader_CommitReportsGTETimestamp_RespectsFinality(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	var finalityDepth int64 = 10
	s, _, onRampAddress := setupGetCommitGTETimestampTest(ctx, t, finalityDepth, false)

	tokenA := common.HexToAddress("123")
	const numReports = 5

	firstReportTs := emitCommitReports(ctx, t, s, numReports, tokenA, onRampAddress)

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var reports []plugintypes.CommitPluginReportWithMeta
	var err error
	// Will not return any reports as the finality depth is not reached.
	require.Never(t, func() bool {
		reports, err = s.reader.CommitReportsGTETimestamp(
			ctx,
			chainD,
			// Skips first report
			//nolint:gosec // this won't overflow
			time.Unix(int64(firstReportTs)+1, 0),
			10,
		)
		require.NoError(t, err)
		return len(reports) == numReports-1
	}, 20*time.Second, 50*time.Millisecond)

	// Commit finality depth number of blocks.
	for i := 0; i < int(finalityDepth); i++ {
		s.sb.Commit()
	}

	require.Eventually(t, func() bool {
		reports, err = s.reader.CommitReportsGTETimestamp(
			ctx,
			chainD,
			// Skips first report
			//nolint:gosec // this won't overflow
			time.Unix(int64(firstReportTs)+1, 0),
			10,
		)
		require.NoError(t, err)
		return len(reports) == numReports-1
	}, 30*time.Second, 50*time.Millisecond)

	assert.Len(t, reports, numReports-1)
	assert.Len(t, reports[0].Report.MerkleRoots, 1)
	assert.Equal(t, chainS1, reports[0].Report.MerkleRoots[0].ChainSel)
	assert.Equal(t, onRampAddress.Bytes(), []byte(reports[0].Report.MerkleRoots[0].OnRampAddress))
	assert.Equal(t, cciptypes.SeqNum(10), reports[0].Report.MerkleRoots[0].SeqNumsRange.Start())
	assert.Equal(t, cciptypes.SeqNum(20), reports[0].Report.MerkleRoots[0].SeqNumsRange.End())
	assert.Equal(t, "0x0200000000000000000000000000000000000000000000000000000000000000",
		reports[0].Report.MerkleRoots[0].MerkleRoot.String())
	assert.Equal(t, tokenA.String(), string(reports[0].Report.PriceUpdates.TokenPriceUpdates[0].TokenID))
	assert.Equal(t, uint64(1000), reports[0].Report.PriceUpdates.TokenPriceUpdates[0].Price.Uint64())
	assert.Equal(t, chainD, reports[0].Report.PriceUpdates.GasPriceUpdates[0].ChainSel)
	assert.Equal(t, uint64(90), reports[0].Report.PriceUpdates.GasPriceUpdates[0].GasPrice.Uint64())
}

func TestCCIPReader_ExecutedMessageRanges(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	s := setupExecutedMessageRangesTest(ctx, t, false)
	_, err := s.contract.EmitExecutionStateChanged(
		s.auth,
		uint64(chainS1),
		14,
		cciptypes.Bytes32{1, 0, 0, 1},
		cciptypes.Bytes32{1, 0, 0, 1, 1, 0, 0, 1},
		1,
		[]byte{1, 2, 3, 4},
		big.NewInt(250_000),
	)
	require.NoError(t, err)
	s.sb.Commit()

	_, err = s.contract.EmitExecutionStateChanged(
		s.auth,
		uint64(chainS1),
		15,
		cciptypes.Bytes32{1, 0, 0, 2},
		cciptypes.Bytes32{1, 0, 0, 2, 1, 0, 0, 2},
		1,
		[]byte{1, 2, 3, 4, 5},
		big.NewInt(350_000),
	)
	require.NoError(t, err)
	s.sb.Commit()

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var executedRanges []cciptypes.SeqNumRange
	require.Eventually(t, func() bool {
		executedRanges, err = s.reader.ExecutedMessageRanges(
			ctx,
			chainS1,
			chainD,
			cciptypes.NewSeqNumRange(14, 15),
		)
		require.NoError(t, err)
		return len(executedRanges) == 2
	}, tests.WaitTimeout(t), 50*time.Millisecond)

	assert.Equal(t, cciptypes.SeqNum(14), executedRanges[0].Start())
	assert.Equal(t, cciptypes.SeqNum(14), executedRanges[0].End())

	assert.Equal(t, cciptypes.SeqNum(15), executedRanges[1].Start())
	assert.Equal(t, cciptypes.SeqNum(15), executedRanges[1].End())
}

func TestCCIPReader_MsgsBetweenSeqNums(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	s := setupMsgsBetweenSeqNumsTest(ctx, t, false)
	_, err := s.contract.EmitCCIPMessageSent(s.auth, uint64(chainD), ccip_reader_tester.InternalEVM2AnyRampMessage{
		Header: ccip_reader_tester.InternalRampMessageHeader{
			MessageId:           [32]byte{1, 0, 0, 0, 0},
			SourceChainSelector: uint64(chainS1),
			DestChainSelector:   uint64(chainD),
			SequenceNumber:      10,
		},
		Sender:         utils.RandomAddress(),
		Data:           make([]byte, 0),
		Receiver:       utils.RandomAddress().Bytes(),
		ExtraArgs:      make([]byte, 0),
		FeeToken:       utils.RandomAddress(),
		FeeTokenAmount: big.NewInt(1),
		FeeValueJuels:  big.NewInt(2),
		TokenAmounts:   []ccip_reader_tester.InternalEVM2AnyTokenTransfer{{Amount: big.NewInt(1)}, {Amount: big.NewInt(2)}},
	})
	require.NoError(t, err)

	_, err = s.contract.EmitCCIPMessageSent(s.auth, uint64(chainD), ccip_reader_tester.InternalEVM2AnyRampMessage{
		Header: ccip_reader_tester.InternalRampMessageHeader{
			MessageId:           [32]byte{1, 0, 0, 0, 1},
			SourceChainSelector: uint64(chainS1),
			DestChainSelector:   uint64(chainD),
			SequenceNumber:      15,
		},
		Sender:         utils.RandomAddress(),
		Data:           make([]byte, 0),
		Receiver:       utils.RandomAddress().Bytes(),
		ExtraArgs:      make([]byte, 0),
		FeeToken:       utils.RandomAddress(),
		FeeTokenAmount: big.NewInt(3),
		FeeValueJuels:  big.NewInt(4),
		TokenAmounts:   []ccip_reader_tester.InternalEVM2AnyTokenTransfer{{Amount: big.NewInt(3)}, {Amount: big.NewInt(4)}},
	})
	require.NoError(t, err)

	s.sb.Commit()

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var msgs []cciptypes.Message
	require.Eventually(t, func() bool {
		msgs, err = s.reader.MsgsBetweenSeqNums(
			ctx,
			chainS1,
			cciptypes.NewSeqNumRange(5, 20),
		)
		require.NoError(t, err)
		return len(msgs) == 2
	}, tests.WaitTimeout(t), 100*time.Millisecond)

	require.Len(t, msgs, 2)
	// sort to ensure ascending order of sequence numbers.
	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Header.SequenceNumber < msgs[j].Header.SequenceNumber
	})
	require.Equal(t, cciptypes.SeqNum(10), msgs[0].Header.SequenceNumber)
	require.Equal(t, big.NewInt(1), msgs[0].FeeTokenAmount.Int)
	require.Equal(t, big.NewInt(2), msgs[0].FeeValueJuels.Int)
	require.Equal(t, int64(1), msgs[0].TokenAmounts[0].Amount.Int64())
	require.Equal(t, int64(2), msgs[0].TokenAmounts[1].Amount.Int64())

	require.Equal(t, cciptypes.SeqNum(15), msgs[1].Header.SequenceNumber)
	require.Equal(t, big.NewInt(3), msgs[1].FeeTokenAmount.Int)
	require.Equal(t, big.NewInt(4), msgs[1].FeeValueJuels.Int)
	require.Equal(t, int64(3), msgs[1].TokenAmounts[0].Amount.Int64())
	require.Equal(t, int64(4), msgs[1].TokenAmounts[1].Amount.Int64())

	for _, msg := range msgs {
		require.Equal(t, chainS1, msg.Header.SourceChainSelector)
		require.Equal(t, chainD, msg.Header.DestChainSelector)
	}
}

func TestCCIPReader_NextSeqNum(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)

	onChainSeqNums := map[cciptypes.ChainSelector]cciptypes.SeqNum{
		chainS1: 10,
		chainS2: 20,
		chainS3: 30,
	}

	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameOffRamp: {
				ContractABI: ccip_reader_tester.CCIPReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.MethodNameGetSourceChainConfig: {
						ChainSpecificName: "getSourceChainConfig",
						ReadType:          evmtypes.Method,
					},
				},
			},
		},
	}

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:        chainD,
		DestChain:          chainD,
		OnChainSeqNums:     onChainSeqNums,
		Cfg:                cfg,
		ToBindContracts:    nil,
		ToMockBindings:     nil,
		BindTester:         true,
		ContractNameToBind: consts.ContractNameOffRamp,
		SimulatedBackend:   sb,
		Auth:               auth,
	})

	seqNums, err := s.reader.NextSeqNum(ctx, []cciptypes.ChainSelector{chainS1, chainS2, chainS3})
	require.NoError(t, err)
	assert.Len(t, seqNums, 3)
	assert.Equal(t, cciptypes.SeqNum(10), seqNums[chainS1])
	assert.Equal(t, cciptypes.SeqNum(20), seqNums[chainS2])
	assert.Equal(t, cciptypes.SeqNum(30), seqNums[chainS3])
}

func TestCCIPReader_GetExpectedNextSequenceNumber(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	//env := NewMemoryEnvironmentContractsOnly(t, logger.TestLogger(t), 2, 4, nil)
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	destChain, srcChain := selectors[0], selectors[1]

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, destChain, srcChain, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, srcChain, destChain, false)

	reader := testSetupRealContracts(
		ctx,
		t,
		destChain,
		map[cciptypes.ChainSelector][]types.BoundContract{
			cciptypes.ChainSelector(srcChain): {
				{
					Address: state.Chains[srcChain].OnRamp.Address().String(),
					Name:    consts.ContractNameOnRamp,
				},
			},
		},
		nil,
		env,
	)

	maxExpectedSeqNum := uint64(10)
	var i uint64
	for i = 1; i < maxExpectedSeqNum; i++ {
		msg := testhelpers.DefaultRouterMessage(state.Chains[destChain].Receiver.Address())
		msgSentEvent := testhelpers.TestSendRequest(t, env.Env, state, srcChain, destChain, false, msg)
		require.Equal(t, uint64(i), msgSentEvent.SequenceNumber)
		require.Equal(t, uint64(i), msgSentEvent.Message.Header.Nonce) // check outbound nonce incremented
		seqNum, err2 := reader.GetExpectedNextSequenceNumber(ctx, cs(srcChain), cs(destChain))
		require.NoError(t, err2)
		require.Equal(t, cciptypes.SeqNum(i+1), seqNum)
	}
}

func TestCCIPReader_Nonces(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	var nonces = map[cciptypes.ChainSelector]map[common.Address]uint64{
		chainS1: {
			utils.RandomAddress(): 10,
			utils.RandomAddress(): 20,
		},
		chainS2: {
			utils.RandomAddress(): 30,
			utils.RandomAddress(): 40,
		},
		chainS3: {
			utils.RandomAddress(): 50,
			utils.RandomAddress(): 60,
		},
	}

	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameNonceManager: {
				ContractABI: ccip_reader_tester.CCIPReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.MethodNameGetInboundNonce: {
						ChainSpecificName: "getInboundNonce",
						ReadType:          evmtypes.Method,
					},
				},
			},
		},
	}

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:        chainD,
		DestChain:          chainD,
		Cfg:                cfg,
		BindTester:         true,
		ContractNameToBind: consts.ContractNameNonceManager,
		SimulatedBackend:   sb,
		Auth:               auth,
	})

	// Add some nonces.
	for chain, addrs := range nonces {
		for addr, nonce := range addrs {
			_, err := s.contract.SetInboundNonce(s.auth, uint64(chain), nonce, common.LeftPadBytes(addr.Bytes(), 32))
			require.NoError(t, err)
		}
	}
	s.sb.Commit()

	for sourceChain, addrs := range nonces {
		var addrQuery []string
		for addr := range addrs {
			addrQuery = append(addrQuery, addr.String())
		}
		addrQuery = append(addrQuery, utils.RandomAddress().String())

		results, err := s.reader.Nonces(ctx, sourceChain, chainD, addrQuery)
		require.NoError(t, err)
		assert.Len(t, results, len(addrQuery))
		for addr, nonce := range addrs {
			assert.Equal(t, nonce, results[addr.String()])
		}
	}
}

func Test_GetChainFeePriceUpdates(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain1, chain2, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain2, chain1, false)

	// Change the gas price for chain2
	feeQuoter := state.Chains[chain1].FeeQuoter
	_, err = feeQuoter.UpdatePrices(
		env.Env.Chains[chain1].DeployerKey, fee_quoter.InternalPriceUpdates{
			GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
				{
					DestChainSelector: chain2,
					UsdPerUnitGas:     defaultGasPrice.ToInt(),
				},
			},
		},
	)
	require.NoError(t, err)
	be := env.Env.Chains[chain1].Client.(*memory.Backend)
	be.Commit()

	gas, err := feeQuoter.GetDestinationChainGasPrice(&bind.CallOpts{}, chain2)
	require.NoError(t, err)
	require.Equal(t, defaultGasPrice.ToInt(), gas.Value)

	reader := testSetupRealContracts(
		ctx,
		t,
		chain1,
		//evmconfig.DestReaderConfig,
		map[cciptypes.ChainSelector][]types.BoundContract{
			cciptypes.ChainSelector(chain1): {
				{
					Address: state.Chains[chain1].FeeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
			},
		},
		nil,
		env,
	)

	updates := reader.GetChainFeePriceUpdate(ctx, []cciptypes.ChainSelector{cs(chain1), cs(chain2)})
	// only chain1 has a bound contract
	require.Len(t, updates, 1)
	require.Equal(t, defaultGasPrice.ToInt(), updates[cs(chain2)].Value.Int)
}

func Test_LinkPriceUSD(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain1, chain2, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain2, chain1, false)

	reader := testSetupRealContracts(
		ctx,
		t,
		chain1,
		map[cciptypes.ChainSelector][]types.BoundContract{
			cciptypes.ChainSelector(chain1): {
				{
					Address: state.Chains[chain1].FeeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
			},
		},
		nil,
		env,
	)

	linkPriceUSD, err := reader.LinkPriceUSD(ctx)
	require.NoError(t, err)
	require.NotNil(t, linkPriceUSD.Int)
	require.Equal(t, testhelpers.DefaultLinkPrice, linkPriceUSD.Int)
}

func Test_GetMedianDataAvailabilityGasConfig(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	env, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(4))
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	destChain, chain1, chain2, chain3 := selectors[0], selectors[1], selectors[2], selectors[3]

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain1, destChain, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain2, destChain, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain3, destChain, false)

	boundContracts := map[cciptypes.ChainSelector][]types.BoundContract{}
	for i, selector := range env.Env.AllChainSelectorsExcluding([]uint64{destChain}) {
		feeQuoter := state.Chains[selector].FeeQuoter
		destChainCfg := changeset.DefaultFeeQuoterDestChainConfig()
		//nolint:gosec // disable G115
		destChainCfg.DestDataAvailabilityOverheadGas = uint32(100 + i)
		//nolint:gosec // disable G115
		destChainCfg.DestGasPerDataAvailabilityByte = uint16(200 + i)
		//nolint:gosec // disable G115
		destChainCfg.DestDataAvailabilityMultiplierBps = uint16(1 + i)
		_, err2 := feeQuoter.ApplyDestChainConfigUpdates(env.Env.Chains[selector].DeployerKey, []fee_quoter.FeeQuoterDestChainConfigArgs{
			{
				DestChainSelector: destChain,
				DestChainConfig:   destChainCfg,
			},
		})
		require.NoError(t, err2)
		be := env.Env.Chains[selector].Client.(*memory.Backend)
		be.Commit()
		boundContracts[cs(selector)] = []types.BoundContract{
			{
				Address: feeQuoter.Address().String(),
				Name:    consts.ContractNameFeeQuoter,
			},
		}
	}

	reader := testSetupRealContracts(
		ctx,
		t,
		destChain,
		boundContracts,
		nil,
		env,
	)

	daConfig, err := reader.GetMedianDataAvailabilityGasConfig(ctx)
	require.NoError(t, err)

	// Verify the results
	require.Equal(t, uint32(101), daConfig.DestDataAvailabilityOverheadGas)
	require.Equal(t, uint16(201), daConfig.DestGasPerDataAvailabilityByte)
	require.Equal(t, uint16(2), daConfig.DestDataAvailabilityMultiplierBps)
}

func Test_GetWrappedNativeTokenPriceUSD(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain1, chain2, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain2, chain1, false)

	reader := testSetupRealContracts(
		ctx,
		t,
		chain1,
		map[cciptypes.ChainSelector][]types.BoundContract{
			cciptypes.ChainSelector(chain1): {
				{
					Address: state.Chains[chain1].FeeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
				{
					Address: state.Chains[chain1].Router.Address().String(),
					Name:    consts.ContractNameRouter,
				},
			},
		},
		nil,
		env,
	)

	prices := reader.GetWrappedNativeTokenPriceUSD(ctx, []cciptypes.ChainSelector{cciptypes.ChainSelector(chain1), cciptypes.ChainSelector(chain2)})

	// Only chainD has reader contracts bound
	require.Len(t, prices, 1)
	require.Equal(t, testhelpers.DefaultWethPrice, prices[cciptypes.ChainSelector(chain1)].Int)
}

// Benchmark Results:
// Benchmark_CCIPReader_CommitReportsGTETimestamp/FirstLogs_0_MatchLogs_0-14             16948      67728 ns/op        30387 B/op          417 allocs/op
// Benchmark_CCIPReader_CommitReportsGTETimestamp/FirstLogs_1_MatchLogs_10-14            1650       741741 ns/op       528334 B/op         9929 allocs/op
// Benchmark_CCIPReader_CommitReportsGTETimestamp/FirstLogs_10_MatchLogs_100-14          195        6096328 ns/op      4739856 B/op        92345 allocs/op
// Benchmark_CCIPReader_CommitReportsGTETimestamp/FirstLogs_100_MatchLogs_10000-14       2          582712583 ns/op    454375304 B/op      8931990 allocs/op
func Benchmark_CCIPReader_CommitReportsGTETimestamp(b *testing.B) {
	tests := []struct {
		logsInsertedFirst    int
		logsInsertedMatching int
	}{
		{0, 0},
		{1, 10},
		{10, 100},
		{100, 10_000},
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("FirstLogs_%d_MatchLogs_%d", tt.logsInsertedMatching, tt.logsInsertedFirst), func(b *testing.B) {
			benchmarkCommitReports(b, tt.logsInsertedFirst, tt.logsInsertedMatching)
		})
	}
}

func benchmarkCommitReports(b *testing.B, logsInsertedFirst int, logsInsertedMatching int) {
	// Initialize test setup
	ctx := tests.Context(b)
	s, _, _ := setupGetCommitGTETimestampTest(ctx, b, 0, true)

	if logsInsertedFirst > 0 {
		populateDatabaseForCommitReportAccepted(ctx, b, s, chainD, chainS1, logsInsertedFirst, 0)
	}

	queryTimestamp := time.Now()

	if logsInsertedMatching > 0 {
		populateDatabaseForCommitReportAccepted(ctx, b, s, chainD, chainS1, logsInsertedMatching, logsInsertedFirst)
	}

	// Reset timer to measure only the query time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		reports, err := s.reader.CommitReportsGTETimestamp(ctx, chainD, queryTimestamp, logsInsertedFirst)
		require.NoError(b, err)
		require.Len(b, reports, logsInsertedFirst)
	}
}

func populateDatabaseForCommitReportAccepted(
	ctx context.Context,
	b *testing.B,
	testEnv *testSetupData,
	destChain cciptypes.ChainSelector,
	sourceChain cciptypes.ChainSelector,
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
		// For first set of logs, set timestamp to 1 hour ago
		timestamp = time.Now().Add(-1 * time.Hour)
	} else {
		// For matching logs, use current time
		timestamp = time.Now()
	}

	for i := 0; i < numOfReports; i++ {
		// Calculate unique BlockNumber and LogIndex
		blockNumber := int64(offset + i + 1) // Offset ensures unique block numbers
		logIndex := int64(offset + i + 1)    // Offset ensures unique log indices

		// Simulate merkleRoots
		merkleRoots := []offramp.InternalMerkleRoot{
			{
				SourceChainSelector: uint64(sourceChain),
				OnRampAddress:       utils.RandomAddress().Bytes(),
				// #nosec G115
				MinSeqNr: uint64(i * 100),
				// #nosec G115
				MaxSeqNr:   uint64(i*100 + 99),
				MerkleRoot: utils.RandomBytes32(),
			},
		}

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
		encodedData, err := commitReportEvent.Inputs.Pack(merkleRoots, priceUpdates)
		require.NoError(b, err)

		// Topics (first one is the event signature)
		topics := [][]byte{
			commitReportEventSig[:],
		}

		// Create log entry
		logs = append(logs, logpoller.Log{
			EvmChainId:     ubig.New(new(big.Int).SetUint64(uint64(destChain))),
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
	require.NoError(b, testEnv.orm.InsertLogs(ctx, logs))
	require.NoError(b, testEnv.orm.InsertBlock(ctx, utils.RandomHash(), int64(offset+numOfReports), timestamp, int64(offset+numOfReports)))
}

// Benchmark Results:
// Benchmark_CCIPReader_ExecutedMessageRanges/LogsInserted_0_StartSeq_0_EndSeq_10-14               13599            93414 ns/op           43389 B/op        654 allocs/op
// Benchmark_CCIPReader_ExecutedMessageRanges/LogsInserted_10_StartSeq_10_EndSeq_20-14             13471            88392 ns/op           43011 B/op        651 allocs/op
// Benchmark_CCIPReader_ExecutedMessageRanges/LogsInserted_10_StartSeq_0_EndSeq_9-14                2799           473396 ns/op          303737 B/op       4535 allocs/op
// Benchmark_CCIPReader_ExecutedMessageRanges/LogsInserted_100_StartSeq_0_EndSeq_100-14              438          2724414 ns/op         2477573 B/op      37468 allocs/op
// Benchmark_CCIPReader_ExecutedMessageRanges/LogsInserted_100000_StartSeq_99744_EndSeq_100000-14     40         29118796 ns/op        12607995 B/op     179396 allocs/op
func Benchmark_CCIPReader_ExecutedMessageRanges(b *testing.B) {
	tests := []struct {
		logsInserted int
		startSeqNum  cciptypes.SeqNum
		endSeqNum    cciptypes.SeqNum
	}{
		{0, 0, 10},                        // no logs
		{10, 10, 20},                      // out of bounds
		{10, 0, 9},                        // get all messages with 10 logs
		{100, 0, 100},                     // get all messages with 100 logs
		{100_000, 100_000 - 256, 100_000}, // get the last 256 messages
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("LogsInserted_%d_StartSeq_%d_EndSeq_%d", tt.logsInserted, tt.startSeqNum, tt.endSeqNum), func(b *testing.B) {
			benchmarkExecutedMessageRanges(b, tt.logsInserted, tt.startSeqNum, tt.endSeqNum)
		})
	}
}

func benchmarkExecutedMessageRanges(b *testing.B, logsInsertedFirst int, startSeqNum, endSeqNum cciptypes.SeqNum) {
	// Initialize test setup
	ctx := tests.Context(b)
	s := setupExecutedMessageRangesTest(ctx, b, true)
	expectedRangeLen := calculateExpectedRangeLen(logsInsertedFirst, startSeqNum, endSeqNum)

	// Insert logs in two phases based on parameters
	if logsInsertedFirst > 0 {
		populateDatabaseForExecutionStateChanged(ctx, b, s, chainS1, chainD, logsInsertedFirst, 0)
	}

	// Reset timer to measure only the query time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		executedRanges, err := s.reader.ExecutedMessageRanges(
			ctx,
			chainS1,
			chainD,
			cciptypes.NewSeqNumRange(startSeqNum, endSeqNum),
		)
		require.NoError(b, err)
		require.Len(b, executedRanges, expectedRangeLen)
	}
}

func populateDatabaseForExecutionStateChanged(
	ctx context.Context,
	b *testing.B,
	testEnv *testSetupData,
	sourceChain cciptypes.ChainSelector,
	destChain cciptypes.ChainSelector,
	numOfEvents int,
	offset int,
) {
	var logs []logpoller.Log
	executionStateEvent, exists := offrampABI.Events[consts.EventNameExecutionStateChanged]
	require.True(b, exists, "Event ExecutionStateChanged not found in ABI")

	executionStateEventSig := executionStateEvent.ID
	executionStateEventAddress := testEnv.contractAddr

	for i := 0; i < numOfEvents; i++ {
		// Calculate unique BlockNumber and LogIndex
		blockNumber := int64(offset + i + 1) // Offset ensures unique block numbers
		logIndex := int64(offset + i + 1)    // Offset ensures unique log indices

		// Populate fields for the event
		sourceChainSelector := uint64(sourceChain)
		// #nosec G115
		sequenceNumber := uint64(offset + i)
		messageID := utils.NewHash()
		messageHash := utils.NewHash()
		state := uint8(1)
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
			EvmChainId:     ubig.New(big.NewInt(0).SetUint64(uint64(destChain))),
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
	require.NoError(b, testEnv.orm.InsertLogs(ctx, logs))
	require.NoError(b, testEnv.orm.InsertBlock(ctx, utils.RandomHash(), int64(offset+numOfEvents), time.Now(), int64(offset+numOfEvents)))
}

// Benchmark Results:
// Benchmark_CCIPReader_MessageSentRanges/LogsInserted_0_StartSeq_0_EndSeq_10-14                     13729             85838 ns/op           43473 B/op        647 allocs/op
// Benchmark_CCIPReader_MessageSentRanges/LogsInserted_10_StartSeq_0_EndSeq_9-14                      870           1405208 ns/op         1156315 B/op      21102 allocs/op
// Benchmark_CCIPReader_MessageSentRanges/LogsInserted_100_StartSeq_0_EndSeq_100-14                    90          12129488 ns/op        10833395 B/op     201076 allocs/op
// Benchmark_CCIPReader_MessageSentRanges/LogsInserted_100000_StartSeq_99744_EndSeq_100000-14          10         105741438 ns/op        49103282 B/op     796213 allocs/op
func Benchmark_CCIPReader_MessageSentRanges(b *testing.B) {
	tests := []struct {
		logsInserted int
		startSeqNum  cciptypes.SeqNum
		endSeqNum    cciptypes.SeqNum
	}{
		{0, 0, 10},                        // No logs
		{10, 0, 9},                        // Get all messages with 10 logs
		{100, 0, 100},                     // Get all messages with 100 logs
		{100_000, 100_000 - 256, 100_000}, // Get the last 256 messages
	}

	for _, tt := range tests {
		b.Run(fmt.Sprintf("LogsInserted_%d_StartSeq_%d_EndSeq_%d", tt.logsInserted, tt.startSeqNum, tt.endSeqNum), func(b *testing.B) {
			benchmarkMessageSentRanges(b, tt.logsInserted, tt.startSeqNum, tt.endSeqNum)
		})
	}
}

func benchmarkMessageSentRanges(b *testing.B, logsInserted int, startSeqNum, endSeqNum cciptypes.SeqNum) {
	// Initialize test setup
	ctx := tests.Context(b)
	s := setupMsgsBetweenSeqNumsTest(ctx, b, true)
	expectedRangeLen := calculateExpectedRangeLen(logsInserted, startSeqNum, endSeqNum)

	err := s.extendedCR.Bind(ctx, []types.BoundContract{
		{
			Address: s.contractAddr.String(),
			Name:    consts.ContractNameOnRamp,
		},
	})
	require.NoError(b, err)

	// Insert logs if needed
	if logsInserted > 0 {
		populateDatabaseForMessageSent(ctx, b, s, chainS1, chainD, logsInserted, 0)
	}

	// Reset timer to measure only the query time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msgs, err := s.reader.MsgsBetweenSeqNums(
			ctx,
			chainS1,
			cciptypes.NewSeqNumRange(startSeqNum, endSeqNum),
		)
		require.NoError(b, err)
		require.Len(b, msgs, expectedRangeLen)
	}
}

func populateDatabaseForMessageSent(
	ctx context.Context,
	b *testing.B,
	testEnv *testSetupData,
	sourceChain cciptypes.ChainSelector,
	destChain cciptypes.ChainSelector,
	numOfEvents int,
	offset int,
) {
	var logs []logpoller.Log
	messageSentEvent, exists := onrampABI.Events[consts.EventNameCCIPMessageSent]
	require.True(b, exists, "Event CCIPMessageSent not found in ABI")

	messageSentEventSig := messageSentEvent.ID
	messageSentEventAddress := testEnv.contractAddr

	for i := 0; i < numOfEvents; i++ {
		// Calculate unique BlockNumber and LogIndex
		blockNumber := int64(offset + i + 1) // Offset ensures unique block numbers
		logIndex := int64(offset + i + 1)    // Offset ensures unique log indices

		// Populate fields for the event
		destChainSelector := uint64(destChain)
		// #nosec G115
		sequenceNumber := uint64(offset + i)

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

		// Create InternalEVM2AnyRampMessage struct
		message := onramp.InternalEVM2AnyRampMessage{
			Header:    header,
			Sender:    utils.RandomAddress(),
			Data:      []byte{0x04, 0x05},
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
			EvmChainId:     ubig.New(big.NewInt(0).SetUint64(uint64(sourceChain))),
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
	require.NoError(b, testEnv.orm.InsertLogs(ctx, logs))
	require.NoError(b, testEnv.orm.InsertBlock(ctx, utils.RandomHash(), int64(offset+numOfEvents), time.Now(), int64(offset+numOfEvents)))
}

func calculateExpectedRangeLen(logsInserted int, startSeq, endSeq cciptypes.SeqNum) int {
	if logsInserted == 0 {
		return 0
	}
	start := uint64(startSeq)
	end := uint64(endSeq)
	// #nosec G115
	logs := uint64(logsInserted)

	if start >= logs {
		return 0
	}

	if end >= logs {
		end = logs - 1
	}

	// #nosec G115
	return int(end - start + 1)
}

func setupSimulatedBackendAndAuth(t testing.TB) (*simulated.Backend, *bind.TransactOpts) {
	privateKey, err := crypto.GenerateKey()
	require.NoError(t, err)

	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	require.True(t, ok)

	alloc := map[common.Address]ethtypes.Account{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := simulated.NewBackend(alloc, simulated.WithBlockGasLimit(8000000))

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	require.NoError(t, err)
	auth.GasLimit = uint64(6000000)

	return simulatedBackend, auth
}

func testSetupRealContracts(
	ctx context.Context,
	t *testing.T,
	destChain uint64,
	toBindContracts map[cciptypes.ChainSelector][]types.BoundContract,
	toMockBindings map[cciptypes.ChainSelector][]types.BoundContract,
	env testhelpers.DeployedEnv,
) ccipreaderpkg.CCIPReader {
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            0,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)

	var crs = make(map[cciptypes.ChainSelector]contractreader.Extended)
	for chain, bindings := range toBindContracts {
		be := env.Env.Chains[uint64(chain)].Client.(*memory.Backend)
		cl := client.NewSimulatedBackendClient(t, be.Sim, big.NewInt(0).SetUint64(uint64(chain)))
		headTracker := headtracker.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
		lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(0).SetUint64(uint64(chain)), db, lggr),
			cl,
			lggr,
			headTracker,
			lpOpts,
		)
		require.NoError(t, lp.Start(ctx))

		var cfg evmtypes.ChainReaderConfig
		if chain == cs(destChain) {
			cfg = evmconfig.DestReaderConfig
		} else {
			cfg = evmconfig.SourceReaderConfig
		}
		cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, cfg)
		require.NoError(t, err)

		extendedCr2 := contractreader.NewExtendedContractReader(cr)
		err = extendedCr2.Bind(ctx, bindings)
		require.NoError(t, err)
		crs[cciptypes.ChainSelector(chain)] = extendedCr2

		err = cr.Start(ctx)
		require.NoError(t, err)

		t.Cleanup(func() {
			require.NoError(t, cr.Close())
			require.NoError(t, lp.Close())
			require.NoError(t, db.Close())
		})
	}

	for chain, bindings := range toMockBindings {
		if _, ok := crs[chain]; ok {
			require.False(t, ok, "chain %d already exists", chain)
		}
		m := readermocks.NewMockContractReaderFacade(t)
		m.EXPECT().Bind(ctx, bindings).Return(nil)
		ecr := contractreader.NewExtendedContractReader(m)
		err := ecr.Bind(ctx, bindings)
		require.NoError(t, err)
		crs[chain] = ecr
	}

	contractReaders := map[cciptypes.ChainSelector]contractreader.Extended{}
	for chain, cr := range crs {
		contractReaders[chain] = cr
	}
	contractWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)
	edc := ccipevm.NewExtraArgsCodec()
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, cciptypes.ChainSelector(destChain), nil, edc)

	return reader
}

func testSetup(
	ctx context.Context,
	t testing.TB,
	params testSetupParams,
) *testSetupData {
	address, _, _, err := ccip_reader_tester.DeployCCIPReaderTester(params.Auth, params.SimulatedBackend.Client())
	assert.NoError(t, err)
	params.SimulatedBackend.Commit()

	// Setup contract client
	contract, err := ccip_reader_tester.NewCCIPReaderTester(address, params.SimulatedBackend.Client())
	assert.NoError(t, err)

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	// Parameterize database selection
	var db *sqlx.DB
	if params.UseHeavyDB {
		_, db = heavyweight.FullTestDBV2(t, nil) // Heavyweight database for benchmarks
	} else {
		db = pgtest.NewSqlxDB(t) // Simple in-memory DB for tests
	}
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            params.FinalityDepth,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, params.SimulatedBackend, big.NewInt(0).SetUint64(uint64(params.ReaderChain)))
	headTracker := headtracker.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	orm := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(params.ReaderChain)), db, lggr)
	lp := logpoller.NewLogPoller(
		orm,
		cl,
		lggr,
		headTracker,
		lpOpts,
	)
	assert.NoError(t, lp.Start(ctx))

	for sourceChain, seqNum := range params.OnChainSeqNums {
		_, err1 := contract.SetSourceChainConfig(params.Auth, uint64(sourceChain), ccip_reader_tester.OffRampSourceChainConfig{
			IsEnabled: true,
			MinSeqNr:  uint64(seqNum),
			OnRamp:    utils.RandomAddress().Bytes(),
		})
		assert.NoError(t, err1)
		params.SimulatedBackend.Commit()
		scc, err1 := contract.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, uint64(sourceChain))
		assert.NoError(t, err1)
		assert.Equal(t, seqNum, cciptypes.SeqNum(scc.MinSeqNr))
	}

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, params.Cfg)
	require.NoError(t, err)

	extendedCr := contractreader.NewExtendedContractReader(cr)

	if params.BindTester {
		err = extendedCr.Bind(ctx, []types.BoundContract{
			{
				Address: address.String(),
				Name:    params.ContractNameToBind,
			},
		})
		require.NoError(t, err)
	}

	var otherCrs = make(map[cciptypes.ChainSelector]contractreader.Extended)
	for chain, bindings := range params.ToBindContracts {
		cl2 := client.NewSimulatedBackendClient(t, params.SimulatedBackend, big.NewInt(0).SetUint64(uint64(chain)))
		headTracker2 := headtracker.NewSimulatedHeadTracker(cl2, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
		lp2 := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(0).SetUint64(uint64(chain)), db, lggr),
			cl2,
			lggr,
			headTracker2,
			lpOpts,
		)
		require.NoError(t, lp2.Start(ctx))

		cr2, err2 := evm.NewChainReaderService(ctx, lggr, lp2, headTracker2, cl2, params.Cfg)
		require.NoError(t, err2)

		extendedCr2 := contractreader.NewExtendedContractReader(cr2)
		err2 = extendedCr2.Bind(ctx, bindings)
		require.NoError(t, err2)
		otherCrs[chain] = extendedCr2
	}

	for chain, bindings := range params.ToMockBindings {
		if _, ok := otherCrs[chain]; ok {
			require.False(t, ok, "chain %d already exists", chain)
		}
		m := readermocks.NewMockContractReaderFacade(t)
		m.EXPECT().Bind(ctx, bindings).Return(nil)
		ecr := contractreader.NewExtendedContractReader(m)
		err = ecr.Bind(ctx, bindings)
		require.NoError(t, err)
		otherCrs[chain] = ecr
	}

	err = cr.Start(ctx)
	require.NoError(t, err)

	contractReaders := map[cciptypes.ChainSelector]contractreader.Extended{params.ReaderChain: extendedCr}
	for chain, cr := range otherCrs {
		contractReaders[chain] = cr
	}
	contractWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)
	edc := ccipevm.NewExtraArgsCodec()
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, params.DestChain, nil, edc)

	t.Cleanup(func() {
		require.NoError(t, cr.Close())
		require.NoError(t, lp.Close())
		require.NoError(t, db.Close())
	})

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           params.SimulatedBackend,
		auth:         params.Auth,
		orm:          orm,
		lp:           lp,
		cl:           cl,
		reader:       reader,
		extendedCR:   extendedCr,
	}
}

type testSetupParams struct {
	ReaderChain        cciptypes.ChainSelector
	DestChain          cciptypes.ChainSelector
	OnChainSeqNums     map[cciptypes.ChainSelector]cciptypes.SeqNum
	Cfg                evmtypes.ChainReaderConfig
	ToBindContracts    map[cciptypes.ChainSelector][]types.BoundContract
	ToMockBindings     map[cciptypes.ChainSelector][]types.BoundContract
	BindTester         bool
	ContractNameToBind string
	SimulatedBackend   *simulated.Backend
	Auth               *bind.TransactOpts
	FinalityDepth      int64
	UseHeavyDB         bool
}

type testSetupData struct {
	contractAddr common.Address
	contract     *ccip_reader_tester.CCIPReaderTester
	sb           *simulated.Backend
	auth         *bind.TransactOpts
	orm          logpoller.ORM
	lp           logpoller.LogPoller
	cl           client.Client
	reader       ccipreaderpkg.CCIPReader
	extendedCR   contractreader.Extended
}

func cs(i uint64) cciptypes.ChainSelector {
	return cciptypes.ChainSelector(i)
}
