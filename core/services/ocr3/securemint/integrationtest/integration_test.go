package integrationtest

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	datastreamsllo "github.com/smartcontractkit/chainlink-data-streams/llo"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-evm/gethwrappers/llo-feeds/generated/configurator"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr3/securemint"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm/llo"
	"github.com/smartcontractkit/freeport"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/gethwrappers2/ocr2aggregator"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/confighelper"
	"github.com/smartcontractkit/libocr/offchainreporting2plus/ocr3confighelper"
	ocr2types "github.com/smartcontractkit/libocr/offchainreporting2plus/types"
	"github.com/smartcontractkit/por_mock_ocr3plugin/por"
)

var (
	fNodes = uint8(1)
	nNodes = 4 // number of nodes (not including bootstrap)
)

// TestIntegration_SecureMint_happy_path tests runs a small DON which runs the secure mint plugin
// and verifies that it can successfully create reports.
//
// Inspired by:
// * core/internal/features/ocr2/features_ocr2_test.go
// * core/services/ocr2/plugins/ocr2keeper/integration_21_test.go
func TestIntegration_SecureMint_happy_path(t *testing.T) {
	const salt = 100

	clientCSAKeys := make([]csakey.KeyV2, nNodes)
	clientPubKeys := make([]ed25519.PublicKey, nNodes)
	for i := range nNodes {
		k := big.NewInt(int64(salt + i))
		key := csakey.MustNewV2XXXTestingOnly(k)
		clientCSAKeys[i] = key
		clientPubKeys[i] = key.PublicKey
	}

	t.Logf("clientPubKeys: %v", clientPubKeys)

	steve, backend := setupBlockchain(t)
	fromBlock, err := backend.Client().BlockNumber(testutils.Context(t))
	require.NoError(t, err)
	t.Logf("Starting from block: %d", fromBlock)

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_securemint", backend, bootstrapCSAKey, nil)
	bootstrapNode := node{app: appBootstrap, keyBundle: bootstrapKb}

	p2pV2Bootstrappers := []commontypes.BootstrapperLocator{
		// Supply the bootstrap IP and port as a V2 peer address
		{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
	}

	// Setup oracle nodes
	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		// inform node about bootstrap node
		c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
	})
	for i, node := range nodes {
		t.Logf("node %d clientPubKey: %x", i, node.clientPubKey)
	}

	allowedSenders := make([]common.Address, len(nodes))
	for i, node := range nodes {
		keys, err := node.app.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		allowedSenders[i] = keys[0].Address // assuming the first key is the transmitter
	}

	// aggregatorAddress := setSecureMintOnchainConfigUsingAggregator(t, steve, backend, nodes, oracles)
	_, configuratorAddress := setSecureMintOnchainConfigUsingOCR3Configurator(t, steve, backend, nodes, oracles)

	t.Logf("Creating bootstrap job with configurator address: %s", configuratorAddress.Hex())
	bootstrapJob := createSecureMintBootstrapJob(t, bootstrapNode, configuratorAddress, testutils.SimulatedChainID.String(), fmt.Sprintf("%d", fromBlock))
	t.Logf("Created bootstrap job: %s with id %d", bootstrapJob.Name.ValueOrZero(), bootstrapJob.ID)

	jobIDs := addSecureMintOCRJobs(t, nodes, configuratorAddress)

	t.Logf("jobIDs: %v", jobIDs)
	validateJobsRunningSuccessfully(t, nodes, jobIDs)

}

func setupBlockchain(t *testing.T) (
	*bind.TransactOpts,
	evmtypes.Backend,
) {
	steve := evmtestutils.MustNewSimTransactor(t) // config contract deployer and owner
	genesisData := gethtypes.GenesisAlloc{steve.From: {Balance: assets.Ether(1000).ToInt()}}
	backend := cltest.NewSimulatedBackend(t, genesisData, ethconfig.Defaults.Miner.GasCeil)
	backend.Commit()
	backend.Commit() // ensure starting block number at least 1

	return steve, backend
}

