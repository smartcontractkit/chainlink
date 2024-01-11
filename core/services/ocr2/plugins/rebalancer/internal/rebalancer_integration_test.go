package internal_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/hashicorp/consul/sdk/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocrtypes "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/test-go/testify/require"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/shared/generated/no_op_ocr3"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/keystest"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
)

func TestRebalancer_Integration(t *testing.T) {
	newTestUniverse(t)
}

type ocr2Node struct {
	app                  *cltest.TestApplication
	peerID               string
	transmitter          common.Address
	effectiveTransmitter common.Address
	keybundle            ocr2key.KeyBundle
	sendingKeys          []string
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
	p2pKey := keystest.NewP2PKeyV2(t)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.Insecure.OCRDevelopmentMode = ptr(true) // Disables ocr spec validation so we can have fast polling for the test.

		c.Feature.LogPoller = ptr(true)

		c.P2P.PeerID = ptr(p2pKey.PeerID())
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
		c.OCR2.ContractPollInterval = models.MustNewDuration(5 * time.Second)
	})

	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, b, p2pKey)

	var sendingKeys []ethkey.KeyV2
	{
		var err error
		sendingKeys, err = app.KeyStore.Eth().EnabledKeysForChain(testutils.SimulatedChainID)
		require.NoError(t, err)
		require.Len(t, sendingKeys, 1)
	}
	transmitter := sendingKeys[0].Address
	effectiveTransmitter := sendingKeys[0].Address

	// Fund the sending keys with some ETH.
	var sendingKeyStrings []string
	for _, k := range sendingKeys {
		sendingKeyStrings = append(sendingKeyStrings, k.Address.String())
		n, err := b.NonceAt(testutils.Context(t), owner.From, nil)
		require.NoError(t, err)

		tx := cltest.NewLegacyTransaction(
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
		sendingKeys:          sendingKeyStrings,
	}
}

