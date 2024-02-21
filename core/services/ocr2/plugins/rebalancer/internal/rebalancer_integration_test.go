package internal_test

import (
	"context"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
	"github.com/hashicorp/consul/sdk/freeport"
	chainsel "github.com/smartcontractkit/chain-selectors"
	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-common/pkg/loop"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/mailbox"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/client"
	v2toml "github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	evmutils "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/legacyevm"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/arm_proxy_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/lock_release_token_pool"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/mock_arm_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/router"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ccip/generated/weth9"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/mock_l1_bridge_adapter"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/rebalancer/generated/rebalancer"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/logger/audit"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/rebalancer/bridge/testonlybridge"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"

	"github.com/smartcontractkit/chainlink/v2/core/utils"
	"github.com/smartcontractkit/chainlink/v2/plugins"
)

var (
	mainChainID = int64(chainsel.GETH_TESTNET.EvmChainID)
)

func TestRebalancer_Integration(t *testing.T) {
	newTestUniverse(t, 2)
}

type ocr3Node struct {
	app          chainlink.Application
	peerID       string
	transmitters map[int64]common.Address
	keybundle    ocr2key.KeyBundle
}

type onchainUniverse struct {
	backend         *backends.SimulatedBackend
	chainID         uint64
	wethToken       *weth9.WETH9
	lockReleasePool *lock_release_token_pool.LockReleaseTokenPool
	rebalancer      *rebalancer.Rebalancer
	bridgeAdapter   *mock_l1_bridge_adapter.MockL1BridgeAdapter
}

func setupNodeOCR3(
	t *testing.T,
	owner *bind.TransactOpts,
	port int,
	chainIDToBackend map[int64]*backends.SimulatedBackend,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
	useForwarders bool,
) *ocr3Node {
	// Do not want to load fixtures as they contain a dummy chainID.
	config, db := heavyweight.FullTestDBNoFixturesV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = config.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = config.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.OCR.Enabled = ptr(false)
		c.OCR.DefaultTransactionQueueDepth = ptr(uint32(200))
		c.OCR2.Enabled = ptr(true)

		c.EVM[0].LogPollInterval = config.MustNewDuration(500 * time.Millisecond)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](3_500_000)
		c.EVM[0].Transactions.ForwardersEnabled = &useForwarders
		c.OCR2.ContractPollInterval = config.MustNewDuration(5 * time.Second)

		var chains v2toml.EVMConfigs
		for chainID := range chainIDToBackend {
			chains = append(chains, createConfigV2Chain(big.NewInt(chainID)))
		}
		c.EVM = chains
		c.OCR2.ContractPollInterval = config.MustNewDuration(5 * time.Second)

	})

	lggr := logger.TestLogger(t)
	clients := make(map[int64]client.Client)

	for chainID, backend := range chainIDToBackend {
		clients[chainID] = client.NewSimulatedBackendClient(t, backend, big.NewInt(chainID))
	}

	master := keystore.New(db, utils.FastScryptParams, lggr, config.Database())

	keystore := KeystoreSim{
		eks: &EthKeystoreSim{
			Eth: master.Eth(),
			t:   t,
		},
		csa: master.CSA(),
	}
	mailMon := mailbox.NewMonitor("Rebalancer", lggr.Named("mailbox"))
	evmOpts := chainlink.EVMFactoryConfig{
		ChainOpts: legacyevm.ChainOpts{
			AppConfig: config,
			GenEthClient: func(i *big.Int) client.Client {
				t.Log("genning eth client for chain id:", i.String())
				client, ok := clients[i.Int64()]
				if !ok {
					t.Fatal("no backend for chainID", i)
				}
				return client
			},
			MailMon: mailMon,
			DB:      db,
		},
		CSAETHKeystore: keystore,
	}
	relayerFactory := chainlink.RelayerFactory{
		Logger:       lggr,
		LoopRegistry: plugins.NewLoopRegistry(lggr.Named("LoopRegistry"), config.Tracing()),
		GRPCOpts:     loop.GRPCOpts{},
	}
	initOps := []chainlink.CoreRelayerChainInitFunc{chainlink.InitEVM(testutils.Context(t), relayerFactory, evmOpts)}
	rci, err := chainlink.NewCoreRelayerChainInteroperators(initOps...)
	require.NoError(t, err)

	app, err := chainlink.NewApplication(chainlink.ApplicationOpts{
		Config:                     config,
		SqlxDB:                     db,
		KeyStore:                   master,
		RelayerChainInteroperators: rci,
		Logger:                     lggr,
		ExternalInitiatorManager:   nil,
		CloseLogger:                lggr.Sync,
		UnrestrictedHTTPClient:     &http.Client{},
		RestrictedHTTPClient:       &http.Client{},
		AuditLogger:                audit.NoopLogger,
		MailMon:                    mailMon,
		LoopRegistry:               plugins.NewLoopRegistry(lggr, config.Tracing()),
	})
	require.NoError(t, err)
	require.NoError(t, app.GetKeyStore().Unlock("password"))
	_, err = app.GetKeyStore().P2P().Create()
	require.NoError(t, err)

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()

	// create a transmitter for each chain
	transmitters := make(map[int64]common.Address)
	for chainID, backend := range chainIDToBackend {
		addrs, err2 := app.GetKeyStore().Eth().EnabledAddressesForChain(big.NewInt(chainID))
		require.NoError(t, err2)
		if len(addrs) == 1 {
			// just fund the address
			fundAddress(t, owner, addrs[0], assets.Ether(10).ToInt(), backend)
			transmitters[chainID] = addrs[0]
		} else {
			// create key and fund it
			_, err3 := app.GetKeyStore().Eth().Create(big.NewInt(chainID))
			require.NoError(t, err3, "failed to create key for chain", chainID)
			sendingKeys, err3 := app.GetKeyStore().Eth().EnabledAddressesForChain(big.NewInt(chainID))
			require.NoError(t, err3)
			require.Len(t, sendingKeys, 1)
			fundAddress(t, owner, sendingKeys[0], assets.Ether(10).ToInt(), backend)
			transmitters[chainID] = sendingKeys[0]
		}
	}
	require.Len(t, transmitters, len(chainIDToBackend))

	keybundle, err := app.GetKeyStore().OCR2().Create(chaintype.EVM)
	require.NoError(t, err)

	return &ocr3Node{
		app:          app,
		peerID:       peerID.Raw(),
		transmitters: transmitters,
		keybundle:    keybundle,
	}
}