func setupNodes(t *testing.T, nNodes int, backend evmtypes.Backend, clientCSAKeys []csakey.KeyV2, f func(*chainlink.Config)) (oracles []confighelper.OracleIdentityExtra, nodes []node) {
	ports := freeport.GetN(t, nNodes)
	for i := range nNodes {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_securemint_%d", i), backend, clientCSAKeys[i], f)

		nodes = append(nodes, node{
			app:          app,
			clientPubKey: transmitter,
			keyBundle:    kb,
			observedLogs: observedLogs,
		})
		offchainPublicKey, err := hex.DecodeString(strings.TrimPrefix(kb.OnChainPublicKey(), "0x"))
		require.NoError(t, err)
		oracles = append(oracles, confighelper.OracleIdentityExtra{
			OracleIdentity: confighelper.OracleIdentity{
				OnchainPublicKey:  offchainPublicKey,
				TransmitAccount:   ocr2types.Account(hex.EncodeToString(transmitter[:])),
				OffchainPublicKey: kb.OffchainPublicKey(),
				PeerID:            peerID,
			},
			ConfigEncryptionPublicKey: kb.ConfigEncryptionPublicKey(),
		})
	}
	return
}

func validateJobsRunningSuccessfully(t *testing.T, nodes []node, jobIDs map[int]int32) {

	// 1. Assert no job spec errors
	for i, node := range nodes {
		jobs, _, err := node.app.JobORM().FindJobs(testutils.Context(t), 0, 1000)
		require.NoErrorf(t, err, "assert error finding jobs for node %d", i)
		t.Logf("%d jobs found for node %d", len(jobs), i)
		for _, j := range jobs {
			t.Logf("job %d on node %d oracle spec: %#v", j.ID, i, j.OCR2OracleSpec)
			t.Logf("job %d on node %d pipeline spec: %#v", j.ID, i, j.PipelineSpec)
		}
		// No spec errors
		for _, j := range jobs {
			ignore := 0
			for _, jse := range j.JobSpecErrors {
				// Non-fatal timing related error, ignore for testing.
				if strings.Contains(jse.Description, "leader's phase conflicts tGrace timeout") {
					ignore++
				} else {
					t.Errorf("assert error: job spec error on node %d: %v", i, jse)
				}
			}
			require.Lenf(t, j.JobSpecErrors, ignore, "assert error: job spec errors on node %d", i)
		}
	}

	t.Logf("No job spec errors identified for any node")

	// time.Sleep(30 * time.Second) // wait for jobs to run

	runs, err := nodes[0].app.PipelineORM().GetAllRuns(testutils.Context(t))
	require.NoError(t, err, "assert error getting all runs")
	t.Logf("Found %d runs", len(runs))
	for _, run := range runs {
		t.Logf("Run ID: %d, Job ID: %d, Status: %s", run.ID, run.JobID, run.Status())
	}

	// 2. Assert that all the Secure Mint jobs get a run with valid values eventually
	var wg sync.WaitGroup
	for i, node := range nodes {
		wg.Add(1)
		go func() {
			defer wg.Done()

			pr := cltest.WaitForPipelineComplete(t, i, jobIDs[i], 1, 0, node.app.JobORM(), 30*time.Second, 1*time.Second)
			outputs, err := pr[0].Outputs.MarshalJSON()
			if !assert.NoError(t, err) {
				t.Logf("assert error marshalling outputs for job %d: %v", jobIDs[i], err)
				return
			}
			t.Logf("Pipeline itself is %+v", pr[0])
			t.Logf("Pipeline run outputs are %s", string(outputs))
		}()
	}
	t.Logf("waiting for pipeline runs to complete")
	wg.Wait()
	t.Logf("All pipeline runs completed successfully")

	// 3. Check that transmissions work
	expectedNumTransmissions := int32(4)
	gomega.NewWithT(t).Eventually(func() bool {
		numTransmissions := securemint.StubTransmissionCounter.Load()
		t.Logf("Number of (stub) report transmissions: %d", numTransmissions)
		return numTransmissions >= expectedNumTransmissions
	}, 30*time.Second, 1*time.Second).Should(
		gomega.BeTrue(),
		fmt.Sprintf("expected at least %d reports transmitted, but got less", expectedNumTransmissions),
	)
}

