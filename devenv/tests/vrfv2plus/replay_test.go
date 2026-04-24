package vrfv2plus

import (
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"
	"github.com/smartcontractkit/chainlink-testing-framework/framework/clclient"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	productvrfv2plus "github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFv2PlusReplayAfterTimeout(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrf2plus-out.toml"

	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load env-vrf2plus-out.toml")

	cfg, err := products.LoadOutput[productvrfv2plus.Configurator](outputFile)
	require.NoError(t, err, "failed to load VRFv2Plus config")
	require.NotEmpty(t, cfg.Config, "vrfv2_plus config must not be empty")

	c := cfg.Config[0]
	require.NotEmpty(t, c.DeployedContracts.Coordinator, "coordinator address must not be empty")
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

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to create CL clients")

	batchCoordAddr := c.DeployedContracts.BatchCoordinator

	t.Run("Timed out request fulfilled after node restart with replay", func(t *testing.T) {
		var isNativeBilling = false

		// Deploy consumer and create an unfunded subscription so the request gets stuck.
		consumer, dErr := contracts.DeployVRFv2PlusLoadTestConsumer(chainClient, coord.Address())
		require.NoError(t, dErr, "failed to deploy load test consumer")

		subTx, sErr := coord.CreateSubscription()
		require.NoError(t, sErr, "failed to create subscription")
		receipt, rErr := chainClient.Client.TransactionReceipt(ctx, subTx.Hash())
		require.NoError(t, rErr)
		subID, pErr := contracts.FindSubscriptionID(receipt)
		require.NoError(t, pErr)

		aErr := coord.AddConsumer(subID, consumer.Address())
		require.NoError(t, aErr, "failed to add consumer")

		// Send the initial request with an unfunded sub — it will time out.
		initialRequest, qErr := consumer.RequestRandomnessWithEvent(
			keyHash, subID, c.MinimumConfirmations,
			defaultCallbackGasLimit, isNativeBilling, defaultNumWords, defaultRequestCount,
		)
		require.NoError(t, qErr, "error requesting randomness")

		// Wait until the underfunded request is observed as pending. This is
		// more robust in CI than a fixed sleep around request timeout boundaries.
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			pendingReqExists, pxErr := coord.PendingRequestsExist(ctx, subID)
			if pxErr != nil {
				return false
			}
			return pendingReqExists
		}, 30*time.Second, time.Second).Should(gomega.BeTrue(),
			"pending request must exist before funding and replay")

		// Fund the subscription so the node could fulfill — but it already timed out.
		linkToken, lErr := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
		require.NoError(t, lErr, "failed to load LINK token")

		nativeWei := products.EtherToWei(big.NewFloat(c.SubFundingAmountNative))
		fErr := coord.FundSubscriptionWithNative(subID, nativeWei)
		require.NoError(t, fErr, "failed to fund subscription with native")

		encodedSubID, eErr := encodeSubID(subID)
		require.NoError(t, eErr)
		_, fErr = linkToken.TransferAndCall(coord.Address(), products.EtherToWei(big.NewFloat(c.SubFundingAmountLink)), encodedSubID)
		require.NoError(t, fErr, "failed to fund subscription with LINK")

		// Delete the current job (which has the 5s timeout).
		resp, delErr := cl[0].DeleteJob(c.VRFKeyData.VRFJobID)
		require.NoError(t, delErr, "error deleting VRF job")
		require.Equal(t, http.StatusNoContent, resp.StatusCode, "unexpected status code deleting job")

		// Build a new pipeline spec identical to the original but using the default "latest" block.
		newPipelineSpec := &productvrfv2plus.TxPipelineSpec{
			Address:               coord.Address(),
			EstimateGasMultiplier: 1.1,
			FromAddress:           c.VRFKeyData.TxKeyAddresses[0],
		}
		observationSource, oErr := newPipelineSpec.String()
		require.NoError(t, oErr, "failed to build pipeline spec")

		// Create the new job with a 1h timeout.
		var newJobID string
		newJobSpec := &productvrfv2plus.JobSpec{
			Name:                          "vrf-v2-plus-replay-" + uuid.NewString(),
			CoordinatorAddress:            coord.Address(),
			BatchCoordinatorAddress:       batchCoordAddr,
			PublicKey:                     c.VRFKeyData.PubKeyCompressed,
			ExternalJobID:                 uuid.New().String(),
			ObservationSource:             observationSource,
			MinIncomingConfirmations:      int(c.MinimumConfirmations),
			FromAddresses:                 c.VRFKeyData.TxKeyAddresses,
			EVMChainID:                    in.Blockchains[0].Out.ChainID,
			BatchFulfillmentEnabled:       false,
			BatchFulfillmentGasMultiplier: 1.1,
			BackOffInitialDelay:           15 * time.Second,
			BackOffMaxDelay:               5 * time.Minute,
			PollPeriod:                    1 * time.Second,
			RequestTimeout:                1 * time.Hour,
		}
		job, jErr := cl[0].MustCreateJob(newJobSpec)
		require.NoError(t, jErr, "error creating replay VRF job")
		newJobID = job.Data.ID

		t.Cleanup(func() {
			if newJobID != "" {
				cl[0].MustDeleteJob(newJobID) //nolint:errcheck // best-effort cleanup in test teardown
			}
		})

		// Wait for the initial stuck request to be fulfilled by the new job with 1h timeout.
		gomega.NewGomegaWithT(t).Eventually(func() bool {
			event, wErr := coord.FilterRandomWordsFulfilled(
				&bind.FilterOpts{Context: ctx},
				initialRequest.RequestId,
			)
			if wErr != nil {
				return false
			}
			return event.Success
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"timed out waiting for initial request to be fulfilled by replayed job")

		event, wErr := coord.FilterRandomWordsFulfilled(&bind.FilterOpts{Context: ctx}, initialRequest.RequestId)
		require.NoError(t, wErr)
		require.False(t, event.OnlyPremium, "RandomWordsFulfilled Event's OnlyPremium field should be false")
		require.Equal(t, isNativeBilling, event.NativePayment, "RandomWordsFulfilled Event's NativePayment should be false")
		require.True(t, event.Success, "RandomWordsFulfilled Event's Success should be true")

		status, sErr := consumer.GetRequestStatus(ctx, initialRequest.RequestId)
		require.NoError(t, sErr, "error getting rand request status")
		require.True(t, status.Fulfilled, "request should be fulfilled")
	})
}