func newTestUniverse(t *testing.T, numChains int) {
	// create chains and deploy contracts
	owner, chains := createChains(t, numChains)
	universes := deployContracts(t, owner, chains)
	createConnectedNetwork(t, owner, chains, universes)
	transferBalances(t, owner, universes)
	mainContract := universes[mainChainID].rebalancer.Address()

	t.Log("Creating bootstrap node")
	bootstrapNodePort := freeport.GetOne(t)
	bootstrapNode := setupNodeOCR3(t, owner, bootstrapNodePort, chains, nil, false)
	numNodes := 4

	t.Log("creating ocr3 nodes")
	var (
		oracles        = make(map[int64][]confighelper2.OracleIdentityExtra)
		transmitters   = make(map[int64][]common.Address)
		onchainPubKeys []common.Address
		kbs            []ocr2key.KeyBundle
		apps           []chainlink.Application
		nodes          []*ocr3Node
	)
	ports := freeport.GetN(t, numNodes)
	for i := 0; i < numNodes; i++ {
		// Supply the bootstrap IP and port as a V2 peer address
		bootstrappers := []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.peerID, Addrs: []string{
				fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort),
			}},
		}
		node := setupNodeOCR3(t, owner, ports[i], chains, bootstrappers, false)

		kbs = append(kbs, node.keybundle)
		apps = append(apps, node.app)
		for chainID, transmitter := range node.transmitters {
			transmitters[chainID] = append(transmitters[chainID], transmitter)
		}
		onchainPubKeys = append(onchainPubKeys, common.BytesToAddress(node.keybundle.PublicKey()))
		for chainID, transmitter := range node.transmitters {
			identity := confighelper2.OracleIdentityExtra{
				OracleIdentity: confighelper2.OracleIdentity{
					OnchainPublicKey:  node.keybundle.PublicKey(),
					TransmitAccount:   ocrtypes.Account(transmitter.Hex()),
					OffchainPublicKey: node.keybundle.OffchainPublicKey(),
					PeerID:            node.peerID,
				},
				ConfigEncryptionPublicKey: node.keybundle.ConfigEncryptionPublicKey(),
			}
			oracles[chainID] = append(oracles[chainID], identity)
		}
		nodes = append(nodes, node)
	}

	t.Log("starting ticker to commit blocks")
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	tickCtx, tickCancel := context.WithCancel(testutils.Context(t))
	defer tickCancel()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-tick.C:
				for _, backend := range chains {
					backend.Commit()
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

	t.Log("setting config")
	blocksBeforeConfig := setRebalancerConfigs(
		t,
		owner,
		universes,
		chains,
		onchainPubKeys,
		transmitters,
		oracles)
	mainFromBlock := blocksBeforeConfig[mainChainID]

	t.Log("adding bootstrap node job")
	err := bootstrapNode.app.Start(testutils.Context(t))
	require.NoError(t, err, "failed to start bootstrap node")
	t.Cleanup(func() {
		require.NoError(t, bootstrapNode.app.Stop())
	})

	evmChains := bootstrapNode.app.GetRelayers().LegacyEVMChains()
	require.NotNil(t, evmChains)
	require.Len(t, evmChains.Slice(), numChains)
	bootstrapJobSpec := fmt.Sprintf(
		`
type = "bootstrap"
name = "bootstrap"
contractConfigTrackerPollInterval = "1s"
relay = "evm"
schemaVersion = 1
contractID = "%s"
[relayConfig]
chainID = 1337
fromBlock = %d
`, mainContract.Hex(), mainFromBlock)
	t.Log("creating bootstrap job with spec:\n", bootstrapJobSpec)
	ocrJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(bootstrapJobSpec)
	require.NoError(t, err, "failed to validate bootstrap job")
	err = bootstrapNode.app.AddJobV2(testutils.Context(t), &ocrJob)
	require.NoError(t, err, "failed to add bootstrap job")

	t.Log("creating ocr3 jobs")
	for i := 0; i < numNodes; i++ {
		err = apps[i].Start(testutils.Context(t))
		require.NoError(t, err)
		tapp := apps[i]
		t.Cleanup(func() {
			require.NoError(t, tapp.Stop())
		})

		mainChain := mustGetChainByEvmID(t, testutils.SimulatedChainID.Int64())

		jobSpec := fmt.Sprintf(
			`
type                 	= "offchainreporting2"
schemaVersion        	= 1
name                 	= "rebalancer-integration-test"
maxTaskDuration      	= "30s"
contractID           	= "%s"
ocrKeyBundleID       	= "%s"
relay                	= "evm"
pluginType           	= "rebalancer"
transmitterID        	= "%s"
forwardingAllowed       = false
contractConfigTrackerPollInterval = "5s"

[relayConfig]
chainID              	= 1337
# This is the fromBlock for the main chain
fromBlock               = %d
[relayConfig.fromBlocks]
# these are the fromBlock values for the follower chains
%s

[pluginConfig]
liquidityManagerAddress = "%s"
liquidityManagerNetwork = "%d"
closePluginTimeoutSec = 10
[pluginConfig.rebalancerConfig]
type = "ping-pong"
`,
			mainContract.Hex(),
			kbs[i].ID(),
			nodes[i].transmitters[1337].Hex(),
			mainFromBlock,
			buildFollowerChainsFromBlocksToml(blocksBeforeConfig),
			mainContract.Hex(),
			mainChain.Selector)
		t.Log("Creating rebalancer job with spec:\n", jobSpec)
		ocrJob2, err2 := validate.ValidatedOracleSpecToml(
			apps[i].GetConfig().OCR2(),
			apps[i].GetConfig().Insecure(),
			jobSpec)
		require.NoError(t, err2, "failed to validate rebalancer job")
		err2 = apps[i].AddJobV2(testutils.Context(t), &ocrJob2)
		require.NoError(t, err2, "failed to add rebalancer job")
	}

	t.Log("waiting for a transmission")
	waitForTransmissions(t, universes)
}

