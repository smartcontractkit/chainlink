package usdcreader

import (
	"context"
	"math/big"
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

	sel "github.com/smartcontractkit/chain-selectors"

	"github.com/smartcontractkit/chainlink-ccip/pkg/contractreader"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"

	evmconfig "github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/configs/evm"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/headtracker"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/usdc_reader_tester"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	evmtypes "github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/types"
)

func Test_USDCReader_MessageHashes(t *testing.T) {
	finalityDepth := 5

	ctx := testutils.Context(t)
	ethereumChain := cciptypes.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)
	ethereumDomainCCTP := reader.CCTPDestDomains[uint64(ethereumChain)]
	avalancheChain := cciptypes.ChainSelector(sel.AVALANCHE_MAINNET.Selector)
	avalancheDomainCCTP := reader.CCTPDestDomains[uint64(avalancheChain)]
	polygonChain := cciptypes.ChainSelector(sel.POLYGON_MAINNET.Selector)
	polygonDomainCCTP := reader.CCTPDestDomains[uint64(polygonChain)]

	ts := testSetup(ctx, t, ethereumChain, evmconfig.USDCReaderConfig, finalityDepth)

	usdcReader, err := reader.NewUSDCMessageReader(
		ctx,
		logger.TestLogger(t),
		map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
			ethereumChain: {
				SourceMessageTransmitterAddr: ts.contractAddr.String(),
			},
		},
		map[cciptypes.ChainSelector]contractreader.ContractReaderFacade{
			ethereumChain: ts.reader,
		})
	require.NoError(t, err)

	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 11)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 21)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 31)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 41)
	emitMessageSent(t, ts, ethereumDomainCCTP, polygonDomainCCTP, 31)
	emitMessageSent(t, ts, ethereumDomainCCTP, polygonDomainCCTP, 41)
	// Finalize events
	for i := 0; i < finalityDepth; i++ {
		ts.sb.Commit()
	}
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 51)

	// Need to replay as sometimes the logs are not picked up by the log poller (?)
	// Maybe another situation where chain reader doesn't register filters as expected.
	require.NoError(t, ts.lp.Replay(ctx, 1))

	tt := []struct {
		name           string
		tokens         map[reader.MessageTokenID]cciptypes.RampTokenAmount
		sourceChain    cciptypes.ChainSelector
		destChain      cciptypes.ChainSelector
		expectedMsgIDs []reader.MessageTokenID
	}{
		{
			name:           "empty messages should return empty response",
			tokens:         map[reader.MessageTokenID]cciptypes.RampTokenAmount{},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "single token message",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{reader.NewMessageTokenID(1, 1)},
		},
		{
			name: "single token message but different chain",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(31, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      polygonChain,
			expectedMsgIDs: []reader.MessageTokenID{reader.NewMessageTokenID(1, 2)},
		},
		{
			name: "message without matching nonce",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(1234, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "message without matching source domain",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, avalancheDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
		{
			name: "message with multiple tokens",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(21, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain: ethereumChain,
			destChain:   avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{
				reader.NewMessageTokenID(1, 1),
				reader.NewMessageTokenID(1, 2),
			},
		},
		{
			name: "message with multiple tokens, one without matching nonce",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(12, ethereumDomainCCTP).ToBytes(),
				},
				reader.NewMessageTokenID(1, 3): {
					ExtraData: reader.NewSourceTokenDataPayload(31, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain: ethereumChain,
			destChain:   avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{
				reader.NewMessageTokenID(1, 1),
				reader.NewMessageTokenID(1, 3),
			},
		},
		{
			name: "not finalized events are not returned",
			tokens: map[reader.MessageTokenID]cciptypes.RampTokenAmount{
				reader.NewMessageTokenID(1, 5): {
					ExtraData: reader.NewSourceTokenDataPayload(51, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []reader.MessageTokenID{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			hashes, err1 := usdcReader.MessageHashes(ctx, tc.sourceChain, tc.destChain, tc.tokens)
			require.NoError(t, err1)

			require.Equal(t, len(tc.expectedMsgIDs), len(hashes))

			for _, msgID := range tc.expectedMsgIDs {
				_, ok := hashes[msgID]
				require.True(t, ok)
			}
		})
	}
}

func emitMessageSent(t *testing.T, testEnv *testSetupData, source, dest uint32, nonce uint64) {
	payload := utils.RandomBytes32()
	_, err := testEnv.contract.EmitMessageSent(
		testEnv.auth,
		reader.CCTPMessageVersion,
		source,
		dest,
		utils.RandomBytes32(),
		utils.RandomBytes32(),
		[32]byte{},
		nonce,
		payload[:],
	)
	require.NoError(t, err)
	testEnv.sb.Commit()
}

func testSetup(ctx context.Context, t *testing.T, readerChain cciptypes.ChainSelector, cfg evmtypes.ChainReaderConfig, depth int) *testSetupData {
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

	address, _, _, err := usdc_reader_tester.DeployUSDCReaderTester(
		auth,
		simulatedBackend,
	)
	require.NoError(t, err)
	simulatedBackend.Commit()

	contract, err := usdc_reader_tester.NewUSDCReaderTester(address, simulatedBackend)
	require.NoError(t, err)

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(zapcore.ErrorLevel)
	db := pgtest.NewSqlxDB(t)
	lpOpts := logpoller.Opts{
		PollPeriod:               time.Millisecond,
		FinalityDepth:            int64(depth),
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
	require.NoError(t, lp.Start(ctx))

	cr, err := evm.NewChainReaderService(ctx, lggr, lp, headTracker, cl, cfg)
	require.NoError(t, err)

	err = cr.Start(ctx)
	require.NoError(t, err)

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
		cl:           cl,
		reader:       cr,
		lp:           lp,
	}
}

type testSetupData struct {
	contractAddr common.Address
	contract     *usdc_reader_tester.USDCReaderTester
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
	cl           client.Client
	reader       types.ContractReader
	lp           logpoller.LogPoller
}
