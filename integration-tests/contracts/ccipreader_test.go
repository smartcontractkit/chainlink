package contracts

import (
	"context"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient/simulated"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-ccip/plugintypes"
	"github.com/smartcontractkit/chainlink/deployment/ccip/changeset"
	"github.com/smartcontractkit/chainlink/deployment/environment/memory"
	"github.com/smartcontractkit/chainlink/integration-tests/utils/pgtest"

	readermocks "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/contractreader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"

	"github.com/smartcontractkit/chainlink-common/pkg/codec"
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

func setupGetCommitGTETimestampTest(ctx context.Context, t *testing.T, finalityDepth int64) (*testSetupData, int64, common.Address) {
	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameOffRamp: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{consts.EventNameCommitReportAccepted},
				},
				ContractABI: ccip_reader_tester.CCIPReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.EventNameCommitReportAccepted: {
						ChainSpecificName: consts.EventNameCommitReportAccepted,
						ReadType:          evmtypes.Event,
					},
				},
			},
		},
	}

	sb, auth := setupSimulatedBackendAndAuth(t)
	onRampAddress := utils.RandomAddress()
	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:    chainD,
		DestChain:      chainD,
		OnChainSeqNums: nil,
		Cfg:            cfg,
		ToMockBindings: map[cciptypes.ChainSelector][]types.BoundContract{
			chainS1: {
				{
					Address: onRampAddress.Hex(),
					Name:    consts.ContractNameOnRamp,
				},
			},
		},
		BindTester:       true,
		SimulatedBackend: sb,
		Auth:             auth,
		FinalityDepth:    finalityDepth,
	})

	return s, finalityDepth, onRampAddress
}

func emitCommitReports(ctx context.Context, t *testing.T, s *testSetupData, numReports int, tokenA common.Address, onRampAddress common.Address) uint64 {
	var firstReportTs uint64
	for i := 0; i < numReports; i++ {
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
					MerkleRoot:          [32]byte{uint8(i) + 1}, //nolint:gosec // this won't overflow
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
		assert.NoError(t, err)
		bh := s.sb.Commit()
		b, err := s.sb.Client().BlockByHash(ctx, bh)
		require.NoError(t, err)
		if firstReportTs == 0 {
			firstReportTs = b.Time()
		}
	}
	return firstReportTs
}

func TestCCIPReader_CommitReportsGTETimestamp(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	s, _, onRampAddress := setupGetCommitGTETimestampTest(ctx, t, 0)

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
	s, _, onRampAddress := setupGetCommitGTETimestampTest(ctx, t, finalityDepth)

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
	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameOffRamp: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{consts.EventNameExecutionStateChanged},
				},
				ContractABI: ccip_reader_tester.CCIPReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.EventNameExecutionStateChanged: {
						ChainSpecificName: consts.EventNameExecutionStateChanged,
						ReadType:          evmtypes.Event,
						EventDefinitions: &evmtypes.EventDefinitions{
							GenericTopicNames: map[string]string{
								"sourceChainSelector": consts.EventAttributeSourceChain,
								"sequenceNumber":      consts.EventAttributeSequenceNumber,
							},
							GenericDataWordDetails: map[string]evmtypes.DataWordDetail{
								consts.EventAttributeState: {
									Name: "state",
								},
							},
						},
					},
				},
			},
		},
	}

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:      chainD,
		DestChain:        chainD,
		OnChainSeqNums:   nil,
		Cfg:              cfg,
		ToBindContracts:  nil,
		ToMockBindings:   nil,
		BindTester:       true,
		SimulatedBackend: sb,
		Auth:             auth,
	})
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
	assert.NoError(t, err)
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
	assert.NoError(t, err)
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

	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameOnRamp: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{consts.EventNameCCIPMessageSent},
				},
				ContractABI: ccip_reader_tester.CCIPReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.EventNameCCIPMessageSent: {
						ChainSpecificName: "CCIPMessageSent",
						ReadType:          evmtypes.Event,
						EventDefinitions: &evmtypes.EventDefinitions{
							GenericDataWordDetails: map[string]evmtypes.DataWordDetail{
								consts.EventAttributeSourceChain:    {Name: "message.header.sourceChainSelector"},
								consts.EventAttributeDestChain:      {Name: "message.header.destChainSelector"},
								consts.EventAttributeSequenceNumber: {Name: "message.header.sequenceNumber"},
							},
						},
						OutputModifications: codec.ModifiersConfig{
							&codec.WrapperModifierConfig{Fields: map[string]string{
								"Message.FeeTokenAmount":      "Int",
								"Message.FeeValueJuels":       "Int",
								"Message.TokenAmounts.Amount": "Int",
							}},
						},
					},
				},
			},
		},
	}

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, testSetupParams{
		ReaderChain:      chainS1,
		DestChain:        chainD,
		OnChainSeqNums:   nil,
		Cfg:              cfg,
		ToBindContracts:  nil,
		ToMockBindings:   nil,
		BindTester:       true,
		SimulatedBackend: sb,
		Auth:             auth,
	})

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
	assert.NoError(t, err)

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
	assert.NoError(t, err)

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
		ReaderChain:      chainD,
		DestChain:        chainD,
		OnChainSeqNums:   onChainSeqNums,
		Cfg:              cfg,
		ToBindContracts:  nil,
		ToMockBindings:   nil,
		BindTester:       true,
		SimulatedBackend: sb,
		Auth:             auth,
	})

	seqNums, err := s.reader.NextSeqNum(ctx, []cciptypes.ChainSelector{chainS1, chainS2, chainS3})
	assert.NoError(t, err)
	assert.Len(t, seqNums, 3)
	assert.Equal(t, cciptypes.SeqNum(10), seqNums[0])
	assert.Equal(t, cciptypes.SeqNum(20), seqNums[1])
	assert.Equal(t, cciptypes.SeqNum(30), seqNums[2])
}