func waitForTransmissions(
	t *testing.T,
	universes map[int64]onchainUniverse,
) {
	start := uint64(1)
	liquidityTransferredSink := make(chan *rebalancer.RebalancerLiquidityTransferred)
	var subs []event.Subscription
	for _, uni := range universes {
		sub, err := uni.rebalancer.WatchLiquidityTransferred(&bind.WatchOpts{
			Start: &start,
		}, liquidityTransferredSink, nil, nil, nil)
		require.NoError(t, err, "failed to create subscription")
		subs = append(subs, sub)
	}
	defer func() {
		for _, sub := range subs {
			sub.Unsubscribe()
		}
	}()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	sentEvents := map[string]struct{}{}
	for {
		select {
		case lt := <-liquidityTransferredSink:
			// determine if it's a send or receive event based on the BridgeReturnData field
			// if it's a send event, then the BridgeReturnData will not be empty
			if len(lt.BridgeReturnData) > 0 {
				// for the test bridges, bridge return data is just a nonce
				nonce, err := testonlybridge.UnpackBridgeSendReturnData(lt.BridgeReturnData)
				require.NoError(t, err)
				t.Log("received send event with nonce:", nonce, "tx hash:", lt.Raw.TxHash.String())
				sentEvents[nonce.String()] = struct{}{}
			} else {
				// for the test bridges, the bridge specific data is an amount and a nonce
				amount, nonce, err := testonlybridge.UnpackFinalizeBridgePayload(lt.BridgeSpecificData)
				require.NoError(t, err)
				t.Log("received receive event with amount:", amount, "nonce:", nonce, "tx hash:", lt.Raw.TxHash.String())
				if _, ok := sentEvents[nonce.String()]; !ok {
					t.Fatal("received receive event without corresponding send event")
				} else {
					t.Log("received corresponding receive event")
					return
				}
			}
		case <-ticker.C:
			t.Log("waiting for transmission or liquidity transferred event")
		}
	}
}

