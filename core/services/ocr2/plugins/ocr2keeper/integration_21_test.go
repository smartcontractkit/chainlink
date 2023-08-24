package ocr2keeper_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/v3/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	automationForwarderLogic "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_forwarder_logic"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/basic_upkeep_contract"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	registrylogica21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_a_wrapper_2_1"
	registrylogicb21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
	registry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestFilterNamesFromSpec21(t *testing.T) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	address := common.HexToAddress(hexutil.Encode(b))

	spec := &job.OCR2OracleSpec{
		PluginType: job.OCR2Keeper,
		ContractID: address.String(), // valid contract addr
	}

	names, err := ocr2keeper.FilterNamesFromSpec21(spec)
	require.NoError(t, err)

	assert.Len(t, names, 2)
	assert.Equal(t, logpoller.FilterName("KeepersRegistry TransmitEventProvider", address), names[0])
	assert.Equal(t, logpoller.FilterName("KeeperRegistry Events", address), names[1])

	spec = &job.OCR2OracleSpec{
		PluginType: job.OCR2Keeper,
		ContractID: "0x5431", // invalid contract addr
	}
	_, err = ocr2keeper.FilterNamesFromSpec21(spec)
	require.ErrorContains(t, err, "not a valid EIP55 formatted address")
}

func TestIntegration_KeeperPluginConditionalUpkeep(t *testing.T) {
	g := gomega.NewWithT(t)
	lggr := logger.TestLogger(t)

	// setup blockchain
	sergey := testutils.MustNewSimTransactor(t) // owns all the link
	steve := testutils.MustNewSimTransactor(t)  // registry owner
	carrol := testutils.MustNewSimTransactor(t) // upkeep owner
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(10000).ToInt()},
		steve.From:  {Balance: assets.Ether(10000).ToInt()},
		carrol.From: {Balance: assets.Ether(10000).ToInt()},
	}
	// Generate 5 keys for nodes (1 bootstrap + 4 ocr nodes) and fund them with ether
	var nodeKeys [5]ethkey.KeyV2
	for i := int64(0); i < 5; i++ {
		nodeKeys[i] = cltest.MustGenerateRandomKey(t)
		genesisData[nodeKeys[i].Address] = core.GenesisAccount{Balance: assets.Ether(1000).ToInt()}
	}

	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	stopMining := cltest.Mine(backend, 3*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	defer stopMining()

	// Deploy registry
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
	require.NoError(t, err)
	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
	require.NoError(t, err)
	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(2000000000000000000))
	require.NoError(t, err)
	registry := deployKeeper21Registry(t, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

	nodes := setupNodes(t, nodeKeys, registry, backend, steve)

	<-time.After(time.Second * 5)

	upkeeps := 1

	_, err = linkToken.Transfer(sergey, carrol.From, big.NewInt(0).Mul(oneHunEth, big.NewInt(int64(upkeeps+1))))
	require.NoError(t, err)

	// Register new upkeep
	upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend)
	require.NoError(t, err)
	registrationTx, err := registry.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, 0, []byte{}, []byte{}, []byte{})
	require.NoError(t, err)
	backend.Commit()
	upkeepID := getUpkeepIdFromTx21(t, registry, registrationTx, backend)

	// Fund the upkeep
	_, err = linkToken.Transfer(sergey, carrol.From, oneHunEth)
	require.NoError(t, err)
	_, err = linkToken.Approve(carrol, registry.Address(), oneHunEth)
	require.NoError(t, err)
	_, err = registry.AddFunds(carrol, upkeepID, oneHunEth)
	require.NoError(t, err)
	backend.Commit()

	// Set upkeep to be performed
	_, err = upkeepContract.SetBytesToSend(carrol, payload1)
	require.NoError(t, err)
	_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
	require.NoError(t, err)
	backend.Commit()

	lggr.Infow("Upkeep registered and funded", "upkeepID", upkeepID.String())

	// keeper job is triggered and payload is received
	receivedBytes := func() []byte {
		received, err2 := upkeepContract.ReceivedBytes(nil)
		require.NoError(t, err2)
		return received
	}
	g.Eventually(receivedBytes, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(payload1))

	checkPipelineRuns(t, nodes, 1)

	// change payload
	_, err = upkeepContract.SetBytesToSend(carrol, payload2)
	require.NoError(t, err)
	_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
	require.NoError(t, err)

	// observe 2nd job run and received payload changes
	g.Eventually(receivedBytes, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(payload2))
}

