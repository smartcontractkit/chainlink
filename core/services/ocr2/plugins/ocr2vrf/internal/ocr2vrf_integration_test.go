package internal_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/libocr/commontypes"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	ocr2dkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/mock_v3_aggregator_contract"
	dkg_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/load_test_beacon_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	vrf_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
)

type ocr2vrfUniverse struct {
	owner   *bind.TransactOpts
	backend *backends.SimulatedBackend

	dkgAddress common.Address
	dkg        *dkg_wrapper.DKG

	beaconAddress      common.Address
	coordinatorAddress common.Address
	beacon             *vrf_beacon.VRFBeacon
	coordinator        *vrf_wrapper.VRFCoordinator

	linkAddress common.Address
	link        *link_token_interface.LinkToken

	consumerAddress common.Address
	consumer        *vrf_beacon_consumer.BeaconVRFConsumer

	loadTestConsumerAddress common.Address
	loadTestConsumer        *load_test_beacon_consumer.LoadTestBeaconVRFConsumer

	feedAddress common.Address
	feed        *mock_v3_aggregator_contract.MockV3AggregatorContract
}

type ocr2Node struct {
	app         *cltest.TestApplication
	peerID      string
	transmitter common.Address
	keybundle   ocr2key.KeyBundle
	config      *configtest.TestGeneralConfig
}

func setupOCR2VRFContracts(
	t *testing.T, beaconPeriod int64, keyID [32]byte, consumerShouldFail bool) ocr2vrfUniverse {
	owner := testutils.MustNewSimTransactor(t)
	genesisData := core.GenesisAlloc{
		owner.From: {
			Balance: assets.Ether(100),
		},
	}
	b := backends.NewSimulatedBackend(genesisData, ethconfig.Defaults.Miner.GasCeil*2)

	// deploy OCR2VRF contracts, which have the following deploy order:
	// * link token
	// * link/eth feed
	// * DKG
	// * VRF
	// * VRF consumer
	linkAddress, _, link, err := link_token_interface.DeployLinkToken(
		owner, b)
	require.NoError(t, err)

	b.Commit()

	feedAddress, _, feed, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
		owner, b, 18, assets.GWei(1e7)) // 0.01 eth per link
	require.NoError(t, err)

	b.Commit()

	dkgAddress, _, dkg, err := dkg_wrapper.DeployDKG(owner, b)
	require.NoError(t, err)

	b.Commit()

	coordinatorAddress, _, coordinator, err := vrf_wrapper.DeployVRFCoordinator(
		owner, b, big.NewInt(beaconPeriod), linkAddress)
	require.NoError(t, err)

	b.Commit()

	beaconAddress, _, beacon, err := vrf_beacon.DeployVRFBeacon(
		owner, b, linkAddress, coordinatorAddress, dkgAddress, keyID)
	require.NoError(t, err)

	b.Commit()

	consumerAddress, _, consumer, err := vrf_beacon_consumer.DeployBeaconVRFConsumer(
		owner, b, coordinatorAddress, consumerShouldFail, big.NewInt(beaconPeriod))
	require.NoError(t, err)

	b.Commit()

	loadTestConsumerAddress, _, loadTestConsumer, err := load_test_beacon_consumer.DeployLoadTestBeaconVRFConsumer(
		owner, b, coordinatorAddress, consumerShouldFail, big.NewInt(beaconPeriod))
	require.NoError(t, err)

	b.Commit()

	_, err = dkg.AddClient(owner, keyID, beaconAddress)
	require.NoError(t, err)

	b.Commit()

	_, err = coordinator.SetProducer(owner, beaconAddress)
	require.NoError(t, err)

	// Achieve finality depth so the CL node can work properly.
	for i := 0; i < 20; i++ {
		b.Commit()
	}

	return ocr2vrfUniverse{
		owner:                   owner,
		backend:                 b,
		dkgAddress:              dkgAddress,
		dkg:                     dkg,
		beaconAddress:           beaconAddress,
		coordinatorAddress:      coordinatorAddress,
		beacon:                  beacon,
		coordinator:             coordinator,
		linkAddress:             linkAddress,
		link:                    link,
		consumerAddress:         consumerAddress,
		consumer:                consumer,
		loadTestConsumerAddress: loadTestConsumerAddress,
		loadTestConsumer:        loadTestConsumer,
		feedAddress:             feedAddress,
		feed:                    feed,
	}
}

