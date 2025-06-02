package ccip

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	chain_selectors "github.com/smartcontractkit/chain-selectors"

	readermocks "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/contractreader"
	typepkgmock "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-ccip/plugintypes"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/types/query/primitives"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"
	cldf_chain "github.com/smartcontractkit/chainlink-deployments-framework/chain"

	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset/testhelpers"
	"github.com/smartcontractkit/chainlink/deployment/ccip/shared/stateview"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/pgtest"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/client"
	"github.com/smartcontractkit/chainlink-evm/pkg/heads/headstest"
	"github.com/smartcontractkit/chainlink-evm/pkg/logpoller"
	evmchaintypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink-evm/pkg/utils"
	ubig "github.com/smartcontractkit/chainlink-evm/pkg/utils/big"

	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_0_0/rmn_proxy_contract"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/ccip_reader_tester"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/fee_quoter"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/offramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/onramp"
	"github.com/smartcontractkit/chainlink-ccip/chains/evm/gobindings/generated/v1_6_0/rmn_remote"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

const (
	chainS1   = cciptypes.ChainSelector(1)
	chainS2   = cciptypes.ChainSelector(2)
	chainS3   = cciptypes.ChainSelector(3)
	chainD    = cciptypes.ChainSelector(4)
	chainSEVM = cciptypes.ChainSelector(5009297550715157269)
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
			chainS2: {
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

func setupExecutedMessagesTest(ctx context.Context, t testing.TB, useHeavyDB bool) *testSetupData {
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

func setupMsgsBetweenSeqNumsTest(ctx context.Context, t testing.TB, useHeavyDB bool, sourceChainSel cciptypes.ChainSelector) *testSetupData {
	sb, auth := setupSimulatedBackendAndAuth(t)
	return testSetup(ctx, t, testSetupParams{
		ReaderChain:        sourceChainSel,
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
			BlessedMerkleRoots: []ccip_reader_tester.InternalMerkleRoot{
				{
					SourceChainSelector: uint64(chainS1),
					MinSeqNr:            10,
					MaxSeqNr:            20,
					MerkleRoot:          [32]byte{i + 1},
					OnRampAddress:       common.LeftPadBytes(onRampAddress.Bytes(), 32),
				},
			},
			UnblessedMerkleRoots: []ccip_reader_tester.InternalMerkleRoot{
				{
					SourceChainSelector: uint64(chainS2),
					MinSeqNr:            20,
					MaxSeqNr:            30,
					MerkleRoot:          [32]byte{i + 2},
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

func TestCCIPReader_GetRMNRemoteConfig(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	sb, auth := setupSimulatedBackendAndAuth(t)

	rmnRemoteAddr, _, _, err := rmn_remote.DeployRMNRemote(auth, sb.Client(), uint64(chainD), utils.RandomAddress())
	require.NoError(t, err)
	sb.Commit()

	proxyAddr, _, _, err := rmn_proxy_contract.DeployRMNProxy(auth, sb.Client(), rmnRemoteAddr)
	require.NoError(t, err)
	sb.Commit()

	t.Logf("Proxy address: %s, rmn remote address: %s", proxyAddr.Hex(), rmnRemoteAddr.Hex())

	proxy, err := rmn_proxy_contract.NewRMNProxy(proxyAddr, sb.Client())
	require.NoError(t, err)

	currARM, err := proxy.GetARM(&bind.CallOpts{
		Context: ctx,
	})
	require.NoError(t, err)
	require.Equal(t, currARM, rmnRemoteAddr)

	rmnRemote, err := rmn_remote.NewRMNRemote(rmnRemoteAddr, sb.Client())
	require.NoError(t, err)

	_, err = rmnRemote.SetConfig(auth, rmn_remote.RMNRemoteConfig{
		RmnHomeContractConfigDigest: utils.RandomBytes32(),
		Signers: []rmn_remote.RMNRemoteSigner{
			{
				OnchainPublicKey: utils.RandomAddress(),
				NodeIndex:        0,
			},
			{
				OnchainPublicKey: utils.RandomAddress(),
				NodeIndex:        1,
			},
			{
				OnchainPublicKey: utils.RandomAddress(),
				NodeIndex:        2,
			},
		},
		FSign: 1, // 2*FSign + 1 == 3
	})
	require.NoError(t, err)
	sb.Commit()

	db := pgtest.NewSqlxDB(t)
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            1,
		BackfillBatchSize:        10,
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, sb, big.NewInt(1337))
	headTracker := headstest.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
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
			Address: proxyAddr.String(),
			Name:    consts.ContractNameRMNRemote,
		},
	})
	require.NoError(t, err)

	mockAddrCodec := newMockAddressCodec(t)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(
		ctx,
		lggr,
		map[cciptypes.ChainSelector]contractreader.Extended{
			chainD: extendedCr,
		},
		nil,
		chainD,
		rmnRemoteAddr.Bytes(),
		mockAddrCodec,
	)

	exp, err := rmnRemote.GetVersionedConfig(&bind.CallOpts{
		Context: ctx,
	})
	require.NoError(t, err)

	rmnRemoteConfig, err := reader.GetRMNRemoteConfig(ctx)
	require.NoError(t, err)
	require.Equal(t, exp.Config.RmnHomeContractConfigDigest[:], rmnRemoteConfig.ConfigDigest[:])
	require.Equal(t, len(exp.Config.Signers), len(rmnRemoteConfig.Signers))
	for i, signer := range exp.Config.Signers {
		require.Equal(t, signer.OnchainPublicKey.Bytes(), []byte(rmnRemoteConfig.Signers[i].OnchainPublicKey))
		require.Equal(t, signer.NodeIndex, rmnRemoteConfig.Signers[i].NodeIndex)
	}
	require.Equal(t, exp.Config.FSign, rmnRemoteConfig.FSign)
}

func TestCCIPReader_GetOffRampConfigDigest(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
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
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, sb, big.NewInt(1337))
	headTracker := headstest.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
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
	mokAddrCodec := newMockAddressCodec(t)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(
		ctx,
		lggr,
		map[cciptypes.ChainSelector]contractreader.Extended{
			chainD: extendedCr,
		},
		nil,
		chainD,
		addr.Bytes(),
		mokAddrCodec,
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
	ctx := t.Context()
	s, _, onRampAddress := setupGetCommitGTETimestampTest(ctx, t, 0, false)

	tokenA := common.HexToAddress("123")
	const numReports = 5

	firstReportTs := emitCommitReports(ctx, t, s, numReports, tokenA, onRampAddress)

	iter, err := s.contract.FilterCommitReportAccepted(&bind.FilterOpts{
		Start: 0,
	})
	require.NoError(t, err)
	var onchainEvents []*ccip_reader_tester.CCIPReaderTesterCommitReportAccepted
	for iter.Next() {
		onchainEvents = append(onchainEvents, iter.Event)
	}
	require.Len(t, onchainEvents, numReports)
	sort.Slice(onchainEvents, func(i, j int) bool {
		return onchainEvents[i].Raw.BlockNumber < onchainEvents[j].Raw.BlockNumber
	})

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var ccipReaderReports []plugintypes.CommitPluginReportWithMeta
	require.Eventually(t, func() bool {
		var err2 error
		ccipReaderReports, err2 = s.reader.CommitReportsGTETimestamp(
			ctx,
			// Skips first report
			//nolint:gosec // this won't overflow
			time.Unix(int64(firstReportTs)+1, 0),
			primitives.Unconfirmed,
			10,
		)
		require.NoError(t, err2)
		return len(ccipReaderReports) == numReports-1
	}, 30*time.Second, 50*time.Millisecond)

	// trim the first report to simulate the timestamp filter above.
	onchainEvents = onchainEvents[1:]
	require.Len(t, onchainEvents, numReports-1)

	require.Len(t, ccipReaderReports, numReports-1)
	for i := range onchainEvents {
		// check blessed roots are deserialized correctly
		requireEqualRoots(t, onchainEvents[i].BlessedMerkleRoots, ccipReaderReports[i].Report.BlessedMerkleRoots)

		// check unblessed roots are deserialized correctly
		requireEqualRoots(t, onchainEvents[i].UnblessedMerkleRoots, ccipReaderReports[i].Report.UnblessedMerkleRoots)

		// check price updates are deserialized correctly
		requireEqualPriceUpdates(t, onchainEvents[i].PriceUpdates, ccipReaderReports[i].Report.PriceUpdates)
	}
}

func requireEqualPriceUpdates(
	t *testing.T,
	onchainPriceUpdates ccip_reader_tester.InternalPriceUpdates,
	ccipReaderPriceUpdates cciptypes.PriceUpdates,
) {
	// token price update equality
	require.Equal(t, len(onchainPriceUpdates.TokenPriceUpdates), len(ccipReaderPriceUpdates.TokenPriceUpdates))
	for i := range onchainPriceUpdates.TokenPriceUpdates {
		require.Equal(t,
			onchainPriceUpdates.TokenPriceUpdates[i].SourceToken.Bytes(),
			hexutil.MustDecode(string(ccipReaderPriceUpdates.TokenPriceUpdates[i].TokenID)))
		require.Equal(t,
			onchainPriceUpdates.TokenPriceUpdates[i].UsdPerToken,
			ccipReaderPriceUpdates.TokenPriceUpdates[i].Price.Int)
	}

	// gas price update equality
	require.Equal(t, len(onchainPriceUpdates.GasPriceUpdates), len(ccipReaderPriceUpdates.GasPriceUpdates))
	for i := range onchainPriceUpdates.GasPriceUpdates {
		require.Equal(t,
			onchainPriceUpdates.GasPriceUpdates[i].DestChainSelector,
			uint64(ccipReaderPriceUpdates.GasPriceUpdates[i].ChainSel))
		require.Equal(t,
			onchainPriceUpdates.GasPriceUpdates[i].UsdPerUnitGas,
			ccipReaderPriceUpdates.GasPriceUpdates[i].GasPrice.Int)
	}
}

func requireEqualRoots(
	t *testing.T,
	onchainRoots []ccip_reader_tester.InternalMerkleRoot,
	ccipReaderRoots []cciptypes.MerkleRootChain,
) {
	require.Equal(t, len(onchainRoots), len(ccipReaderRoots))
	for i := 0; i < len(onchainRoots); i++ {
		require.Equal(t,
			onchainRoots[i].SourceChainSelector,
			uint64(ccipReaderRoots[i].ChainSel),
		)

		// onchain emits the padded address but ccip reader currently sets the unpadded address
		// TODO: fix this!
		require.Equal(t,
			onchainRoots[i].OnRampAddress,
			common.LeftPadBytes([]byte(ccipReaderRoots[i].OnRampAddress), 32),
		)
		require.Equal(t,
			onchainRoots[i].MinSeqNr,
			uint64(ccipReaderRoots[i].SeqNumsRange.Start()),
		)
		require.Equal(t,
			onchainRoots[i].MaxSeqNr,
			uint64(ccipReaderRoots[i].SeqNumsRange.End()),
		)
		require.Equal(t,
			onchainRoots[i].MerkleRoot,
			[32]byte(ccipReaderRoots[i].MerkleRoot),
		)
	}
}
func commitSqNrs(
	s *testSetupData,
	chainSel cciptypes.ChainSelector,
	seqNums []uint64,
	state uint8) error {
	for _, sqnr := range seqNums {
		_, err := s.contract.EmitExecutionStateChanged(
			s.auth,
			uint64(chainSel),
			sqnr,
			cciptypes.Bytes32{1, 0, 0, 1},
			cciptypes.Bytes32{1, 0, 0, 1, 1, 0, 0, 1},
			state,
			[]byte{1, 2, 3, 4},
			big.NewInt(250_000),
		)
		if err != nil {
			return err
		}
	}
	return nil
}
func TestCCIPReader_ExecutedMessages_SingleChain(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	s := setupExecutedMessagesTest(ctx, t, false)
	err := commitSqNrs(s, chainS1, []uint64{14}, 1)
	require.NoError(t, err)
	s.sb.Commit()

	err = commitSqNrs(s, chainS1, []uint64{15}, 1)
	require.NoError(t, err)
	s.sb.Commit()

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var executedMsgs map[cciptypes.ChainSelector][]cciptypes.SeqNum
	require.Eventually(t, func() bool {
		executedMsgs, err = s.reader.ExecutedMessages(
			ctx,
			map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{
				chainS1: {
					cciptypes.NewSeqNumRange(14, 15),
				},
			},
			primitives.Unconfirmed,
		)
		require.NoError(t, err)
		return len(executedMsgs[chainS1]) == 2
	}, tests.WaitTimeout(t), 50*time.Millisecond)

	assert.Equal(t, []cciptypes.SeqNum{14, 15}, executedMsgs[chainS1])
}

func TestCCIPReader_ExecutedMessages_MultiChain(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	s := setupExecutedMessagesTest(ctx, t, false)
	err := commitSqNrs(s, chainS1, []uint64{15}, 1)
	require.NoError(t, err)
	s.sb.Commit()

	err = commitSqNrs(s, chainS2, []uint64{15}, 1)
	require.NoError(t, err)
	s.sb.Commit()

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var executedMsgs map[cciptypes.ChainSelector][]cciptypes.SeqNum
	require.Eventually(t, func() bool {
		executedMsgs, err = s.reader.ExecutedMessages(
			ctx,
			map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{
				chainS1: {
					cciptypes.NewSeqNumRange(14, 16),
				},
				chainS2: {
					cciptypes.NewSeqNumRange(15, 15),
				},
				chainS3: {}, // empty, should not affect query
			},
			primitives.Unconfirmed,
		)
		require.NoError(t, err)
		return executedMsgs[chainS1][0] == 15 && executedMsgs[chainS2][0] == 15
	}, tests.WaitTimeout(t), 50*time.Millisecond)

	assert.Equal(t, []cciptypes.SeqNum{15}, executedMsgs[chainS1])
	assert.Len(t, executedMsgs, 2)
	assert.Equal(t, []cciptypes.SeqNum{15}, executedMsgs[chainS1])
	assert.Equal(t, []cciptypes.SeqNum{15}, executedMsgs[chainS2])
}

func TestCCIPReader_ExecutedMessages_MultiChainDisjoint(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	s := setupExecutedMessagesTest(ctx, t, false)
	err := commitSqNrs(s, chainS1, []uint64{15, 17, 70}, 1)
	require.NoError(t, err)
	s.sb.Commit()

	err = commitSqNrs(s, chainS2, []uint64{15, 16}, 1)
	require.NoError(t, err)
	s.sb.Commit()

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var executedMsgs map[cciptypes.ChainSelector][]cciptypes.SeqNum
	require.Eventually(t, func() bool {
		executedMsgs, err = s.reader.ExecutedMessages(
			ctx,
			map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{
				chainS1: {
					cciptypes.NewSeqNumRange(10, 20),
					cciptypes.NewSeqNumRange(70, 70),
				},
				chainS2: {
					cciptypes.NewSeqNumRange(15, 16),
				},
			},
			primitives.Unconfirmed,
		)
		require.NoError(t, err)
		return len(executedMsgs) == 2
	}, tests.WaitTimeout(t), 50*time.Millisecond)

	assert.Len(t, executedMsgs[chainS1], 3)
	assert.Len(t, executedMsgs[chainS2], 2)
	assert.Equal(t, []cciptypes.SeqNum{15, 17, 70}, executedMsgs[chainS1])
	assert.Equal(t, []cciptypes.SeqNum{15, 16}, executedMsgs[chainS2])
}

func TestCCIPReader_MsgsBetweenSeqNums(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

	s := setupMsgsBetweenSeqNumsTest(ctx, t, false, chainSEVM)
	_, err := s.contract.EmitCCIPMessageSent(s.auth, uint64(chainD), ccip_reader_tester.InternalEVM2AnyRampMessage{
		Header: ccip_reader_tester.InternalRampMessageHeader{
			MessageId:           [32]byte{1, 0, 0, 0, 0},
			SourceChainSelector: uint64(chainSEVM),
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
			SourceChainSelector: uint64(chainSEVM),
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
			chainSEVM,
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
		require.Equal(t, chainSEVM, msg.Header.SourceChainSelector)
		require.Equal(t, chainD, msg.Header.DestChainSelector)
	}
}

func TestCCIPReader_NextSeqNum(t *testing.T) {
	t.Parallel()
	ctx := t.Context()

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
	ctx := t.Context()
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)

	var selectors = env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
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
					Address: state.MustGetEVMChainState(srcChain).OnRamp.Address().String(),
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
		msg := testhelpers.DefaultRouterMessage(state.MustGetEVMChainState(destChain).Receiver.Address())
		msgSentEvent := testhelpers.TestSendRequest(t, env.Env, state, srcChain, destChain, false, msg)
		require.Equal(t, uint64(i), msgSentEvent.SequenceNumber)
		require.Equal(t, uint64(i), msgSentEvent.Message.Header.Nonce) // check outbound nonce incremented
		seqNum, err2 := reader.GetExpectedNextSequenceNumber(ctx, cs(srcChain))
		require.NoError(t, err2)
		require.Equal(t, cciptypes.SeqNum(i+1), seqNum)
	}
}

func TestCCIPReader_Nonces(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
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

	request := make(map[cciptypes.ChainSelector][]string)
	for chain, addresses := range nonces {
		request[chain] = make([]string, 0, len(addresses))
		for address := range addresses {
			request[chain] = append(request[chain], address.String())
		}
		request[chain] = append(request[chain], utils.RandomAddress().String())
	}

	results, err := s.reader.Nonces(ctx, request)
	require.NoError(t, err)

	for chain, addresses := range nonces {
		assert.Len(t, results[chain], len(request[chain]))
		for address, nonce := range addresses {
			assert.Equal(t, nonce, results[chain][address.String()])
		}
	}
}

func TestCCIPReader_GetContractAddress(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	sb, auth := setupSimulatedBackendAndAuth(t)

	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:        chainS1,
		DestChain:          chainD,
		OnChainSeqNums:     nil,
		Cfg:                evmconfig.DestReaderConfig,
		BindTester:         true,
		ContractNameToBind: consts.ContractNameOffRamp,
		SimulatedBackend:   sb,
		Auth:               auth,
		UseHeavyDB:         false,
	})

	t.Run("success - single bound address", func(t *testing.T) {
		myContractName := consts.ContractNameOffRamp
		myAddress := s.contractAddr

		err := s.extendedCR.Bind(ctx, []types.BoundContract{
			{
				Address: myAddress.String(),
				Name:    myContractName,
			},
		})
		require.NoError(t, err)

		gotBytes, err := s.reader.GetContractAddress(myContractName, chainS1)
		require.NoError(t, err)

		require.Equal(t, myAddress.Bytes(), gotBytes, "expected the bound contract address to match")
	})

	t.Run("error - no bindings found", func(t *testing.T) {
		_, err := s.reader.GetContractAddress("UnboundContract", chainS1)
		require.Error(t, err)
		require.Contains(t, err.Error(), "expected one binding for the UnboundContract contract, got 0")
	})

	t.Run("success - multiple bindings, return override binding", func(t *testing.T) {
		myContractName := consts.ContractNameOffRamp
		addr1 := s.contractAddr
		addr2, _, _, err := ccip_reader_tester.DeployCCIPReaderTester(auth, sb.Client())
		require.NoError(t, err)
		sb.Commit()

		err = s.extendedCR.Bind(ctx, []types.BoundContract{
			{
				Address: addr1.String(),
				Name:    myContractName,
			},
			{
				Address: addr2.String(),
				Name:    myContractName,
			},
		})
		require.NoError(t, err)

		gotBytes, err := s.reader.GetContractAddress(myContractName, chainS1)
		require.NoError(t, err)

		require.Equal(t, addr2.Bytes(), gotBytes, "expected the bound contract override address to match")
	})

	t.Run("error - chain not supported", func(t *testing.T) {
		// Suppose chainS2 is not set up in this test environment (no contract reader).
		// The call should fail with "contract reader not found for chain".
		_, err := s.reader.GetContractAddress("TestContract", chainS2)
		require.Error(t, err)
		require.Contains(t, err.Error(), "contract reader not found for chain 2")
	})
}

func TestCCIPReader_DiscoverContracts(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	sb, auth := setupSimulatedBackendAndAuth(t)

	//--------------------------------Setup--------------------------------//
	onRampS1StaticConfig := onramp.OnRampStaticConfig{
		ChainSelector:      uint64(chainS1),
		RmnRemote:          utils.RandomAddress(),
		NonceManager:       utils.RandomAddress(),
		TokenAdminRegistry: utils.RandomAddress(),
	}

	onRampS1DynamicConfig := onramp.OnRampDynamicConfig{
		FeeQuoter:              utils.RandomAddress(),
		ReentrancyGuardEntered: false,
		MessageInterceptor:     utils.ZeroAddress,
		FeeAggregator:          utils.RandomAddress(),
		AllowlistAdmin:         utils.RandomAddress(),
	}

	destinationChainConfigArgs := []onramp.OnRampDestChainConfigArgs{
		{
			DestChainSelector: uint64(chainD),
			Router:            utils.RandomAddress(),
			AllowlistEnabled:  false,
		},
	}
	onRampS1Addr, _, _, err := onramp.DeployOnRamp(auth, sb.Client(), onRampS1StaticConfig, onRampS1DynamicConfig, destinationChainConfigArgs)
	require.NoError(t, err)
	sb.Commit()

	offRampDStaticConfig := offramp.OffRampStaticConfig{
		ChainSelector:        uint64(chainD),
		GasForCallExactCheck: 0,
		RmnRemote:            utils.RandomAddress(),
		TokenAdminRegistry:   utils.RandomAddress(),
		NonceManager:         utils.RandomAddress(),
	}

	offRampDDynamicConfig := offramp.OffRampDynamicConfig{
		FeeQuoter:                               utils.RandomAddress(),
		PermissionLessExecutionThresholdSeconds: 1,
		MessageInterceptor:                      utils.ZeroAddress,
	}

	offRampDSourceChainConfigArgs := []offramp.OffRampSourceChainConfigArgs{
		{
			Router:                    destinationChainConfigArgs[0].Router,
			SourceChainSelector:       onRampS1StaticConfig.ChainSelector,
			IsEnabled:                 true,
			IsRMNVerificationDisabled: true,
			OnRamp:                    common.LeftPadBytes(onRampS1Addr.Bytes(), 32),
		},
	}
	offRampDestAddr, _, _, err := offramp.DeployOffRamp(auth, sb.Client(), offRampDStaticConfig, offRampDDynamicConfig, offRampDSourceChainConfigArgs)
	require.NoError(t, err)
	sb.Commit()

	clS1 := client.NewSimulatedBackendClient(t, sb, big.NewInt(0).SetUint64(uint64(chainS1)))
	headTrackerS1 := headstest.NewSimulatedHeadTracker(clS1, true, 1)
	ormS1 := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(chainS1)), pgtest.NewSqlxDB(t), logger.TestLogger(t))
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            0,
		BackfillBatchSize:        10,
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	lpS1 := logpoller.NewLogPoller(
		ormS1,
		clS1,
		logger.TestLogger(t),
		headTrackerS1,
		lpOpts,
	)
	require.NoError(t, lpS1.Start(ctx))

	clD := client.NewSimulatedBackendClient(t, sb, big.NewInt(0).SetUint64(uint64(chainD)))
	headTrackerD := headstest.NewSimulatedHeadTracker(clD, true, 1)
	ormD := logpoller.NewORM(big.NewInt(0).SetUint64(uint64(chainD)), pgtest.NewSqlxDB(t), logger.TestLogger(t))
	lpD := logpoller.NewLogPoller(
		ormD,
		clD,
		logger.TestLogger(t),
		headTrackerD,
		lpOpts,
	)
	require.NoError(t, lpD.Start(ctx))

	crS1, err := evm.NewChainReaderService(ctx, logger.TestLogger(t), lpS1, headTrackerS1, clS1, evmconfig.SourceReaderConfig)
	require.NoError(t, err)
	extendedCrS1 := contractreader.NewExtendedContractReader(crS1)

	crD, err := evm.NewChainReaderService(ctx, logger.TestLogger(t), lpD, headTrackerD, clD, evmconfig.DestReaderConfig)
	require.NoError(t, err)
	extendedCrD := contractreader.NewExtendedContractReader(crD)
	err = extendedCrD.Bind(ctx, []types.BoundContract{
		{
			Address: offRampDestAddr.String(),
			Name:    consts.ContractNameOffRamp,
		},
	})
	require.NoError(t, err)

	err = crS1.Start(ctx)
	require.NoError(t, err)
	err = crD.Start(ctx)
	require.NoError(t, err)

	contractReaders := map[cciptypes.ChainSelector]contractreader.Extended{}
	contractReaders[chainS1] = extendedCrS1
	contractReaders[chainD] = extendedCrD

	contractWriters := make(map[cciptypes.ChainSelector]types.ContractWriter)

	mokAddrCodec := newMockAddressCodec(t)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, logger.TestLogger(t), contractReaders, contractWriters, chainD, offRampDestAddr.Bytes(), mokAddrCodec)

	t.Cleanup(func() {
		assert.NoError(t, crS1.Close())
		assert.NoError(t, lpS1.Close())
		assert.NoError(t, crD.Close())
		assert.NoError(t, lpD.Close())
	})
	//--------------------------------Setup done--------------------------------//

	// Call the ccip chain reader with DiscoverContracts for test
	contractAddresses, err := reader.DiscoverContracts(ctx, []cciptypes.ChainSelector{chainS1, chainD})
	require.NoError(t, err)

	require.Equal(t, contractAddresses[consts.ContractNameOnRamp][chainS1], cciptypes.UnknownAddress(common.LeftPadBytes(onRampS1Addr.Bytes(), 32)))
	require.Equal(t, contractAddresses[consts.ContractNameRouter][chainD], cciptypes.UnknownAddress(destinationChainConfigArgs[0].Router.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameRMNRemote][chainD], cciptypes.UnknownAddress(offRampDStaticConfig.RmnRemote.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameNonceManager][chainD], cciptypes.UnknownAddress(offRampDStaticConfig.NonceManager.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameFeeQuoter][chainD], cciptypes.UnknownAddress(offRampDDynamicConfig.FeeQuoter.Bytes()))

	// Now Sync the CCIP Reader's S1 chain contract reader with OnRamp binding
	onRampContractMapping := make(ccipreaderpkg.ContractAddresses)
	onRampContractMapping[consts.ContractNameOnRamp] = make(map[cciptypes.ChainSelector]cciptypes.UnknownAddress)
	onRampContractMapping[consts.ContractNameOnRamp][chainS1] = onRampS1Addr.Bytes()

	err = reader.Sync(ctx, onRampContractMapping)
	require.NoError(t, err)

	// Since config poller has default refresh interval of 30s, we need to wait for the contract to be discovered
	require.Eventually(t, func() bool {
		contractAddresses, err = reader.DiscoverContracts(ctx, []cciptypes.ChainSelector{chainS1, chainD})
		if err != nil {
			return false
		}

		// Check if router and FeeQuoter addresses on source chain are now discovered
		routerS1, routerExists := contractAddresses[consts.ContractNameRouter][chainS1]
		feeQuoterS1, feeQuoterExists := contractAddresses[consts.ContractNameFeeQuoter][chainS1]

		return routerExists && feeQuoterExists &&
			bytes.Equal(routerS1, destinationChainConfigArgs[0].Router.Bytes()) &&
			bytes.Equal(feeQuoterS1, onRampS1DynamicConfig.FeeQuoter.Bytes())
	}, tests.WaitTimeout(t), 100*time.Millisecond, "Router and FeeQuoter addresses were not discovered on source chain in time")

	// Final assertions again for completeness:
	contractAddresses, err = reader.DiscoverContracts(ctx, []cciptypes.ChainSelector{chainS1, chainD})
	require.NoError(t, err)

	require.Equal(t, contractAddresses[consts.ContractNameOnRamp][chainS1], cciptypes.UnknownAddress(common.LeftPadBytes(onRampS1Addr.Bytes(), 32)))
	require.Equal(t, contractAddresses[consts.ContractNameRouter][chainD], cciptypes.UnknownAddress(destinationChainConfigArgs[0].Router.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameRMNRemote][chainD], cciptypes.UnknownAddress(offRampDStaticConfig.RmnRemote.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameNonceManager][chainD], cciptypes.UnknownAddress(offRampDStaticConfig.NonceManager.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameFeeQuoter][chainD], cciptypes.UnknownAddress(offRampDDynamicConfig.FeeQuoter.Bytes()))

	// Final assert to confirm source chain addresses discovered
	require.Equal(t, contractAddresses[consts.ContractNameRouter][chainS1], cciptypes.UnknownAddress(destinationChainConfigArgs[0].Router.Bytes()))
	require.Equal(t, contractAddresses[consts.ContractNameFeeQuoter][chainS1], cciptypes.UnknownAddress(onRampS1DynamicConfig.FeeQuoter.Bytes()))
}