func TestIntegration_KeeperPluginLogUpkeep(t *testing.T) {
	t.Skip() // TODO: fix test (fails in CI)
	g := gomega.NewWithT(t)

	// setup blockchain
	sergey := testutils.MustNewSimTransactor(t) // owns all the link
	steve := testutils.MustNewSimTransactor(t)  // registry owner
	carrol := testutils.MustNewSimTransactor(t) // upkeep owner
	genesisData := core.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(10000).ToInt()},
		steve.From:  {Balance: assets.Ether(10000).ToInt()},
		carrol.From: {Balance: assets.Ether(10000).ToInt()},
	}
	// Generate 5 keys for nodes (1 bootstrap + 4 ocr nodes) and fund them with ether
	var nodeKeys [5]ethkey.KeyV2
	for i := int64(0); i < 5; i++ {
		nodeKeys[i] = cltest.MustGenerateRandomKey(t)
		genesisData[nodeKeys[i].Address] = core.GenesisAccount{Balance: assets.Ether(1000).ToInt()}
	}

	backend := cltest.NewSimulatedBackend(t, genesisData, uint32(ethconfig.Defaults.Miner.GasCeil))
	stopMining := cltest.Mine(backend, 3*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	defer stopMining()

	// Deploy registry
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend)
	require.NoError(t, err)
	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(60000000000))
	require.NoError(t, err)
	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend, 18, big.NewInt(2000000000000000000))
	require.NoError(t, err)
	registry := deployKeeper21Registry(t, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

	nodes := setupNodes(t, nodeKeys, registry, backend, steve)
	// wait for nodes to start
	// TODO: find a better way to do this
	<-time.After(time.Second * 10)

	upkeeps := 1

	_, err = linkToken.Transfer(sergey, carrol.From, big.NewInt(0).Mul(oneHunEth, big.NewInt(int64(upkeeps+1))))
	require.NoError(t, err)

	backend.Commit()

	ids, addrs, contracts := deployUpkeeps(t, backend, carrol, steve, linkToken, registry, upkeeps)
	require.Equal(t, upkeeps, len(ids))
	require.Equal(t, len(ids), len(contracts))
	require.Equal(t, len(ids), len(addrs))

	backend.Commit()

	emits := 10
	go emitEvents(testutils.Context(t), t, emits, contracts, carrol, func() {
		backend.Commit()
		time.Sleep(3 * time.Second)
	})

	listener, done := listenPerformed(t, backend, registry, ids, int64(1))
	g.Eventually(listener, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.BeTrue())
	done()

	runs := checkPipelineRuns(t, nodes, 1*len(nodes)) // TODO: TBD

	t.Run("recover logs", func(t *testing.T) {
		t.Skip() // TODO: fix test (fails in CI)

		addr, contract := addrs[0], contracts[0]
		upkeepID := registerUpkeep(t, registry, addr, carrol, steve, backend)
		backend.Commit()
		t.Logf("Registered new upkeep %s for address %s", upkeepID.String(), addr.String())
		// blockBeforeEmits := backend.Blockchain().CurrentBlock().Number.Uint64()
		// Emit 100 logs in a burst
		emits := 100
		i := 0
		emitEvents(testutils.Context(t), t, 100, []*log_upkeep_counter_wrapper.LogUpkeepCounter{contract}, carrol, func() {
			i++
			if i%(emits/4) == 0 {
				backend.Commit()
				time.Sleep(time.Millisecond * 250) // otherwise we get "invalid transaction nonce" errors
			}
		})
		// Mine enough blocks to ensre these logs don't fall into log provider range
		dummyBlocks := 500
		for i := 0; i < dummyBlocks; i++ {
			backend.Commit()
			time.Sleep(time.Millisecond * 10)
		}
		t.Logf("Mined %d blocks", dummyBlocks)

		// listener, done := listenPerformed(t, backend, registry, []*big.Int{upkeepID}, int64(blockBeforeEmits))
		// defer done()
		// g.Eventually(listener, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.BeTrue())

		expectedPostRecover := runs + emits // TODO: TBD
		waitPipelineRuns(t, nodes, expectedPostRecover, testutils.WaitTimeout(t), cltest.DBPollingInterval)

	})
}

