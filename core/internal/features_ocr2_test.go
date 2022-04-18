package internal_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	testoffchainaggregator2 "github.com/smartcontractkit/libocr/gethwrappers2/testocr2aggregator"
	ocrnetworking "github.com/smartcontractkit/libocr/networking"
	confighelper2 "github.com/smartcontractkit/libocr/offchainreporting2/confighelper"
	ocrtypes2 "github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/bridges"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/link_token_interface"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ocr2key"
	"github.com/smartcontractkit/chainlink/core/services/ocr2/validate"
	"github.com/smartcontractkit/chainlink/core/services/ocrbootstrap"
	"github.com/smartcontractkit/chainlink/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

func setupOCR2Contracts(t *testing.T) (*bind.TransactOpts, *backends.SimulatedBackend, common.Address, *ocr2aggregator.OCR2Aggregator) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err, "failed to generate ethereum identity")
	owner := cltest.MustNewSimulatedBackendKeyedTransactor(t, key)
	sb := new(big.Int)
	sb, _ = sb.SetString("100000000000000000000", 10) // 1 eth
	genesisData := core.GenesisAlloc{owner.From: {Balance: sb}}
	gasLimit := ethconfig.Defaults.Miner.GasCeil * 2
	b := backends.NewSimulatedBackend(genesisData, gasLimit)
	linkTokenAddress, _, linkContract, err := link_token_interface.DeployLinkToken(owner, b)
	require.NoError(t, err)
	accessAddress, _, _, err := testoffchainaggregator2.DeploySimpleWriteAccessController(owner, b)
	require.NoError(t, err, "failed to deploy test access controller contract")
	b.Commit()

	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))
	ocrContractAddress, _, ocrContract, err := ocr2aggregator.DeployOCR2Aggregator(
		owner,
		b,
		linkTokenAddress, //_link common.Address,
		minAnswer,        // -2**191
		maxAnswer,        // 2**191 - 1
		accessAddress,
		accessAddress,
		9,
		"TEST",
	)

	require.NoError(t, err)
	_, err = linkContract.Transfer(owner, ocrContractAddress, big.NewInt(1000))
	require.NoError(t, err)
	b.Commit()
	return owner, b, ocrContractAddress, ocrContract
}

func setupNodeOCR2(t *testing.T, owner *bind.TransactOpts, port uint16, dbName string, b *backends.SimulatedBackend) (*cltest.TestApplication, string, common.Address, ocr2key.KeyBundle, *configtest.TestGeneralConfig) {
	config, _ := heavyweight.FullTestDB(t, fmt.Sprintf("%s%d", dbName, port))
	config.Overrides.FeatureOffchainReporting = null.BoolFrom(false)
	config.Overrides.FeatureOffchainReporting2 = null.BoolFrom(true)
	config.Overrides.P2PEnabled = null.BoolFrom(true)
	config.Overrides.P2PNetworkingStack = ocrnetworking.NetworkingStackV2
	config.Overrides.P2PListenPort = null.NewInt(0, true)
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

	sendingKeys, err := app.KeyStore.Eth().SendingKeys(nil)
	require.NoError(t, err)
	require.Len(t, sendingKeys, 1)
	transmitter := sendingKeys[0].Address.Address()

	// Fund the transmitter address with some ETH
	n, err := b.NonceAt(context.Background(), owner.From, nil)
	require.NoError(t, err)

	tx := types.NewTransaction(n, transmitter, big.NewInt(1000000000000000000), 21000, big.NewInt(1000000000), nil)
	signedTx, err := owner.Signer(owner.From, tx)
	require.NoError(t, err)
	err = b.SendTransaction(context.Background(), signedTx)
	require.NoError(t, err)
	b.Commit()

	kb, err := app.GetKeyStore().OCR2().Create("evm")
	require.NoError(t, err)
	return app, peerID.Raw(), transmitter, kb, config
}

