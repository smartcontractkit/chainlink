package internal_test

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/kyber/v3"

	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2plus/types"

	"github.com/smartcontractkit/chainlink-vrf/altbn_128"
	ocr2dkg "github.com/smartcontractkit/chainlink-vrf/dkg"
	"github.com/smartcontractkit/chainlink-vrf/ocr2vrf"
	ocr2vrftypes "github.com/smartcontractkit/chainlink-vrf/types"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	commonutils "github.com/smartcontractkit/chainlink-common/pkg/utils"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/forwarders"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils"
	ubig "github.com/smartcontractkit/chainlink/v2/core/chains/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/authorized_forwarder"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/mock_v3_aggregator_contract"
	dkg_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/load_test_beacon_consumer"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	vrf_wrapper "github.com/smartcontractkit/chainlink/v2/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
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

	subID *big.Int
}

const (
	fundingAmount int64 = 5e18
)

type ocr2Node struct {
	app                  *cltest.TestApplication
	peerID               string
	transmitter          common.Address
	effectiveTransmitter common.Address
	keybundle            ocr2key.KeyBundle
	sendingKeys          []string
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
	// * VRF (coordinator, and beacon)
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

	coordinatorAddress, _, coordinator, err := vrf_wrapper.DeployVRFCoordinator(
		owner, b, big.NewInt(beaconPeriod), linkAddress)
	require.NoError(t, err)
	b.Commit()

	require.NoError(t, commonutils.JustError(coordinator.SetCallbackConfig(owner, vrf_wrapper.VRFCoordinatorCallbackConfig{
		MaxCallbackGasLimit:        2.5e6,
		MaxCallbackArgumentsLength: 160, // 5 EVM words
	})))
	b.Commit()

	require.NoError(t, commonutils.JustError(coordinator.SetCoordinatorConfig(owner, vrf_wrapper.VRFBeaconTypesCoordinatorConfig{
		RedeemableRequestGasOverhead: 50_000,
		CallbackRequestGasOverhead:   50_000,
		StalenessSeconds:             60,
		FallbackWeiPerUnitLink:       assets.GWei(int(1e7)).ToInt(),
	})))
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

	// Set up coordinator subscription for billing.
	require.NoError(t, commonutils.JustError(coordinator.CreateSubscription(owner)))
	b.Commit()

	fopts := &bind.FilterOpts{}

	subscriptionIterator, err := coordinator.FilterSubscriptionCreated(fopts, nil, []common.Address{owner.From})
	require.NoError(t, err)

	require.True(t, subscriptionIterator.Next())
	subID := subscriptionIterator.Event.SubId

	require.NoError(t, commonutils.JustError(coordinator.AddConsumer(owner, subID, consumerAddress)))
	b.Commit()
	require.NoError(t, commonutils.JustError(coordinator.AddConsumer(owner, subID, loadTestConsumerAddress)))
	b.Commit()
	data, err := utils.ABIEncode(`[{"type":"uint256"}]`, subID)
	require.NoError(t, err)
	require.NoError(t, commonutils.JustError(link.TransferAndCall(owner, coordinatorAddress, big.NewInt(5e18), data)))
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
		subID:                   subID,
	}
}