func newTestUniverse(t *testing.T) {
	ctx := testutils.Context(t)
	owner := testutils.MustNewSimTransactor(t)
	mainBackend := backends.NewSimulatedBackend(core.GenesisAlloc{
		owner.From: core.GenesisAccount{
			Balance: assets.Ether(1000).ToInt(),
		},
	}, 30e6)

	// deploy the ocr3 contract
	addr, _, _, err := no_op_ocr3.DeployNoOpOCR3(owner, mainBackend)
	require.NoError(t, err, "failed to deploy NoOpOCR3 contract")
	mainBackend.Commit()
	wrapper, err := no_op_ocr3.NewNoOpOCR3(addr, mainBackend)
	require.NoError(t, err, "failed to create NoOpOCR3 wrapper")

	t.Log("Creating bootstrap node")
	bootstrapNodePort := freeport.GetOne(t)
	bootstrapNode := setupNodeOCR2(t, owner, bootstrapNodePort, "bootstrap", mainBackend, false, nil)
	numNodes := 4

	t.Log("creating ocr3 nodes")
	var (
		oracles               []confighelper2.OracleIdentityExtra
		transmitters          []common.Address
		effectiveTransmitters []common.Address
		onchainPubKeys        []common.Address
		kbs                   []ocr2key.KeyBundle
		apps                  []*cltest.TestApplication
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
		node := setupNodeOCR2(t, owner, ports[i], fmt.Sprintf("ocr2vrforacle%d", i), mainBackend, false, bootstrappers)
		sendingKeys = append(sendingKeys, node.sendingKeys)

		kbs = append(kbs, node.keybundle)
		apps = append(apps, node.app)
		transmitters = append(transmitters, node.transmitter)
		effectiveTransmitters = append(effectiveTransmitters, node.effectiveTransmitter)
		onchainPubKeys = append(onchainPubKeys, common.BytesToAddress(node.keybundle.PublicKey()))
		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  node.keybundle.PublicKey(),
				TransmitAccount:   ocrtypes.Account(node.transmitter.String()),
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
			mainBackend.Commit()
		}
	}()

	blockBeforeConfig, err := mainBackend.BlockByNumber(ctx, nil)
	require.NoError(t, err)

	t.Log("setting config")
	setRebalancerConfig(t, owner, wrapper, mainBackend, onchainPubKeys, effectiveTransmitters, oracles)

	t.Log("adding bootstrap node job")
	err = bootstrapNode.app.Start(ctx)
	require.NoError(t, err, "failed to start bootstrap node")

	evmChains := bootstrapNode.app.GetRelayers().LegacyEVMChains()
	require.NotNil(t, evmChains)
	require.Len(t, evmChains.Slice(), 1)
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
`, addr.Hex(), blockBeforeConfig.Number().Uint64())
	t.Log("creating bootstrap job with spec:\n", bootstrapJobSpec)
	ocrJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(bootstrapJobSpec)
	require.NoError(t, err, "failed to validate bootstrap job")
	err = bootstrapNode.app.AddJobV2(ctx, &ocrJob)
	require.NoError(t, err, "failed to add bootstrap job")

	t.Log("creating ocr3 jobs")
	for i := 0; i < numNodes; i++ {
		var sendingKeysString = fmt.Sprintf(`"%s"`, sendingKeys[i][0])
		for x := 1; x < len(sendingKeys[i]); x++ {
			sendingKeysString = fmt.Sprintf(`%s,"%s"`, sendingKeysString, sendingKeys[i][x])
		}
		err = apps[i].Start(testutils.Context(t))
		require.NoError(t, err)

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
fromBlock               = %d

[pluginConfig]
liquidityManagerAddress = "%s"
liquidityManagerNetwork = %d
closePluginTimeoutSec = 10
[pluginConfig.rebalancerConfig]
type = "random"
[pluginConfig.rebalancerConfig.randomRebalancerConfig]
maxNumTransfers = 5
checkSourceDestEqual = false
`,
			addr.Hex(),
			kbs[i].ID(),
			transmitters[i].Hex(),
			blockBeforeConfig.Number().Uint64(),
			addr.Hex(),
			testutils.SimulatedChainID)
		t.Log("Creating rebalancer job with spec:\n", jobSpec)
		ocrJob2, err2 := validate.ValidatedOracleSpecToml(apps[i].Config.OCR2(), apps[i].Config.Insecure(), jobSpec)
		require.NoError(t, err2, "failed to validate rebalancer job")
		err2 = apps[i].AddJobV2(ctx, &ocrJob2)
		require.NoError(t, err2, "failed to add rebalancer job")
	}

	t.Log("waiting for a transmission")
	start := uint64(1)
	sink := make(chan *no_op_ocr3.NoOpOCR3Transmitted)
	sub, err := wrapper.WatchTransmitted(&bind.WatchOpts{
		Start: &start,
	}, sink)
	require.NoError(t, err, "failed to create subscription")
	defer sub.Unsubscribe()
	ticker := time.NewTicker(1 * time.Second)
outer:
	for {
		select {
		case te := <-sink:
			t.Log("got transmission event, config digest:", hexutil.Encode(te.ConfigDigest[:]), "seqNr:", te.SequenceNumber)
			break outer
		case <-ticker.C:
			t.Log("waiting for transmission event")
		}
	}

	t.Log("done")
}

func setRebalancerConfig(
	t *testing.T,
	owner *bind.TransactOpts,
	wrapper *no_op_ocr3.NoOpOCR3,
	mainBackend *backends.SimulatedBackend,
	onchainPubKeys,
	effectiveTransmitters []common.Address,
	oracles []confighelper2.OracleIdentityExtra) {
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
		3,
		schedule,
		oracles,
		offchainConfig,
		50*time.Millisecond,  // maxDurationQuery
		5*time.Second,        // maxDurationObservation
		10*time.Second,       // maxDurationShouldAcceptAttestedReport
		100*time.Millisecond, // maxDurationShouldTransmitAcceptedReport
		int(f),
		onchainConfig)
	require.NoError(t, err, "failed to create contract config")
	t.Log("onchain config:", hexutil.Encode(onchainConfig), "offchain config:", hexutil.Encode(offchainConfig))
	_, err = wrapper.SetOCR3Config(
		owner,
		onchainPubKeys,
		effectiveTransmitters,
		f,
		onchainConfig,
		offchainConfigVersion,
		offchainConfig)
	require.NoError(t, err, "failed to set config")
	mainBackend.Commit()
}

func ptr[T any](v T) *T { return &v }
