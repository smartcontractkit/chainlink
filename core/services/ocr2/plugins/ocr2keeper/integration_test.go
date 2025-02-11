package ocr2keeper_test

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/umbracle/ethgo/abi"

	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrTypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-automation/pkg/v2/config"
	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-common/pkg/types"
	"github.com/smartcontractkit/chainlink-common/pkg/utils/tests"

	"github.com/smartcontractkit/chainlink-integrations/evm/assets"
	evmtypes "github.com/smartcontractkit/chainlink-integrations/evm/types"
	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/logpoller"
	"github.com/smartcontractkit/chainlink/v2/core/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/basic_upkeep_contract"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_logic2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/keeper_registry_wrapper2_0"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/chaintype"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/ocr2keeper/evmregistry/v21/mercury"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

const (
	MercuryCredName = "cred1"
)

var (
	oneEth    = big.NewInt(1000000000000000000)
	oneHunEth = big.NewInt(0).Mul(oneEth, big.NewInt(100))

	payload1 = common.Hex2Bytes("1234")
	payload2 = common.Hex2Bytes("ABCD")
)

func deployKeeper20Registry(
	t *testing.T,
	auth *bind.TransactOpts,
	backend evmtypes.Backend,
	linkAddr, linkFeedAddr,
	gasFeedAddr common.Address,
) *keeper_registry_wrapper2_0.KeeperRegistry {
	logicAddr, _, _, err := keeper_registry_logic2_0.DeployKeeperRegistryLogic(
		auth,
		backend.Client(),
		0, // Payment model
		linkAddr,
		linkFeedAddr,
		gasFeedAddr)
	require.NoError(t, err)
	backend.Commit()

	regAddr, _, _, err := keeper_registry_wrapper2_0.DeployKeeperRegistry(
		auth,
		backend.Client(),
		logicAddr,
	)
	require.NoError(t, err)
	backend.Commit()

	registry, err := keeper_registry_wrapper2_0.NewKeeperRegistry(regAddr, backend.Client())
	require.NoError(t, err)

	return registry
}

func setupNode(
	t *testing.T,
	port int,
	nodeKey ethkey.KeyV2,
	backend evmtypes.Backend,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
	mercury mercury.MercuryEndpointMock,
) (chainlink.Application, string, common.Address, ocr2key.KeyBundle) {
	ctx := testutils.Context(t)
	p2pKey := p2pkey.MustNewV2XXXTestingOnly(big.NewInt(int64(port)))
	p2paddresses := []string{fmt.Sprintf("127.0.0.1:%d", port)}
	cfg, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Feature.LogPoller = ptr(true)

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)
		c.P2P.V2.AnnounceAddresses = &p2paddresses
		c.P2P.V2.ListenAddresses = &p2paddresses
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.EVM[0].Transactions.ForwardersEnabled = ptr(true)
		c.EVM[0].GasEstimator.Mode = ptr("FixedPrice")
		s.Mercury.Credentials = map[string]toml.MercuryCredentials{
			MercuryCredName: {
				LegacyURL: models.MustSecretURL(mercury.URL()),
				URL:       models.MustSecretURL(mercury.URL()),
				Username:  models.NewSecret(mercury.Username()),
				Password:  models.NewSecret(mercury.Password()),
			},
		}
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, cfg, backend, nodeKey, p2pKey)
	kb, err := app.GetKeyStore().OCR2().Create(ctx, chaintype.EVM)
	require.NoError(t, err)

	err = app.Start(ctx)
	require.NoError(t, err)

	t.Cleanup(func() {
		assert.NoError(t, app.Stop())
	})

	return app, p2pKey.PeerID().Raw(), nodeKey.Address, kb
}

type Node struct {
	App         chainlink.Application
	Transmitter common.Address
	KeyBundle   ocr2key.KeyBundle
}

func (node *Node) AddJob(t *testing.T, spec string) {
	c := node.App.GetConfig()
	jb, err := validate.ValidatedOracleSpecToml(testutils.Context(t), c.OCR2(), c.Insecure(), spec, nil)
	require.NoError(t, err)
	err = node.App.AddJobV2(testutils.Context(t), &jb)
	require.NoError(t, err)
}