func TestIntegration_OCR2(t *testing.T) {
	owner, b, ocrContractAddress, ocrContract := setupOCR2Contracts(t)

	lggr := logger.TestLogger(t)
	// Note it's plausible these ports could be occupied on a CI machine.
	// May need a port randomize + retry approach if we observe collisions.
	bootstrapNodePort := uint16(29999)
	appBootstrap, bootstrapPeerID, _, _, _ := setupNodeOCR2(t, owner, bootstrapNodePort, "bootstrap", b)

	var (
		oracles      []confighelper2.OracleIdentityExtra
		transmitters []common.Address
		kbs          []ocr2key.KeyBundle
		apps         []*cltest.TestApplication
	)
	for i := uint16(0); i < 4; i++ {
		app, peerID, transmitter, kb, cfg := setupNodeOCR2(t, owner, bootstrapNodePort+1+i, fmt.Sprintf("oracle%d", i), b)
		// Supply the bootstrap IP and port as a V2 peer address
		cfg.Overrides.P2PV2Bootstrappers = []commontypes.BootstrapperLocator{
			{PeerID: bootstrapPeerID, Addrs: []string{
				fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort),
			}},
		}

		kbs = append(kbs, kb)
		apps = append(apps, app)
		transmitters = append(transmitters, transmitter)

		oracles = append(oracles, confighelper2.OracleIdentityExtra{
			OracleIdentity: confighelper2.OracleIdentity{
				OnchainPublicKey:  kb.PublicKey(),
				TransmitAccount:   ocrtypes2.Account(transmitter.String()),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}

	tick := time.NewTicker(1 * time.Second)
	defer tick.Stop()
	go func() {
		for range tick.C {
			b.Commit()
		}
	}()

	lggr.Debugw("Setting Payees on OraclePlugin Contract", "transmitters", transmitters)
	_, err := ocrContract.SetPayees(
		owner,
		transmitters,
		transmitters,
	)
	require.NoError(t, err)
	signers, transmitters, threshold, onchainConfig, encodedConfigVersion, encodedConfig, err := confighelper2.ContractSetConfigArgsForEthereumIntegrationTest(
		oracles,
		1,
		1000000000/100, // threshold PPB
	)
	require.NoError(t, err)
	lggr.Debugw("Setting Config on Oracle Contract",
		"signers", signers,
		"transmitters", transmitters,
		"threshold", threshold,
		"onchainConfig", onchainConfig,
		"encodedConfigVersion", encodedConfigVersion,
	)
	_, err = ocrContract.SetConfig(
		owner,
		signers,
		transmitters,
		threshold,
		onchainConfig,
		encodedConfigVersion,
		encodedConfig,
	)
	require.NoError(t, err)
	b.Commit()

	err = appBootstrap.Start(testutils.Context(t))
	require.NoError(t, err)
	defer appBootstrap.Stop()

	chainSet := appBootstrap.GetChains().EVM
	require.NotNil(t, chainSet)
	ocrJob, err := ocrbootstrap.ValidatedBootstrapSpecToml(fmt.Sprintf(`
type				= "bootstrap"
name				= "bootstrap"
relay				= "evm"
schemaVersion		= 1
contractID			= "%s"
[relayConfig]
chainID 			= 1337
`, ocrContractAddress))
	require.NoError(t, err)
	err = appBootstrap.AddJobV2(context.Background(), &ocrJob)
	require.NoError(t, err)

	var jids []int32
	var servers, slowServers = make([]*httptest.Server, 4), make([]*httptest.Server, 4)
	// We expect metadata of:
	//  latestAnswer:nil // First call
	//  latestAnswer:0
	//  latestAnswer:10
	//  latestAnswer:20
	//  latestAnswer:30
	var metaLock sync.Mutex
	expectedMeta := map[string]struct{}{
		"0": {}, "10": {}, "20": {}, "30": {},
	}
	for i := 0; i < 4; i++ {
		err = apps[i].Start(testutils.Context(t))
		require.NoError(t, err)
		defer apps[i].Stop()

		// API speed is > observation timeout set in ContractSetConfigArgsForIntegrationTest
		slowServers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(5 * time.Second)
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"data":10}`))
		}))
		defer slowServers[i].Close()
		servers[i] = httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			b, err := ioutil.ReadAll(req.Body)
			require.NoError(t, err)
			var m bridges.BridgeMetaDataJSON
			require.NoError(t, json.Unmarshal(b, &m))
			if m.Meta.LatestAnswer != nil && m.Meta.UpdatedAt != nil {
				metaLock.Lock()
				delete(expectedMeta, m.Meta.LatestAnswer.String())
				metaLock.Unlock()
			}
			res.WriteHeader(http.StatusOK)
			res.Write([]byte(`{"data":10}`))
		}))
		defer servers[i].Close()
		u, _ := url.Parse(servers[i].URL)
		apps[i].BridgeORM().CreateBridgeType(&bridges.BridgeType{
			Name: bridges.BridgeName(fmt.Sprintf("bridge%d", i)),
			URL:  models.WebURL(*u),
		})

		ocrJob, err := validate.ValidatedOracleSpecToml(apps[i].Config, fmt.Sprintf(`
type               = "offchainreporting2"
relay              = "evm"
schemaVersion      = 1
pluginType         = "median"
name               = "web oracle spec"
contractID         = "%s"
ocrKeyBundleID     = "%s"
transmitterID      = "%s"
contractConfigConfirmations = 1
contractConfigTrackerPollInterval = "1s"
observationSource  = """
    // data source 1
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="data"];
    ds1_multiply [type=multiply times=%d];

    // data source 2
    ds2          [type=http method=GET url="%s"];
    ds2_parse    [type=jsonparse path="data"];
    ds2_multiply [type=multiply times=%d];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""
[relayConfig]
chainID = 1337
[pluginConfig]
juelsPerFeeCoinSource = """
		// data source 1
		ds1          [type=bridge name="%s"];
		ds1_parse    [type=jsonparse path="data"];
		ds1_multiply [type=multiply times=%d];

		// data source 2
		ds2          [type=http method=GET url="%s"];
		ds2_parse    [type=jsonparse path="data"];
		ds2_multiply [type=multiply times=%d];

		ds1 -> ds1_parse -> ds1_multiply -> answer1;
		ds2 -> ds2_parse -> ds2_multiply -> answer1;

	answer1 [type=median index=0];
"""
`, ocrContractAddress, kbs[i].ID(), transmitters[i], fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i, fmt.Sprintf("bridge%d", i), i, slowServers[i].URL, i))
		require.NoError(t, err)
		err = apps[i].AddJobV2(context.Background(), &ocrJob)
		require.NoError(t, err)
		jids = append(jids, ocrJob.ID)
	}

	// Assert that all the OCR jobs get a run with valid values eventually.
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		ic := i
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Want at least 2 runs so we see all the metadata.
			pr := cltest.WaitForPipelineComplete(t, ic, jids[ic], 2, 7, apps[ic].JobORM(), 1*time.Minute, 1*time.Second)
			jb, err := pr[0].Outputs.MarshalJSON()
			require.NoError(t, err)
			assert.Equal(t, []byte(fmt.Sprintf("[\"%d\"]", 10*ic)), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1])
			require.NoError(t, err)
		}()
	}
	wg.Wait()

	// 4 oracles reporting 0, 10, 20, 30. Answer should be 20 (results[4/2]).
	gomega.NewGomegaWithT(t).Eventually(func() string {
		answer, err := ocrContract.LatestAnswer(nil)
		require.NoError(t, err)
		return answer.String()
	}, 1*time.Minute, 200*time.Millisecond).Should(gomega.Equal("20"))

	for _, app := range apps {
		jobs, _, err := app.JobORM().FindJobs(0, 1000)
		require.NoError(t, err)
		// No spec errors
		for _, j := range jobs {
			ignore := 0
			for i := range j.JobSpecErrors {
				// Non-fatal timing related error, ignore for testing.
				if strings.Contains(j.JobSpecErrors[i].Description, "leader's phase conflicts tGrace timeout") {
					ignore++
				}
			}
			require.Len(t, j.JobSpecErrors, ignore)
		}
	}
	assert.Len(t, expectedMeta, 0, "expected metadata %v", expectedMeta)

	// Assert we can read the latest config digest and epoch after a report has been submitted.
	contractABI, err := abi.JSON(strings.NewReader(ocr2aggregator.OCR2AggregatorABI))
	require.NoError(t, err)
	ct := evm.NewOCRContractTransmitter(ocrContractAddress, apps[0].Chains.EVM.Chains()[0].Client(), contractABI, nil, lggr)
	configDigest, epoch, err := ct.LatestConfigDigestAndEpoch(context.Background())
	require.NoError(t, err)
	details, err := ocrContract.LatestConfigDetails(nil)
	require.NoError(t, err)
	assert.True(t, bytes.Equal(configDigest[:], details.ConfigDigest[:]))
	digestAndEpoch, err := ocrContract.LatestConfigDigestAndEpoch(nil)
	require.NoError(t, err)
	assert.Equal(t, digestAndEpoch.Epoch, epoch)
}