func Test_GetChainFeePriceUpdates(t *testing.T) {
	t.Parallel()
	env, _ := testhelpers.NewMemoryEnvironment(t, testhelpers.WithNumOfChains(3))
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	dest, source1, source2 := selectors[0], selectors[1], selectors[2]

	// Setup: Add lanes and default configs (This sets default prices)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, source1, dest, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, source2, dest, false)

	// Setup: Explicitly change the gas prices for source1 and source2 on dest's FeeQuoter
	feeQuoterDest := state.MustGetEVMChainState(dest).FeeQuoter
	source1GasPrice := big.NewInt(987654321) // Use a distinct value for source1
	source2GasPrice := big.NewInt(123456789) // Use a distinct value for source2
	_, err = feeQuoterDest.UpdatePrices(
		env.Env.BlockChains.EVMChains()[dest].DeployerKey, fee_quoter.InternalPriceUpdates{
			GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
				{
					DestChainSelector: source1, // Corresponds to sending message *to* source1
					UsdPerUnitGas:     source1GasPrice,
				},
				{
					DestChainSelector: source2, // Corresponds to sending message *to* source2
					UsdPerUnitGas:     source2GasPrice,
				},
			},
		},
	)
	require.NoError(t, err)
	be := env.Env.BlockChains.EVMChains()[dest].Client.(*memory.Backend)
	be.Commit()

	// Verify the updates took effect on-chain (optional sanity check)
	gas1, err := feeQuoterDest.GetDestinationChainGasPrice(&bind.CallOpts{}, source1)
	require.NoError(t, err)
	require.Equal(t, source1GasPrice, gas1.Value)
	gas2, err := feeQuoterDest.GetDestinationChainGasPrice(&bind.CallOpts{}, source2)
	require.NoError(t, err)
	require.Equal(t, source2GasPrice, gas2.Value)

	// Setup: Create the reader instance configured for dest (destination)
	// Note: The testSetupRealContracts binds the FeeQuoter contract for the *destination* chain (dest here)
	reader := testSetupRealContracts(
		t.Context(),
		t,
		dest, // Reader is configured for dest
		map[cciptypes.ChainSelector][]types.BoundContract{
			cciptypes.ChainSelector(dest): { // Binding for the reader's chain (dest)
				{
					Address: state.MustGetEVMChainState(dest).FeeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
			},
			// Note: No bindings needed for source chains (source1, source2) for this specific reader function
		},
		nil,
		env,
	)

	t.Run("happy path - fetch prices for multiple source chains", func(t *testing.T) {
		// Act: Query for both source chains
		updates := reader.GetChainFeePriceUpdate(t.Context(), []cciptypes.ChainSelector{cs(source1), cs(source2)})

		// Assert: Expect updates for both source1 and source2
		require.Len(t, updates, 2, "Should get updates for both source chains")

		// Check source1 price (should be the explicitly set value)
		require.Contains(t, updates, cs(source1))
		assert.NotNil(t, updates[cs(source1)].Value)
		assert.Equal(t, 0, updates[cs(source1)].Value.Cmp(source1GasPrice), "Source1 price mismatch")
		assert.NotZero(t, updates[cs(source1)].Timestamp, "Source1 timestamp should be non-zero")

		// Check source2 price (should be the explicitly set value)
		require.Contains(t, updates, cs(source2))
		assert.NotNil(t, updates[cs(source2)].Value)
		assert.Equal(t, 0, updates[cs(source2)].Value.Cmp(source2GasPrice), "Source2 price mismatch")
		assert.NotZero(t, updates[cs(source2)].Timestamp, "Source2 timestamp should be non-zero")
	})

	t.Run("query non-existent chain", func(t *testing.T) {
		nonExistentChain := cciptypes.ChainSelector(99999)
		// Act: Query for existing (source1, source2) and non-existent chains. Also query for dest itself.
		updates := reader.GetChainFeePriceUpdate(t.Context(), []cciptypes.ChainSelector{cs(dest), cs(source1), cs(source2), nonExistentChain})

		// Assert: Expect updates only for source1 and source2.
		require.Len(t, updates, 2, "Should only get updates for source1 and source2")
		require.NotContains(t, updates, cs(dest)) // Dest itself shouldn't have an entry
		require.Contains(t, updates, cs(source1))
		require.Contains(t, updates, cs(source2))
		require.NotContains(t, updates, nonExistentChain)
	})

	t.Run("query empty selectors", func(t *testing.T) {
		// Act: Query with an empty slice
		updates := reader.GetChainFeePriceUpdate(t.Context(), []cciptypes.ChainSelector{})

		// Assert: Expect an empty map
		require.Empty(t, updates)
	})
}

