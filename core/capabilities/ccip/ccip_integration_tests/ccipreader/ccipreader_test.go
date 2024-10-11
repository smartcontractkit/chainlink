package ccipreader

import (
	"context"
	"math/big"
	"sort"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/ccip_reader_tester"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

	readermocks "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	ccipreaderpkg "github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/plugintypes"
)

const (
	chainS1 = cciptypes.ChainSelector(1)
	chainS2 = cciptypes.ChainSelector(2)
	chainS3 = cciptypes.ChainSelector(3)
	chainD  = cciptypes.ChainSelector(4)
)

func TestCCIPReader_CommitReportsGTETimestamp(t *testing.T) {
	ctx := testutils.Context(t)

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

	onRampAddress := utils.RandomAddress()
	s := testSetup(ctx, t, chainD, chainD, nil, cfg, map[cciptypes.ChainSelector][]types.BoundContract{
		chainS1: {
			{
				Address: onRampAddress.Hex(),
				Name:    consts.ContractNameOnRamp,
			},
		},
	})

	tokenA := common.HexToAddress("123")
	const numReports = 5

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
			RmnRawVs: big.NewInt(100),
		})
		assert.NoError(t, err)
		s.sb.Commit()
	}

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, s.lp.Replay(ctx, 1))

	var reports []plugintypes.CommitPluginReportWithMeta
	var err error
	require.Eventually(t, func() bool {
		reports, err = s.reader.CommitReportsGTETimestamp(
			ctx,
			chainD,
			time.Unix(30, 0), // Skips first report, simulated backend report timestamps are [20, 30, 40, ...]
			10,
		)
		require.NoError(t, err)
		return len(reports) == numReports-1
	}, tests.WaitTimeout(t), 50*time.Millisecond)

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

	// TODO assert once chainlink-ccip changes are done
	// assert.Len(t, reports[0].Report.RMNSignatures, 2)
	// assert.Equal(t, reports[0].Report.RMNSignatures[0].R, [32]byte{1})
	// assert.Equal(t, reports[0].Report.RMNSignatures[0].S, [32]byte{2})
}

func TestCCIPReader_ExecutedMessageRanges(t *testing.T) {
	ctx := testutils.Context(t)
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
					},
				},
			},
		},
	}

	s := testSetup(ctx, t, chainD, chainD, nil, cfg, nil)

	_, err := s.contract.EmitExecutionStateChanged(
		s.auth,
		uint64(chainS1),
		14,
		cciptypes.Bytes32{1, 0, 0, 1},
		1,
		[]byte{1, 2, 3, 4},
	)
	assert.NoError(t, err)
	s.sb.Commit()

	_, err = s.contract.EmitExecutionStateChanged(
		s.auth,
		uint64(chainS1),
		15,
		cciptypes.Bytes32{1, 0, 0, 2},
		1,
		[]byte{1, 2, 3, 4, 5},
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
	}, testutils.WaitTimeout(t), 50*time.Millisecond)

	assert.Equal(t, cciptypes.SeqNum(14), executedRanges[0].Start())
	assert.Equal(t, cciptypes.SeqNum(14), executedRanges[0].End())

	assert.Equal(t, cciptypes.SeqNum(15), executedRanges[1].Start())
	assert.Equal(t, cciptypes.SeqNum(15), executedRanges[1].End())
}

func TestCCIPReader_MsgsBetweenSeqNums(t *testing.T) {
	ctx := testutils.Context(t)

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
					},
				},
			},
		},
	}

	s := testSetup(ctx, t, chainS1, chainD, nil, cfg, nil)

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
		FeeTokenAmount: big.NewInt(0),
		FeeValueJuels:  big.NewInt(0),
		TokenAmounts:   make([]ccip_reader_tester.InternalEVM2AnyTokenTransfer, 0),
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
		FeeTokenAmount: big.NewInt(0),
		FeeValueJuels:  big.NewInt(0),
		TokenAmounts:   make([]ccip_reader_tester.InternalEVM2AnyTokenTransfer, 0),
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
	require.Equal(t, cciptypes.SeqNum(15), msgs[1].Header.SequenceNumber)
	for _, msg := range msgs {
		require.Equal(t, chainS1, msg.Header.SourceChainSelector)
		require.Equal(t, chainD, msg.Header.DestChainSelector)
	}
}

func TestCCIPReader_NextSeqNum(t *testing.T) {
	ctx := testutils.Context(t)

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

	s := testSetup(ctx, t, chainD, chainD, onChainSeqNums, cfg, nil)

	seqNums, err := s.reader.NextSeqNum(ctx, []cciptypes.ChainSelector{chainS1, chainS2, chainS3})
	assert.NoError(t, err)
	assert.Len(t, seqNums, 3)
	assert.Equal(t, cciptypes.SeqNum(10), seqNums[0])
	assert.Equal(t, cciptypes.SeqNum(20), seqNums[1])
	assert.Equal(t, cciptypes.SeqNum(30), seqNums[2])
}