func (node *Node) AddBootstrapJob(t *testing.T, spec string) {
	jb, err := ocrbootstrap.ValidatedBootstrapSpecToml(spec)
	require.NoError(t, err)
	err = node.App.AddJobV2(testutils.Context(t), &jb)
	require.NoError(t, err)
}

func accountsToAddress(accounts []ocrTypes.Account) (addresses []common.Address, err error) {
	for _, signer := range accounts {
		bytes, err := hexutil.Decode(string(signer))
		if err != nil {
			return []common.Address{}, errors.Wrap(err, fmt.Sprintf("given address is not valid %s", signer))
		}
		if len(bytes) != 20 {
			return []common.Address{}, errors.Errorf("address is not 20 bytes %s", signer)
		}
		addresses = append(addresses, common.BytesToAddress(bytes))
	}
	return addresses, nil
}

func getUpkeepIDFromTx(t *testing.T, registry *keeper_registry_wrapper2_0.KeeperRegistry, registrationTx *gethtypes.Transaction, backend evmtypes.Backend) *big.Int {
	receipt, err := backend.Client().TransactionReceipt(testutils.Context(t), registrationTx.Hash())
	require.NoError(t, err)
	parsedLog, err := registry.ParseUpkeepRegistered(*receipt.Logs[0])
	require.NoError(t, err)
	return parsedLog.Id
}

func TestIntegration_KeeperPluginBasic(t *testing.T) {
	tests.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/AUTO-11072")
	runKeeperPluginBasic(t)
}

