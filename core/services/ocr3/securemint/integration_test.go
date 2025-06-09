package llo_test

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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/data-feeds/generated/data_feeds_cache"
	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	evmtestutils "github.com/smartcontractkit/chainlink-evm/pkg/testutils"
	evmtypes "github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/csakey"
	"github.com/smartcontractkit/chainlink/v2/core/services/ocr2/testhelpers"
	"github.com/smartcontractkit/chainlink/v2/core/services/relay/evm"
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

// TODO(gg) see also:
// https://github.com/smartcontractkit/mercury-pipeline/blob/9f0bc5d457d57d5807122446cb936306ecf1b263/e2e_tests/mercuryhelpers/helpers.go#L308 for example of onchain config

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

	steve, backend := setupBlockchain(t)
	fromBlock, err := backend.Client().BlockNumber(testutils.Context(t))
	require.NoError(t, err)
	t.Logf("Starting from block: %d", fromBlock)

	// Setup bootstrap
	bootstrapCSAKey := csakey.MustNewV2XXXTestingOnly(big.NewInt(salt - 1))
	bootstrapNodePort := freeport.GetOne(t)
	appBootstrap, bootstrapPeerID, _, bootstrapKb, _ := setupNode(t, bootstrapNodePort, "bootstrap_securemint", backend, bootstrapCSAKey, nil)
	bootstrapNode := Node{App: appBootstrap, KeyBundle: bootstrapKb}

	p2pV2Bootstrappers := []commontypes.BootstrapperLocator{
		// Supply the bootstrap IP and port as a V2 peer address
		{PeerID: bootstrapPeerID, Addrs: []string{fmt.Sprintf("127.0.0.1:%d", bootstrapNodePort)}},
	}

	// Setup oracle nodes
	oracles, nodes := setupNodes(t, nNodes, backend, clientCSAKeys, func(c *chainlink.Config) {
		// TODO(gg): something like this + extra config
		// c.Feature.SecureMint.Enabled = true

		// inform node about bootstrap node
		c.P2P.V2.DefaultBootstrappers = &p2pV2Bootstrappers
	})

	// pluginConfig := fmt.Sprintf(`servers = { "%s" = "%x" }
	// donID = %d
	// channelDefinitionsContractAddress = "0x%x"
	// channelDefinitionsContractFromBlock = %d`, serverURL, serverPubKey, donID, configStoreAddress, fromBlock)
	// addOCRJobsEVMPremiumLegacy(t, streams, serverPubKey, serverURL, legacyVerifierAddr, bootstrapPeerID, bootstrapNodePort, nodes, configStoreAddress, clientPubKeys, pluginConfig, relayType, relayConfig)

	allowedSenders := make([]common.Address, len(nodes))
	for i, node := range nodes {
		keys, err := node.App.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
		require.NoError(t, err)
		allowedSenders[i] = keys[0].Address // assuming the first key is the transmitter
	}

	aggregatorAddress := setSecureMintOnchainConfigUsingAggregator(t, steve, backend, nodes, oracles)

	t.Logf("Creating bootstrap job with aggregator address: %s", aggregatorAddress.Hex())
	bootstrapJob := createSecureMintBootstrapJob(t, bootstrapNode, aggregatorAddress, testutils.SimulatedChainID.String(), fmt.Sprintf("%d", fromBlock))
	t.Logf("Created bootstrap job: %s with id %d", bootstrapJob.Name.ValueOrZero(), bootstrapJob.ID)

	// TODO(gg): enable this for writing step
	// TODO(gg): deduplicate
	// feedIDBytes := [16]byte{}
	// copy(feedIDBytes[:], common.FromHex("0xA1B2C3D4E5F600010203040506070809"))

	// dfCacheAddress, dfCacheContract := setupDataFeedsCacheContract(t, steve, backend, allowedSenders, steve.From.Hex(), "securemint")
	// t.Logf("Deployed and configured DataFeedsCache contract at: %s", dfCacheAddress.Hex())
	// desc, err := dfCacheContract.GetDescription(&bind.CallOpts{}, feedIDBytes)
	// require.NoError(t, err)
	// t.Logf("DataFeedsCache description: %s", desc)

	jobIDs := addSecureMintOCRJobs(t, nodes, aggregatorAddress)

	t.Logf("jobIDs: %v", jobIDs)
	validateJobsRunningSuccessfully(t, nodes, jobIDs)

	// wait for a minute for the jobs to run and collect data
	// time.Sleep(1 * time.Minute)
}

