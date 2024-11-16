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
	"golang.org/x/exp/maps"

	"github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink/deployment"

	"github.com/smartcontractkit/chainlink/v2/core/capabilities"
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
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay"
	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

type Node struct {
	App chainlink.Application
	// Transmitter key/OCR keys for this node
	Chains     []uint64 // chain selectors
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

// Creates a CL node which is:
// - Configured for OCR
// - Configured for the chains specified
// - Transmitter keys funded.
func NewNode(
	t *testing.T,
	port int, // Port for the P2P V2 listener.
	chains map[uint64]deployment.Chain,
	logLevel zapcore.Level,
	bootstrap bool,
	registryConfig deployment.CapabilityRegistryConfig,
) *Node {
	evmchains := make(map[uint64]EVMChain)
	for _, chain := range chains {
		evmChainID, err := chainsel.ChainIdFromSelector(chain.Selector)
		if err != nil {
			t.Fatal(err)
		}
		evmchains[evmChainID] = EVMChain{
			Backend:     chain.Client.(*Backend).Sim,
			DeployerKey: chain.DeployerKey,
		}
	}

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
		for chainID := range evmchains {
			chainConfigs = append(chainConfigs, createConfigV2Chain(chainID))
		}
		c.EVM = chainConfigs
	})

	// Set logging.
	lggr := logger.TestLogger(t)
	lggr.SetLogLevel(logLevel)

	// Create clients for the core node backed by sim.
	clients := make(map[uint64]client.Client)
	for chainID, chain := range evmchains {
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

	// Build Beholder auth
	ctx := tests.Context(t)
	require.NoError(t, master.Unlock(ctx, "password"))
	require.NoError(t, master.CSA().EnsureKey(ctx))
	beholderAuthHeaders, csaPubKeyHex, err := keystore.BuildBeholderAuth(master)
	require.NoError(t, err)

	// Build relayer factory with EVM.
	relayerFactory := chainlink.RelayerFactory{
		Logger:               lggr,
		LoopRegistry:         plugins.NewLoopRegistry(lggr.Named("LoopRegistry"), cfg.Tracing(), cfg.Telemetry(), beholderAuthHeaders, csaPubKeyHex),
		GRPCOpts:             loop.GRPCOpts{},
		CapabilitiesRegistry: capabilities.NewRegistry(lggr),
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
		LoopRegistry:               plugins.NewLoopRegistry(lggr, cfg.Tracing(), cfg.Telemetry(), beholderAuthHeaders, csaPubKeyHex),
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, db.Close())
	})
	keys := CreateKeys(t, app, chains)

	return &Node{
		App:        app,
		Chains:     maps.Keys(chains),
		Keys:       keys,
		Addr:       net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: port},
		IsBoostrap: bootstrap,
	}
}

type Keys struct {
	PeerID                   p2pkey.PeerID
	CSA                      csakey.KeyV2
	TransmittersByEVMChainID map[uint64]common.Address
	OCRKeyBundles            map[chaintype.ChainType]ocr2key.KeyBundle
}

func CreateKeys(t *testing.T,
	app chainlink.Application, chains map[uint64]deployment.Chain) Keys {
	ctx := tests.Context(t)
	_, err := app.GetKeyStore().P2P().Create(ctx)
	require.NoError(t, err)

	err = app.GetKeyStore().CSA().EnsureKey(ctx)
	require.NoError(t, err)
	csaKeys, err := app.GetKeyStore().CSA().GetAll()
	require.NoError(t, err)
	csaKey := csaKeys[0]

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()
	// create a transmitter for each chain
	transmitters := make(map[uint64]common.Address)
	keybundles := make(map[chaintype.ChainType]ocr2key.KeyBundle)
	for _, chain := range chains {
		family, err := chainsel.GetSelectorFamily(chain.Selector)
		require.NoError(t, err)

		var ctype chaintype.ChainType
		switch family {
		case chainsel.FamilyEVM:
			ctype = chaintype.EVM
		case chainsel.FamilySolana:
			ctype = chaintype.Solana
		case chainsel.FamilyStarknet:
			ctype = chaintype.StarkNet
		case chainsel.FamilyCosmos:
			ctype = chaintype.Cosmos
		case chainsel.FamilyAptos:
			ctype = chaintype.Aptos
		default:
			panic(fmt.Sprintf("Unsupported chain family %v", family))
		}

		keybundle, err := app.GetKeyStore().OCR2().Create(ctx, ctype)
		require.NoError(t, err)
		keybundles[ctype] = keybundle

		if family != chainsel.FamilyEVM {
			// TODO: only support EVM transmission keys for now
			continue
		}

		evmChainID, err := chainsel.ChainIdFromSelector(chain.Selector)
		require.NoError(t, err)

		cid := big.NewInt(int64(evmChainID))
		addrs, err2 := app.GetKeyStore().Eth().EnabledAddressesForChain(ctx, cid)
		require.NoError(t, err2)
		if len(addrs) == 1 {
			// just fund the address
			transmitters[evmChainID] = addrs[0]
		} else {
			// create key and fund it
			_, err3 := app.GetKeyStore().Eth().Create(ctx, cid)
			require.NoError(t, err3, "failed to create key for chain", evmChainID)
			sendingKeys, err3 := app.GetKeyStore().Eth().EnabledAddressesForChain(ctx, cid)
			require.NoError(t, err3)
			require.Len(t, sendingKeys, 1)
			transmitters[evmChainID] = sendingKeys[0]
		}
		backend := chain.Client.(*Backend).Sim
		fundAddress(t, chain.DeployerKey, transmitters[evmChainID], assets.Ether(1000).ToInt(), backend)
	}

	return Keys{
		PeerID:                   peerID,
		CSA:                      csaKey,
		TransmittersByEVMChainID: transmitters,
		OCRKeyBundles:            keybundles,
	}
}

func createConfigV2Chain(chainID uint64) *v2toml.EVMConfig {
	chainIDBig := evmutils.NewI(int64(chainID))
	chain := v2toml.Defaults(chainIDBig)
	chain.GasEstimator.LimitDefault = ptr(uint64(5e6))
	chain.LogPollInterval = config.MustNewDuration(500 * time.Millisecond)
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