func setSecureMintOnchainConfigUsingOCR3Configurator(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []node, oracles []confighelper.OracleIdentityExtra) (*configurator.Configurator, common.Address) {

	// 1. Deploy configurator contract
	configuratorAddress, _, configurator, err := configurator.DeployConfigurator(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	// Ensure we have finality depth worth of blocks to start.
	for range 5 {
		backend.Commit()
	}
	t.Logf("Deployed OCR3Configurator contract at: %s", configuratorAddress.Hex())

	smPluginConfig := por.PorOffchainConfig{MaxChains: 5}
	smPluginConfigBytes, err := smPluginConfig.Serialize()
	require.NoError(t, err)

	onchainConfig, err := (&datastreamsllo.EVMOnchainConfigCodec{}).Encode(datastreamsllo.OnchainConfig{
		Version:                 1,
		PredecessorConfigDigest: nil,
	})
	require.NoError(t, err)

	signers, transmitters, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress,
		20*time.Second,       // deltaResend,
		400*time.Millisecond, // deltaInitial,
		500*time.Millisecond, // deltaRound,
		250*time.Millisecond, // deltaGrace,
		300*time.Millisecond, // deltaCertifiedCommitRequest,
		1*time.Minute,        // deltaStage,
		100,                  // rMax,
		[]int{len(oracles)},  // s,
		oracles,              // oracles,
		smPluginConfigBytes,  // reportingPluginConfig,
		nil,                  // maxDurationInitialization,
		250*time.Millisecond, // maxDurationQuery,
		1*time.Second,        // maxDurationObservation,
		1*time.Second,        // maxDurationShouldAcceptAttestedReport,
		1*time.Second,        // maxDurationShouldTransmitAcceptedReport,
		int(fNodes),          // f,
		onchainConfig,        // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
	)
	require.NoError(t, err)

	// clientPubKeys := make([]ed25519.PublicKey, nNodes)
	// for i := 0; i < nNodes; i++ {
	// 	k := big.NewInt(int64(salt + i))
	// 	key := csakey.MustNewV2XXXTestingOnly(k)
	// 	clientCSAKeys[i] = key
	// 	clientPubKeys[i] = key.PublicKey
	// }

	// func (r wsrpcRequest) TransmitterID() ocr2types.Account {
	// 	return ocr2types.Account(fmt.Sprintf("%x", r.pk))
	// }

	// TransmitAccount:   ocr2types.Account(hex.EncodeToString(transmitter[:])),

	// 3. Set config on the contract
	var signerKeys [][]byte
	for _, signer := range signers {
		signerKeys = append(signerKeys, signer)
	}

	transmitterAddresses := make([]common.Address, len(nodes))
	for i := range nodes {
		keys, err := nodes[i].app.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		transmitterAddresses[i] = keys[0].Address // assuming the first key is the transmitter
	}
	t.Logf("transmitterAddresses: %v", transmitterAddresses)

	transmitterAddrs := make([][32]byte, len(transmitterAddresses))
	for i := range transmitterAddresses {
		copy(transmitterAddrs[i][:], transmitterAddresses[i][:])
	}
	t.Logf("transmitterAddrs: %v", transmitterAddrs)

	offchainTransmitters := make([][32]byte, len(transmitters))
	for i := range transmitters {
		copy(offchainTransmitters[i][:], transmitters[i][:])
	}
	t.Logf("offchainTransmitters: %v", offchainTransmitters)

	// transmitters should be the nodes' csa keys
	offchainTransmitters2 := make([][32]byte, len(nodes))
	for i := range nodes {
		copy(offchainTransmitters2[i][:], nodes[i].clientPubKey[:])
	}
	t.Logf("offchainTransmitters2: %v", offchainTransmitters2)

	offchainTransmitters3 := make([][32]byte, nNodes)
	for i := 0; i < nNodes; i++ {
		offchainTransmitters3[i] = nodes[i].clientPubKey // use csa keys as transmitters
	}
	t.Logf("offchainTransmitters3: %v", offchainTransmitters3)

	configID := [32]byte{}
	copy(configID[:], common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000001"))

	x := func(donID uint32) [32]byte {
		var b [32]byte
		copy(b[:], common.LeftPadBytes(big.NewInt(int64(donID)).Bytes(), 32))
		return b
	}
	donIDBytes32 := x(1)
	t.Logf("donIDBytes32: %x", donIDBytes32)
	t.Logf("configID: %x", configID)

	_, err = configurator.SetProductionConfig(steve, configID, signerKeys, offchainTransmitters3, f, outOnchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		t.Logf("Error: %s", err)
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to configure contract: %s %s", errString, err)
	}

	// make sure config is finalized
	for range 5 {
		backend.Commit()
	}

	var topic common.Hash
	topic = llo.ProductionConfigSet

	logs, err := backend.Client().FilterLogs(testutils.Context(t), ethereum.FilterQuery{Addresses: []common.Address{configuratorAddress}, Topics: [][]common.Hash{[]common.Hash{topic, configID}}})
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(logs), 1)
	cfg, err := llo.DecodeProductionConfigSetLog(logs[len(logs)-1].Data)
	require.NoError(t, err)

	t.Logf("Configurator config digest: 0x%x", cfg.ConfigDigest)

	return configurator, configuratorAddress
}

