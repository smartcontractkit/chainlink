package memory

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zapcore"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	v2toml "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	configv2 "github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

func Context(tb testing.TB) context.Context {
	ctx := context.Background()
	var cancel func()
	switch t := tb.(type) {
	case *testing.T:
		if d, ok := t.Deadline(); ok {
			ctx, cancel = context.WithDeadline(ctx, d)
		}
	}
	if cancel == nil {
		ctx, cancel = context.WithCancel(ctx)
	}
	tb.Cleanup(cancel)
	return ctx
}

type Node struct {
	App chainlink.Application
	// Transmitter key/OCR keys for this node
	Keys       Keys
	Addr       net.TCPAddr
	IsBoostrap bool
}

func (n Node) ReplayLogs(chains map[uint64]uint64) error {
	for sel, block := range chains {
		chainID, _ := chainsel.ChainIdFromSelector(sel)
		if err := n.App.ReplayFromBlock(big.NewInt(int64(chainID)), block, false); err != nil {
			return err
		}
	}
	return nil
}

type RegistryConfig struct {
	EVMChainID uint64
	Contract   common.Address
}

// Creates a CL node which is:
// - Configured for OCR
// - Configured for the chains specified
// - Transmitter keys funded.
func NewNode(
	t *testing.T,
	port int, // Port for the P2P V2 listener.
	chains map[uint64]EVMChain,
	logLevel zapcore.Level,
	bootstrap bool,
	registryConfig RegistryConfig,
) *Node {
	// Do not want to load fixtures as they contain a dummy chainID.
	// Create database and initial configuration.
	cfg, db := heavyweight.FullTestDBNoFixturesV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		// P2P V2 configs.
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = config.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = config.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}

		// Enable Capabilities, This is a pre-requisite for registrySyncer to work.
		if registryConfig.Contract != common.HexToAddress("0x0") {
			c.Capabilities.ExternalRegistry.NetworkID = ptr(relay.NetworkEVM)
			c.Capabilities.ExternalRegistry.ChainID = ptr(strconv.FormatUint(uint64(registryConfig.EVMChainID), 10))
			c.Capabilities.ExternalRegistry.Address = ptr(registryConfig.Contract.String())
		}

		// OCR configs
		c.OCR.Enabled = ptr(false)
		c.OCR.DefaultTransactionQueueDepth = ptr(uint32(200))
		c.OCR2.Enabled = ptr(true)
		c.OCR2.ContractPollInterval = config.MustNewDuration(5 * time.Second)

		c.Log.Level = ptr(configv2.LogLevel(logLevel))

		var chainConfigs v2toml.EVMConfigs
		for chainID := range chains {
			chainConfigs = append(chainConfigs, createConfigV2Chain(chainID))
		}
		c.EVM = chainConfigs
	})

	// Set logging.
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(logLevel)

	// Create clients for the core node backed by sim.
	clients := make(map[uint64]client.Client)
	for chainID, chain := range chains {
		clients[chainID] = client.NewSimulatedBackendClient(t, chain.Backend, big.NewInt(int64(chainID)))
	}

	// Create keystore
	master := keystore.New(db, utils.FastScryptParams, lggr)
	kStore := KeystoreSim{
		eks: &EthKeystoreSim{
			Eth: master.Eth(),
		},
		csa: master.CSA(),
	}

	// Build evm factory using clients + keystore.
	mailMon := mailbox.NewMonitor("node", lggr.Named("mailbox"))
	evmOpts := chainlink.EVMFactoryConfig{
		ChainOpts: legacyevm.ChainOpts{
			AppConfig: cfg,
			GenEthClient: func(i *big.Int) client.Client {
				ethClient, ok := clients[i.Uint64()]
				if !ok {
					t.Fatal("no backend for chainID", i)
				}
				return ethClient
			},
			MailMon: mailMon,
			DS:      db,
		},
		CSAETHKeystore: kStore,
	}

	// Build relayer factory with EVM.
	relayerFactory := chainlink.RelayerFactory{
		Logger:       lggr,
		LoopRegistry: plugins.NewLoopRegistry(lggr.Named("LoopRegistry"), cfg.Tracing()),
		GRPCOpts:     loop.GRPCOpts{},
	}
	initOps := []chainlink.CoreRelayerChainInitFunc{chainlink.InitEVM(context.Background(), relayerFactory, evmOpts)}
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
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	keys := CreateKeys(t, app, chains)

	return &Node{
		App:        app,
		Keys:       keys,
		Addr:       net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: port},
		IsBoostrap: bootstrap,
	}
}

type Keys struct {
	PeerID                   p2pkey.PeerID
	TransmittersByEVMChainID map[uint64]common.Address
	OCRKeyBundle             ocr2key.KeyBundle
}

func CreateKeys(t *testing.T,
	app chainlink.Application, chains map[uint64]EVMChain) Keys {
	ctx := Context(t)
	require.NoError(t, app.GetKeyStore().Unlock(ctx, "password"))
	_, err := app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()
	// create a transmitter for each chain
	transmitters := make(map[uint64]common.Address)
	for chainID, chain := range chains {
		cid := big.NewInt(int64(chainID))
		addrs, err2 := app.GetKeyStore().Eth().EnabledAddressesForChain(Context(t), cid)
		require.NoError(t, err2)
		if len(addrs) == 1 {
			// just fund the address
			fundAddress(t, chain.DeployerKey, addrs[0], assets.Ether(10).ToInt(), chain.Backend)
			transmitters[chainID] = addrs[0]
		} else {
			// create key and fund it
			_, err3 := app.GetKeyStore().Eth().Create(Context(t), cid)
			require.NoError(t, err3, "failed to create key for chain", chainID)
			sendingKeys, err3 := app.GetKeyStore().Eth().EnabledAddressesForChain(Context(t), cid)
			require.NoError(t, err3)
			require.Len(t, sendingKeys, 1)
			fundAddress(t, chain.DeployerKey, sendingKeys[0], assets.Ether(10).ToInt(), chain.Backend)
			transmitters[chainID] = sendingKeys[0]
		}
	}
	require.Len(t, transmitters, len(chains))

	keybundle, err := app.GetKeyStore().OCR2().Create(ctx, chaintype.EVM)
	require.NoError(t, err)
	return Keys{
		PeerID:                   peerID,
		TransmittersByEVMChainID: transmitters,
		OCRKeyBundle:             keybundle,
	}
}

func createConfigV2Chain(chainID uint64) *v2toml.EVMConfig {
	chainIDBig := evmutils.NewI(int64(chainID))
	chain := v2toml.Defaults(chainIDBig)
	chain.GasEstimator.LimitDefault = ptr(uint64(5e6))
	chain.LogPollInterval = config.MustNewDuration(1000 * time.Millisecond)
	chain.Transactions.ForwardersEnabled = ptr(false)
	chain.FinalityDepth = ptr(uint32(2))
	return &v2toml.EVMConfig{
		ChainID: chainIDBig,
		Enabled: ptr(true),
		Chain:   chain,
		Nodes:   v2toml.EVMNodes{&v2toml.Node{}},
	}
}

func ptr[T any](v T) *T { return &v }

var _ keystore.Eth = &EthKeystoreSim{}

type EthKeystoreSim struct {
	keystore.Eth
}

// override
func (e *EthKeystoreSim) SignTx(ctx context.Context, address common.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error) {
	// always sign with chain id 1337 for the simulated backend
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