func setupNodeOCR2(
	t *testing.T,
	owner *bind.TransactOpts,
	port uint16,
	dbName string,
	b *backends.SimulatedBackend,
) *ocr2Node {
	config, _ := heavyweight.FullTestDB(t, fmt.Sprintf("%s%d", dbName, port))
	config.Overrides.FeatureOffchainReporting = null.BoolFrom(false)
	config.Overrides.FeatureOffchainReporting2 = null.BoolFrom(true)
	config.Overrides.FeatureLogPoller = null.BoolFrom(true)
	poll := 500 * time.Millisecond
	config.Overrides.GlobalEvmLogPollInterval = &poll
	config.Overrides.P2PEnabled = null.BoolFrom(true)
	config.Overrides.P2PNetworkingStack = ocrnetworking.NetworkingStackV2
	config.Overrides.P2PListenPort = null.NewInt(0, true)
	config.Overrides.GlobalEvmGasLimitDefault = null.NewInt(3_500_000, true)
	config.Overrides.SetP2PV2DeltaDial(500 * time.Millisecond)
	config.Overrides.SetP2PV2DeltaReconcile(5 * time.Second)
	p2paddresses := []string{
		fmt.Sprintf("127.0.0.1:%d", port),
	}
	config.Overrides.P2PV2ListenAddresses = p2paddresses
	// Disables ocr spec validation so we can have fast polling for the test.
	config.Overrides.Dev = null.BoolFrom(true)

	app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, b)
	_, err := app.GetKeyStore().P2P().Create()
	require.NoError(t, err)

	p2pIDs, err := app.GetKeyStore().P2P().GetAll()
	require.NoError(t, err)
	require.Len(t, p2pIDs, 1)
	peerID := p2pIDs[0].PeerID()

	config.Overrides.P2PPeerID = peerID

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address

	// Fund the transmitter address with some ETH
	n, err := b.NonceAt(testutils.Context(t), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(
		n, transmitter,
		assets.Ether(1),
		21000,
		assets.GWei(1),
		nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(testutils.Context(t), signedTx)
	require.NoError(t, err)
	b.Commit()

	kb, err := app.GetKeyStore().OCR2().Create("evm")
	require.NoError(t, err)

	return &ocr2Node{
		app:         app,
		peerID:      peerID.Raw(),
		transmitter: transmitter,
		keybundle:   kb,
		config:      config,
	}
}

func TestIntegration_OCR2VRF(t *testing.T) {
	keyID := randomKeyID(t)
	uni := setupOCR2VRFContracts(t, 5, keyID, false)

	t.Log("Creating bootstrap node")

	bootstrapNodePort := getFreePort(t)
	bootstrapNode := setupNodeOCR2(t, uni.owner, bootstrapNodePort, "bootstrap", uni.backend)
	numNodes := 5

	t.Log("Creating OCR2 nodes")
	var (
		oracles        []confighelper2.OracleIdentityExtra
		transmitters   []common.Address
		onchainPubKeys []common.Address
		kbs            []ocr2key.KeyBundle
		apps           []*cltest.TestApplication
		dkgEncrypters  []dkgencryptkey.Key
		dkgSigners     []dkgsignkey.Key
	)
	for i := 0; i < numNodes; i++ {
		node := setupNodeOCR2(t, uni.owner, getFreePort(t), fmt.Sprintf("ocr2vrforacle%d", i), uni.backend)
		// Supply the bootstrap IP and port as a V2 peer address
		node.config.Overrides.P2PV2Bootstrappers = []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.peerID, Addrs: []string{
				fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort),
			}},
		}

		dkgSignKey, err := node.app.GetKeyStore().DKGSign().Create()
		require.NoError(t, err)

		dkgEncryptKey, err := node.app.GetKeyStore().DKGEncrypt().Create()
		require.NoError(t, err)

		kbs = append(kbs, node.keybundle)
		apps = append(apps, node.app)
		transmitters = append(transmitters, node.transmitter)
		dkgEncrypters = append(dkgEncrypters, dkgEncryptKey)
		dkgSigners = append(dkgSigners, dkgSignKey)
		onchainPubKeys = append(onchainPubKeys, common.BytesToAddress(node.keybundle.PublicKey()))
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  node.keybundle.PublicKey(),
				TransmitAccount:   ocrtypes2.Account(node.transmitter.String()),
				OffchainPublicKey: node.keybundle.OffchainPublicKey(),
				PeerID:            node.peerID,
			},
			ConfigEncryptionPublicKey: node.keybundle.ConfigEncryptionPublicKey(),
		})
	}

	t.Log("starting ticker to commit blocks")
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go func() {
		for range tick.C {
			uni.backend.Commit()
		}
	}()

	blockBeforeConfig, err := uni.backend.BlockByNumber(context.Background(), nil)
	require.NoError(t, err)

	t.Log("Setting DKG config before block:", blockBeforeConfig.Number().String())

	// set config for dkg
	setDKGConfig(
		t,
		uni,
		onchainPubKeys,
		transmitters,
		1,
		oracles,
		dkgSigners,
		dkgEncrypters,
		keyID,
	)

	t.Log("Adding bootstrap node job")
	err = bootstrapNode.app.Start(testutils.Context(t))
	require.NoError(t, err)

	chainSet := bootstrapNode.app.GetChains().EVM
	require.NotNil(t, chainSet)
	bootstrapJobSpec := fmt.Sprintf(`
type				= "bootstrap"
name				= "bootstrap"
relay				= "evm"
schemaVersion		= 1
contractID			= "%s"
[relayConfig]
chainID 			= 1337
`, uni.dkgAddress.Hex())
	t.Log("Creating bootstrap job:", bootstrapJobSpec)
	ocrJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(bootstrapJobSpec)
	require.NoError(t, err)
	err = bootstrapNode.app.AddJobV2(context.Background(), &ocrJob)
	require.NoError(t, err)

	t.Log("Creating OCR2VRF jobs")
	for i := 0; i < numNodes; i++ {
		err = apps[i].Start(testutils.Context(t))
		require.NoError(t, err)

		jobSpec := fmt.Sprintf(`
type                 	= "offchainreporting2"
schemaVersion        	= 1
name                 	= "ocr2 vrf integration test"
maxTaskDuration      	= "30s"
contractID           	= "%s"
ocrKeyBundleID       	= "%s"
relay                	= "evm"
pluginType           	= "ocr2vrf"
transmitterID        	= "%s"

[relayConfig]
chainID              	= 1337

[pluginConfig]
dkgEncryptionPublicKey 	= "%s"
dkgSigningPublicKey    	= "%s"
dkgKeyID               	= "%s"
dkgContractAddress     	= "%s"

vrfCoordinatorAddress   = "%s"
linkEthFeedAddress     	= "%s"
confirmationDelays     	= %s # This is an array
lookbackBlocks         	= %d # This is an integer
`, uni.beaconAddress.String(),
			kbs[i].ID(),
			transmitters[i],
			dkgEncrypters[i].PublicKeyString(),
			dkgSigners[i].PublicKeyString(),
			hex.EncodeToString(keyID[:]),
			uni.dkgAddress.String(),
			uni.coordinatorAddress.String(),
			uni.feedAddress.String(),
			"[1, 2, 3, 4, 5, 6, 7, 8]", // conf delays
			1000,                       // lookback blocks
		)
		t.Log("Creating OCR2VRF job with spec:", jobSpec)
		ocrJob, err := validate.ValidatedOracleSpecToml(apps[i].Config, jobSpec)
		require.NoError(t, err)
		err = apps[i].AddJobV2(context.Background(), &ocrJob)
		require.NoError(t, err)
	}

	t.Log("jobs added, running log poller replay")

	// Once all the jobs are added, replay to ensure we have the configSet logs.
	for _, app := range apps {
		require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Replay(context.Background(), blockBeforeConfig.Number().Int64()))
	}
	require.NoError(t, bootstrapNode.app.Chains.EVM.Chains()[0].LogPoller().Replay(context.Background(), blockBeforeConfig.Number().Int64()))

	t.Log("Waiting for DKG key to get written")
	// poll until a DKG key is written to the contract
	// at that point we can start sending VRF requests
	var emptyKH [32]byte
	emptyHash := crypto.Keccak256Hash(emptyKH[:])
	gomega.NewWithT(t).Eventually(func() bool {
		kh, err := uni.beacon.SProvingKeyHash(&bind.CallOpts{
			Context: testutils.Context(t),
		})
		require.NoError(t, err)
		t.Log("proving keyhash:", hexutil.Encode(kh[:]))
		return crypto.Keccak256Hash(kh[:]) != emptyHash
	}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())

	t.Log("DKG key written, setting VRF config")

	// set config for vrf now that dkg is ready
	setVRFConfig(
		t,
		uni,
		onchainPubKeys,
		transmitters,
		1,
		oracles,
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		keyID)

	t.Log("Sending VRF request")

	// Send a VRF request and mine it
	_, err = uni.consumer.TestRequestRandomness(uni.owner, 2, 1, big.NewInt(1))
	require.NoError(t, err)
	_, err = uni.consumer.TestRequestRandomnessFulfillment(uni.owner, 1, 1, big.NewInt(2), 50_000, []byte{})
	require.NoError(t, err)

	_, err = uni.loadTestConsumer.TestRequestRandomnessFulfillmentBatch(uni.owner, 1, 1, big.NewInt(2), 200_000, []byte{}, big.NewInt(2))
	require.NoError(t, err)

	uni.backend.Commit()

	t.Log("waiting for fulfillment")

	// poll until we're able to redeem the randomness without reverting
	// at that point, it's been fulfilled
	gomega.NewWithT(t).Eventually(func() bool {
		_, err1 := uni.consumer.TestRedeemRandomness(uni.owner, big.NewInt(0))
		t.Logf("TestRedeemRandomness err: %+v", err1)
		return err1 == nil
	}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())

	// Mine block after redeeming randomness
	uni.backend.Commit()

	// poll until we're able to verify that consumer contract has stored randomness as expected
	// First arg is the request ID, which starts at zero, second is the index into
	// the random words.
	gomega.NewWithT(t).Eventually(func() bool {

		var errs []error
		rw1, err := uni.consumer.SReceivedRandomnessByRequestID(nil, big.NewInt(0), big.NewInt(0))
		t.Logf("TestRedeemRandomness 1st word err: %+v", err)
		errs = append(errs, err)
		rw2, err := uni.consumer.SReceivedRandomnessByRequestID(nil, big.NewInt(0), big.NewInt(1))
		t.Logf("TestRedeemRandomness 2nd word err: %+v", err)
		errs = append(errs, err)
		rw3, err := uni.consumer.SReceivedRandomnessByRequestID(nil, big.NewInt(1), big.NewInt(0))
		t.Logf("FulfillRandomness 1st word err: %+v", err)
		errs = append(errs, err)
		rw4, err := uni.loadTestConsumer.SReceivedRandomnessByRequestID(nil, big.NewInt(2), big.NewInt(0))
		t.Logf("Batch FulfillRandomness 1st word err: %+v", err)
		errs = append(errs, err)
		rw5, err := uni.loadTestConsumer.SReceivedRandomnessByRequestID(nil, big.NewInt(3), big.NewInt(0))
		t.Logf("Batch FulfillRandomness 2nd word err: %+v", err)
		errs = append(errs, err)
		batchTotalRequests, err := uni.loadTestConsumer.STotalRequests(nil)
		t.Logf("Batch FulfillRandomness total requests err: %+v", err)
		errs = append(errs, err)
		batchTotalFulfillments, err := uni.loadTestConsumer.STotalFulfilled(nil)
		t.Logf("Batch FulfillRandomness total fulfillments err: %+v", err)
		errs = append(errs, err)
		err = nil
		if batchTotalRequests.Int64() != batchTotalFulfillments.Int64() {
			err = errors.New("batchTotalRequests is not equal to batchTotalFulfillments")
			errs = append(errs, err)
		}
		t.Logf("Batch FulfillRandomness total requests/fulfillments equal err: %+v", err)

		t.Logf("randomness from redeemRandomness: %s %s", rw1.String(), rw2.String())
		t.Logf("randomness from fulfillRandomness: %s", rw3.String())
		t.Logf("randomness from batch fulfillRandomness: %s %s", rw4.String(), rw5.String())
		t.Logf("total batch requested and fulfilled: %d %d", batchTotalRequests, batchTotalFulfillments)

		for _, err := range errs {
			if err != nil {
				return false
			}
		}
		return true
	}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())
}