func setSecureMintOnchainConfigUsingAggregator(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []node, oracles []confighelper.OracleIdentityExtra) common.Address {

	// 1. Deploy aggregator contract

	// these min and max answers are not used by the secure mint oracle but they're needed for validation in aggregator.setConfig()
	minAnswer := big.NewInt(0)
	maxAnswer := big.NewInt(999999)
	aggregatorAddress, _, aggregatorContract, err := ocr2aggregator.DeployOCR2Aggregator(
		steve,
		backend.Client(),
		common.Address{}, // LINK address
		minAnswer,
		maxAnswer,
		common.Address{},   // billingAccessController
		common.Address{},   // requesterAccessController
		9,                  // decimals
		"secure mint test", // description
	)
	if err != nil {
		rPCError, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to deploy OCR2Aggregator contract: %s", rPCError)
	}
	// Ensure we have finality depth worth of blocks to start.
	for range 20 {
		backend.Commit()
	}
	t.Logf("Deployed OCR2Aggregator contract at: %s", aggregatorAddress.Hex())

	// 2. Create config
	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer)
	require.NoError(t, err)

	smPluginConfig := por.PorOffchainConfig{MaxChains: 5}
	smPluginConfigBytes, err := smPluginConfig.Serialize()
	require.NoError(t, err)

	signers, _, f, outOnchainConfig, offchainConfigVersion, offchainConfig, err := ocr3confighelper.ContractSetConfigArgsForTests(
		2*time.Second,        // deltaProgress,
		20*time.Second,       // deltaResend,
		400*time.Millisecond, // deltaInitial,
		500*time.Millisecond, // deltaRound,
		250*time.Millisecond, // deltaGrace,
		300*time.Millisecond, // deltaCertifiedCommitRequest,
		1*time.Minute,        // deltaStage,
		100,                  // rMax,
		[]int{len(oracles)},  // s,
		oracles,              // oracles,
		smPluginConfigBytes,  // reportingPluginConfig,
		nil,                  // maxDurationInitialization,
		250*time.Millisecond, // maxDurationQuery,
		1*time.Second,        // maxDurationObservation,
		1*time.Second,        // maxDurationShouldAcceptAttestedReport,
		1*time.Second,        // maxDurationShouldTransmitAcceptedReport,
		int(fNodes),          // f,
		onchainConfig,        // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
	)
	require.NoError(t, err)

	// 3. Set config on the contract
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)

	transmitterAddresses := make([]common.Address, len(nodes))
	for i := range nodes {
		keys, err := nodes[i].app.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		transmitterAddresses[i] = keys[0].Address // assuming the first key is the transmitter
	}

	_, err = aggregatorContract.SetConfig(steve, signerAddresses, transmitterAddresses, f, outOnchainConfig, offchainConfigVersion, offchainConfig)
	if err != nil {
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to configure contract: %s", errString)
	}

	// make sure config is finalized
	for range 20 {
		backend.Commit()
	}

	aggregatorConfigDigest, err := aggregatorContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
	if err != nil {
		rPCError, err := rPCErrorFromError(err)
		require.NoError(t, err)
		t.Fatalf("Failed to get latest config digest: %s", rPCError)
	}
	t.Logf("Aggregator config digest: 0x%x", aggregatorConfigDigest.ConfigDigest)

	return aggregatorAddress
}

