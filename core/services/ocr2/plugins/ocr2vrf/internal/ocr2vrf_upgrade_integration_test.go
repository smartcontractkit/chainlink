package internal_test

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/generated/mock_v3_aggregator_contract"
	dkg_wrapper "github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/dkg"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/test_vrf_beacon_v1"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/test_vrf_beacon_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/test_vrf_coordinator_v1"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/test_vrf_coordinator_v2"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_consumer"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_beacon_proxy"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_coordinator_proxy"
	"github.com/smartcontractkit/chainlink/core/gethwrappers/ocr2vrf/generated/vrf_proxy_admin"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgencryptkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/dkgsignkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupOCR2VRFV2Contracts(
	t *testing.T,
	beaconPeriod int64,
	uni ocr2vrfUniverse) ocr2vrfUniverse {

	coordinatorImplAddress, _, _, err := test_vrf_coordinator_v2.DeployTestVRFCoordinatorV2(uni.owner, uni.backend)
	require.NoError(t, err)
	uni.backend.Commit()

	beaconImplAddress, _, _, err := test_vrf_beacon_v2.DeployTestVRFBeaconV2(uni.owner, uni.backend)
	require.NoError(t, err)
	uni.backend.Commit()

	coordinatorAbi, err := test_vrf_coordinator_v2.TestVRFCoordinatorV2MetaData.GetAbi()
	require.NoError(t, err)
	coordinatorCalldata, err := coordinatorAbi.Pack("initialize", big.NewInt(beaconPeriod), uni.linkAddress, uni.owner.From)
	require.NoError(t, err)

	beaconAbi, err := test_vrf_beacon_v2.TestVRFBeaconV2MetaData.GetAbi()
	require.NoError(t, err)
	beaconCalldata, err := beaconAbi.Pack("initialize", uni.linkAddress, uni.coordinatorAddress, uni.dkgAddress, uni.keyID, uni.owner.From, true)
	require.NoError(t, err)

	_, err = uni.proxyAdmin.VrfUpgradeAndCall(
		uni.owner,
		uni.coordinatorAddress,
		uni.beaconAddress,
		coordinatorImplAddress,
		beaconImplAddress,
		coordinatorCalldata,
		beaconCalldata,
	)
	require.NoError(t, err)

	uni.backend.Commit()

	beacon, err := test_vrf_beacon_v2.NewTestVRFBeaconV2(uni.beaconAddress, uni.backend)
	require.NoError(t, err)
	coordinator, err := test_vrf_coordinator_v2.NewTestVRFCoordinatorV2(uni.coordinatorAddress, uni.backend)
	require.NoError(t, err)

	callOpts := &bind.CallOpts{
		Context: testutils.Context(t),
	}

	tokenAddr := common.Address{19: 0x1}
	_, err = coordinator.RegisterToken(uni.owner, tokenAddr)
	require.NoError(t, err)
	uni.backend.Commit()

	tokenAddress, err := coordinator.AcceptedTokens(callOpts, big.NewInt(0))
	require.NoError(t, err)

	allowOnchainVerification, err := beacon.AllowOnchainVerification(callOpts)
	require.NoError(t, err)

	// assert allowOnchainVerification was set as expected through initializer
	assert.True(t, allowOnchainVerification)
	// assert RegisterToken successfully added an to a newly introduced array state variable
	assert.Equal(t, tokenAddr, tokenAddress)

	return ocr2vrfUniverse{
		owner:              uni.owner,
		backend:            uni.backend,
		dkgAddress:         uni.dkgAddress,
		dkg:                uni.dkg,
		beaconAddress:      uni.beaconAddress,
		coordinatorAddress: uni.coordinatorAddress,
		beacon:             uni.beacon,
		coordinator:        uni.coordinator,
		linkAddress:        uni.linkAddress,
		link:               uni.link,
		consumerAddress:    uni.consumerAddress,
		consumer:           uni.consumer,
		feedAddress:        uni.feedAddress,
		feed:               uni.feed,
		keyID:              uni.keyID,
	}
}