func setRebalancerConfig(
	t *testing.T,
	owner *bind.TransactOpts,
	wrapper *rebalancer.Rebalancer,
	chain *backends.SimulatedBackend,
	onchainPubKeys []common.Address,
	transmitters []common.Address,
	oracles []confighelper2.OracleIdentityExtra,
) (blockBeforeConfig int64) {
	beforeConfig, err := chain.BlockByNumber(testutils.Context(t), nil)
	require.NoError(t, err)

	// most of the config on the follower chains does not matter
	// except for signers + transmitters
	var schedule []int
	for range oracles {
		schedule = append(schedule, 1)
	}
	offchainConfig, onchainConfig := []byte{}, []byte{}
	f := uint8(1)
	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		30*time.Second, // deltaProgress
		10*time.Second, // deltaResend
		20*time.Second, // deltaInitial
		2*time.Second,  // deltaRound
		20*time.Second, // deltaGrace
		10*time.Second, // deltaCertifiedCommitRequest
		10*time.Second, // deltaStage
		3,              // rmax
		schedule,
		oracles,
		offchainConfig,
		50*time.Millisecond, // maxDurationQuery
		5*time.Second,       // maxDurationObservation
		10*time.Second,      // maxDurationShouldAcceptAttestedReport
		10*time.Second,      // maxDurationShouldTransmitAcceptedReport
		int(f),
		onchainConfig)
	require.NoError(t, err, "failed to create contract config")
	_, err = wrapper.SetOCR3Config(
		owner,
		onchainPubKeys,
		transmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig)
	require.NoError(t, err, "failed to set config")
	chain.Commit()

	iter, err := wrapper.FilterConfigSet(&bind.FilterOpts{
		Start: beforeConfig.Number().Uint64(),
	})
	require.NoError(t, err, "failed to create ConfigSet filter")
	require.True(t, iter.Next())
	e := iter.Event
	require.Equal(t, onchainPubKeys, e.Signers, "signers do not match")
	require.Equal(t, transmitters, e.Transmitters, "transmitters do not match")
	t.Log("config digest for rebalancer at address: ", wrapper.Address(), ", is:", hexutil.Encode(e.ConfigDigest[:]))

	return beforeConfig.Number().Int64()
}

func setRebalancerConfigs(
	t *testing.T,
	owner *bind.TransactOpts,
	universes map[int64]onchainUniverse,
	chains map[int64]*backends.SimulatedBackend,
	onchainPubKeys []common.Address,
	transmitters map[int64][]common.Address,
	oracles map[int64][]confighelper2.OracleIdentityExtra) (blocksBeforeConfig map[int64]int64) {
	blocksBeforeConfig = make(map[int64]int64)
	for chainID, uni := range universes {
		blocksBeforeConfig[chainID] = setRebalancerConfig(
			t,
			owner,
			uni.rebalancer,
			chains[chainID],
			onchainPubKeys,
			transmitters[chainID],
			oracles[chainID],
		)
	}
	return
}

