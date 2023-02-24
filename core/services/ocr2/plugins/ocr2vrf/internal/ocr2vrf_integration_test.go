package internal_test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"net"
	"os"
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
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/ocr2vrf/altbn_128"
	ocr2dkg "github.com/smartcontractkit/ocr2vrf/dkg"
	"github.com/smartcontractkit/ocr2vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/ocr2vrf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/mock_v3_aggregator_contract"
	dkg_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/load_test_beacon_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	vrf_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_router"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/utils"
)

type ocr2vrfUniverse struct {
	owner   *bind.TransactOpts
	backend *backends.SimulatedBackend

	dkgAddress common.Address
	dkg        *dkg_wrapper.DKG

	beaconAddress      common.Address
	coordinatorAddress common.Address
	routerAddress      common.Address
	beacon             *vrf_beacon.VRFBeacon
	coordinator        *vrf_wrapper.VRFCoordinator
	router             *vrf_router.VRFRouter

	linkAddress common.Address
	link        *link_token_interface.LinkToken

	consumerAddress common.Address
	consumer        *vrf_beacon_consumer.BeaconVRFConsumer

	loadTestConsumerAddress common.Address
	loadTestConsumer        *load_test_beacon_consumer.LoadTestBeaconVRFConsumer

	feedAddress common.Address
	feed        *mock_v3_aggregator_contract.MockV3AggregatorContract

	subID *big.Int
}

type ocr2Node struct {
	app                  *cltest.TestApplication
	peerID               string
	transmitter          common.Address
	effectiveTransmitter common.Address
	keybundle            ocr2key.KeyBundle
	config               config.GeneralConfig
}

func setupOCR2VRFContracts(
	t *testing.T, beaconPeriod int64, keyID [32]byte, consumerShouldFail bool) ocr2vrfUniverse {
	owner := testutils.MustNewSimTransactor(t)
	owner.GasPrice = assets.GWei(1).ToInt()
	genesisData := core.GenesisAlloc{
		owner.From: {
			Balance: assets.Ether(100).ToInt(),
		},
	}
	b := backends.NewSimulatedBackend(genesisData, ethconfig.Defaults.Miner.GasCeil*2)

	// deploy OCR2VRF contracts, which have the following deploy order:
	// * link token
	// * link/eth feed
	// * DKG
	// * VRF (router, coordinator, and beacon)
	// * VRF consumer
	linkAddress, _, link, err := link_token_interface.DeployLinkToken(
		owner, b)
	require.NoError(t, err)
	b.Commit()

	feedAddress, _, feed, err := mock_v3_aggregator_contract.DeployMockV3AggregatorContract(
		owner, b, 18, assets.GWei(int(1e7)).ToInt()) // 0.01 eth per link
	require.NoError(t, err)
	b.Commit()

	dkgAddress, _, dkg, err := dkg_wrapper.DeployDKG(owner, b)
	require.NoError(t, err)
	b.Commit()

	routerAddress, _, router, err := vrf_router.DeployVRFRouter(owner, b)
	require.NoError(t, err)
	b.Commit()

	coordinatorAddress, _, coordinator, err := vrf_wrapper.DeployVRFCoordinator(
		owner, b, big.NewInt(beaconPeriod), linkAddress, feedAddress, routerAddress)
	require.NoError(t, err)
	b.Commit()

	require.NoError(t, utils.JustError(coordinator.SetBillingConfig(owner, vrf_wrapper.VRFBeaconTypesBillingConfig{
		RedeemableRequestGasOverhead: 50_000,
		CallbackRequestGasOverhead:   50_000,
		StalenessSeconds:             60,
		PremiumPercentage:            0,
		FallbackWeiPerUnitLink:       assets.GWei(int(1e7)).ToInt(),
	})))
	b.Commit()

	require.NoError(t, utils.JustError(router.RegisterCoordinator(owner, coordinatorAddress)))
	b.Commit()

	beaconAddress, _, beacon, err := vrf_beacon.DeployVRFBeacon(
		owner, b, linkAddress, coordinatorAddress, dkgAddress, keyID)
	require.NoError(t, err)
	b.Commit()

	consumerAddress, _, consumer, err := vrf_beacon_consumer.DeployBeaconVRFConsumer(
		owner, b, routerAddress, consumerShouldFail, big.NewInt(beaconPeriod))
	require.NoError(t, err)
	b.Commit()

	loadTestConsumerAddress, _, loadTestConsumer, err := load_test_beacon_consumer.DeployLoadTestBeaconVRFConsumer(
		owner, b, routerAddress, consumerShouldFail, big.NewInt(beaconPeriod))
	require.NoError(t, err)
	b.Commit()

	// Set up coordinator subscription for billing.
	require.NoError(t, utils.JustError(coordinator.CreateSubscription(owner)))
	b.Commit()

	fopts := &bind.FilterOpts{}

	subscriptionIterator, err := coordinator.FilterSubscriptionCreated(fopts, nil, []common.Address{owner.From})
	require.NoError(t, err)

	require.True(t, subscriptionIterator.Next())
	subID := subscriptionIterator.Event.SubId

	require.NoError(t, utils.JustError(coordinator.AddConsumer(owner, subID, consumerAddress)))
	b.Commit()
	require.NoError(t, utils.JustError(coordinator.AddConsumer(owner, subID, loadTestConsumerAddress)))
	b.Commit()
	data, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	require.NoError(t, err)
	require.NoError(t, utils.JustError(link.TransferAndCall(owner, coordinatorAddress, big.NewInt(5e18), data)))
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
		routerAddress:           routerAddress,
		beacon:                  beacon,
		coordinator:             coordinator,
		router:                  router,
		linkAddress:             linkAddress,
		link:                    link,
		consumerAddress:         consumerAddress,
		consumer:                consumer,
		loadTestConsumerAddress: loadTestConsumerAddress,
		loadTestConsumer:        loadTestConsumer,
		feedAddress:             feedAddress,
		feed:                    feed,
		subID:                   subID,
	}
}

