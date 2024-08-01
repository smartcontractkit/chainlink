package ccip_integration_tests

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities/ccip/validate"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/jmoiron/sqlx"

	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	v2toml "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	configv2 "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/plugins"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"
)

type ocr3Node struct {
	app          chainlink.Application
	peerID       string
	transmitters map[uint64]common.Address
	keybundle    ocr2key.KeyBundle
	db           *sqlx.DB
}

// setupNodeOCR3 creates a chainlink node and any associated keys in order to run
// ccip.
func setupNodeOCR3(
	t *testing.T,
	port int,
	universes map[uint64]onchainUniverse,
	homeChainUniverse homeChain,
	logLevel zapcore.Level,
) *ocr3Node {
	// Do not want to load fixtures as they contain a dummy chainID.
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		// P2P V2 configs.
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = config.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = config.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}

		// Enable Capabilities, This is a pre-requisite for registrySyncer to work.
		c.Capabilities.ExternalRegistry.NetworkID = ptr(relay.NetworkEVM)
		c.Capabilities.ExternalRegistry.ChainID = ptr(strconv.FormatUint(homeChainUniverse.chainID, 10))
		c.Capabilities.ExternalRegistry.Address = ptr(homeChainUniverse.capabilityRegistry.Address().String())

		// OCR configs
		c.OCR.Enabled = ptr(false)
		c.OCR.DefaultTransactionQueueDepth = ptr(uint32(200))
		c.OCR2.Enabled = ptr(true)
		c.OCR2.ContractPollInterval = config.MustNewDuration(5 * time.Second)

		c.Log.Level = ptr(configv2.LogLevel(logLevel))

		var chains v2toml.EVMConfigs
		for chainID := range universes {
			chains = append(chains, createConfigV2Chain(uBigInt(chainID)))
		}
		c.EVM = chains
	})

	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(logLevel)
	ctx := testutils.Context(t)
	clients := make(map[uint64]client.Client)

	for chainID, uni := range universes {
		clients[chainID] = client.NewSimulatedBackendClient(t, uni.backend, uBigInt(chainID))
	}

	master := keystore.New(db, utils.FastScryptParams, lggr)

	kStore := KeystoreSim{
		eks: &EthKeystoreSim{
			Eth: master.Eth(),
			t:   t,
		},
		csa: master.CSA(),
	}
	mailMon := mailbox.NewMonitor("ccip", lggr.Named("mailbox"))
	evmOpts := chainlink.EVMFactoryConfig{
		ChainOpts: legacyevm.ChainOpts{
			AppConfig: cfg,
			GenEthClient: func(i *big.Int) client.Client {
				t.Log("genning eth client for chain id:", i.String())
				client, ok := clients[i.Uint64()]
				if !ok {
					t.Fatal("no backend for chainID", i)
				}
				return client
			},
			MailMon: mailMon,
			DS:      db,
		},
		CSAETHKeystore: kStore,
	}
	relayerFactory := chainlink.RelayerFactory{
		Logger:       lggr,
		LoopRegistry: plugins.NewLoopRegistry(lggr.Named("LoopRegistry"), cfg.Tracing()),
		GRPCOpts:     loop.GRPCOpts{},
	}
	initOps := []chainlink.CoreRelayerChainInitFunc{chainlink.InitEVM(testutils.Context(t), relayerFactory, evmOpts)}
	rci, err := chainlink.NewCoreRelayerChainInteroperators(initOps...)
	require.NoError(t, err)

	app, err := chainlink.NewApplication(chainlink.ApplicationOpts{
		Config:                     cfg,
		DS:                         db,
		KeyStore:                   master,
		RelayerChainInteroperators: rci,
		Logger:                     lggr,
		ExternalInitiatorManager:   nil,
		CloseLogger:                lggr.Sync,
		UnrestrictedHTTPClient:     &http.Client{},
		RestrictedHTTPClient:       &http.Client{},
		AuditLogger:                audit.NoopLogger,
		MailMon:                    mailMon,
		LoopRegistry:               plugins.NewLoopRegistry(lggr, cfg.Tracing()),
	})
	require.NoError(t, err)
	require.NoError(t, app.GetKeyStore().Unlock(ctx, "password"))
	_, err = app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()
	// create a transmitter for each chain
	transmitters := make(map[uint64]common.Address)
	for chainID, uni := range universes {
		backend := uni.backend
		owner := uni.owner
		cID := uBigInt(chainID)
		addrs, err2 := app.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), cID)
		require.NoError(t, err2)
		if len(addrs) == 1 {
			// just fund the address
			fundAddress(t, owner, addrs[0], assets.Ether(10).ToInt(), backend)
			transmitters[chainID] = addrs[0]
		} else {
			// create key and fund it
			_, err3 := app.GetKeyStore().Eth().Create(testutils.Context(t), cID)
			require.NoError(t, err3, "failed to create key for chain", chainID)
			sendingKeys, err3 := app.GetKeyStore().Eth().EnabledAddressesForChain(testutils.Context(t), cID)
			require.NoError(t, err3)
			require.Len(t, sendingKeys, 1)
			fundAddress(t, owner, sendingKeys[0], assets.Ether(10).ToInt(), backend)
			transmitters[chainID] = sendingKeys[0]
		}
	}
	require.Len(t, transmitters, len(universes))

	keybundle, err := app.GetKeyStore().OCR2().Create(ctx, chaintype.EVM)
	require.NoError(t, err)

	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})

	return &ocr3Node{
		// can't use this app because it doesn't have the right toml config
		// missing bootstrapp
		app:          app,
		peerID:       peerID.Raw(),
		transmitters: transmitters,
		keybundle:    keybundle,
		db:           db,
	}
}

