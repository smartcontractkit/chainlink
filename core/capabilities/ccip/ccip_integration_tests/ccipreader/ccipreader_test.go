package ccipreader

import (
	"context"
	"encoding/hex"
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

	readermocks "github.com/smartcontractkit/chainlink-ccip/mocks/pkg/contractreader"
	cciptypes "github.com/smartcontractkit/chainlink-ccip/pkg/types/ccipocr3"
	"github.com/smartcontractkit/chainlink-common/pkg/services/servicetest"

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
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"

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

var (
	defaultGasPrice  = assets.GWei(10)
	InitialLinkPrice = e18Mult(20)
	InitialWethPrice = e18Mult(4000)
	linkAddress      = utils.RandomAddress()
	wethAddress      = utils.RandomAddress()
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

	sb, auth := setupSimulatedBackendAndAuth(t)
	onRampAddress := utils.RandomAddress()
	s := testSetup(ctx, t, chainD, chainD, nil, cfg, nil, map[cciptypes.ChainSelector][]types.BoundContract{
		chainS1: {
			{
				Address: onRampAddress.Hex(),
				Name:    consts.ContractNameOnRamp,
			},
		},
	},
		true,
		sb,
		auth,
	)

	tokenA := common.HexToAddress("123")
	const numReports = 5

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
	s := testSetup(ctx, t, chainD, chainD, nil, cfg, nil, nil, true, sb, auth)

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
	s := testSetup(ctx, t, chainS1, chainD, nil, cfg, nil, nil, true, sb, auth)

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

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, chainD, chainD, onChainSeqNums, cfg, nil, nil, true, sb, auth)

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

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, chainS1, chainD, nil, cfg, nil, nil, true, sb, auth)

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

	sb, auth := setupSimulatedBackendAndAuth(t)
	s := testSetup(ctx, t, chainD, chainD, nil, cfg, nil, nil, true, sb, auth)

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
	ctx := testutils.Context(t)
	sb, auth := setupSimulatedBackendAndAuth(t)
	feeQuoter := deployFeeQuoterWithPrices(t, auth, sb, chainS1)

	s := testSetup(ctx, t, chainD, chainD, nil, evmconfig.DestReaderConfig,
		map[cciptypes.ChainSelector][]types.BoundContract{
			chainD: {
				{
					Address: feeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
			},
		},
		nil,
		false,
		sb,
		auth,
	)

	updates := s.reader.GetChainFeePriceUpdate(ctx, []cciptypes.ChainSelector{chainS1, chainS2})
	// only chainS1 has a bound contract
	require.Len(t, updates, 1)
	require.Equal(t, defaultGasPrice.ToInt(), updates[chainS1].Value.Int)
}