func setDKGConfig(
	t *testing.T,
	uni ocr2vrfUniverse,
	onchainPubKeys []common.Address,
	transmitters []common.Address,
	f uint8,
	oracleIdentities []confighelper2.OracleIdentityExtra,
	signKeys []dkgsignkey.Key,
	encryptKeys []dkgencryptkey.Key,
	keyID [32]byte,
) {
	var (
		signingPubKeys []kyber.Point
		encryptPubKeys []kyber.Point
	)
	for i := range signKeys {
		signingPubKeys = append(signingPubKeys, signKeys[i].PublicKey)
		encryptPubKeys = append(encryptPubKeys, encryptKeys[i].PublicKey)
	}

	offchainConfig, err := ocr2dkg.OffchainConfig(
		encryptPubKeys,
		signingPubKeys,
		&altbn_128.G1{},
		&ocr2vrftypes.PairingTranslation{
			Suite: &altbn_128.PairingSuite{},
		})
	require.NoError(t, err)
	onchainConfig, err := ocr2dkg.OnchainConfig(keyID)
	require.NoError(t, err)

	var schedule []int
	for range oracleIdentities {
		schedule = append(schedule, 1)
	}

	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		30*time.Second,
		10*time.Second,
		10*time.Second,
		20*time.Second,
		20*time.Second,
		3,
		schedule,
		oracleIdentities,
		offchainConfig,
		50*time.Millisecond,
		10*time.Second,
		10*time.Second,
		100*time.Millisecond,
		1*time.Second,
		int(f),
		onchainConfig)
	require.NoError(t, err)

	_, err = uni.dkg.SetConfig(uni.owner, onchainPubKeys, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err)

	uni.backend.Commit()
}

