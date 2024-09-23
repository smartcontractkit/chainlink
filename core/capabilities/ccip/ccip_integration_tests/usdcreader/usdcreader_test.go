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
	sel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/chainlink-ccip/execute/exectypes"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

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
	ctx := testutils.Context(t)
	ethereumChain := cciptypes.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)
	ethereumDomainCCTP := reader.CCTPDestDomains[uint64(ethereumChain)]
	avalancheChain := cciptypes.ChainSelector(sel.AVALANCHE_MAINNET.Selector)
	avalancheDomainCCTP := reader.CCTPDestDomains[uint64(avalancheChain)]
	polygonChain := cciptypes.ChainSelector(sel.POLYGON_MAINNET.Selector)
	polygonDomainCCTP := reader.CCTPDestDomains[uint64(polygonChain)]

	cfg := evmtypes.ChainReaderConfig{
		Contracts: map[string]evmtypes.ChainContractReader{
			consts.ContractNameCCTPMessageTransmitter: {
				ContractPollingFilter: evmtypes.ContractPollingFilter{
					GenericEventNames: []string{consts.EventNameCCTPMessageSent},
				},
				ContractABI: usdc_reader_tester.USDCReaderTesterABI,
				Configs: map[string]*evmtypes.ChainReaderDefinition{
					consts.EventNameCCTPMessageSent: {
						ChainSpecificName: consts.EventNameCCTPMessageSent,
						ReadType:          evmtypes.Event,
					},
				},
			},
		},
	}

	ts := testSetup(ctx, t, ethereumChain, cfg)

	usdcReader, err := reader.NewUSDCMessageReader(
		map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
			ethereumChain: {
				SourceMessageTransmitterAddr: ts.contractAddr.String(),
			},
		},
		map[cciptypes.ChainSelector]types.ContractReader{
			ethereumChain: ts.reader,
		})
	require.NoError(t, err)

	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 11)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 21)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 31)
	emitMessageSent(t, ts, ethereumDomainCCTP, avalancheDomainCCTP, 41)
	emitMessageSent(t, ts, ethereumDomainCCTP, polygonDomainCCTP, 31)
	emitMessageSent(t, ts, ethereumDomainCCTP, polygonDomainCCTP, 41)

	tt := []struct {
		name           string
		tokens         map[exectypes.MessageTokenID]cciptypes.RampTokenAmount
		sourceChain    cciptypes.ChainSelector
		destChain      cciptypes.ChainSelector
		expectedMsgIDs []exectypes.MessageTokenID
	}{
		{
			name:           "empty messages should return empty response",
			tokens:         map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []exectypes.MessageTokenID{},
		},
		{
			name: "single token message",
			tokens: map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
				exectypes.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []exectypes.MessageTokenID{exectypes.NewMessageTokenID(1, 1)},
		},
		{
			name: "single token message but different chain",
			tokens: map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
				exectypes.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(31, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      polygonChain,
			expectedMsgIDs: []exectypes.MessageTokenID{exectypes.NewMessageTokenID(1, 2)},
		},
		{
			name: "message without matching nonce",
			tokens: map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
				exectypes.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(1234, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []exectypes.MessageTokenID{},
		},
		{
			name: "message without matching source domain",
			tokens: map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
				exectypes.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, avalancheDomainCCTP).ToBytes(),
				},
			},
			sourceChain:    ethereumChain,
			destChain:      avalancheChain,
			expectedMsgIDs: []exectypes.MessageTokenID{},
		},
		{
			name: "message with multiple tokens",
			tokens: map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
				exectypes.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
				exectypes.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(21, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain: ethereumChain,
			destChain:   avalancheChain,
			expectedMsgIDs: []exectypes.MessageTokenID{
				exectypes.NewMessageTokenID(1, 1),
				exectypes.NewMessageTokenID(1, 2),
			},
		},
		{
			name: "message with multiple tokens, one without matching nonce",
			tokens: map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
				exectypes.NewMessageTokenID(1, 1): {
					ExtraData: reader.NewSourceTokenDataPayload(11, ethereumDomainCCTP).ToBytes(),
				},
				exectypes.NewMessageTokenID(1, 2): {
					ExtraData: reader.NewSourceTokenDataPayload(12, ethereumDomainCCTP).ToBytes(),
				},
				exectypes.NewMessageTokenID(1, 3): {
					ExtraData: reader.NewSourceTokenDataPayload(31, ethereumDomainCCTP).ToBytes(),
				},
			},
			sourceChain: ethereumChain,
			destChain:   avalancheChain,
			expectedMsgIDs: []exectypes.MessageTokenID{
				exectypes.NewMessageTokenID(1, 1),
				exectypes.NewMessageTokenID(1, 3),
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			require.Eventually(t, func() bool {
				hashes, err1 := usdcReader.MessageHashes(ctx, tc.sourceChain, tc.destChain, tc.tokens)
				require.NoError(t, err1)

				if len(tc.expectedMsgIDs) != len(hashes) {
					return false
				}

				for _, msgID := range tc.expectedMsgIDs {
					if _, ok := hashes[msgID]; !ok {
						return false
					}
				}
				return true
			}, 2*time.Second, 100*time.Millisecond)
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

func testSetup(ctx context.Context, t *testing.T, readerChain cciptypes.ChainSelector, cfg evmtypes.ChainReaderConfig) *testSetupData {
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
	}
}

type testSetupData struct {
	contractAddr common.Address
	contract     *usdc_reader_tester.USDCReaderTester
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
	cl           client.Client
	reader       types.ContractReader
}