func Test_LinkPriceUSD(t *testing.T) {
	ctx := testutils.Context(t)
	sb, auth := setupSimulatedBackendAndAuth(t)
	feeQuoter := deployFeeQuoterWithPrices(t, auth, sb, chainS1)

	s := testSetup(ctx, t, chainD, chainD, nil, evmconfig.DestReaderConfig,
		map[cciptypes.ChainSelector][]types.BoundContract{
			chainD: {
				{
					Address: feeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
			},
		},
		nil,
		false,
		sb,
		auth,
	)

	linkPriceUSD, err := s.reader.LinkPriceUSD(ctx)
	require.NoError(t, err)
	require.NotNil(t, linkPriceUSD.Int)
	require.Equal(t, InitialLinkPrice, linkPriceUSD.Int)
}

func Test_GetMedianDataAvailabilityGasConfig(t *testing.T) {
	ctx := testutils.Context(t)

	sb, auth := setupSimulatedBackendAndAuth(t)

	// All fee quoters using same auth and simulated backend for simplicity
	feeQuoter1 := deployFeeQuoterWithPrices(t, auth, sb, chainD)
	feeQuoter2 := deployFeeQuoterWithPrices(t, auth, sb, chainD)
	feeQuoter3 := deployFeeQuoterWithPrices(t, auth, sb, chainD)
	feeQuoters := []*fee_quoter.FeeQuoter{feeQuoter1, feeQuoter2, feeQuoter3}

	// Update the dest chain config for each fee quoter
	for i, fq := range feeQuoters {
		destChainCfg := defaultFeeQuoterDestChainConfig()
		//nolint:gosec // disable G115
		destChainCfg.DestDataAvailabilityOverheadGas = uint32(100 + i)
		//nolint:gosec // disable G115
		destChainCfg.DestGasPerDataAvailabilityByte = uint16(200 + i)
		//nolint:gosec // disable G115
		destChainCfg.DestDataAvailabilityMultiplierBps = uint16(1 + i)
		_, err := fq.ApplyDestChainConfigUpdates(auth, []fee_quoter.FeeQuoterDestChainConfigArgs{
			{
				DestChainSelector: uint64(chainD),
				DestChainConfig:   destChainCfg,
			},
		})
		sb.Commit()
		require.NoError(t, err)
	}

	s := testSetup(ctx, t, chainD, chainD, nil, evmconfig.DestReaderConfig, map[cciptypes.ChainSelector][]types.BoundContract{
		chainS1: {
			{
				Address: feeQuoter1.Address().String(),
				Name:    consts.ContractNameFeeQuoter,
			},
		},
		chainS2: {
			{
				Address: feeQuoter2.Address().String(),
				Name:    consts.ContractNameFeeQuoter,
			},
		},
		chainS3: {
			{
				Address: feeQuoter3.Address().String(),
				Name:    consts.ContractNameFeeQuoter,
			},
		},
	}, nil,
		false,
		sb,
		auth,
	)

	daConfig, err := s.reader.GetMedianDataAvailabilityGasConfig(ctx)
	require.NoError(t, err)

	// Verify the results
	require.Equal(t, uint32(101), daConfig.DestDataAvailabilityOverheadGas)
	require.Equal(t, uint16(201), daConfig.DestGasPerDataAvailabilityByte)
	require.Equal(t, uint16(2), daConfig.DestDataAvailabilityMultiplierBps)
}

func Test_GetWrappedNativeTokenPriceUSD(t *testing.T) {
	ctx := testutils.Context(t)
	sb, auth := setupSimulatedBackendAndAuth(t)
	feeQuoter := deployFeeQuoterWithPrices(t, auth, sb, chainS1)

	// Mock the routerContract to return a native token address
	routerContract := deployRouterWithNativeToken(t, auth, sb)

	s := testSetup(ctx, t, chainD, chainD, nil, evmconfig.DestReaderConfig,
		map[cciptypes.ChainSelector][]types.BoundContract{
			chainD: {
				{
					Address: feeQuoter.Address().String(),
					Name:    consts.ContractNameFeeQuoter,
				},
				{
					Address: routerContract.Address().String(),
					Name:    consts.ContractNameRouter,
				},
			},
		},
		nil,
		false,
		sb,
		auth,
	)

	prices := s.reader.GetWrappedNativeTokenPriceUSD(ctx, []cciptypes.ChainSelector{chainD, chainS1})

	// Only chainD has reader contracts bound
	require.Len(t, prices, 1)
	require.Equal(t, InitialWethPrice, prices[chainD].Int)
}

func deployRouterWithNativeToken(t *testing.T, auth *bind.TransactOpts, sb *simulated.Backend) *router.Router {
	address, _, _, err := router.DeployRouter(
		auth,
		sb.Client(),
		wethAddress,
		utils.RandomAddress(), // armProxy address
	)
	require.NoError(t, err)
	sb.Commit()

	routerContract, err := router.NewRouter(address, sb.Client())
	require.NoError(t, err)

	return routerContract
}

func deployFeeQuoterWithPrices(t *testing.T, auth *bind.TransactOpts, sb *simulated.Backend, destChain cciptypes.ChainSelector) *fee_quoter.FeeQuoter {
	address, _, _, err := fee_quoter.DeployFeeQuoter(
		auth,
		sb.Client(),
		fee_quoter.FeeQuoterStaticConfig{
			MaxFeeJuelsPerMsg:            big.NewInt(0).Mul(big.NewInt(2e2), big.NewInt(1e18)),
			LinkToken:                    linkAddress,
			TokenPriceStalenessThreshold: uint32(24 * 60 * 60),
		},
		[]common.Address{auth.From},
		[]common.Address{wethAddress, linkAddress},
		[]fee_quoter.FeeQuoterTokenPriceFeedUpdate{},
		[]fee_quoter.FeeQuoterTokenTransferFeeConfigArgs{},
		[]fee_quoter.FeeQuoterPremiumMultiplierWeiPerEthArgs{},
		[]fee_quoter.FeeQuoterDestChainConfigArgs{
			{

				DestChainSelector: uint64(destChain),
				DestChainConfig:   defaultFeeQuoterDestChainConfig(),
			},
		},
	)

	require.NoError(t, err)
	sb.Commit()

	feeQuoter, err := fee_quoter.NewFeeQuoter(address, sb.Client())
	require.NoError(t, err)

	_, err = feeQuoter.UpdatePrices(
		auth, fee_quoter.InternalPriceUpdates{
			GasPriceUpdates: []fee_quoter.InternalGasPriceUpdate{
				{
					DestChainSelector: uint64(chainS1),
					UsdPerUnitGas:     defaultGasPrice.ToInt(),
				},
			},
			TokenPriceUpdates: []fee_quoter.InternalTokenPriceUpdate{
				{
					SourceToken: linkAddress,
					UsdPerToken: InitialLinkPrice,
				},
				{
					SourceToken: wethAddress,
					UsdPerToken: InitialWethPrice,
				},
			},
		},
	)
	require.NoError(t, err)
	sb.Commit()

	gas, err := feeQuoter.GetDestinationChainGasPrice(&bind.CallOpts{}, uint64(chainS1))
	require.NoError(t, err)
	require.Equal(t, defaultGasPrice.ToInt(), gas.Value)

	return feeQuoter
}

func defaultFeeQuoterDestChainConfig() fee_quoter.FeeQuoterDestChainConfig {
	// https://github.com/smartcontractkit/ccip/blob/c4856b64bd766f1ddbaf5d13b42d3c4b12efde3a/contracts/src/v0.8/ccip/libraries/Internal.sol#L337-L337
	/*
		```Solidity
			// bytes4(keccak256("CCIP ChainFamilySelector EVM"))
			bytes4 public constant CHAIN_FAMILY_SELECTOR_EVM = 0x2812d52c;
		```
	*/
	evmFamilySelector, _ := hex.DecodeString("2812d52c")
	return fee_quoter.FeeQuoterDestChainConfig{
		IsEnabled:                         true,
		MaxNumberOfTokensPerMsg:           10,
		MaxDataBytes:                      256,
		MaxPerMsgGasLimit:                 3_000_000,
		DestGasOverhead:                   50_000,
		DefaultTokenFeeUSDCents:           1,
		DestGasPerPayloadByte:             10,
		DestDataAvailabilityOverheadGas:   100,
		DestGasPerDataAvailabilityByte:    100,
		DestDataAvailabilityMultiplierBps: 1,
		DefaultTokenDestGasOverhead:       125_000,
		DefaultTxGasLimit:                 200_000,
		GasMultiplierWeiPerEth:            1,
		NetworkFeeUSDCents:                1,
		ChainFamilySelector:               [4]byte(evmFamilySelector),
	}
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

func testSetup(
	ctx context.Context,
	t *testing.T,
	readerChain,
	destChain cciptypes.ChainSelector,
	onChainSeqNums map[cciptypes.ChainSelector]cciptypes.SeqNum,
	cfg evmtypes.ChainReaderConfig,
	toBindContracts map[cciptypes.ChainSelector][]types.BoundContract,
	toMockBindings map[cciptypes.ChainSelector][]types.BoundContract,
	bindTester bool,
	simulatedBackend *simulated.Backend,
	auth *bind.TransactOpts,
) *testSetupData {
	address, _, _, err := ccip_reader_tester.DeployCCIPReaderTester(auth, simulatedBackend.Client())
	assert.NoError(t, err)
	simulatedBackend.Commit()

	// Setup contract client
	contract, err := ccip_reader_tester.NewCCIPReaderTester(address, simulatedBackend.Client())
	assert.NoError(t, err)

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	db := pgtest.NewSqlxDB(t)
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
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
	servicetest.Run(t, lp)

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

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, cfg)
	require.NoError(t, err)

	extendedCr := contractreader.NewExtendedContractReader(cr)

	if bindTester {
		err = extendedCr.Bind(ctx, []types.BoundContract{
			{
				Address: address.String(),
				Name:    contractNames[0],
			},
		})
		require.NoError(t, err)
	}

	var otherCrs = make(map[cciptypes.ChainSelector]contractreader.Extended)
	for chain, bindings := range toBindContracts {
		cl2 := client.NewSimulatedBackendClient(t, simulatedBackend, big.NewInt(0).SetUint64(uint64(chain)))
		headTracker2 := headtracker.NewSimulatedHeadTracker(cl2, lpOpts.UseFinalityTag, lpOpts.FinalityDepth)
		lp2 := logpoller.NewLogPoller(logpoller.NewORM(big.NewInt(0).SetUint64(uint64(chain)), db, lggr),
			cl2,
			lggr,
			headTracker2,
			lpOpts,
		)
		servicetest.Run(t, lp2)

		cr2, err2 := evm.NewChainReaderService(ctx, lggr, lp2, headTracker2, cl2, cfg)
		require.NoError(t, err2)

		extendedCr2 := contractreader.NewExtendedContractReader(cr2)
		err2 = extendedCr2.Bind(ctx, bindings)
		require.NoError(t, err2)
		otherCrs[chain] = extendedCr2
	}

	for chain, bindings := range toMockBindings {
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

	servicetest.Run(t, cr)

	contractReaders := map[cciptypes.ChainSelector]contractreader.Extended{readerChain: extendedCr}
	for chain, cr := range otherCrs {
		contractReaders[chain] = cr
	}
	contractWriters := make(map[cciptypes.ChainSelector]types.ChainWriter)
	reader := ccipreaderpkg.NewCCIPReaderWithExtendedContractReaders(ctx, lggr, contractReaders, contractWriters, destChain, nil)

	return &testSetupData{
		contractAddr: address,
		contract:     contract,
		sb:           simulatedBackend,
		auth:         auth,
		lp:           lp,
		cl:           cl,
		reader:       reader,
		extendedCR:   extendedCr,
	}
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

func uBigInt(i uint64) *big.Int {
	return new(big.Int).SetUint64(i)
}

func e18Mult(amount uint64) *big.Int {
	return new(big.Int).Mul(uBigInt(amount), uBigInt(1e18))
}