func waitPipelineRuns(t *testing.T, nodes []Node, n int, timeout, interval time.Duration) {
	ctx, cancel := context.WithTimeout(testutils.Context(t), timeout)
	defer cancel()
	var allRuns []pipeline.Run
	for len(allRuns) < n && ctx.Err() == nil {
		allRuns = []pipeline.Run{}
		for _, node := range nodes {
			runs, err := node.App.PipelineORM().GetAllRuns()
			require.NoError(t, err)
			allRuns = append(allRuns, runs...)
		}
		time.Sleep(interval)
	}
	runs := len(allRuns)
	t.Logf("found %d pipeline runs", runs)
	require.GreaterOrEqual(t, runs, n)
}

func checkPipelineRuns(t *testing.T, nodes []Node, n int) int {
	var allRuns []pipeline.Run
	for _, node := range nodes {
		runs, err2 := node.App.PipelineORM().GetAllRuns()
		require.NoError(t, err2)
		allRuns = append(allRuns, runs...)
	}
	runs := len(allRuns)
	t.Logf("found %d pipeline runs", runs)
	require.GreaterOrEqual(t, runs, n)
	return runs
}

func emitEvents(ctx context.Context, t *testing.T, n int, contracts []*log_upkeep_counter_wrapper.LogUpkeepCounter, carrol *bind.TransactOpts, afterEmit func()) {
	for i := 0; i < n && ctx.Err() == nil; i++ {
		for _, contract := range contracts {
			// t.Logf("[automation-ocr3 | EvmRegistry] calling upkeep contracts to emit events. run: %d; contract addr: %s", i+1, contract.Address().Hex())
			_, err := contract.Start(carrol)
			require.NoError(t, err)
		}
		afterEmit()
	}
}

func mapListener(m *sync.Map, n int) func() bool {
	return func() bool {
		count := 0
		m.Range(func(key, value interface{}) bool {
			count++
			return true
		})
		return count > n
	}
}

func listenPerformed(t *testing.T, backend *backends.SimulatedBackend, registry *iregistry21.IKeeperRegistryMaster, ids []*big.Int, startBlock int64) (func() bool, func()) {
	cache := &sync.Map{}
	ctx, cancel := context.WithCancel(testutils.Context(t))
	start := startBlock
	go func() {
		for ctx.Err() == nil {
			bl := backend.Blockchain().CurrentBlock().Number.Uint64()
			sc := make([]bool, len(ids))
			for i := range sc {
				sc[i] = true
			}
			iter, err := registry.FilterUpkeepPerformed(&bind.FilterOpts{
				Start:   uint64(start),
				End:     &bl,
				Context: ctx,
			}, ids, sc)
			if ctx.Err() != nil {
				return
			}
			require.NoError(t, err)
			for iter.Next() {
				if iter.Event != nil {
					t.Logf("[automation-ocr3 | EvmRegistry] upkeep performed event emitted for id %s", iter.Event.Id.String())
					cache.Store(iter.Event.Id.String(), true)
				}
			}
			require.NoError(t, iter.Close())
			time.Sleep(time.Second)
		}
	}()

	return mapListener(cache, 0), cancel
}