func ptr[T any](v T) *T { return &v }

var _ keystore.Eth = &EthKeystoreSim{}

type EthKeystoreSim struct {
	keystore.Eth
	t *testing.T
}

// override
func (e *EthKeystoreSim) SignTx(ctx context.Context, address common.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error) {
	// always sign with chain id 1337 for the simulated backend
	e.t.Log("always signing tx for chain id:", chainID.String(), "with chain id 1337, tx hash:", tx.Hash())
	return e.Eth.SignTx(ctx, address, tx, big.NewInt(1337))
}

type KeystoreSim struct {
	eks keystore.Eth
	csa keystore.CSA
}

func (e KeystoreSim) Eth() keystore.Eth {
	return e.eks
}

func (e KeystoreSim) CSA() keystore.CSA {
	return e.csa
}

func fundAddress(t *testing.T, from *bind.TransactOpts, to common.Address, amount *big.Int, backend *backends.SimulatedBackend) {
	nonce, err := backend.PendingNonceAt(testutils.Context(t), from.From)
	require.NoError(t, err)
	gp, err := backend.SuggestGasPrice(testutils.Context(t))
	require.NoError(t, err)
	rawTx := gethtypes.NewTx(&gethtypes.LegacyTx{
		Nonce:    nonce,
		GasPrice: gp,
		Gas:      21000,
		To:       &to,
		Value:    amount,
	})
	signedTx, err := from.Signer(from.From, rawTx)
	require.NoError(t, err)
	err = backend.SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	backend.Commit()
}

func createConfigV2Chain(chainID *big.Int) *v2toml.EVMConfig {
	chain := v2toml.Defaults((*evmutils.Big)(chainID))
	chain.GasEstimator.LimitDefault = ptr(uint64(5e6))
	chain.LogPollInterval = config.MustNewDuration(100 * time.Millisecond)
	chain.Transactions.ForwardersEnabled = ptr(false)
	chain.FinalityDepth = ptr(uint32(2))
	return &v2toml.EVMConfig{
		ChainID: (*evmutils.Big)(chainID),
		Enabled: ptr(true),
		Chain:   chain,
		Nodes:   v2toml.EVMNodes{&v2toml.Node{}},
	}
}

// Commit blocks periodically in the background for all chains
func commitBlocksBackground(t *testing.T, universes map[uint64]onchainUniverse, tick *time.Ticker) {
	t.Log("starting ticker to commit blocks")
	tickCtx, tickCancel := context.WithCancel(testutils.Context(t))
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-tick.C:
				for _, uni := range universes {
					uni.backend.Commit()
				}
			case <-tickCtx.Done():
				return
			}
		}
	}()
	t.Cleanup(func() {
		tickCancel()
		wg.Wait()
	})
}

// p2pKeyID: nodes p2p id
// ocrKeyBundleID: nodes ocr key bundle id
func mustGetJobSpec(t *testing.T, bootstrapP2PID p2pkey.PeerID, bootstrapPort int, p2pKeyID string, ocrKeyBundleID string) job.Job {
	specArgs := validate.SpecArgs{
		P2PV2Bootstrappers: []string{
			fmt.Sprintf("%s@127.0.0.1:%d", bootstrapP2PID.Raw(), bootstrapPort),
		},
		CapabilityVersion:      CapabilityVersion,
		CapabilityLabelledName: CapabilityLabelledName,
		OCRKeyBundleIDs: map[string]string{
			relay.NetworkEVM: ocrKeyBundleID,
		},
		P2PKeyID:     p2pKeyID,
		PluginConfig: map[string]any{},
	}
	specToml, err := validate.NewCCIPSpecToml(specArgs)
	require.NoError(t, err)
	jb, err := validate.ValidatedCCIPSpec(specToml)
	require.NoError(t, err)
	return jb
}
