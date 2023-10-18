package ocr2keeper_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
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
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/ocr2keepers/pkg/v3/config"

	relaytypes "github.com/smartcontractkit/chainlink-relay/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	automationForwarderLogic "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/automation_forwarder_logic"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/basic_upkeep_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/dummy_protocol_wrapper"
	iregistry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/i_keeper_registry_master_wrapper_2_1"
	registrylogica21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_a_wrapper_2_1"
	registrylogicb21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic_b_wrapper_2_1"
	registry21 "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper_2_1"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_triggered_streams_lookup_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/log_upkeep_counter_wrapper"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	evm21 "github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evm21"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
)

func TestFilterNamesFromSpec21(t *testing.T) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	address := common.HexToAddress(hexutil.Encode(b))

	spec := &job.OCR2OracleSpec{
		PluginType: relaytypes.OCR2Keeper,
		ContractID: address.String(), // valid contract addr
	}

	names, err := ocr2keeper.FilterNamesFromSpec21(spec)
	require.NoError(t, err)

	assert.Len(t, names, 2)
	assert.Equal(t, logpoller.FilterName("KeepersRegistry TransmitEventProvider", address), names[0])
	assert.Equal(t, logpoller.FilterName("KeeperRegistry Events", address), names[1])

	spec = &job.OCR2OracleSpec{
		PluginType: relaytypes.OCR2Keeper,
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

	nodes, _ := setupNodes(t, nodeKeys, registry, backend, steve)

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
	nodes, _ := setupNodes(t, nodeKeys, registry, backend, steve)
	upkeeps := 1

	_, err = linkToken.Transfer(sergey, carrol.From, big.NewInt(0).Mul(oneHunEth, big.NewInt(int64(upkeeps+1))))
	require.NoError(t, err)

	backend.Commit()

	ids, addrs, contracts := deployUpkeeps(t, backend, carrol, steve, linkToken, registry, upkeeps)
	require.Equal(t, upkeeps, len(ids))
	require.Equal(t, len(ids), len(contracts))
	require.Equal(t, len(ids), len(addrs))

	backend.Commit()

	emits := 1
	go emitEvents(testutils.Context(t), t, emits, contracts, carrol, func() {
		backend.Commit()
	})

	listener, done := listenPerformed(t, backend, registry, ids, int64(1))
	g.Eventually(listener, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.BeTrue())
	done()

	runs := checkPipelineRuns(t, nodes, 1)

	t.Run("recover logs", func(t *testing.T) {

		addr, contract := addrs[0], contracts[0]
		upkeepID := registerUpkeep(t, registry, addr, carrol, steve, backend)
		backend.Commit()
		t.Logf("Registered new upkeep %s for address %s", upkeepID.String(), addr.String())
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
		t.Logf("Mined %d blocks, waiting for logs to be recovered", dummyBlocks)

		expectedPostRecover := runs + emits
		waitPipelineRuns(t, nodes, expectedPostRecover, testutils.WaitTimeout(t), cltest.DBPollingInterval)

	})
}