func setVRFConfig(
	t *testing.T,
	uni ocr2vrfUniverse,
	onchainPubKeys []common.Address,
	transmitters []common.Address,
	f uint8,
	oracleIdentities []confighelper2.OracleIdentityExtra,
	confDelaysSl []int,
	keyID [32]byte,
) {
	offchainConfig := ocr2vrf.OffchainConfig()

	confDelays := make(map[uint32]struct{})
	for _, c := range confDelaysSl {
		confDelays[uint32(c)] = struct{}{}
	}

	onchainConfig := ocr2vrf.OnchainConfig(confDelays)

	var schedule []int
	for range oracleIdentities {
		schedule = append(schedule, 1)
	}

	_, _, f, onchainConfig, offchainConfigVersion, offchainConfig, err := confighelper2.ContractSetConfigArgsForTests(
		30*time.Second,
		10*time.Second,
		10*time.Second,
		20*time.Second,
		20*time.Second,
		3,
		schedule,
		oracleIdentities,
		offchainConfig,
		50*time.Millisecond,
		10*time.Second,
		10*time.Second,
		100*time.Millisecond,
		1*time.Second,
		int(f),
		onchainConfig)
	require.NoError(t, err)

	_, err = uni.beacon.SetConfig(
		uni.owner, onchainPubKeys, transmitters, f, onchainConfig, offchainConfigVersion, offchainConfig)
	require.NoError(t, err)

	uni.backend.Commit()
}

func randomKeyID(t *testing.T) (r [32]byte) {
	_, err := rand.Read(r[:])
	require.NoError(t, err)
	return
}

func getFreePort(t *testing.T) uint16 {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	require.NoError(t, err)

	l, err := net.ListenTCP("tcp", addr)
	require.NoError(t, err)
	defer l.Close()

	return uint16(l.Addr().(*net.TCPAddr).Port)
}
