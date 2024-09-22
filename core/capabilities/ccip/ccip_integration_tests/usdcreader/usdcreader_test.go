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
	"github.com/smartcontractkit/chainlink-ccip/execute/exectypes"
	"github.com/smartcontractkit/chainlink-ccip/pkg/consts"
	"github.com/smartcontractkit/chainlink-ccip/pkg/reader"
	"github.com/smartcontractkit/chainlink-ccip/pluginconfig"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	cciptypes "github.com/smartcontractkit/chainlink-common/pkg/types/ccipocr3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	sel "github.com/smartcontractkit/chain-selectors"

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

func Test_USDCReader(t *testing.T) {
	ctx := testutils.Context(t)
	sourceChain := cciptypes.ChainSelector(sel.ETHEREUM_MAINNET_OPTIMISM_1.Selector)
	sourceDomainCCTP := uint32(2)
	destChain := cciptypes.ChainSelector(sel.AVALANCHE_MAINNET.Selector)
	destDomainCCTP := uint32(1)

	cctpVersion := uint32(0)
	nonce := uint64(20)

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

	testEnv := testSetup(ctx, t, sourceChain, cfg, cctpVersion, sourceDomainCCTP)

	configs := map[cciptypes.ChainSelector]pluginconfig.USDCCCTPTokenConfig{
		sourceChain: {
			SourcePoolAddress:            "0xA",
			SourceMessageTransmitterAddr: testEnv.contractAddr.String(),
		},
	}

	perChainReader := map[cciptypes.ChainSelector]types.ContractReader{
		sourceChain: testEnv.reader,
	}

	token := reader.NewSourceTokenDataPayload(nonce, sourceDomainCCTP).ToBytes()

	payload := utils.RandomBytes32()

	usdcReader, err := reader.NewUSDCMessageReader(configs, perChainReader)
	require.NoError(t, err)

	_, err = testEnv.contract.EmitMessageSent(testEnv.auth, destDomainCCTP, utils.RandomBytes32(), utils.RandomBytes32(), [32]byte{}, nonce, payload[:])
	require.NoError(t, err)
	testEnv.sb.Commit()

	time.Sleep(5 * time.Second)

	hashes, err := usdcReader.MessageHashes(ctx, sourceChain, destChain, map[exectypes.MessageTokenID]cciptypes.RampTokenAmount{
		exectypes.NewMessageTokenID(10, 10): cciptypes.RampTokenAmount{
			ExtraData:         token,
			SourcePoolAddress: []byte("0xA"),
		},
	})
	assert.Len(t, hashes, 1)
}

func testSetup(ctx context.Context, t *testing.T, readerChain cciptypes.ChainSelector, cfg evmtypes.ChainReaderConfig, cctpVersion, localDomain uint32) *testSetupData {
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
		cctpVersion,
		localDomain,
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
		lp:           lp,
		cl:           cl,
		reader:       cr,
	}
}

type testSetupData struct {
	contractAddr common.Address
	contract     *usdc_reader_tester.USDCReaderTester
	sb           *backends.SimulatedBackend
	auth         *bind.TransactOpts
	lp           logpoller.LogPoller
	cl           client.Client
	reader       types.ContractReader
}