func ptr[T any](v T) *T { return &v }

func createConfigV2Chain(chainID *big.Int) *v2toml.EVMConfig {
	chain := v2toml.Defaults((*evmutils.Big)(chainID))
	chain.GasEstimator.LimitDefault = ptr(uint32(4e6))
	chain.LogPollInterval = config.MustNewDuration(500 * time.Millisecond)
	chain.Transactions.ForwardersEnabled = ptr(false)
	chain.FinalityDepth = ptr(uint32(2))
	return &v2toml.EVMConfig{
		ChainID: (*evmutils.Big)(chainID),
		Enabled: ptr(true),
		Chain:   chain,
		Nodes:   v2toml.EVMNodes{&v2toml.Node{}},
	}
}

var _ keystore.Eth = &EthKeystoreSim{}

type EthKeystoreSim struct {
	keystore.Eth
	t *testing.T
}

// override
func (e *EthKeystoreSim) SignTx(address common.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error) {
	// always sign with chain id 1337 for the simulated backend
	e.t.Log("always signing tx for chain id:", chainID.String(), "with chain id 1337, tx hash:", tx.Hash())
	return e.Eth.SignTx(address, tx, big.NewInt(1337))
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

func createChains(t *testing.T, numChains int) (owner *bind.TransactOpts, chains map[int64]*backends.SimulatedBackend) {
	owner = testutils.MustNewSimTransactor(t)
	chains = make(map[int64]*backends.SimulatedBackend)

	chains[mainChainID] = backends.NewSimulatedBackend(core.GenesisAlloc{
		owner.From: core.GenesisAccount{
			Balance: assets.Ether(10_000).ToInt(),
		},
	}, 30e6)

	for chainID := int64(chainsel.TEST_90000001.EvmChainID); chainID < int64(chainsel.TEST_90000020.EvmChainID); chainID++ {
		chains[chainID] = backends.NewSimulatedBackend(core.GenesisAlloc{
			owner.From: core.GenesisAccount{
				Balance: assets.Ether(10000).ToInt(),
			},
		}, 30e6)

		if len(chains) == numChains {
			break
		}
	}
	return
}

func deployContracts(
	t *testing.T,
	owner *bind.TransactOpts,
	chains map[int64]*backends.SimulatedBackend,
) (
	universes map[int64]onchainUniverse,
) {
	universes = make(map[int64]onchainUniverse)
	for chainID, backend := range chains {
		// Deploy wrapped ether contract
		// will act as the ERC-20 being bridged
		wethAddress, _, _, err := weth9.DeployWETH9(owner, backend)
		require.NoError(t, err, "failed to deploy WETH9 contract")
		backend.Commit()
		wethToken, err := weth9.NewWETH9(wethAddress, backend)
		require.NoError(t, err, "failed to create WETH9 wrapper")

		// deposit some eth into the weth contract
		_, err = wethToken.Deposit(&bind.TransactOpts{
			From:    owner.From,
			Signer:  owner.Signer,
			Value:   assets.Ether(100).ToInt(),
			Context: testutils.Context(t),
		})
		require.NoError(t, err, "failed to deposit eth into weth contract")

		// deploy arm and arm proxy.
		// required by the token pool
		// otherwise not used by this test.
		armAddress, _, _, err := mock_arm_contract.DeployMockARMContract(owner, backend)
		require.NoError(t, err, "failed to deploy MockARMContract contract")
		backend.Commit()
		armProxyAddress, _, _, err := arm_proxy_contract.DeployARMProxyContract(owner, backend, armAddress)
		require.NoError(t, err, "failed to deploy ARMProxyContract contract")
		backend.Commit()

		routerAddress, _, _, err := router.DeployRouter(owner, backend, wethAddress, armProxyAddress)
		require.NoError(t, err, "failed to deploy Router contract")
		backend.Commit()

		// deploy lock/release pool targeting the weth9 contract
		lockReleasePoolAddress, _, _, err := lock_release_token_pool.DeployLockReleaseTokenPool(
			owner, backend, wethAddress, []common.Address{}, armProxyAddress, true, routerAddress)
		require.NoError(t, err, "failed to deploy LockReleaseTokenPool contract")
		backend.Commit()
		lockReleasePool, err := lock_release_token_pool.NewLockReleaseTokenPool(lockReleasePoolAddress, backend)
		require.NoError(t, err)

		// deploy the rebalancer and set the liquidity container to be the lock release pool
		ch := mustGetChainByEvmID(t, chainID)
		rebalancerAddr, _, _, err := rebalancer.DeployRebalancer(owner, backend, wethAddress, ch.Selector, lockReleasePoolAddress)
		require.NoError(t, err, "failed to deploy Rebalancer contract")
		rebalancer, err := rebalancer.NewRebalancer(rebalancerAddr, backend)
		require.NoError(t, err, "failed to create Rebalancer wrapper")

		// set the rebalancer of the lock release pool to be the just deployed rebalancer
		_, err = lockReleasePool.SetRebalancer(owner, rebalancerAddr)
		require.NoError(t, err, "failed to set rebalancer on lock/release pool")
		backend.Commit()
		actualRebalancer, err := lockReleasePool.GetRebalancer(&bind.CallOpts{Context: testutils.Context(t)})
		require.NoError(t, err)
		require.Equal(t, rebalancerAddr, actualRebalancer)

		// deploy the bridge adapter to point to the weth contract address
		bridgeAdapterAddress, _, _, err := mock_l1_bridge_adapter.DeployMockL1BridgeAdapter(owner, backend, wethAddress)
		require.NoError(t, err, "failed to deploy mock l1 bridge adapter")
		backend.Commit()
		bridgeAdapter, err := mock_l1_bridge_adapter.NewMockL1BridgeAdapter(bridgeAdapterAddress, backend)
		require.NoError(t, err)

		universes[chainID] = onchainUniverse{
			backend:         backend,
			chainID:         uint64(chainID),
			wethToken:       wethToken,
			lockReleasePool: lockReleasePool,
			rebalancer:      rebalancer,
			bridgeAdapter:   bridgeAdapter,
		}

		t.Log("deployed contracts for chain:", chainID,
			"weth:", wethAddress.Hex(),
			"lockReleasePool:", lockReleasePool.Address().Hex(),
			"rebalancer:", rebalancerAddr.Hex(),
			"bridgeAdapter:", bridgeAdapterAddress.Hex(),
		)
	}
	return
}

func buildFollowerChainsFromBlocksToml(fromBlocks map[int64]int64) string {
	var s string
	for chainID, fromBlock := range fromBlocks {
		if chainID == mainChainID {
			continue
		}
		s += fmt.Sprintf("%d = %d\n", chainID, fromBlock)
	}
	return s
}

func transferBalances(
	t *testing.T,
	owner *bind.TransactOpts,
	universes map[int64]onchainUniverse,
) {
	for _, uni := range universes {
		// move some weth to the bridge adapters
		// so that they can transfer it to the rebalancer
		// when it calls finalizeWithdrawal
		_, err := uni.wethToken.Transfer(owner, uni.bridgeAdapter.Address(), assets.Ether(5).ToInt())
		require.NoError(t, err, "failed to transfer weth to bridge adapter")
		uni.backend.Commit()
		// confirm balance
		bal, err := uni.wethToken.BalanceOf(&bind.CallOpts{Context: testutils.Context(t)}, uni.bridgeAdapter.Address())
		require.NoError(t, err)
		require.Equal(t, assets.Ether(5).ToInt(), bal)

		// move some weth to the lock/release pool
		// the LM will pull from this pool in order to send cross-chain
		_, err = uni.wethToken.Transfer(owner, uni.lockReleasePool.Address(), assets.Ether(5).ToInt())
		require.NoError(t, err, "failed to transfer weth to lock/release pool")
		uni.backend.Commit()
		// confirm balance
		bal, err = uni.wethToken.BalanceOf(&bind.CallOpts{Context: testutils.Context(t)}, uni.lockReleasePool.Address())
		require.NoError(t, err)
		require.Equal(t, assets.Ether(5).ToInt(), bal)

		// check the balance of the token pool through the rebalancer,
		// should be the same as the balance of the lock/release pool
		// retrieved above.
		bal, err = uni.rebalancer.GetLiquidity(&bind.CallOpts{Context: testutils.Context(t)})
		require.NoError(t, err)
		require.Equal(t, assets.Ether(5).ToInt(), bal)
	}
}

// create a connection from the main chain to all follower chains
// and from all follower chains to the main chain
// this is analogous to the main chain being an L1 and all other
// chains being L2's.
func createConnectedNetwork(
	t *testing.T,
	owner *bind.TransactOpts,
	chains map[int64]*backends.SimulatedBackend,
	universes map[int64]onchainUniverse,
) {
	for chainID, uni := range universes {
		if chainID == mainChainID {
			continue
		}
		// follower -> main connection
		remoteChain := mustGetChainByEvmID(t, mainChainID)

		_, err := uni.rebalancer.SetCrossChainRebalancer(
			owner,
			rebalancer.IRebalancerCrossChainRebalancerArgs{
				RemoteRebalancer:    universes[mainChainID].rebalancer.Address(),
				RemoteChainSelector: remoteChain.Selector,
				Enabled:             true,
				LocalBridge:         uni.bridgeAdapter.Address(),
				RemoteToken:         universes[mainChainID].wethToken.Address(),
			})
		require.NoError(t, err, "failed to SetCrossChainRebalancer from follower to main chain")
		chains[chainID].Commit()

		mgr, err := uni.rebalancer.GetCrossChainRebalancer(&bind.CallOpts{Context: testutils.Context(t)}, remoteChain.Selector)
		require.NoError(t, err)
		require.Equal(t, universes[mainChainID].rebalancer.Address(), mgr.RemoteRebalancer)
		require.Equal(t, uni.bridgeAdapter.Address(), mgr.LocalBridge)
		require.Equal(t, universes[mainChainID].wethToken.Address(), mgr.RemoteToken)
		require.True(t, mgr.Enabled)

		// main -> follower connection
		remoteChain = mustGetChainByEvmID(t, chainID)

		_, err = universes[mainChainID].rebalancer.SetCrossChainRebalancer(
			owner,
			rebalancer.IRebalancerCrossChainRebalancerArgs{
				RemoteRebalancer:    uni.rebalancer.Address(),
				RemoteChainSelector: remoteChain.Selector,
				Enabled:             true,
				LocalBridge:         universes[mainChainID].bridgeAdapter.Address(),
				RemoteToken:         uni.wethToken.Address(),
			})
		require.NoError(t, err, "failed to add neighbor from main to follower chain")
		chains[mainChainID].Commit()

		mgr, err = universes[mainChainID].rebalancer.GetCrossChainRebalancer(
			&bind.CallOpts{Context: testutils.Context(t)}, remoteChain.Selector)
		require.NoError(t, err)
		require.Equal(t, uni.rebalancer.Address(), mgr.RemoteRebalancer)
		require.Equal(t, universes[mainChainID].bridgeAdapter.Address(), mgr.LocalBridge)
		require.Equal(t, uni.wethToken.Address(), mgr.RemoteToken)
		require.True(t, mgr.Enabled)
	}

	// sanity check connections
	for chainID, uni := range universes {
		destChains, err := uni.rebalancer.GetSupportedDestChains(&bind.CallOpts{Context: testutils.Context(t)})
		require.NoError(t, err, "couldn't get supported dest chains")
		t.Log("num dest chains:", len(destChains), "dest chains:", destChains)
		if chainID == mainChainID {
			require.Len(t, destChains, len(universes)-1)
		} else {
			require.Len(t, destChains, 1)
		}
		mgrs, err := uni.rebalancer.GetAllCrossChainRebalancers(&bind.CallOpts{
			Context: testutils.Context(t),
		})
		require.NoError(t, err, "couldn't get all cross-chain liquidity managers")
		t.Log("chainID:", chainID, "num neighbors:", len(mgrs))
		if chainID == mainChainID {
			// should be connected to all follower chains
			require.Len(t, mgrs, len(universes)-1, "unexpected number of neighbors on main chain")
		} else {
			// should be connected to just the main chain
			require.Len(t, mgrs, 1, "unexpected number of neighbors on follower chain")
		}
	}
}

func mustGetChainByEvmID(t *testing.T, chainID int64) chainsel.Chain {
	ch, exists := chainsel.ChainByEvmChainID(uint64(chainID))
	require.True(t, exists)
	return ch
}