func TestCCIPReader_GetExpectedNextSequenceNumber(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	//env := NewMemoryEnvironmentContractsOnly(t, logger.TestLogger(t), 2, 4, nil)
	env := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), 2, 4, nil)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	destChain, srcChain := selectors[0], selectors[1]

	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, destChain, srcChain, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, srcChain, destChain, false))

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
		msg := changeset.DefaultRouterMessage(state.Chains[destChain].Receiver.Address())
		msgSentEvent := changeset.TestSendRequest(t, env.Env, state, srcChain, destChain, false, msg)
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
		ReaderChain:      chainD,
		DestChain:        chainD,
		Cfg:              cfg,
		BindTester:       true,
		SimulatedBackend: sb,
		Auth:             auth,
	})

	// Add some nonces.
	for chain, addrs := range nonces {
		for addr, nonce := range addrs {
			_, err := s.contract.SetInboundNonce(s.auth, uint64(chain), nonce, common.LeftPadBytes(addr.Bytes(), 32))
			assert.NoError(t, err)
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
		assert.NoError(t, err)
		assert.Len(t, results, len(addrQuery))
		for addr, nonce := range addrs {
			assert.Equal(t, nonce, results[addr.String()])
		}
	}
}

func Test_GetChainFeePriceUpdates(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	env := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), 2, 4, nil)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain1, chain2, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain2, chain1, false))

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
	env := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), 2, 4, nil)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain1, chain2, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain2, chain1, false))

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
	require.Equal(t, changeset.DefaultInitialPrices.LinkPrice, linkPriceUSD.Int)
}

func Test_GetMedianDataAvailabilityGasConfig(t *testing.T) {
	t.Parallel()
	ctx := tests.Context(t)
	env := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), 4, 4, nil)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	destChain, chain1, chain2, chain3 := selectors[0], selectors[1], selectors[2], selectors[3]

	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain1, destChain, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain2, destChain, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain3, destChain, false))

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
	env := changeset.NewMemoryEnvironmentWithJobsAndContracts(t, logger.TestLogger(t), 2, 4, nil)
	state, err := changeset.LoadOnchainState(env.Env)
	require.NoError(t, err)

	selectors := env.Env.AllChainSelectors()
	chain1, chain2 := selectors[0], selectors[1]

	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain1, chain2, false))
	require.NoError(t, changeset.AddLaneWithDefaultPricesAndFeeQuoterConfig(env.Env, state, chain2, chain1, false))

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
	require.Equal(t, changeset.DefaultInitialPrices.WethPrice, prices[cciptypes.ChainSelector(chain1)].Int)
}

func setupSimulatedBackendAndAuth(t *testing.T) (*simulated.Backend, *bind.TransactOpts) {
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
	env changeset.DeployedEnv,
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
	contractWriters := make(map[cciptypes.ChainSelector]types.ChainWriter)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, cciptypes.ChainSelector(destChain), nil)

	return reader
}

func testSetup(
	ctx context.Context,
	t *testing.T,
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
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            params.FinalityDepth,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, params.SimulatedBackend, big.NewInt(0).SetUint64(uint64(params.ReaderChain)))
	headTracker := headtracker.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(0).SetUint64(uint64(params.ReaderChain)), db, lggr),
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

	contractNames := maps.Keys(params.Cfg.Contracts)

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, params.Cfg)
	require.NoError(t, err)

	extendedCr := contractreader.NewExtendedContractReader(cr)

	if params.BindTester {
		err = extendedCr.Bind(ctx, []types.BoundContract{
			{
				Address: address.String(),
				Name:    contractNames[0],
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
	contractWriters := make(map[cciptypes.ChainSelector]types.ChainWriter)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, params.DestChain, nil)

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
		lp:           lp,
		cl:           cl,
		reader:       reader,
		extendedCR:   extendedCr,
	}
}

type testSetupParams struct {
	ReaderChain      cciptypes.ChainSelector
	DestChain        cciptypes.ChainSelector
	OnChainSeqNums   map[cciptypes.ChainSelector]cciptypes.SeqNum
	Cfg              evmtypes.ChainReaderConfig
	ToBindContracts  map[cciptypes.ChainSelector][]types.BoundContract
	ToMockBindings   map[cciptypes.ChainSelector][]types.BoundContract
	BindTester       bool
	SimulatedBackend *simulated.Backend
	Auth             *bind.TransactOpts
	FinalityDepth    int64
}

type testSetupData struct {
	contractAddr common.Address
	contract     *ccip_reader_tester.CCIPReaderTester
	sb           *simulated.Backend
	auth         *bind.TransactOpts
	lp           logpoller.LogPoller
	cl           client.Client
	reader       ccipreaderpkg.CCIPReader
	extendedCR   contractreader.Extended
}

func cs(i uint64) cciptypes.ChainSelector {
	return cciptypes.ChainSelector(i)
}
