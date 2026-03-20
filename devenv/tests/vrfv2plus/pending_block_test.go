package vrfv2plus

import (
	"fmt"
	"strconv"
	"testing"

	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-testing-framework/framework"

	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/smartcontractkit/chainlink/devenv/products"
	"github.com/smartcontractkit/chainlink/devenv/products/vrfv2plus"
)

func TestVRFv2PlusPendingBlockSimulationAndZeroConfirmationDelays(t *testing.T) {
	t.Cleanup(func() {
		_, cErr := framework.SaveContainerLogs(fmt.Sprintf("%s-%s", framework.DefaultCTFLogsDir, t.Name()))
		require.NoError(t, cErr)
	})

	outputFile := "../../env-vrf2plus-out.toml"

	in, err := de.LoadOutput[de.Cfg](outputFile)
	require.NoError(t, err, "failed to load env-vrf2plus-out.toml")

	cfg, err := products.LoadOutput[vrfv2plus.Configurator](outputFile)
	require.NoError(t, err, "failed to load VRFv2Plus config")
	require.NotEmpty(t, cfg.Config, "vrfv2_plus config must not be empty")

	c := cfg.Config[0]
	require.NotEmpty(t, c.DeployedContracts.Coordinator, "coordinator address must not be empty")
	require.NotEmpty(t, c.VRFKeyData.KeyHash, "key hash must not be empty")

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

	// Deploy consumer, create funded subscription, add consumer.
	consumer, dErr := contracts.DeployVRFv2PlusLoadTestConsumer(chainClient, coord.Address())
	require.NoError(t, dErr, "failed to deploy load test consumer")

	subID, sErr := createAndFundSub(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, c.SubFundingAmountNative)
	require.NoError(t, sErr, "failed to create and fund subscription")

	aErr := coord.AddConsumer(subID, consumer.Address())
	require.NoError(t, aErr, "failed to add consumer to subscription")

	// Request randomness with 0 confirmations (env is configured with minimum_confirmations=0
	// and vrf_job_simulation_block="pending").
	requested, qErr := consumer.RequestRandomnessWithEvent(
		keyHash, subID, 0, defaultCallbackGasLimit, true, defaultNumWords, defaultRequestCount,
	)
	require.NoError(t, qErr, "RequestRandomness failed")
	require.NotNil(t, requested)

	gomega.NewGomegaWithT(t).Eventually(func() bool {
		event, fErr := coord.FilterRandomWordsFulfilled(
			&bind.FilterOpts{Context: ctx},
			requested.RequestId,
		)
		if fErr != nil {
			return false
		}
		return event.Success
	}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
		"timed out waiting for RandomWordsFulfilled event")

	status, sErr := consumer.GetRequestStatus(ctx, requested.RequestId)
	require.NoError(t, sErr, "error getting rand request status")
	require.True(t, status.Fulfilled, "request should be fulfilled")
}