func setupNodeOCR2(
	t *testing.T,
	owner *bind.TransactOpts,
	port int,
	dbName string,
	b *backends.SimulatedBackend,
	useForwarders bool,
	p2pV2Bootstrappers []commontypes.BootstrapperLocator,
) *ocr2Node {
	ctx := testutils.Context(t)
	p2pKey := keystest.NewP2PKeyV2(t)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
		c.P2P.V2.Enabled = ptr(true)
		c.P2P.V2.DeltaDial = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.P2P.V2.DeltaReconcile = commonconfig.MustNewDuration(5 * time.Second)
		c.P2P.V2.ListenAddresses = &[]string{fmt.Sprintf("127.0.0.1:%d", port)}
		if len(p2pV2Bootstrappers) > 0 {
			c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
		}

		c.OCR.Enabled = ptr(false)
		c.OCR2.Enabled = ptr(true)

		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(500 * time.Millisecond)
		c.EVM[0].GasEstimator.LimitDefault = ptr[uint64](3_500_000)
		c.EVM[0].Transactions.ForwardersEnabled = &useForwarders
		c.OCR2.ContractPollInterval = commonconfig.MustNewDuration(10 * time.Second)
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	var sendingKeys []ethkey.KeyV2
	{
		var err error
		sendingKeys, err = app.KeyStore.Eth().EnabledKeysForChain(ctx, testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Len(t, sendingKeys, 1)
	}
	transmitter := sendingKeys[0].Address
	effectiveTransmitter := sendingKeys[0].Address

	if useForwarders {
		sendingKeysAddresses := []common.Address{sendingKeys[0].Address}

		// Add new sending key.
		k, err := app.KeyStore.Eth().Create(ctx)
		require.NoError(t, err)
		require.NoError(t, app.KeyStore.Eth().Add(ctx, k.Address, testutils.SimulatedChainID))
		require.NoError(t, app.KeyStore.Eth().Enable(ctx, k.Address, testutils.SimulatedChainID))
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
		forwarderORM := forwarders.NewORM(app.GetDB())
		chainID := ubig.Big(*b.Blockchain().Config().ChainID)
		_, err = forwarderORM.CreateForwarder(testutils.Context(t), faddr, chainID)
		require.NoError(t, err)
		effectiveTransmitter = faddr
	}

	// Fund the sending keys with some ETH.
	var sendingKeyStrings []string
	for _, k := range sendingKeys {
		sendingKeyStrings = append(sendingKeyStrings, k.Address.String())
		n, err := b.NonceAt(ctx, owner.From, nil)
		require.NoError(t, err)

		tx := cltest.NewLegacyTransaction(
			n, k.Address,
			assets.Ether(1).ToInt(),
			21000,
			assets.GWei(1).ToInt(),
			nil)
		signedTx, err := owner.Signer(owner.From, tx)
		require.NoError(t, err)
		err = b.SendTransaction(ctx, signedTx)
		require.NoError(t, err)
		b.Commit()
	}

	kb, err := app.GetKeyStore().OCR2().Create(ctx, "evm")
	require.NoError(t, err)

	return &ocr2Node{
		app:                  app,
		peerID:               p2pKey.PeerID().Raw(),
		transmitter:          transmitter,
		effectiveTransmitter: effectiveTransmitter,
		keybundle:            kb,
		sendingKeys:          sendingKeyStrings,
	}
}

func TestIntegration_OCR2VRF_ForwarderFlow(t *testing.T) {
	testutils.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/VRF-688")
	runOCR2VRFTest(t, true)
}

func TestIntegration_OCR2VRF(t *testing.T) {
	testutils.SkipFlakey(t, "https://smartcontract-it.atlassian.net/browse/VRF-688")
	runOCR2VRFTest(t, false)
}

func runOCR2VRFTest(t *testing.T, useForwarders bool) {
	ctx := testutils.Context(t)
	keyID := randomKeyID(t)
	uni := setupOCR2VRFContracts(t, 5, keyID, false)

	t.Log("Creating bootstrap node")

	bootstrapNodePort := freeport.GetOne(t)
	bootstrapNode := setupNodeOCR2(t, uni.owner, bootstrapNodePort, "bootstrap", uni.backend, false, nil)
	numNodes := 5

	t.Log("Creating OCR2 nodes")
	var (
		oracles               []confighelper2.OracleIdentityExtra
		transmitters          []common.Address
		payees                []common.Address
		payeeTransactors      []*bind.TransactOpts
		effectiveTransmitters []common.Address
		onchainPubKeys        []common.Address
		kbs                   []ocr2key.KeyBundle
		apps                  []*cltest.TestApplication
		dkgEncrypters         []dkgencryptkey.Key
		dkgSigners            []dkgsignkey.Key
		sendingKeys           [][]string
	)
	ports := freeport.GetN(t, numNodes)
	for i := 0; i < numNodes; i++ {
		// Supply the bootstrap IP and port as a V2 peer address
		bootstrappers := []commontypes.BootstrapperLocator{
			{PeerID: bootstrapNode.peerID, Addrs: []string{
				fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort),
			}},
		}
		node := setupNodeOCR2(t, uni.owner, ports[i], fmt.Sprintf("ocr2vrforacle%d", i), uni.backend, useForwarders, bootstrappers)
		sendingKeys = append(sendingKeys, node.sendingKeys)

		dkgSignKey, err := node.app.GetKeyStore().DKGSign().Create(ctx)
		require.NoError(t, err)

		dkgEncryptKey, err := node.app.GetKeyStore().DKGEncrypt().Create(ctx)
		require.NoError(t, err)

		kbs = append(kbs, node.keybundle)
		apps = append(apps, node.app)
		transmitters = append(transmitters, node.transmitter)
		payeeTransactor := testutils.MustNewSimTransactor(t)
		payeeTransactors = append(payeeTransactors, payeeTransactor)
		payees = append(payees, payeeTransactor.From)
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

	_, err := uni.beacon.SetPayees(uni.owner, transmitters, payees)
	require.NoError(t, err)

	t.Log("starting ticker to commit blocks")
	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go func() {
		for range tick.C {
			uni.backend.Commit()
		}
	}()

	blockBeforeConfig, err := uni.backend.BlockByNumber(ctx, nil)
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
	err = bootstrapNode.app.Start(ctx)
	require.NoError(t, err)

	evmChains := bootstrapNode.app.GetRelayers().LegacyEVMChains()
	require.NotNil(t, evmChains)
	bootstrapJobSpec := fmt.Sprintf(`
type				= "bootstrap"
name				= "bootstrap"
contractConfigTrackerPollInterval = "15s"
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
	err = bootstrapNode.app.AddJobV2(ctx, &ocrJob)
	require.NoError(t, err)

	t.Log("Creating OCR2VRF jobs")
	for i := 0; i < numNodes; i++ {
		var sendingKeysString = fmt.Sprintf(`"%s"`, sendingKeys[i][0])
		for x := 1; x < len(sendingKeys[i]); x++ {
			sendingKeysString = fmt.Sprintf(`%s,"%s"`, sendingKeysString, sendingKeys[i][x])
		}
		err = apps[i].Start(ctx)
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
contractConfigTrackerPollInterval = "15s"

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
		ocrJob2, err2 := validate.ValidatedOracleSpecToml(testutils.Context(t), apps[i].Config.OCR2(), apps[i].Config.Insecure(), jobSpec, nil)
		require.NoError(t, err2)
		err2 = apps[i].AddJobV2(ctx, &ocrJob2)
		require.NoError(t, err2)
	}

	t.Log("Waiting for DKG key to get written")
	// poll until a DKG key is written to the contract
	// at that point we can start sending VRF requests
	var emptyKH [32]byte
	emptyHash := crypto.Keccak256Hash(emptyKH[:])
	gomega.NewWithT(t).Eventually(func() bool {
		kh, err2 := uni.beacon.SProvingKeyHash(&bind.CallOpts{
			Context: ctx,
		})
		require.NoError(t, err2)
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

	redemptionRequestID, err := uni.consumer.SMostRecentRequestID(nil)
	require.NoError(t, err)

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

	fulfillmentRequestID, err := uni.consumer.SMostRecentRequestID(nil)
	require.NoError(t, err)

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

	batchFulfillmentRequestID1, err := uni.loadTestConsumer.SRequestIDs(nil, big.NewInt(0), big.NewInt(0))
	require.NoError(t, err)

	batchFulfillmentRequestID2, err := uni.loadTestConsumer.SRequestIDs(nil, big.NewInt(0), big.NewInt(1))
	require.NoError(t, err)

	// There is no premium on these requests, so the cost of the requests should have been:
	// = ((request overhead + callback gas allowance) * (gas price) / (LINK/ETH ratio)) * batch size
	// = (((50_000 + 200_000) * 1 Gwei) / .01) * 2
	// = 50_000_000 GJuels
	subAfterBatchFulfillmentRequest, err := uni.coordinator.GetSubscription(nil, uni.subID)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(subAfterFulfillmentRequest.Balance.Int64()-assets.GWei(50_000_000).Int64()), subAfterBatchFulfillmentRequest.Balance)

	t.Logf("sub balance after batch fulfillment request: %d", subAfterBatchFulfillmentRequest.Balance)

	t.Log("waiting for fulfillment")

	var balanceAfterRefund *big.Int
	// poll until we're able to redeem the randomness without reverting
	// at that point, it's been fulfilled
	gomega.NewWithT(t).Eventually(func() bool {
		_, err2 := uni.consumer.TestRedeemRandomness(uni.owner, uni.subID, redemptionRequestID)
		t.Logf("TestRedeemRandomness err: %+v", err2)
		return err2 == nil
	}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())

	gomega.NewWithT(t).Eventually(func() bool {
		// Ensure a refund is provided. Refund amount comes out to ~15_700_000 GJuels.
		// We use an upper and lower bound such that this part of the test is not excessively brittle to upstream tweaks.
		refundUpperBound := big.NewInt(0).Add(assets.GWei(17_000_000).ToInt(), subAfterBatchFulfillmentRequest.Balance)
		refundLowerBound := big.NewInt(0).Add(assets.GWei(15_000_000).ToInt(), subAfterBatchFulfillmentRequest.Balance)
		subAfterRefund, err2 := uni.coordinator.GetSubscription(nil, uni.subID)
		require.NoError(t, err2)
		balanceAfterRefund = subAfterRefund.Balance
		if ok := ((balanceAfterRefund.Cmp(refundUpperBound) == -1) && (balanceAfterRefund.Cmp(refundLowerBound) == 1)); !ok {
			t.Logf("unexpected sub balance after refund: %d", balanceAfterRefund)
			return false
		}
		return true
	}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())

	// Mine block after redeeming randomness
	uni.backend.Commit()

	// ensure that total sub balance is updated correctly
	totalSubBalance, err := uni.coordinator.GetSubscriptionLinkBalance(nil)
	require.NoError(t, err)
	require.True(t, totalSubBalance.Cmp(balanceAfterRefund) == 0)
	// ensure total link balance is correct before any payout
	totalLinkBalance, err := uni.link.BalanceOf(nil, uni.coordinatorAddress)
	require.NoError(t, err)
	require.True(t, totalLinkBalance.Cmp(big.NewInt(fundingAmount)) == 0)

	// get total owed amount to NOPs and ensure linkAvailableForPayment (CLL profit) calculation is correct
	nopOwedAmount := new(big.Int)
	for _, transmitter := range transmitters {
		owedAmount, err2 := uni.beacon.OwedPayment(nil, transmitter)
		require.NoError(t, err2)
		nopOwedAmount = new(big.Int).Add(nopOwedAmount, owedAmount)
	}
	linkAvailable, err := uni.beacon.LinkAvailableForPayment(nil)
	require.NoError(t, err)
	debt := new(big.Int).Add(totalSubBalance, nopOwedAmount)
	profit := new(big.Int).Sub(totalLinkBalance, debt)
	require.True(t, linkAvailable.Cmp(profit) == 0)

	// test cancel subscription
	linkBalanceBeforeCancel, err := uni.link.BalanceOf(nil, uni.owner.From)
	require.NoError(t, err)
	_, err = uni.coordinator.CancelSubscription(uni.owner, uni.subID, uni.owner.From)
	require.NoError(t, err)
	uni.backend.Commit()
	linkBalanceAfterCancel, err := uni.link.BalanceOf(nil, uni.owner.From)
	require.NoError(t, err)
	require.True(t, new(big.Int).Add(linkBalanceBeforeCancel, totalSubBalance).Cmp(linkBalanceAfterCancel) == 0)
	totalSubBalance, err = uni.coordinator.GetSubscriptionLinkBalance(nil)
	require.NoError(t, err)
	require.True(t, totalSubBalance.Cmp(big.NewInt(0)) == 0)
	totalLinkBalance, err = uni.link.BalanceOf(nil, uni.coordinatorAddress)
	require.NoError(t, err)
	require.True(t, totalLinkBalance.Cmp(new(big.Int).Sub(big.NewInt(fundingAmount), balanceAfterRefund)) == 0)

	// payout node operators
	totalNopPayout := new(big.Int)
	for idx, payeeTransactor := range payeeTransactors {
		// Fund the payee with some ETH.
		n, err2 := uni.backend.NonceAt(ctx, uni.owner.From, nil)
		require.NoError(t, err2)
		tx := cltest.NewLegacyTransaction(
			n, payeeTransactor.From,
			assets.Ether(1).ToInt(),
			21000,
			assets.GWei(1).ToInt(),
			nil)
		signedTx, err2 := uni.owner.Signer(uni.owner.From, tx)
		require.NoError(t, err2)
		err2 = uni.backend.SendTransaction(ctx, signedTx)
		require.NoError(t, err2)

		_, err2 = uni.beacon.WithdrawPayment(payeeTransactor, transmitters[idx])
		require.NoError(t, err2)
		uni.backend.Commit()
		payoutAmount, err2 := uni.link.BalanceOf(nil, payeeTransactor.From)
		require.NoError(t, err2)
		totalNopPayout = new(big.Int).Add(totalNopPayout, payoutAmount)
		owedAmountAfter, err2 := uni.beacon.OwedPayment(nil, transmitters[idx])
		require.NoError(t, err2)
		require.True(t, owedAmountAfter.Cmp(big.NewInt(0)) == 0)
	}
	require.True(t, nopOwedAmount.Cmp(totalNopPayout) == 0)

	// check total link balance after NOP payout
	totalLinkBalanceAfterNopPayout, err := uni.link.BalanceOf(nil, uni.coordinatorAddress)
	require.NoError(t, err)
	require.True(t, totalLinkBalanceAfterNopPayout.Cmp(new(big.Int).Sub(totalLinkBalance, totalNopPayout)) == 0)
	totalSubBalance, err = uni.coordinator.GetSubscriptionLinkBalance(nil)
	require.NoError(t, err)
	require.True(t, totalSubBalance.Cmp(big.NewInt(0)) == 0)

	// withdraw remaining profits after NOP payout
	linkAvailable, err = uni.beacon.LinkAvailableForPayment(nil)
	require.NoError(t, err)
	linkBalanceBeforeWithdraw, err := uni.link.BalanceOf(nil, uni.owner.From)
	require.NoError(t, err)
	_, err = uni.beacon.WithdrawFunds(uni.owner, uni.owner.From, linkAvailable)
	require.NoError(t, err)
	uni.backend.Commit()
	linkBalanceAfterWithdraw, err := uni.link.BalanceOf(nil, uni.owner.From)
	require.NoError(t, err)
	require.True(t, linkBalanceAfterWithdraw.Cmp(new(big.Int).Add(linkBalanceBeforeWithdraw, linkAvailable)) == 0)
	linkAvailable, err = uni.beacon.LinkAvailableForPayment(nil)
	require.NoError(t, err)
	require.True(t, linkAvailable.Cmp(big.NewInt(0)) == 0)

	// poll until we're able to verify that consumer contract has stored randomness as expected
	// First arg is the request ID, which starts at zero, second is the index into
	// the random words.
	gomega.NewWithT(t).Eventually(func() bool {
		var errs []error
		rw1, err2 := uni.consumer.SReceivedRandomnessByRequestID(nil, redemptionRequestID, big.NewInt(0))
		t.Logf("TestRedeemRandomness 1st word err: %+v", err2)
		errs = append(errs, err2)
		rw2, err2 := uni.consumer.SReceivedRandomnessByRequestID(nil, redemptionRequestID, big.NewInt(1))
		t.Logf("TestRedeemRandomness 2nd word err: %+v", err2)
		errs = append(errs, err2)
		rw3, err2 := uni.consumer.SReceivedRandomnessByRequestID(nil, fulfillmentRequestID, big.NewInt(0))
		t.Logf("FulfillRandomness 1st word err: %+v", err2)
		errs = append(errs, err2)
		rw4, err2 := uni.loadTestConsumer.SReceivedRandomnessByRequestID(nil, batchFulfillmentRequestID1, big.NewInt(0))
		t.Logf("Batch FulfillRandomness 1st word err: %+v", err2)
		errs = append(errs, err2)
		rw5, err2 := uni.loadTestConsumer.SReceivedRandomnessByRequestID(nil, batchFulfillmentRequestID2, big.NewInt(0))
		t.Logf("Batch FulfillRandomness 2nd word err: %+v", err2)
		errs = append(errs, err2)
		batchTotalRequests, err2 := uni.loadTestConsumer.STotalRequests(nil)
		t.Logf("Batch FulfillRandomness total requests err: %+v", err2)
		errs = append(errs, err2)
		batchTotalFulfillments, err2 := uni.loadTestConsumer.STotalFulfilled(nil)
		t.Logf("Batch FulfillRandomness total fulfillments err: %+v", err2)
		errs = append(errs, err2)
		err2 = nil
		if batchTotalRequests.Int64() != batchTotalFulfillments.Int64() {
			err2 = errors.New("batchTotalRequests is not equal to batchTotalFulfillments")
			errs = append(errs, err2)
		}
		t.Logf("Batch FulfillRandomness total requests/fulfillments equal err: %+v", err2)

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
		20*time.Second,
		2*time.Second,
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
		20*time.Second,
		2*time.Second,
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

func ptr[T any](v T) *T { return &v }