func runKeeperPluginBasic(t *testing.T) {
	g := gomega.NewWithT(t)
	lggr := logger.TestLogger(t)

	// setup blockchain
	sergey := testutils.MustNewSimTransactor(t) // owns all the link
	steve := testutils.MustNewSimTransactor(t)  // registry owner
	carrol := testutils.MustNewSimTransactor(t) // upkeep owner
	genesisData := gethtypes.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(1000).ToInt()},
		steve.From:  {Balance: assets.Ether(1000).ToInt()},
		carrol.From: {Balance: assets.Ether(1000).ToInt()},
	}
	// Generate 5 keys for nodes (1 bootstrap + 4 ocr nodes) and fund them with ether
	var nodeKeys [5]ethkey.KeyV2
	for i := int64(0); i < 5; i++ {
		nodeKeys[i] = cltest.MustGenerateRandomKey(t)
		genesisData[nodeKeys[i].Address] = gethtypes.Account{Balance: assets.Ether(1000).ToInt()}
	}

	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	_, stopMining := cltest.Mine(backend, 3*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	defer stopMining()

	// Deploy contracts
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend.Client())
	require.NoError(t, err)
	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend.Client(), 18, big.NewInt(60000000000))
	require.NoError(t, err)
	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend.Client(), 18, big.NewInt(2000000000000000000))
	require.NoError(t, err)
	registry := deployKeeper20Registry(t, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

	// Setup bootstrap + oracle nodes
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, bootstrapTransmitter, bootstrapKb := setupNode(t, bootstrapNodePort, nodeKeys[0], backend, nil, mercury.NewSimulatedMercuryServer())
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
		app, peerID, transmitter, kb := setupNode(t, ports[i], nodeKeys[i+1], backend, []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, mercury.NewSimulatedMercuryServer())

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
		`, i, registry.Address(), node.KeyBundle.ID(), node.Transmitter, fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort), MercuryCredName))
	}

	// Setup config on contract
	configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address registrar)")
	onchainConfig, err := abi.Encode(map[string]interface{}{
		"paymentPremiumPPB":    uint32(0),
		"flatFeeMicroLink":     uint32(0),
		"checkGasLimit":        uint32(6500000),
		"stalenessSeconds":     uint32(90000),
		"gasCeilingMultiplier": uint16(2),
		"minUpkeepSpend":       uint32(0),
		"maxPerformGas":        uint32(5000000),
		"maxCheckDataSize":     uint32(5000),
		"maxPerformDataSize":   uint32(5000),
		"fallbackGasPrice":     big.NewInt(60000000000),
		"fallbackLinkPrice":    big.NewInt(2000000000000000000),
		"transcoder":           testutils.NewAddress(),
		"registrar":            testutils.NewAddress(),
	}, configType)
	require.NoError(t, err)

	offC, err := json.Marshal(config.OffchainConfig{
		PerformLockoutWindow: 100 * 3 * 1000, // ~100 block lockout (on goerli)
		MinConfirmations:     1,
	})
	if err != nil {
		t.Logf("error creating off-chain config: %s", err)
		t.Fail()
	}

	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		10*time.Second,        // deltaProgress time.Duration,
		10*time.Second,        // deltaResend time.Duration,
		2500*time.Millisecond, // deltaRound time.Duration,
		40*time.Millisecond,   // deltaGrace time.Duration,
		15*time.Second,        // deltaStage time.Duration,
		3,                     // rMax uint8,
		[]int{1, 1, 1, 1},
		oracles,
		offC, // reportingPluginConfig []byte,
		nil,
		20*time.Millisecond,   // Max duration query
		1600*time.Millisecond, // Max duration observation
		800*time.Millisecond,
		20*time.Millisecond,
		20*time.Millisecond,
		1, // f
		onchainConfig,
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
		steve,
		signerAddresses,
		transmitterAddresses,
		threshold,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	backend.Commit()

	// Register new upkeep
	upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend.Client())
	require.NoError(t, err)
	registrationTx, err := registry.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{}, []byte{})
	require.NoError(t, err)
	backend.Commit()
	upkeepID := getUpkeepIDFromTx(t, registry, registrationTx, backend)

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

	// change payload
	_, err = upkeepContract.SetBytesToSend(carrol, payload2)
	require.NoError(t, err)
	_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
	require.NoError(t, err)

	// observe 2nd job run and received payload changes
	g.Eventually(receivedBytes, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(payload2))
}

func setupForwarderForNode(
	t *testing.T,
	app chainlink.Application,
	caller *bind.TransactOpts,
	backend evmtypes.Backend,
	recipient common.Address,
	linkAddr common.Address) common.Address {
	ctx := testutils.Context(t)
	faddr, _, authorizedForwarder, err := authorized_forwarder.DeployAuthorizedForwarder(caller, backend.Client(), linkAddr, caller.From, recipient, []byte{})
	require.NoError(t, err)
	backend.Commit()

	// set EOA as an authorized sender for the forwarder
	_, err = authorizedForwarder.SetAuthorizedSenders(caller, []common.Address{recipient})
	require.NoError(t, err)
	backend.Commit()

	// add forwarder address to be tracked in db
	forwarderORM := forwarders.NewORM(app.GetDB())
	chainID, err := backend.Client().ChainID(testutils.Context(t))
	require.NoError(t, err)
	_, err = forwarderORM.CreateForwarder(ctx, faddr, ubig.Big(*chainID))
	require.NoError(t, err)

	chain, err := app.GetRelayers().LegacyEVMChains().Get(chainID.String())
	require.NoError(t, err)
	fwdr, err := chain.TxManager().GetForwarderForEOA(ctx, recipient)
	require.NoError(t, err)
	require.Equal(t, faddr, fwdr)

	return faddr
}

func TestIntegration_KeeperPluginForwarderEnabled(t *testing.T) {
	g := gomega.NewWithT(t)
	lggr := logger.TestLogger(t)

	// setup blockchain
	sergey := testutils.MustNewSimTransactor(t) // owns all the link
	steve := testutils.MustNewSimTransactor(t)  // registry owner
	carrol := testutils.MustNewSimTransactor(t) // upkeep owner
	genesisData := gethtypes.GenesisAlloc{
		sergey.From: {Balance: assets.Ether(1000).ToInt()},
		steve.From:  {Balance: assets.Ether(1000).ToInt()},
		carrol.From: {Balance: assets.Ether(1000).ToInt()},
	}
	// Generate 5 keys for nodes (1 bootstrap + 4 ocr nodes) and fund them with ether
	var nodeKeys [5]ethkey.KeyV2
	for i := int64(0); i < 5; i++ {
		nodeKeys[i] = cltest.MustGenerateRandomKey(t)
		genesisData[nodeKeys[i].Address] = gethtypes.Account{Balance: assets.Ether(1000).ToInt()}
	}

	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	_, stopMining := cltest.Mine(backend, 6*time.Second) // Should be greater than deltaRound since we cannot access old blocks on simulated blockchain
	defer stopMining()

	// Deploy contracts
	linkAddr, _, linkToken, err := link_token_interface.DeployLinkToken(sergey, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	gasFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend.Client(), 18, big.NewInt(60000000000))
	require.NoError(t, err)
	backend.Commit()
	linkFeedAddr, _, _, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(steve, backend.Client(), 18, big.NewInt(2000000000000000000))
	require.NoError(t, err)
	backend.Commit()
	registry := deployKeeper20Registry(t, steve, backend, linkAddr, linkFeedAddr, gasFeedAddr)

	effectiveTransmitters := make([]common.Address, 0)
	// Setup bootstrap + oracle nodes
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, bootstrapTransmitter, bootstrapKb := setupNode(t, bootstrapNodePort, nodeKeys[0], backend, nil, mercury.NewSimulatedMercuryServer())

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
		app, peerID, transmitter, kb := setupNode(t, ports[i], nodeKeys[i+1], backend, []commontypes.BootstrapperLocator{
			// Supply the bootstrap IP and port as a V2 peer address
			{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
		}, mercury.NewSimulatedMercuryServer())
		nodeForwarder := setupForwarderForNode(t, app, sergey, backend, transmitter, linkAddr)
		effectiveTransmitters = append(effectiveTransmitters, nodeForwarder)

		nodes = append(nodes, Node{
			app, transmitter, kb,
		})
		offchainPublicKey, _ := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocrTypes.Account(nodeForwarder.String()),
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
		forwardingAllowed = true

		[relayConfig]
		chainID = 1337

		[pluginConfig]
		cacheEvictionInterval = "1s"
		maxServiceWorkers = 100
		mercuryCredentialName = "%s"
		`, i, registry.Address(), node.KeyBundle.ID(), node.Transmitter, fmt.Sprintf("%s@127.0.0.1:%d", bootstrapPeerID, bootstrapNodePort), MercuryCredName))
	}

	// Setup config on contract
	configType := abi.MustNewType("tuple(uint32 paymentPremiumPPB,uint32 flatFeeMicroLink,uint32 checkGasLimit,uint24 stalenessSeconds,uint16 gasCeilingMultiplier,uint96 minUpkeepSpend,uint32 maxPerformGas,uint32 maxCheckDataSize,uint32 maxPerformDataSize,uint256 fallbackGasPrice,uint256 fallbackLinkPrice,address transcoder,address registrar)")
	onchainConfig, err := abi.Encode(map[string]interface{}{
		"paymentPremiumPPB":    uint32(0),
		"flatFeeMicroLink":     uint32(0),
		"checkGasLimit":        uint32(6500000),
		"stalenessSeconds":     uint32(90000),
		"gasCeilingMultiplier": uint16(2),
		"minUpkeepSpend":       uint32(0),
		"maxPerformGas":        uint32(5000000),
		"maxCheckDataSize":     uint32(5000),
		"maxPerformDataSize":   uint32(5000),
		"fallbackGasPrice":     big.NewInt(60000000000),
		"fallbackLinkPrice":    big.NewInt(2000000000000000000),
		"transcoder":           testutils.NewAddress(),
		"registrar":            testutils.NewAddress(),
	}, configType)
	require.NoError(t, err)

	offC, err := json.Marshal(config.OffchainConfig{
		PerformLockoutWindow: 100 * 12 * 1000, // ~100 block lockout (on goerli)
	})
	if err != nil {
		t.Logf("error creating off-chain config: %s", err)
		t.FailNow()
	}

	signers, transmitters, threshold, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper.ContractSetConfigArgsForTests(
		10*time.Second,       // deltaProgress time.Duration,
		10*time.Second,       // deltaResend time.Duration,
		5*time.Second,        // deltaRound time.Duration,
		500*time.Millisecond, // deltaGrace time.Duration,
		2*time.Second,        // deltaStage time.Duration,
		3,                    // rMax uint8,
		[]int{1, 1, 1, 1},
		oracles,
		offC, // reportingPluginConfig []byte,
		nil,
		50*time.Millisecond, // Max duration query
		1*time.Second,       // Max duration observation
		1*time.Second,
		1*time.Second,
		1*time.Second,
		1, // f
		onchainConfig,
	)
	require.NoError(t, err)
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)
	transmitterAddresses, err := accountsToAddress(transmitters)
	require.NoError(t, err)

	// Make sure we are using forwarders and not node keys as transmitters on chain.
	eoaList := make([]common.Address, 0)
	for _, n := range nodes {
		eoaList = append(eoaList, n.Transmitter)
	}
	require.Equal(t, effectiveTransmitters, transmitterAddresses)
	require.NotEqual(t, eoaList, effectiveTransmitters)
	lggr.Infow("Setting Config on Oracle Contract",
		"signerAddresses", signerAddresses,
		"transmitterAddresses", transmitterAddresses,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", offchainConfigVersion,
		"offchainConfig", offchainConfig,
	)
	_, err = registry.SetConfig(
		steve,
		signerAddresses,
		transmitterAddresses,
		threshold,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig,
	)
	require.NoError(t, err)
	backend.Commit()

	// Register new upkeep
	upkeepAddr, _, upkeepContract, err := basic_upkeep_contract.DeployBasicUpkeepContract(carrol, backend.Client())
	require.NoError(t, err)
	backend.Commit()
	registrationTx, err := registry.RegisterUpkeep(steve, upkeepAddr, 2_500_000, carrol.From, []byte{}, []byte{})
	require.NoError(t, err)
	backend.Commit()
	upkeepID := getUpkeepIDFromTx(t, registry, registrationTx, backend)

	// Fund the upkeep
	_, err = linkToken.Transfer(sergey, carrol.From, oneHunEth)
	require.NoError(t, err)
	backend.Commit()
	_, err = linkToken.Approve(carrol, registry.Address(), oneHunEth)
	require.NoError(t, err)
	backend.Commit()
	_, err = registry.AddFunds(carrol, upkeepID, oneHunEth)
	require.NoError(t, err)
	backend.Commit()

	// Set upkeep to be performed
	_, err = upkeepContract.SetBytesToSend(carrol, payload1)
	require.NoError(t, err)
	backend.Commit()
	_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
	require.NoError(t, err)
	backend.Commit()

	lggr.Infow("Upkeep registered and funded")

	// keeper job is triggered and payload is received
	receivedBytes := func() []byte {
		received, err2 := upkeepContract.ReceivedBytes(nil)
		require.NoError(t, err2)
		return received
	}
	g.Eventually(receivedBytes, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(payload1))

	// change payload
	_, err = upkeepContract.SetBytesToSend(carrol, payload2)
	require.NoError(t, err)
	backend.Commit()
	_, err = upkeepContract.SetShouldPerformUpkeep(carrol, true)
	require.NoError(t, err)
	backend.Commit()

	// observe 2nd job run and received payload changes
	g.Eventually(receivedBytes, testutils.WaitTimeout(t), cltest.DBPollingInterval).Should(gomega.Equal(payload2))
}

func ptr[T any](v T) *T { return &v }

func TestFilterNamesFromSpec20(t *testing.T) {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	require.NoError(t, err)
	address := common.HexToAddress(hexutil.Encode(b))

	spec := &job.OCR2OracleSpec{
		PluginType: types.OCR2Keeper,
		ContractID: address.String(), // valid contract addr
	}

	names, err := ocr2keeper.FilterNamesFromSpec20(spec)
	require.NoError(t, err)

	assert.Len(t, names, 2)
	assert.Equal(t, logpoller.FilterName("OCR2KeeperRegistry - LogProvider", address), names[0])
	assert.Equal(t, logpoller.FilterName("EvmRegistry - Upkeep events for", address), names[1])

	spec = &job.OCR2OracleSpec{
		PluginType: types.OCR2Keeper,
		ContractID: "0x5431", // invalid contract addr
	}
	_, err = ocr2keeper.FilterNamesFromSpec20(spec)
	require.ErrorContains(t, err, "not a valid EIP55 formatted address")
}