func setupNodes(t *testing.T, nNodes int, backend evmtypes.Backend, clientCSAKeys []csakey.KeyV2, f func(*chainlink.Config)) (oracles []confighelper.OracleIdentityExtra, nodes []Node) {
	ports := freeport.GetN(t, nNodes)
	for i := range nNodes {
		app, peerID, transmitter, kb, observedLogs := setupNode(t, ports[i], fmt.Sprintf("oracle_securemint_%d", i), backend, clientCSAKeys[i], f)

		nodes = append(nodes, Node{
			App:          app,
			ClientPubKey: transmitter,
			KeyBundle:    kb,
			ObservedLogs: observedLogs,
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

func validateJobsRunningSuccessfully(t *testing.T, nodes []Node, jobIDs map[int]int32) {

	// 1. Assert no job spec errors
	for i, node := range nodes {
		jobs, _, err := node.App.JobORM().FindJobs(testutils.Context(t), 0, 1000)
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

	runs, err := nodes[0].App.PipelineORM().GetAllRuns(testutils.Context(t))
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
			// t.Logf("finding pipeline runs for job %d on node %d", jobIDs[i], i)
			// completedRuns, err := node.App.JobORM().FindPipelineRunIDsByJobID(testutils.Context(t), jobIDs[i], 0, 10)
			// if !assert.NoError(t, err) {
			// 	t.Logf("assert error finding pipeline runs for job %d: %v", jobIDs[i], err)
			// 	return
			// }
			// t.Logf("found pipeline runs for job %d on node %d: %v", jobIDs[i], i, completedRuns)

			// Want at least 2 runs so we see all the metadata.

			// TODO(gg): fix this, the pipeline completes now
			/**
				cltest.go:969: Found pipeline run 9 with status completed on node 3 for job 1 with task runs: []pipeline.TaskRun(nil)
			    cltest.go:969: Found pipeline run 8 with status completed on node 3 for job 1 with task runs: []pipeline.TaskRun(nil)
			    cltest.go:969: Found pipeline run 7 with status completed on node 3 for job 1 with task runs: []pipeline.TaskRun(nil)
			*/

			pr := cltest.WaitForPipelineComplete(t, i, jobIDs[i], 1, 0, node.App.JobORM(), 30*time.Second, 1*time.Second)
			outputs, err := pr[0].Outputs.MarshalJSON()
			if !assert.NoError(t, err) {
				t.Logf("assert error marshalling outputs for job %d: %v", jobIDs[i], err)
				return
			}
			t.Logf("Pipeline itself is %+v", pr[0])
			t.Logf("Pipeline run outputs are %s", string(outputs))

			// assert.Equalf(t, []byte(fmt.Sprintf("[\"%d\"]", 1000*i)), jb, "pr[0] %+v pr[1] %+v", pr[0], pr[1], "assert error: something unexpected happened")
		}()
	}
	t.Logf("waiting for pipeline runs to complete")
	wg.Wait()
	t.Logf("All pipeline runs completed successfully")
}

// TODO(gg): to set config on DF Cache contract
// func setSecureMintOnchainConfigOnDFCacheContract(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []Node, oracles []confighelper.OracleIdentityExtra, dfCacheAddress common.Address, dfCacheContract *data_feeds_cache.DataFeedsCache) [32]byte {

// 	_, err = dfCacheContract.SetConfig(steve, signerAddresses, transmitterAddresses, f, outOnchainConfig, offchainConfigVersion, offchainConfig)
// 	if err != nil {
// 		errString, err := rPCErrorFromError(err)
// 		require.NoError(t, err)

// 		t.Fatalf("Failed to configure contract: %s", errString)
// 	}

// 	// libocr requires a few confirmations to accept the config
// 	backend.Commit()
// 	backend.Commit()
// 	backend.Commit()
// 	backend.Commit()

// 	l, err := dfCacheContract.LatestConfigDigestAndEpoch(&bind.CallOpts{})
// 	require.NoError(t, err)

// 	return l.ConfigDigest
// }

func setSecureMintOnchainConfigUsingAggregator(t *testing.T, steve *bind.TransactOpts, backend evmtypes.Backend, nodes []Node, oracles []confighelper.OracleIdentityExtra) common.Address {

	// 1. Deploy aggregator contract

	// these min and max answers are not used by the secure mint oracle but they're needed for validation in aggregator.setConfig()
	// TODO(gg): maybe these could be 0 and max int?
	minAnswer, maxAnswer := new(big.Int), new(big.Int)
	minAnswer.Exp(big.NewInt(-2), big.NewInt(191), nil)
	maxAnswer.Exp(big.NewInt(2), big.NewInt(191), nil)
	maxAnswer.Sub(maxAnswer, big.NewInt(1))

	aggregatorAddress, _, aggregatorContract, err := ocr2aggregator.DeployOCR2Aggregator(
		steve,
		backend.Client(),
		common.Address{},   // _link common.Address,
		minAnswer,          // -2**191
		maxAnswer,          // 2**191 - 1
		common.Address{},   // accessAddress
		common.Address{},   // accessAddress
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
	onchainConfig, err := testhelpers.GenerateDefaultOCR2OnchainConfig(minAnswer, maxAnswer) // TODO(gg): this uses the median codec, not sure if this is correct
	require.NoError(t, err)

	smPluginConfig := por.PorOffchainConfig{MaxChains: 5} // TODO(gg): set config values
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
		0,                    // maxDurationQuery,
		250*time.Millisecond, // maxDurationObservation,
		0,                    // maxDurationShouldAcceptAttestedReport,
		0,                    // maxDurationShouldTransmitAcceptedReport,
		int(fNodes),          // f,
		onchainConfig,        // onchainConfig (binary blob containing configuration passed through to the ReportingPlugin and also available to the contract. Unlike ReportingPluginConfig which is only available offchain.)
	)
	require.NoError(t, err)

	// 3. Set config on the contract
	signerAddresses, err := evm.OnchainPublicKeyToAddress(signers)
	require.NoError(t, err)

	transmitterAddresses := make([]common.Address, len(nodes))
	for i := range nodes {
		keys, err := nodes[i].App.GetKeyStore().Eth().EnabledKeysForChain(testutils.Context(t), testutils.SimulatedChainID)
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
