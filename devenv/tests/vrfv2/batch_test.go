package vrfv2

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2 "github.com/smartcontractkit/chainlink/devenv/products/vrfv2"
)

func TestVRFv2BatchFulfillmentEnabledDisabled(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr, "failed to save container logs")
	})

	outputFile := "../../env-vrfv2-out.toml"
	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load devenv env-out from %s", outputFile)

	cfg, err := products.LoadOutput[productvrfv2.Configurator](outputFile)
	require.NoError(t, err, "failed to load vrfv2 product config from env-out")
	c := cfg.Config[0]

	keyHash := mustKeyHash(c)
	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err, "failed to parse chain ID from env-out")
	ctx := t.Context()
	chainClient, err := products.InitSeth(in.Blockchains[0].Out.Nodes[0].ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	coord, err := contracts.LoadVRFCoordinatorV2(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load VRF coordinator v2")
	batchCoordAddr := c.DeployedContracts.BatchCoordinator
	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to connect to Chainlink nodes")

	callbackGas := c.BatchCallbackGasLimit
	if callbackGas == 0 {
		callbackGas = 500_000
	}
	batchBudget := c.BatchTxGasBudget
	if batchBudget == 0 {
		batchBudget = c.MaxGasLimitCoordinator + 400_000
	}
	randRequestCountU32 := (batchBudget / callbackGas) - 1
	require.Greater(t, randRequestCountU32, uint32(1), "batch test needs randRequestCount > 1 (check batch budget vs callback gas)")
	require.LessOrEqual(t, randRequestCountU32, uint32(^uint16(0)), "randRequestCount must fit in uint16")
	randRequestCount := uint16(randRequestCountU32) //nolint:gosec // bounded by require.LessOrEqual to max uint16

	fulfillTimeout := parseFulfillTimeout(c.RandomWordsFulfilledEventTimeout)

	pollPeriod, _ := time.ParseDuration(c.VRFJobPollPeriod)
	if pollPeriod <= 0 {
		pollPeriod = time.Second
	}
	requestTimeout, _ := time.ParseDuration(c.VRFJobRequestTimeout)
	if requestTimeout <= 0 {
		requestTimeout = 24 * time.Hour
	}

	buildPipeline := func() (string, error) {
		ps := &productvrfv2.TxPipelineSpec{
			Address:               coord.Address(),
			EstimateGasMultiplier: c.VRFJobEstimateGasMultiplier,
			FromAddress:           c.VRFKeyData.TxKeyAddresses[0],
		}
		if c.VRFJobSimulationBlock != "" {
			sb := c.VRFJobSimulationBlock
			ps.SimulationBlock = &sb
		}
		return ps.String()
	}

	currentJobID := c.VRFKeyData.VRFJobID
	switchJob := func(t *testing.T, batchOn bool) {
		t.Helper()
		if currentJobID != "" {
			resp, dErr := cl[0].DeleteJob(currentJobID)
			require.NoError(t, dErr, "failed to delete existing VRF job before switch")
			require.Equal(t, http.StatusNoContent, resp.StatusCode, "delete job should return 204, got %d", resp.StatusCode)
		}
		obs, oErr := buildPipeline()
		require.NoError(t, oErr, "failed to build observation pipeline spec")
		gasMult := c.VRFJobBatchFulfillmentGasMultiplier
		if gasMult == 0 {
			gasMult = 1.1
		}
		prefix := "enabled"
		if !batchOn {
			prefix = "disabled"
		}
		spec := &productvrfv2.JobSpec{
			Name:                          fmt.Sprintf("vrf-v2-batch-%s-%s", prefix, uuid.NewString()),
			CoordinatorAddress:            coord.Address(),
			BatchCoordinatorAddress:       batchCoordAddr,
			PublicKey:                     c.VRFKeyData.PubKeyCompressed,
			ExternalJobID:                 uuid.New().String(),
			ObservationSource:             obs,
			MinIncomingConfirmations:      int(c.MinimumConfirmations),
			FromAddresses:                 c.VRFKeyData.TxKeyAddresses,
			EVMChainID:                    in.Blockchains[0].Out.ChainID,
			ForwardingAllowed:             c.VRFJobForwardingAllowed,
			BatchFulfillmentEnabled:       batchOn,
			BatchFulfillmentGasMultiplier: gasMult,
			BackOffInitialDelay:           15 * time.Second,
			BackOffMaxDelay:               5 * time.Minute,
			PollPeriod:                    pollPeriod,
			RequestTimeout:                requestTimeout,
		}
		job, jErr := cl[0].MustCreateJob(spec)
		require.NoError(t, jErr, "failed to create VRF job (batch=%v)", batchOn)
		currentJobID = job.Data.ID
	}
	t.Cleanup(func() {
		if currentJobID != "" {
			_, _ = cl[0].DeleteJob(currentJobID)
		}
	})

	t.Run("Batch Fulfillment Enabled", func(t *testing.T) {
		require.NoError(t, deleteAllJobs(cl[0]), "failed to clear jobs before batch-enabled run")
		currentJobID = ""
		switchJob(t, true)

		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, callbackGas, c.NumberOfWords,
			randRequestCount, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err, "error requesting randomness and waiting for fulfillment (batch on)")

		var wg sync.WaitGroup
		wg.Add(1)
		_, _, err = waitForRequestCountEqualToFulfillmentCount(ctx, consumers[0], 2*time.Minute, &wg)
		require.NoError(t, err, "consumer request/fulfillment counts did not converge (batch on)")
		wg.Wait()

		txs, _, err := cl[0].ReadTransactions()
		require.NoError(t, err, "error reading node transactions")
		var batchTxs []string
		for _, tx := range txs.Data {
			if stringsEqualFoldAddr(tx.Attributes.To, batchCoordAddr) {
				batchTxs = append(batchTxs, tx.Attributes.Hash)
			}
		}
		require.Len(t, batchTxs, 1, "expected exactly one tx from node to batch coordinator")

		fulfillTx, _, err := chainClient.Client.TransactionByHash(ctx, fulfilled.Raw.TxHash)
		require.NoError(t, err, "failed to load fulfillment transaction")
		require.NotNil(t, fulfillTx.To(), "fulfillment tx must have a To address")
		require.True(t, stringsEqualFoldAddr(fulfillTx.To().Hex(), batchCoordAddr),
			"fulfillment should go to batch coordinator %s, got %s", batchCoordAddr, fulfillTx.To().Hex())

		receipt, err := chainClient.Client.TransactionReceipt(ctx, fulfillTx.Hash())
		require.NoError(t, err, "failed to load fulfillment receipt")
		logs, err := contracts.ParseRandomWordsFulfilledLogs(coord, receipt.Logs)
		require.NoError(t, err, "failed to parse RandomWordsFulfilled logs from receipt")
		require.Len(t, logs, int(randRequestCount), "expected %d RandomWordsFulfilled logs in batch receipt", randRequestCount)
	})

	t.Run("Batch Fulfillment Disabled", func(t *testing.T) {
		require.NoError(t, deleteAllJobs(cl[0]), "failed to clear jobs before batch-disabled run")
		currentJobID = ""
		switchJob(t, false)

		consumers, subIDs, err := deployConsumersAndFundSubs(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, 1, 1)
		require.NoError(t, err, "error setting up new consumers and subs")
		subID := subIDs[0]

		_, fulfilled, err := requestRandomnessAndWaitForFulfillment(ctx, consumers[0], coord, keyHash, subID,
			c.MinimumConfirmations, callbackGas, c.NumberOfWords,
			randRequestCount, c.RandomnessRequestCountPerRequestDeviation,
			fulfillTimeout, 0)
		require.NoError(t, err, "error requesting randomness and waiting for fulfillment (batch off)")

		var wg sync.WaitGroup
		wg.Add(1)
		_, _, err = waitForRequestCountEqualToFulfillmentCount(ctx, consumers[0], 2*time.Minute, &wg)
		require.NoError(t, err, "consumer request/fulfillment counts did not converge (batch off)")
		wg.Wait()

		fulfillTx, _, err := chainClient.Client.TransactionByHash(ctx, fulfilled.Raw.TxHash)
		require.NoError(t, err, "failed to load fulfillment transaction")
		require.NotNil(t, fulfillTx.To(), "fulfillment tx must have a To address")
		require.True(t, stringsEqualFoldAddr(fulfillTx.To().Hex(), coord.Address()),
			"fulfillment should go to coordinator %s, got %s", coord.Address(), fulfillTx.To().Hex())

		txs, _, err := cl[0].ReadTransactions()
		require.NoError(t, err, "error reading node transactions")
		var coordTxs int
		for _, tx := range txs.Data {
			if stringsEqualFoldAddr(tx.Attributes.To, coord.Address()) {
				coordTxs++
			}
		}
		require.Equal(t, int(randRequestCount), coordTxs,
			"expected %d txs from node to coordinator (one per word path), got %d", randRequestCount, coordTxs)
	})
}

func stringsEqualFoldAddr(a, b string) bool {
	return common.HexToAddress(a) == common.HexToAddress(b)
}
