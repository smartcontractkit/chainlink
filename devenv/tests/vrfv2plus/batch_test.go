package vrfv2plus

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-evm/gethwrappers/generated/vrf_coordinator_v2_5"
	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2plus "github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFv2PlusBatchFulfillment(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrf2plus-out.toml"

	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load env-out.toml")

	cfg, err := products.LoadOutput[productvrfv2plus.Configurator](outputFile)
	require.NoError(t, err, "failed to load VRFv2Plus config from env-out.toml")
	require.NotEmpty(t, cfg.Config, "vrfv2_plus config must not be empty")

	c := cfg.Config[0]
	require.NotZero(t, c.MaxGasLimitCoordinator,
		"max_gas_limit_coordinator is zero in env output; recreate the vrfv2plus batch environment")

	require.NotEmpty(t, c.DeployedContracts.Coordinator, "coordinator address must not be empty")
	require.NotEmpty(t, c.DeployedContracts.BatchCoordinator, "batch coordinator address must not be empty")
	require.NotEmpty(t, c.VRFKeyData.KeyHash, "key hash must not be empty")
	require.NotEmpty(t, c.VRFKeyData.VRFJobID, "VRF job ID must not be empty")

	keyHashBytes := common.HexToHash(c.VRFKeyData.KeyHash)
	var keyHash [32]byte
	copy(keyHash[:], keyHashBytes[:])

	chainID, err := strconv.ParseUint(in.Blockchains[0].Out.ChainID, 10, 64)
	require.NoError(t, err)

	bcNode := in.Blockchains[0].Out.Nodes[0]
	ctx := t.Context()
	chainClient, err := products.InitSeth(bcNode.ExternalWSUrl, []string{products.NetworkPrivateKey()}, &chainID)
	require.NoError(t, err, "failed to init Seth client")

	coord, err := contracts.LoadVRFCoordinatorV2_5(chainClient, c.DeployedContracts.Coordinator)
	require.NoError(t, err, "failed to load coordinator")

	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to connect to CL nodes")

	// Mirror original integration test semantics.
	// Legacy defaults: callbackGasLimit=500k, batchTxGasBudget=maxGasLimitCoordinator+400k.
	callbackGasLimit := c.BatchCallbackGasLimit
	if callbackGasLimit == 0 {
		callbackGasLimit = 500_000
	}
	batchTxGasBudget := c.BatchTxGasBudget
	if batchTxGasBudget == 0 {
		batchTxGasBudget = c.MaxGasLimitCoordinator + 400_000
	}
	expectedCountU32 := (batchTxGasBudget / callbackGasLimit) - 1
	require.Greater(t, expectedCountU32, uint32(1), "expected batched fulfillment count should be > 1")
	require.LessOrEqual(t, expectedCountU32, uint32(^uint16(0)), "expected count must fit uint16")
	expectedCount := uint16(expectedCountU32) //nolint:gosec // bounded by explicit <= uint16 max assertion above

	consumer, subID := newConsumerAndSub(ctx, t, chainClient, coord, linkToken, c)
	requestID, rErr := consumer.RequestRandomness(
		keyHash,
		subID,
		c.MinimumConfirmations,
		callbackGasLimit,
		true, // original test used native billing
		defaultNumWords,
		expectedCount,
	)
	require.NoError(t, rErr, "RequestRandomness failed")
	require.NotNil(t, requestID)

	// Wait until all requested fulfillments are reported by the consumer.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		reqCount, reqErr := consumer.RequestCount(ctx)
		if reqErr != nil {
			return false
		}
		respCount, respErr := consumer.ResponseCount(ctx)
		if respErr != nil {
			return false
		}
		return reqCount.Cmp(new(big.Int).SetUint64(uint64(expectedCount))) == 0 &&
			respCount.Cmp(new(big.Int).SetUint64(uint64(expectedCount))) == 0
	}, 2*time.Minute, 5*time.Second).Should(gomega.BeTrue(),
		"timed out waiting for request/fulfillment counts to reach %d", expectedCount)

	// Grab a fulfillment event and use its tx hash for assertions.
	var event *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		ev, fErr := coord.FilterRandomWordsFulfilled(
			&bind.FilterOpts{Context: ctx},
			requestID,
		)
		if fErr != nil {
			return false
		}
		event = ev
		return true
	}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
		"timed out waiting for RandomWordsFulfilled event for requestID %s", requestID)
	require.True(t, event.Success, "fulfillment should succeed")

	// Verify fulfillment tx was sent to BatchCoordinator.
	batchCoord, err := contracts.LoadBatchVRFCoordinatorV2Plus(chainClient, c.DeployedContracts.BatchCoordinator)
	require.NoError(t, err, "failed to load BatchVRFCoordinatorV2Plus")

	pollPeriod, pErr := time.ParseDuration(c.VRFJobPollPeriod)
	if pErr != nil || pollPeriod <= 0 {
		pollPeriod = time.Second
	}
	requestTimeout, rtErr := time.ParseDuration(c.VRFJobRequestTimeout)
	if rtErr != nil || requestTimeout <= 0 {
		requestTimeout = 24 * time.Hour
	}

	newPipelineSpec := &productvrfv2plus.TxPipelineSpec{
		Address:               coord.Address(),
		EstimateGasMultiplier: 1.1,
		FromAddress:           c.VRFKeyData.TxKeyAddresses[0],
	}
	observationSource, oErr := newPipelineSpec.String()
	require.NoError(t, oErr, "failed to build pipeline spec")

	currentJobID := c.VRFKeyData.VRFJobID
	switchJob := func(t *testing.T, batchFulfillmentEnabled bool) {
		t.Helper()
		if currentJobID != "" {
			resp, dErr := cl[0].DeleteJob(currentJobID)
			require.NoError(t, dErr, "failed deleting previous VRF job")
			require.Equal(t, http.StatusNoContent, resp.StatusCode, "unexpected status deleting VRF job")
		}

		namePrefix := "enabled"
		if !batchFulfillmentEnabled {
			namePrefix = "disabled"
		}
		newJobSpec := &productvrfv2plus.JobSpec{
			Name:                          fmt.Sprintf("vrf-v2-plus-batch-%s-%s", namePrefix, uuid.NewString()),
			CoordinatorAddress:            coord.Address(),
			BatchCoordinatorAddress:       c.DeployedContracts.BatchCoordinator,
			PublicKey:                     c.VRFKeyData.PubKeyCompressed,
			ExternalJobID:                 uuid.New().String(),
			ObservationSource:             observationSource,
			MinIncomingConfirmations:      int(c.MinimumConfirmations),
			FromAddresses:                 c.VRFKeyData.TxKeyAddresses,
			EVMChainID:                    in.Blockchains[0].Out.ChainID,
			BatchFulfillmentEnabled:       batchFulfillmentEnabled,
			BatchFulfillmentGasMultiplier: 1.1,
			BackOffInitialDelay:           15 * time.Second,
			BackOffMaxDelay:               5 * time.Minute,
			PollPeriod:                    pollPeriod,
			RequestTimeout:                requestTimeout,
		}
		job, cErr := cl[0].MustCreateJob(newJobSpec)
		require.NoError(t, cErr, "failed creating VRF job")
		currentJobID = job.Data.ID
	}

	t.Cleanup(func() {
		if currentJobID != "" {
			cl[0].MustDeleteJob(currentJobID) //nolint:errcheck // best-effort cleanup in test teardown
		}
	})

	runCase := func(t *testing.T, batchFulfillmentEnabled bool, expectedTo string, expectedLogsInFulfillmentTx int) {
		t.Helper()
		switchJob(t, batchFulfillmentEnabled)

		consumer, subID := newConsumerAndSub(ctx, t, chainClient, coord, linkToken, c)
		requestID, rErr := consumer.RequestRandomness(
			keyHash,
			subID,
			c.MinimumConfirmations,
			callbackGasLimit,
			true,
			defaultNumWords,
			expectedCount,
		)
		require.NoError(t, rErr, "RequestRandomness failed")
		require.NotNil(t, requestID)

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			reqCount, reqErr := consumer.RequestCount(ctx)
			if reqErr != nil {
				return false
			}
			respCount, respErr := consumer.ResponseCount(ctx)
			if respErr != nil {
				return false
			}
			return reqCount.Cmp(new(big.Int).SetUint64(uint64(expectedCount))) == 0 &&
				respCount.Cmp(new(big.Int).SetUint64(uint64(expectedCount))) == 0
		}, 2*time.Minute, 5*time.Second).Should(gomega.BeTrue(),
			"timed out waiting for request/fulfillment counts to reach %d", expectedCount)

		var ev *vrf_coordinator_v2_5.VRFCoordinatorV25RandomWordsFulfilled
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			found, fErr := coord.FilterRandomWordsFulfilled(
				&bind.FilterOpts{Context: ctx},
				requestID,
			)
			if fErr != nil {
				return false
			}
			ev = found
			return true
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"timed out waiting for RandomWordsFulfilled event for requestID %s", requestID)
		require.True(t, ev.Success, "fulfillment should succeed")

		tx, _, txErr := chainClient.Client.TransactionByHash(ctx, ev.Raw.TxHash)
		require.NoError(t, txErr, "failed to get fulfillment transaction")
		require.NotNil(t, tx.To(), "fulfillment tx should have a To address")
		require.Equal(t,
			strings.ToLower(expectedTo),
			strings.ToLower(tx.To().Hex()),
			"unexpected fulfillment tx destination")

		fulfillmentLogsCount, lErr := coord.CountRandomWordsFulfilledLogsInTx(ctx, ev.Raw.TxHash)
		require.NoError(t, lErr, "failed to count RandomWordsFulfilled logs in tx")
		require.Equal(t, expectedLogsInFulfillmentTx, fulfillmentLogsCount,
			"unexpected number of RandomWordsFulfilled logs in fulfillment tx")
	}

	t.Run("Batch Fulfillment Enabled", func(t *testing.T) {
		runCase(t, true, batchCoord.Address(), int(expectedCount))
	})

	t.Run("Batch Fulfillment Disabled", func(t *testing.T) {
		runCase(t, false, coord.Address(), 1)
	})
}