func TestCCIPReader_GetExpectedNextSequenceNumber(t *testing.T) {
	ctx := testutils.Context(t)

	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameOnRamp: {
				ContractABI: ccip_reader_tester.CCIPReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.MethodNameGetExpectedNextSequenceNumber: {
						ChainSpecificName: "getExpectedNextSequenceNumber",
						ReadType:          evmtypes.Method,
					},
				},
			},
		},
	}

	s := testSetup(ctx, t, chainS1, chainD, nil, cfg, nil)

	_, err := s.contract.SetDestChainSeqNr(s.auth, uint64(chainD), 10)
	require.NoError(t, err)
	s.sb.Commit()

	seqNum, err := s.reader.GetExpectedNextSequenceNumber(ctx, chainS1, chainD)
	require.NoError(t, err)
	require.Equal(t, cciptypes.SeqNum(10)+1, seqNum)

	_, err = s.contract.SetDestChainSeqNr(s.auth, uint64(chainD), 25)
	require.NoError(t, err)
	s.sb.Commit()

	seqNum, err = s.reader.GetExpectedNextSequenceNumber(ctx, chainS1, chainD)
	require.NoError(t, err)
	require.Equal(t, cciptypes.SeqNum(25)+1, seqNum)
}

func TestCCIPReader_Nonces(t *testing.T) {
	ctx := testutils.Context(t)
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

	s := testSetup(ctx, t, chainD, chainD, nil, cfg, nil)

	// Add some nonces.
	for chain, addrs := range nonces {
		for addr, nonce := range addrs {
			_, err := s.contract.SetInboundNonce(s.auth, uint64(chain), nonce, addr.Bytes())
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

func testSetup(
	ctx context.Context,
	t *testing.T,
	readerChain,
	destChain cciptypes.ChainSelector,
	onChainSeqNums map[cciptypes.ChainSelector]cciptypes.SeqNum,
	cfg evmtypes.ChainReaderConfig,
	otherBindings map[cciptypes.ChainSelector][]types.BoundContract,
) *testSetupData {
	const chainID = 1337

	// Generate a new key pair for the simulated account
	privateKey, err := crypto.GenerateKey()
	assert.NoError(t, err)
	// Set up the genesis account with balance
	blnc, ok := big.NewInt(0).SetString("999999999999999999999999999999999999", 10)
	assert.True(t, ok)
	alloc := map[common.Address]core.GenesisAccount{crypto.PubkeyToAddress(privateKey.PublicKey): {Balance: blnc}}
	simulatedBackend := backends.NewSimulatedBackend(alloc, 0)
	// Create a transactor

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	assert.NoError(t, err)
	auth.GasLimit = uint64(0)

	// Deploy the contract
	address, _, _, err := ccip_reader_tester.DeployCCIPReaderTester(auth, simulatedBackend)
	assert.NoError(t, err)
	simulatedBackend.Commit()

	// Setup contract client
	contract, err := ccip_reader_tester.NewCCIPReaderTester(address, simulatedBackend)
	assert.NoError(t, err)

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            0,
		BackfillBatchSize:        10,
		RpcBatchSize:             10,
		KeepFinalizedBlocksDepth: 100000,
	}
	cl := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(0).SetUint64(uint64(readerChain)))
	headTracker := headtracker.NewSimulatedHeadTracker(cl, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
	lp := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(0).SetUint64(uint64(readerChain)), db, lggr),
		cl,
		lggr,
		headTracker,
		lpOpts,
	)
	assert.NoError(t, lp.Start(ctx))

	for sourceChain, seqNum := range onChainSeqNums {
		_, err1 := contract.SetSourceChainConfig(auth, uint64(sourceChain), ccip_reader_tester.OffRampSourceChainConfig{
			IsEnabled: true,
			MinSeqNr:  uint64(seqNum),
			OnRamp:    utils.RandomAddress().Bytes(),
		})
		assert.NoError(t, err1)
		simulatedBackend.Commit()
		scc, err1 := contract.GetSourceChainConfig(&bind.CallOpts{Context: ctx}, uint64(sourceChain))
		assert.NoError(t, err1)
		assert.Equal(t, seqNum, cciptypes.SeqNum(scc.MinSeqNr))
	}

	contractNames := maps.Keys(cfg.Contracts)
	assert.Len(t, contractNames, 1, "test setup assumes there is only one contract")

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, cfg)
	require.NoError(t, err)

	extendedCr := contractreader.NewExtendedContractReader(cr)
	err = extendedCr.Bind(ctx, []types.BoundContract{
		{
			Address: address.String(),
			Name:    contractNames[0],
		},
	})
	require.NoError(t, err)
	var otherCrs = make(map[cciptypes.ChainSelector]contractreader.Extended)
	for chain, bindings := range otherBindings {
		m := readermocks.NewMockContractReaderFacade(t)
		m.EXPECT().Bind(ctx, bindings).Return(nil)
		ecr := contractreader.NewExtendedContractReader(m)
		err = ecr.Bind(ctx, bindings)
		require.NoError(t, err)
		otherCrs[chain] = ecr
	}

	err = cr.Start(ctx)
	require.NoError(t, err)

	contractReaders := map[cciptypes.ChainSelector]contractreader.Extended{readerChain: extendedCr}
	for chain, cr := range otherCrs {
		contractReaders[chain] = cr
	}
	contractWriters := make(map[cciptypes.ChainSelector]types.ChainWriter)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(lggr, contractReaders, contractWriters, destChain, nil)

	t.Cleanup(func() {
		require.NoError(t, cr.Close())
		require.NoError(t, lp.Close())
		require.NoError(t, db.Close())
	})

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         auth,
		lp:           lp,
		cl:           cl,
		reader:       reader,
	}
}

type testSetupData struct {
	contractAddr common.Address
	contract     *ccip_reader_tester.CCIPReaderTester
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
	lp           logpoller.LogPoller
	cl           client.Client
	reader       ccipreaderpkg.CCIPReader
}