func setupNodes(t *testing.T, nodeKeys [5]ethkey.KeyV2, registry *iregistry21.IKeeperRegistryMaster, backend *backends.SimulatedBackend, usr *bind.TransactOpts) []Node {
	lggr := logger.TestLogger(t)
	// Setup bootstrap + oracle nodes
	bootstrapNodePort := int64(19599)
	appBootstrap, bootstrapPeerID, bootstrapTransmitter, bootstrapKb := setupNode(t, bootstrapNodePort, "bootstrap_keeper_ocr", nodeKeys[0], backend, nil)
	bootstrapNode := Node{
		appBootstrap, bootstrapTransmitter, bootstrapKb,
	}
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	// Set up the minimum 4 oracles all funded
	for i := int64(0); i < 4; i++ {
		app, peerID, transmitter, kb := setupNode(t, bootstrapNodePort+i+1, fmt.Sprintf("oracle_keeper%d", i), nodeKeys[i+1], backend, []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		})

		nodes = append(nodes, Node{
			app, transmitter, kb,
		})
		offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocrTypes.Account(transmitter.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}
	// Add the bootstrap job
	bootstrapNode.AddBootstrapJob(t, fmt.Sprintf(`
		type                              = "bootstrap"
		relay                             = "evm"
		schemaVersion                     = 1
		name                              = "boot"
		contractID                        = "%s"
		contractConfigTrackerPollInterval = "15s"

		[relayConfig]
		chainID = 1337
	`, registry.Address()))

	// Add OCR jobs
	for i, node := range nodes {
		node.AddJob(t, fmt.Sprintf(`
		type = "offchainreporting2"
		pluginType = "ocr2automation"
		relay = "evm"
		name = "ocr2keepers-%d"
		schemaVersion = 1
		contractID = "%s"
		contractConfigTrackerPollInterval = "15s"
		ocrKeyBundleID = "%s"
		transmitterID = "%s"
		p2pv2Bootstrappers = [
		  "%s"
		]

		[relayConfig]
		chainID = 1337

		[pluginConfig]
		maxServiceWorkers = 100
		cacheEvictionInterval = "1s"
		mercuryCredentialName = "%s"
		contractVersion = "v2.1"
		`, i, registry.Address(), node.KeyBundle.ID(), node.Transmitter, fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort), MercuryCredName))
	}

	// Setup config on contract
	configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint32 maxRevertDataSize, uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address[] registrars, address upkeepPrivilegeManager)")
	onchainConfig, err := abi.Encode(map[string]interface{}{
		"paymentPremiumPPB":      uint32(0),
		"flatFeeMicroLink":       uint32(0),
		"checkGasLimit":          uint32(6500000),
		"stalenessSeconds":       uint32(90000),
		"gasCeilingMultiplier":   uint16(2),
		"minUpkeepSpend":         uint32(0),
		"maxPerformGas":          uint32(5000000),
		"maxCheckDataSize":       uint32(5000),
		"maxPerformDataSize":     uint32(5000),
		"maxRevertDataSize":      uint32(5000),
		"fallbackGasPrice":       big.NewInt(60000000000),
		"fallbackLinkPrice":      big.NewInt(2000000000000000000),
		"transcoder":             testutils.NewAddress(),
		"registrars":             []common.Address{testutils.NewAddress()},
		"upkeepPrivilegeManager": testutils.NewAddress(),
	}, configType)
	require.NoError(t, err)
	rawCfg, err := json.Marshal(config.OffchainConfig{
		PerformLockoutWindow: 100 * 12 * 1000, // ~100 block lockout (on goerli)
		MinConfirmations:     1,
	})
	if err != nil {
		t.Logf("error creating off-chain config: %s", err)
		t.FailNow()
	}

	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		5*time.Second,         // deltaProgress time.Duration,
		10*time.Second,        // deltaResend time.Duration,
		100*time.Millisecond,  // deltaInitial time.Duration,
		1000*time.Millisecond, // deltaRound time.Duration,
		40*time.Millisecond,   // deltaGrace time.Duration,
		200*time.Millisecond,  // deltaRequestCertifiedCommit time.Duration,
		30*time.Second,        // deltaStage time.Duration,
		uint64(50),            // rMax uint8,
		[]int{1, 1, 1, 1},     // s []int,
		oracles,               // oracles []OracleIdentityExtra,
		rawCfg,                // reportingPluginConfig []byte,
		20*time.Millisecond,   // maxDurationQuery time.Duration,
		1600*time.Millisecond, // maxDurationObservation time.Duration,
		20*time.Millisecond,   // maxDurationShouldAcceptFinalizedReport time.Duration,
		20*time.Millisecond,   // maxDurationShouldTransmitAcceptedReport time.Duration,
		1,                     // f int,
		onchainConfig,         // onchainConfig []byte,
	)

	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := accountsToAddress(transmitters)
	require.NoError(t, err)

	lggr.Infow("Setting Config on Oracle Contract",
		"signerAddresses", signerAddresses,
		"transmitterAddresses", transmitterAddresses,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", offchainConfigVersion,
		"offchainConfig", offchainConfig,
	)
	_, err = registry.SetConfig(
		usr,
		signerAddresses,
		transmitterAddresses,
		threshold,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	backend.Commit()

	return nodes
}