func setupOCR2VRFV1Contracts(
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

	proxyAdminAddress, _, proxyAdmin, err := vrf_proxy_admin.DeployVRFProxyAdmin(owner, b)
	b.Commit()

	coordinatorImplAddress, _, _, err := test_vrf_coordinator_v1.DeployTestVRFCoordinatorV1(owner, b)
	require.NoError(t, err)

	b.Commit()

	coordinatorAbi, err := vrf_coordinator.VRFCoordinatorMetaData.GetAbi()
	require.NoError(t, err)
	coordinatorCalldata, err := coordinatorAbi.Pack("initialize", big.NewInt(beaconPeriod), linkAddress, owner.From)
	require.NoError(t, err)

	coordinatorAddress, _, _, err := vrf_coordinator_proxy.DeployVRFCoordinatorProxy(
		owner, b, coordinatorImplAddress, proxyAdminAddress, coordinatorCalldata)
	require.NoError(t, err)

	b.Commit()

	beaconImplAddress, _, _, err := test_vrf_beacon_v1.DeployTestVRFBeaconV1(owner, b)
	require.NoError(t, err)

	b.Commit()

	beaconAbi, err := vrf_beacon.VRFBeaconMetaData.GetAbi()
	require.NoError(t, err)
	beaconCalldata, err := beaconAbi.Pack("initialize", linkAddress, coordinatorAddress, dkgAddress, keyID, owner.From)
	require.NoError(t, err)

	beaconAddress, _, _, err := vrf_beacon_proxy.DeployVRFBeaconProxy(
		owner, b, beaconImplAddress, proxyAdminAddress, beaconCalldata)

	b.Commit()

	consumerAddress, _, consumer, err := vrf_beacon_consumer.DeployBeaconVRFConsumer(
		owner, b, coordinatorAddress, consumerShouldFail, big.NewInt(beaconPeriod))
	require.NoError(t, err)

	b.Commit()

	_, err = dkg.AddClient(owner, keyID, beaconAddress)
	require.NoError(t, err)

	b.Commit()

	coordinator, err := vrf_coordinator.NewVRFCoordinator(coordinatorAddress, b)
	require.NoError(t, err)

	beacon, err := vrf_beacon.NewVRFBeacon(beaconAddress, b)
	require.NoError(t, err)

	_, err = coordinator.SetProducer(owner, beaconAddress)
	require.NoError(t, err)

	// Achieve finality depth so the CL node can work properly.
	for i := 0; i < 20; i++ {
		b.Commit()
	}

	return ocr2vrfUniverse{
		owner:              owner,
		backend:            b,
		dkgAddress:         dkgAddress,
		dkg:                dkg,
		beaconAddress:      beaconAddress,
		coordinatorAddress: coordinatorAddress,
		beacon:             beacon,
		coordinator:        coordinator,
		linkAddress:        linkAddress,
		link:               link,
		consumerAddress:    consumerAddress,
		consumer:           consumer,
		feedAddress:        feedAddress,
		feed:               feed,
		proxyAdminAddress:  proxyAdminAddress,
		proxyAdmin:         proxyAdmin,
		keyID:              keyID,
	}
}

func TestUpgrade_OCR2VRF(t *testing.T) {
	keyID := randomKeyID(t)
	uni := setupOCR2VRFV1Contracts(t, 5, keyID, false)

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
	defer bootstrapNode.app.Stop()

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
	var jobIDs []int32
	for i := 0; i < numNodes; i++ {
		err = apps[i].Start(testutils.Context(t))
		require.NoError(t, err)
		defer apps[i].Stop()

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
		jobIDs = append(jobIDs, ocrJob.ID)
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
	// read variables before upgrade.
	// after the upgrade, these variables are compared to make sure storage layout isn't changed in unexpected ways
	callOpts := &bind.CallOpts{
		Context: testutils.Context(t),
	}
	configDetails, err := uni.beacon.LatestConfigDetails(callOpts)
	require.NoError(t, err)
	billingDetails, err := uni.beacon.GetBilling(callOpts)
	require.NoError(t, err)
	confDelays, err := uni.coordinator.GetConfirmationDelays(callOpts)
	require.NoError(t, err)
	beaconPeriodBlocks, err := uni.coordinator.IBeaconPeriodBlocks(callOpts)
	require.NoError(t, err)

	// perform upgrade
	newBeaconPeriod := 10
	uniV2 := setupOCR2VRFV2Contracts(t, int64(newBeaconPeriod), uni)
	sendVRFRequestsAndVerify(t, uniV2)

	configDetailsV2, err := uniV2.beacon.LatestConfigDetails(callOpts)
	require.NoError(t, err)
	billingDetailsV2, err := uniV2.beacon.GetBilling(callOpts)
	require.NoError(t, err)
	confDelaysV2, err := uniV2.coordinator.GetConfirmationDelays(callOpts)
	require.NoError(t, err)
	beaconPeriodBlocksV2, err := uniV2.coordinator.IBeaconPeriodBlocks(callOpts)
	require.NoError(t, err)

	// assert below state variables are not changed
	assert.True(t, reflect.DeepEqual(configDetails, configDetailsV2))
	assert.True(t, reflect.DeepEqual(billingDetails, billingDetailsV2))
	assert.Equal(t, confDelays, confDelaysV2)
	// assert beaconPeriodBlocks was changed during the upgrade
	assert.Equal(t, int64(5), beaconPeriodBlocks.Int64())
	assert.Equal(t, int64(10), beaconPeriodBlocksV2.Int64())
}

func sendVRFRequestsAndVerify(t *testing.T, uni ocr2vrfUniverse) {
	t.Log("Sending VRF request")
	// Send a VRF request and mine it
	_, err := uni.consumer.TestRequestRandomness(uni.owner, 2, 1, big.NewInt(1))
	_, err = uni.consumer.TestRequestRandomnessFulfillment(uni.owner, 1, 1, big.NewInt(2), 50_000, []byte{})
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
		rw1, err1 := uni.consumer.SReceivedRandomnessByRequestID(nil, big.NewInt(0), big.NewInt(0))
		t.Logf("TestRedeemRandomness 1st word err: %+v", err1)
		rw2, err2 := uni.consumer.SReceivedRandomnessByRequestID(nil, big.NewInt(0), big.NewInt(1))
		t.Logf("TestRedeemRandomness 2nd word err: %+v", err2)
		rw3, err3 := uni.consumer.SReceivedRandomnessByRequestID(nil, big.NewInt(1), big.NewInt(0))
		t.Logf("FulfillRandomness 1st word err: %+v", err3)
		t.Log("randomness from redeemRandomness:", rw1.String(), rw2.String())
		t.Log("randomness from fulfillRandomness:", rw3.String())
		return err1 == nil && err2 == nil && err3 == nil
	}, testutils.WaitTimeout(t), 5*time.Second).Should(gomega.BeTrue())
}