func TestIntegration_KeeperPluginLogUpkeep_Retry(t *testing.T) {
	g := gomega.NewWithT(t)

	// setup blockchain
	linkOwner := testutils.MustNewSimTransactor(t)     // owns all the link
	registryOwner := testutils.MustNewSimTransactor(t) // registry owner
	upkeepOwner := testutils.MustNewSimTransactor(t)   // upkeep owner
	genesisData := core.GenesisAlloc{
		linkOwner.From:     {Balance: assets.Ether(10000).ToInt()},
		registryOwner.From: {Balance: assets.Ether(10000).ToInt()},
		upkeepOwner.From:   {Balance: assets.Ether(10000).ToInt()},
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
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(linkOwner, backend)
	require.NoError(t, err)

	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(registryOwner, backend, 18, big.NewInt(60000000000))
	require.NoError(t, err)

	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(registryOwner, backend, 18, big.NewInt(2000000000000000000))
	require.NoError(t, err)

	registry := deployKeeper21Registry(t, registryOwner, backend, linkAddr, linkFeedAddr, gasFeedAddr)

	nodes, mercuryServer := setupNodes(t, nodeKeys, registry, backend, registryOwner)

	const upkeepCount = 10
	const mercuryFailCount = upkeepCount * 3 * 2

	// testing with the mercury server involves mocking responses. currently,
	// there is not a way to connect a mercury call to an upkeep id (though we
	// could add custom headers) so the test must be fairly basic and just
	// count calls before switching to successes
	var (
		mu    sync.Mutex
		count int
	)

	mercuryServer.RegisterHandler(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		count++

		_ = r.ParseForm()

		t.Logf("MercuryHTTPServe:RequestURI: %s", r.RequestURI)

		for key, value := range r.Form {
			t.Logf("MercuryHTTPServe:FormValue: key: %s; value: %s;", key, value)
		}

		// the streams lookup retries against the remote server 3 times before
		// returning a result as retryable.
		// the simulation here should force the streams lookup process to return
		// retryable 2 times.
		// the total count of failures should be (upkeepCount * 3 * tryCount)
		if count <= mercuryFailCount {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		// start sending success messages
		output := `{"chainlinkBlob":"0x0001c38d71fed6c320b90e84b6f559459814d068e2a1700adc931ca9717d4fe70000000000000000000000000000000000000000000000000000000001a80b52b4bf1233f9cb71144a253a1791b202113c4ab4a92fa1b176d684b4959666ff8200000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000260000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001004254432d5553442d415242495452554d2d544553544e4554000000000000000000000000000000000000000000000000000000000000000000000000645570be000000000000000000000000000000000000000000000000000002af2b818dc5000000000000000000000000000000000000000000000000000002af2426faf3000000000000000000000000000000000000000000000000000002af32dc209700000000000000000000000000000000000000000000000000000000012130f8df0a9745bb6ad5e2df605e158ba8ad8a33ef8a0acf9851f0f01668a3a3f2b68600000000000000000000000000000000000000000000000000000000012130f60000000000000000000000000000000000000000000000000000000000000002c4a7958dce105089cf5edb68dad7dcfe8618d7784eb397f97d5a5fade78c11a58275aebda478968e545f7e3657aba9dcbe8d44605e4c6fde3e24edd5e22c94270000000000000000000000000000000000000000000000000000000000000002459c12d33986018a8959566d145225f0c4a4e61a9a3f50361ccff397899314f0018162cf10cd89897635a0bb62a822355bd199d09f4abe76e4d05261bb44733d"}`

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(output))
	})

	defer mercuryServer.Stop()

	_, err = linkToken.Transfer(linkOwner, upkeepOwner.From, big.NewInt(0).Mul(oneHunEth, big.NewInt(int64(upkeepCount+1))))
	require.NoError(t, err)

	backend.Commit()

	feeds, err := newFeedLookupUpkeepController(backend, registryOwner)
	require.NoError(t, err, "no error expected from creating a feed lookup controller")

	// deploy multiple upkeeps that listen to a log emitter and need to be
	// performed for each log event
	_ = feeds.DeployUpkeeps(t, backend, upkeepOwner, upkeepCount)
	_ = feeds.RegisterAndFund(t, registry, registryOwner, backend, linkToken)
	_ = feeds.EnableMercury(t, backend, registry, registryOwner)
	_ = feeds.VerifyEnv(t, backend, registry, registryOwner)

	// start emitting events in a separate go-routine
	// feed lookup relies on a single contract event log to perform multiple
	// listener contracts
	go func() {
		// only 1 event is necessary to make all 10 upkeeps eligible
		_ = feeds.EmitEvents(t, backend, 1, func() {
			// pause per emit for expected block production time
			time.Sleep(3 * time.Second)
		})
	}()

	listener, done := listenPerformed(t, backend, registry, feeds.UpkeepsIds(), int64(1))
	g.Eventually(listener, testutils.WaitTimeout(t)-(5*time.Second), cltest.DBPollingInterval).Should(gomega.BeTrue())

	done()

	_ = checkPipelineRuns(t, nodes, 1*len(nodes)) // TODO: TBD
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

func setupNodes(t *testing.T, nodeKeys [5]ethkey.KeyV2, registry *iregistry21.IKeeperRegistryMaster, backend *backends.SimulatedBackend, usr *bind.TransactOpts) ([]Node, *SimulatedMercuryServer) {
	lggr := logger.TestLogger(t)
	mServer := NewSimulatedMercuryServer()
	mServer.Start()

	// Setup bootstrap + oracle nodes
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, bootstrapTransmitter, bootstrapKb := setupNode(t, bootstrapNodePort, "bootstrap_keeper_ocr", nodeKeys[0], backend, nil, mServer)
	bootstrapNode := Node{
		appBootstrap, bootstrapTransmitter, bootstrapKb,
	}
	var (
		oracles []confighelper.OracleIdentityExtra
		nodes   []Node
	)
	// Set up the minimum 4 oracles all funded
	ports := freeport.GetN(t, 4)
	for i := 0; i < 4; i++ {
		app, peerID, transmitter, kb := setupNode(t, ports[i], fmt.Sprintf("oracle_keeper%d", i), nodeKeys[i+1], backend, []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, mServer)

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
		"upkeepPrivilegeManager": usr.From,
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

	return nodes, mServer
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

// ------- below this line could be added to a test helpers package
type registerAndFundFunc func(*testing.T, common.Address, *bind.TransactOpts, uint8, []byte) *big.Int

func registerAndFund(
	registry *iregistry21.IKeeperRegistryMaster,
	registryOwner *bind.TransactOpts,
	backend *backends.SimulatedBackend,
	linkToken *link_token_interface.LinkToken,
) registerAndFundFunc {
	return func(t *testing.T, upkeepAddr common.Address, upkeepOwner *bind.TransactOpts, trigger uint8, config []byte) *big.Int {
		// register the upkeep on the host registry contract
		registrationTx, err := registry.RegisterUpkeep(
			registryOwner,
			upkeepAddr,
			2_500_000,
			upkeepOwner.From,
			trigger,
			[]byte{},
			config,
			[]byte{},
		)
		require.NoError(t, err)

		backend.Commit()

		receipt, err := backend.TransactionReceipt(testutils.Context(t), registrationTx.Hash())
		require.NoError(t, err)

		parsedLog, err := registry.ParseUpkeepRegistered(*receipt.Logs[0])
		require.NoError(t, err)

		upkeepID := parsedLog.Id

		// Fund the upkeep
		_, err = linkToken.Approve(upkeepOwner, registry.Address(), oneHunEth)
		require.NoError(t, err)

		_, err = registry.AddFunds(upkeepOwner, upkeepID, oneHunEth)
		require.NoError(t, err)

		backend.Commit()

		return upkeepID
	}
}

type feedLookupUpkeepController struct {
	// address for dummy protocol
	logSrcAddr common.Address
	// dummy protocol is a log event source
	protocol      *dummy_protocol_wrapper.DummyProtocol
	protocolOwner *bind.TransactOpts
	// log trigger listener contracts react to logs produced from protocol
	count          int
	upkeepIds      []*big.Int
	addresses      []common.Address
	contracts      []*log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookup
	contractsOwner *bind.TransactOpts
}

func newFeedLookupUpkeepController(
	backend *backends.SimulatedBackend,
	protocolOwner *bind.TransactOpts,
) (*feedLookupUpkeepController, error) {
	addr, _, contract, err := dummy_protocol_wrapper.DeployDummyProtocol(protocolOwner, backend)
	if err != nil {
		return nil, err
	}

	backend.Commit()

	return &feedLookupUpkeepController{
		logSrcAddr:    addr,
		protocol:      contract,
		protocolOwner: protocolOwner,
	}, nil
}

func (c *feedLookupUpkeepController) DeployUpkeeps(
	t *testing.T,
	backend *backends.SimulatedBackend,
	owner *bind.TransactOpts,
	count int,
) error {
	addresses := make([]common.Address, count)
	contracts := make([]*log_triggered_streams_lookup_wrapper.LogTriggeredStreamsLookup, count)

	// deploy n upkeep contracts
	for x := 0; x < count; x++ {
		addr, _, contract, err := log_triggered_streams_lookup_wrapper.DeployLogTriggeredStreamsLookup(
			owner,
			backend,
			false,
			false,
		)

		if err != nil {
			require.NoError(t, err, "test dependent on contract deployment")

			return err
		}

		addresses[x] = addr
		contracts[x] = contract
	}

	backend.Commit()

	c.count = count
	c.addresses = addresses
	c.contracts = contracts
	c.contractsOwner = owner

	return nil
}

func (c *feedLookupUpkeepController) RegisterAndFund(
	t *testing.T,
	registry *iregistry21.IKeeperRegistryMaster,
	registryOwner *bind.TransactOpts,
	backend *backends.SimulatedBackend,
	linkToken *link_token_interface.LinkToken,
) error {
	ids := make([]*big.Int, len(c.contracts))

	t.Logf("address: %s", c.logSrcAddr.Hex())

	logTriggerConfigType := abi.MustNewType("tuple(address contractAddress, uint8 filterSelector, bytes32 topic0, bytes32 topic1, bytes32 topic2, bytes32 topic3)")
	config, err := abi.Encode(map[string]interface{}{
		"contractAddress": c.logSrcAddr,
		"filterSelector":  0,                                                                    // no indexed topics filtered
		"topic0":          "0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd", // LimitOrderExecuted event for dummy protocol
		"topic1":          "0x",
		"topic2":          "0x",
		"topic3":          "0x",
	}, logTriggerConfigType)

	require.NoError(t, err)

	registerFunc := registerAndFund(registry, registryOwner, backend, linkToken)

	for x := range c.contracts {
		ids[x] = registerFunc(t, c.addresses[x], c.contractsOwner, 1, config)
	}

	c.upkeepIds = ids

	return nil
}

func (c *feedLookupUpkeepController) EnableMercury(
	t *testing.T,
	backend *backends.SimulatedBackend,
	registry *iregistry21.IKeeperRegistryMaster,
	registryOwner *bind.TransactOpts,
) error {
	adminBytes, _ := json.Marshal(evm21.UpkeepPrivilegeConfig{
		MercuryEnabled: true,
	})

	for _, id := range c.upkeepIds {
		if _, err := registry.SetUpkeepPrivilegeConfig(registryOwner, id, adminBytes); err != nil {
			require.NoError(t, err)

			return err
		}

		callOpts := &bind.CallOpts{
			Pending: true,
			From:    registryOwner.From,
			Context: context.Background(),
		}

		bts, err := registry.GetUpkeepPrivilegeConfig(callOpts, id)
		if err != nil {
			require.NoError(t, err)

			return err
		}

		var checkBytes evm21.UpkeepPrivilegeConfig
		if err := json.Unmarshal(bts, &checkBytes); err != nil {
			require.NoError(t, err)

			return err
		}

		require.True(t, checkBytes.MercuryEnabled)
	}

	bl, _ := backend.BlockByHash(testutils.Context(t), backend.Commit())
	t.Logf("block number after mercury enabled: %d", bl.NumberU64())

	return nil
}

func (c *feedLookupUpkeepController) VerifyEnv(
	t *testing.T,
	backend *backends.SimulatedBackend,
	registry *iregistry21.IKeeperRegistryMaster,
	registryOwner *bind.TransactOpts,
) error {
	t.Log("verifying number of active upkeeps")

	ids, err := registry.GetActiveUpkeepIDs(&bind.CallOpts{
		Context: testutils.Context(t),
		From:    registryOwner.From,
	}, big.NewInt(0), big.NewInt(100))

	require.NoError(t, err)
	require.Len(t, ids, c.count, "active upkeep ids does not match count")
	require.Len(t, ids, len(c.upkeepIds))

	t.Log("verifying total number of contracts")
	require.Len(t, c.contracts, len(c.upkeepIds), "one contract for each upkeep id expected")

	// call individual contracts to see that they revert
	for _, contract := range c.contracts {
		_, err := contract.CheckLog(c.contractsOwner, log_triggered_streams_lookup_wrapper.Log{
			Index:       big.NewInt(0),
			Timestamp:   big.NewInt(123),
			TxHash:      common.HexToHash("0x1"),
			BlockNumber: big.NewInt(0),
			BlockHash:   common.HexToHash("0x14"),
			Source:      common.HexToAddress("0x2"),
			Topics: [][32]byte{
				common.HexToHash("0xd1ffe9e45581c11d7d9f2ed5f75217cd4be9f8b7eee6af0f6d03f46de53956cd"), // matches executedSig and should result in a feedlookup revert
				common.HexToHash("0x"),
				common.HexToHash("0x"),
				common.HexToHash("0x"),
			},
			Data: []byte{},
		}, []byte("0x"))

		require.Error(t, err, "check log contract call should revert: %s", err)
	}

	return nil
}

func (c *feedLookupUpkeepController) EmitEvents(
	t *testing.T,
	backend *backends.SimulatedBackend,
	count int,
	afterEmit func(),
) error {
	ctx := testutils.Context(t)

	for i := 0; i < count && ctx.Err() == nil; i++ {
		_, err := c.protocol.ExecuteLimitOrder(c.protocolOwner, big.NewInt(1000), big.NewInt(10000), c.logSrcAddr)
		require.NoError(t, err, "no error expected from limit order exec")

		if err != nil {
			return err
		}

		backend.Commit()

		// verify event was emitted
		block, _ := backend.BlockByHash(context.Background(), backend.Commit())
		t.Logf("block number after emit event: %d", block.NumberU64())

		iter, _ := c.protocol.FilterLimitOrderExecuted(
			&bind.FilterOpts{
				Context: testutils.Context(t),
				Start:   block.NumberU64() - 1,
			},
			[]*big.Int{big.NewInt(1000)},
			[]*big.Int{big.NewInt(10000)},
			[]common.Address{c.logSrcAddr},
		)

		var eventEmitted bool
		for iter.Next() {
			if iter.Event != nil {
				eventEmitted = true
			}
		}

		require.True(t, eventEmitted, "event expected on backend")
		if !eventEmitted {
			return fmt.Errorf("event was not emitted")
		}

		afterEmit()
	}

	return nil
}

func (c *feedLookupUpkeepController) UpkeepsIds() []*big.Int {
	return c.upkeepIds
}