func deployUpkeeps(t *testing.T, backend *backends.SimulatedBackend, carrol, steve *bind.TransactOpts, linkToken *link_token_interface.LinkToken, registry *iregistry21.IKeeperRegistryMaster, n int) ([]*big.Int, []common.Address, []*log_upkeep_counter_wrapper.LogUpkeepCounter) {
	ids := make([]*big.Int, n)
	addrs := make([]common.Address, n)
	contracts := make([]*log_upkeep_counter_wrapper.LogUpkeepCounter, n)
	for i := 0; i < n; i++ {
		backend.Commit()
		time.Sleep(1 * time.Second)
		upkeepAddr, _, upkeepContract, err := log_upkeep_counter_wrapper.DeployLogUpkeepCounter(
			carrol, backend,
			big.NewInt(100000),
		)
		require.NoError(t, err)

		upkeepID := registerUpkeep(t, registry, upkeepAddr, carrol, steve, backend)

		// Fund the upkeep
		_, err = linkToken.Approve(carrol, registry.Address(), oneHunEth)
		require.NoError(t, err)
		_, err = registry.AddFunds(carrol, upkeepID, oneHunEth)
		require.NoError(t, err)
		backend.Commit()

		ids[i] = upkeepID
		contracts[i] = upkeepContract
		addrs[i] = upkeepAddr
	}
	return ids, addrs, contracts
}

func registerUpkeep(t *testing.T, registry *iregistry21.IKeeperRegistryMaster, upkeepAddr common.Address, carrol, steve *bind.TransactOpts, backend *backends.SimulatedBackend) *big.Int {
	logTriggerConfigType := abi.MustNewType("tuple(address contractAddress, uint8 filterSelector, bytes32 topic0, bytes32 topic1, bytes32 topic2, bytes32 topic3)")
	logTriggerConfig, err := abi.Encode(map[string]interface{}{
		"contractAddress": upkeepAddr,
		"filterSelector":  0,                                                                    // no indexed topics filtered
		"topic0":          "0x3d53a39550e04688065827f3bb86584cb007ab9ebca7ebd528e7301c9c31eb5d", // event sig for Trigger()
		"topic1":          "0x",
		"topic2":          "0x",
		"topic3":          "0x",
	}, logTriggerConfigType)
	require.NoError(t, err)

	registrationTx, err := registry.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, 1, []byte{}, logTriggerConfig, []byte{})
	require.NoError(t, err)
	backend.Commit()
	upkeepID := getUpkeepIdFromTx21(t, registry, registrationTx, backend)

	return upkeepID
}

func deployKeeper21Registry(
	t *testing.T,
	auth *bind.TransactOpts,
	backend *backends.SimulatedBackend,
	linkAddr, linkFeedAddr,
	gasFeedAddr common.Address,
) *iregistry21.IKeeperRegistryMaster {
	automationForwarderLogicAddr, _, _, err := automationForwarderLogic.DeployAutomationForwarderLogic(auth, backend)
	require.NoError(t, err)
	backend.Commit()
	registryLogicBAddr, _, _, err := registrylogicb21.DeployKeeperRegistryLogicB(
		auth,
		backend,
		0, // Payment model
		linkAddr,
		linkFeedAddr,
		gasFeedAddr,
		automationForwarderLogicAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registryLogicAAddr, _, _, err := registrylogica21.DeployKeeperRegistryLogicA(
		auth,
		backend,
		registryLogicBAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registryAddr, _, _, err := registry21.DeployKeeperRegistry(
		auth,
		backend,
		registryLogicAAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registryMaster, err := iregistry21.NewIKeeperRegistryMaster(registryAddr, backend)
	require.NoError(t, err)

	return registryMaster
}

func getUpkeepIdFromTx21(t *testing.T, registry *iregistry21.IKeeperRegistryMaster, registrationTx *types.Transaction, backend *backends.SimulatedBackend) *big.Int {
	receipt, err := backend.TransactionReceipt(testutils.Context(t), registrationTx.Hash())
	require.NoError(t, err)
	parsedLog, err := registry.ParseUpkeepRegistered(*receipt.Logs[0])
	require.NoError(t, err)
	return parsedLog.Id
}