func Test_LinkPriceUSD(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
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
					Address: state.MustGetEVMChainState(chain1).FeeQuoter.Address().String(),
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

func Test_GetWrappedNativeTokenPriceUSD(t *testing.T) {
	t.Parallel()
	ctx := t.Context()
	env, _ := testhelpers.NewMemoryEnvironment(t)
	state, err := stateview.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.BlockChains.ListChainSelectors(cldf_chain.WithFamily(chain_selectors.FamilyEVM))
	chain1, chain2 := selectors[0], selectors[1]

	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain1, chain2, false)
	testhelpers.AddLaneWithDefaultPricesAndFeeQuoterConfig(t, &env, state, chain2, chain1, false)

	reader := testSetupRealContracts(
		ctx,
		t,
		chain1,
		map[cciptypes.ChainSelector][]types.BoundContract{
			cciptypes.ChainSelector(chain2): {
				{
					Address: state.MustGetEVMChainState(chain2).FeeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
				{
					Address: state.MustGetEVMChainState(chain2).Router.Address().String(),
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
	require.Equal(t, testhelpers.DefaultWethPrice, prices[cciptypes.ChainSelector(chain2)].Int)
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
	ctx := b.Context()
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
		reports, err := s.reader.CommitReportsGTETimestamp(ctx, queryTimestamp, primitives.Unconfirmed, logsInsertedFirst)
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
	require.NoError(b, testEnv.orm.InsertLogs(ctx, logs))
	require.NoError(b, testEnv.orm.InsertBlock(ctx, utils.RandomHash(), int64(offset+numOfReports), timestamp, int64(offset+numOfReports)))
}

// Benchmark Results:
// Benchmark_CCIPReader_ExecutedMessages/LogsInserted_0_StartSeq_0_EndSeq_10-14               13599            93414 ns/op           43389 B/op        654 allocs/op
// Benchmark_CCIPReader_ExecutedMessages/LogsInserted_10_StartSeq_10_EndSeq_20-14             13471            88392 ns/op           43011 B/op        651 allocs/op
// Benchmark_CCIPReader_ExecutedMessages/LogsInserted_10_StartSeq_0_EndSeq_9-14                2799           473396 ns/op          303737 B/op       4535 allocs/op
// Benchmark_CCIPReader_ExecutedMessages/LogsInserted_100_StartSeq_0_EndSeq_100-14              438          2724414 ns/op         2477573 B/op      37468 allocs/op
// Benchmark_CCIPReader_ExecutedMessages/LogsInserted_100000_StartSeq_99744_EndSeq_100000-14     40         29118796 ns/op        12607995 B/op     179396 allocs/op
func Benchmark_CCIPReader_ExecutedMessages(b *testing.B) {
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
			benchmarkExecutedMessages(b, tt.logsInserted, tt.startSeqNum, tt.endSeqNum)
		})
	}
}

func benchmarkExecutedMessages(b *testing.B, logsInsertedFirst int, startSeqNum, endSeqNum cciptypes.SeqNum) {
	// Initialize test setup
	ctx := b.Context()
	s := setupExecutedMessagesTest(ctx, b, true)
	expectedRangeLen := calculateExpectedRangeLen(logsInsertedFirst, startSeqNum, endSeqNum)

	// Insert logs in two phases based on parameters
	if logsInsertedFirst > 0 {
		populateDatabaseForExecutionStateChanged(ctx, b, s, chainS1, chainD, logsInsertedFirst, 0)
	}

	// Reset timer to measure only the query time
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		executedRanges, err := s.reader.ExecutedMessages(
			ctx,
			map[cciptypes.ChainSelector][]cciptypes.SeqNumRange{
				chainS1: {
					cciptypes.NewSeqNumRange(startSeqNum, endSeqNum),
				},
			},
			primitives.Unconfirmed,
		)
		require.NoError(b, err)
		require.Len(b, executedRanges[chainS1], expectedRangeLen)
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
	ctx := b.Context()
	s := setupMsgsBetweenSeqNumsTest(ctx, b, true, chainS1)
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
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)

	var crs = make(map[cciptypes.ChainSelector]contractreader.Extended)
	for chain, bindings := range toBindContracts {
		be := env.Env.BlockChains.EVMChains()[uint64(chain)].Client.(*memory.Backend)
		cl := client.NewSimulatedBackendClient(t, be.Sim, big.NewInt(0).SetUint64(uint64(chain)))
		headTracker := headstest.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
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
	mokAddrCodec := newMockAddressCodec(t)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, cciptypes.ChainSelector(destChain), nil, mokAddrCodec)

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
		RPCBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, params.SimulatedBackend, big.NewInt(0).SetUint64(uint64(params.ReaderChain)))
	headTracker := headstest.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
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
		headTracker2 := headstest.NewSimulatedHeadTracker(cl2, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
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

	mokAddrCodec := newMockAddressCodec(t)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, params.DestChain, nil, mokAddrCodec)

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

func newMockAddressCodec(t testing.TB) *typepkgmock.MockAddressCodec {
	mockAddrCodec := typepkgmock.NewMockAddressCodec(t)
	mockAddrCodec.On("AddressBytesToString", mock.Anything, mock.Anything).
		Return(func(addr cciptypes.UnknownAddress, _ cciptypes.ChainSelector) string {
			return "0x" + hex.EncodeToString(addr)
		}, nil).Maybe()
	mockAddrCodec.On("AddressStringToBytes", mock.Anything, mock.Anything).
		Return(func(addr string, _ cciptypes.ChainSelector) (cciptypes.UnknownAddress, error) {
			addrBytes, err := hex.DecodeString(strings.ToLower(strings.TrimPrefix(addr, "0x")))
			if err != nil {
				return nil, err
			}
			return addrBytes, nil
		}).Maybe()
	return mockAddrCodec
}