func rPCErrorFromError(txError error) (string, error) {
	errBytes, err := json.Marshal(txError)
	if err != nil {
		return "", err
	}
	var callErr struct {
		Code    int
		Data    string `json:"data"`
		Message string `json:"message"`
	}
	err = json.Unmarshal(errBytes, &callErr)
	if err != nil {
		return "", err
	}
	// If the error data is blank
	if len(callErr.Data) == 0 {
		return callErr.Data, nil
	}
	// Some nodes prepend "Reverted " and we also remove the 0x
	trimmed := strings.TrimPrefix(callErr.Data, "Reverted ")[2:]
	data, err := hex.DecodeString(trimmed)
	if err != nil {
		return "", err
	}
	revert, err := abi.UnpackRevert(data)
	// If we can't decode the revert reason, return the raw data
	if err != nil {
		return callErr.Data, nil
	}
	return revert, nil
}

func setupDataFeedsCacheContract(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, allowedSenders []common.Address, workflowOwner, workflowName string) (
	common.Address, *data_feeds_cache.DataFeedsCache) {

	addr, _, dataFeedsCache, err := data_feeds_cache.DeployDataFeedsCache(steve, backend.Client())
	require.NoError(t, err)
	backend.Commit()

	var nameBytes [10]byte
	copy(nameBytes[:], workflowName)

	ownerAddr := common.HexToAddress(workflowOwner)

	_, err = dataFeedsCache.SetFeedAdmin(steve, ownerAddr, true)
	require.NoError(t, err)

	backend.Commit()

	metadatas := make([]data_feeds_cache.DataFeedsCacheWorkflowMetadata, len(allowedSenders))
	for i, sender := range allowedSenders {
		metadatas[i] =
			data_feeds_cache.DataFeedsCacheWorkflowMetadata{
				AllowedSender:        sender,
				AllowedWorkflowOwner: ownerAddr,
				AllowedWorkflowName:  nameBytes,
			}
	}

	feedIDBytes := [16]byte{}
	copy(feedIDBytes[:], common.FromHex("0xA1B2C3D4E5F600010203040506070809"))

	_, err = dataFeedsCache.SetDecimalFeedConfigs(steve, [][16]byte{feedIDBytes}, []string{"securemint"}, metadatas)
	if err != nil {
		errString, err := rPCErrorFromError(err)
		require.NoError(t, err)

		t.Fatalf("Failed to configure contract: %s", errString)
	}

	backend.Commit()

	return addr, dataFeedsCache
}