func setupNodeOCR2(
	t *testing.T,
	owner *bind.TransactOpts,
	port uint16,
	dbName string,
	b *backends.SimulatedBackend,
	useForwarders bool,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
) *ocr2Node {
	p2pKey, err := p2pkey.NewV2()
	require.NoError(t, err)
	config, _ := heavyweight.FullTestDBV2(t, fmt.Sprintf("%s%d", dbName, port), func(c *chainlink.Config, s *chainlink.Secrets) {
		c.DevMode = true // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V1.Enabled = ptr(false)
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = models.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = models.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.EVM[0].LogPollInterval = models.MustNewDuration(500 * time.Millisecond)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint32](3_500_000)
		c.EVM[0].Transactions.ForwardersEnabled = &useForwarders
		c.OCR2.ContractPollInterval = models.MustNewDuration(10 * time.Second)
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	sendingKeys, err := app.KeyStore.Eth().EnabledKeysForChain(testutils.SimulatedChainID)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address
	effectiveTransmitter := sendingKeys[0].Address

	if useForwarders {
		sendingKeysAddresses := []common.Address{sendingKeys[0].Address}

		// Add new sending key.
		k, err := app.KeyStore.Eth().Create()
		require.NoError(t, err)
		require.NoError(t, app.KeyStore.Eth().Enable(k.Address, testutils.SimulatedChainID))
		sendingKeys = append(sendingKeys, k)
		sendingKeysAddresses = append(sendingKeysAddresses, k.Address)

		require.Len(t, sendingKeys, 2)

		// Deploy a forwarder.
		faddr, _, authorizedForwarder, err := authorized_forwarder.DeployAuthorizedForwarder(owner, b, common.HexToAddress("0x326C977E6efc84E512bB9C30f76E30c160eD06FB"), owner.From, common.Address{}, []byte{})
		require.NoError(t, err)

		// Set the node's sending keys as authorized senders.
		_, err = authorizedForwarder.SetAuthorizedSenders(owner, sendingKeysAddresses)
		require.NoError(t, err)
		b.Commit()

		// Add the forwarder to the node's forwarder manager.
		forwarderORM := forwarders.NewORM(app.GetSqlxDB(), logger.TestLogger(t), config)
		chainID := utils.Big(*b.Blockchain().Config().ChainID)
		_, err = forwarderORM.CreateForwarder(faddr, chainID)
		require.NoError(t, err)
		effectiveTransmitter = faddr
	}

	// Fund the sending keys with some ETH.
	for _, k := range sendingKeys {
		n, err := b.NonceAt(testutils.Context(t), owner.From, nil)
		require.NoError(t, err)

		tx := types.NewTransaction(
			n, k.Address,
			assets.Ether(1).ToInt(),
			21000,
			assets.GWei(1).ToInt(),
			nil)
		signedTx, err := owner.Signer(owner.From, tx)
		require.NoError(t, err)
		err = b.SendTransaction(testutils.Context(t), signedTx)
		require.NoError(t, err)
		b.Commit()
	}

	kb, err := app.GetKeyStore().OCR2().Create("evm")
	require.NoError(t, err)

	return &ocr2Node{
		app:                  app,
		peerID:               p2pKey.PeerID().Raw(),
		transmitter:          transmitter,
		effectiveTransmitter: effectiveTransmitter,
		keybundle:            kb,
		config:               config,
	}
}

func TestIntegration_OCR2VRF_ForwarderFlow(t *testing.T) {
	if os.Getenv("CI") == "" && os.Getenv("VRF_LOCAL_TESTING") == "" {
		t.Skip("Skipping test locally.")
	}
	runOCR2VRFTest(t, true)
}

func TestIntegration_OCR2VRF(t *testing.T) {
	if os.Getenv("CI") == "" && os.Getenv("VRF_LOCAL_TESTING") == "" {
		t.Skip("Skipping test locally.")
	}
	runOCR2VRFTest(t, false)
}

func runOCR2VRFTest(t *testing.T, useForwarders bool) {
	keyID := randomKeyID(t)
	uni := setupOCR2VRFContracts(t, 5, keyID, false)

	t.Log("Creating bootstrap node")

	bootstrapNodePort := getFreePort(t)
	bootstrapNode := setupNodeOCR2(t, uni.owner, bootstrapNodePort, "bootstrap", uni.backend, false, nil)
	numNodes := 5

	t.Log("Creating OCR2 nodes")
	var (
		oracles               []confighelper2.OracleIdentityExtra
		transmitters          []common.Address
		effectiveTransmitters []common.Address
		onchainPubKeys        []common.Address
		kbs                   []ocr2key.KeyBundle
		apps                  []*cltest.TestApplication
		dkgEncrypters         []dkgencryptkey.Key
		dkgSigners            []dkgsignkey.Key
	)
	for i := 0; i < numNodes; i++ {
		// Supply the bootstrap IP and port as a V2 peer address
		bootstrappers := []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.peerID, Addrs: []string{
				fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort),
			}},
		}
		node := setupNodeOCR2(t, uni.owner, getFreePort(t), fmt.Sprintf("ocr2vrforacle%d", i), uni.backend, useForwarders, bootstrappers)

		dkgSignKey, err := node.app.GetKeyStore().DKGSign().Create()
		require.NoError(t, err)

		dkgEncryptKey, err := node.app.GetKeyStore().DKGEncrypt().Create()
		require.NoError(t, err)

		kbs = append(kbs, node.keybundle)
		apps = append(apps, node.app)
		transmitters = append(transmitters, node.transmitter)
		effectiveTransmitters = append(effectiveTransmitters, node.effectiveTransmitter)
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
		effectiveTransmitters,
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
fromBlock           = %d
`, uni.dkgAddress.Hex(), blockBeforeConfig.Number().Int64())
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
forwardingAllowed       = %t

[relayConfig]
chainID              	= 1337
fromBlock               = %d

[pluginConfig]
dkgEncryptionPublicKey 	= "%s"
dkgSigningPublicKey    	= "%s"
dkgKeyID               	= "%s"
dkgContractAddress     	= "%s"

vrfCoordinatorAddress   = "%s"
linkEthFeedAddress     	= "%s"
`, uni.beaconAddress.String(),
			kbs[i].ID(),
			transmitters[i],
			useForwarders,
			blockBeforeConfig.Number().Int64(),
			dkgEncrypters[i].PublicKeyString(),
			dkgSigners[i].PublicKeyString(),
			hex.EncodeToString(keyID[:]),
			uni.dkgAddress.String(),
			uni.coordinatorAddress.String(),
			uni.feedAddress.String(),
		)
		t.Log("Creating OCR2VRF job with spec:", jobSpec)
		ocrJob, err := validate.ValidatedOracleSpecToml(apps[i].Config, jobSpec)
		require.NoError(t, err)
		err = apps[i].AddJobV2(context.Background(), &ocrJob)
		require.NoError(t, err)
	}

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
		effectiveTransmitters,
		1,
		oracles,
		[]int{1, 2, 3, 4, 5, 6, 7, 8},
		keyID)

	t.Log("Sending VRF request")

	initialSub, err := uni.coordinator.GetSubscription(nil, uni.subID)
	require.NoError(t, err)
	require.Equal(t, assets.Ether(5).ToInt(), initialSub.Balance)

	// Send a beacon VRF request and mine it
	_, err = uni.consumer.TestRequestRandomness(uni.owner, 2, uni.subID, big.NewInt(1))
	require.NoError(t, err)
	uni.backend.Commit()

	// There is no premium on this request, so the cost of the request should have been:
	// = (request overhead) * (gas price) / (LINK/ETH ratio)
	// = (50_000 * 1 Gwei) / .01
	// = 5_000_000 GJuels
	subAfterBeaconRequest, err := uni.coordinator.GetSubscription(nil, uni.subID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(initialSub.Balance.Int64()-assets.GWei(5_000_000).Int64()), subAfterBeaconRequest.Balance)

	// Send a fulfillment VRF request and mine it
	_, err = uni.consumer.TestRequestRandomnessFulfillment(uni.owner, uni.subID, 1, big.NewInt(2), 100_000, []byte{})
	require.NoError(t, err)
	uni.backend.Commit()

	// There is no premium on this request, so the cost of the request should have been:
	// = (request overhead + callback gas allowance) * (gas price) / (LINK/ETH ratio)
	// = ((50_000 + 100_000) * 1 Gwei) / .01
	// = 15_000_000 GJuels
	subAfterFulfillmentRequest, err := uni.coordinator.GetSubscription(nil, uni.subID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(subAfterBeaconRequest.Balance.Int64()-assets.GWei(15_000_000).Int64()), subAfterFulfillmentRequest.Balance)

	// Send two batched fulfillment VRF requests and mine them
	_, err = uni.loadTestConsumer.TestRequestRandomnessFulfillmentBatch(uni.owner, uni.subID, 1, big.NewInt(2), 200_000, []byte{}, big.NewInt(2))
	require.NoError(t, err)
	uni.backend.Commit()

	// There is no premium on these requests, so the cost of the requests should have been:
	// = ((request overhead + callback gas allowance) * (gas price) / (LINK/ETH ratio)) * batch size
	// = (((50_000 + 200_000) * 1 Gwei) / .01) * 2
	// = 50_000_000 GJuels
	subAfterBatchFulfillmentRequest, err := uni.coordinator.GetSubscription(nil, uni.subID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(subAfterFulfillmentRequest.Balance.Int64()-assets.GWei(50_000_000).Int64()), subAfterBatchFulfillmentRequest.Balance)

	t.Logf("sub balance after batch fulfillment request: %d", subAfterBatchFulfillmentRequest.Balance)

	t.Log("waiting for fulfillment")

	// poll until we're able to redeem the randomness without reverting
	// at that point, it's been fulfilled
	gomega.NewWithT(t).Eventually(func() bool {
		// Ensure a refund is provided. Refund amount comes out to ~20_500_000 GJuels.
		// We use an upper and lower bound such that this part of the test is not excessively brittle to upstream tweaks.
		refundUpperBound := big.NewInt(0).Add(assets.GWei(21_500_000).ToInt(), subAfterBatchFulfillmentRequest.Balance)
		refundLowerBound := big.NewInt(0).Add(assets.GWei(19_500_000).ToInt(), subAfterBatchFulfillmentRequest.Balance)
		subAfterRefund, err := uni.coordinator.GetSubscription(nil, uni.subID)
		require.NoError(t, err)

		_, err1 := uni.consumer.TestRedeemRandomness(uni.owner, uni.subID, big.NewInt(0))
		t.Logf("TestRedeemRandomness err: %+v", err1)
		if err1 != nil {
			return false
		}

		if ok := ((subAfterRefund.Balance.Cmp(refundUpperBound) == -1) && (subAfterRefund.Balance.Cmp(refundLowerBound) == 1)); !ok {
			t.Logf("unexpected sub balance after refund: %d", subAfterRefund.Balance)
			return false
		}

		return true
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
	offchainConfig := ocr2vrf.OffchainConfig(&ocr2vrftypes.CoordinatorConfig{
		CacheEvictionWindowSeconds: 1,
		BatchGasLimit:              5_000_000,
		CoordinatorOverhead:        50_000,
		CallbackOverhead:           50_000,
		BlockGasOverhead:           50_000,
		LookbackBlocks:             1_000,
	})

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
	defer func() { assert.NoError(t, l.Close()) }()

	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func ptr[T any](v T) *T { return &v }
