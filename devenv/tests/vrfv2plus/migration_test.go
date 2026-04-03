package vrfv2plus

import (
	"context"
	"fmt"
	"math/big"
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

// migrationStalenessSeconds and migrationGasAfterPayment mirror the constants
// used during env setup in core.go, which are not exported.
const (
	migrationStalenessSeconds    = uint32(86400)
	migrationGasAfterPayment     = uint32(33285)
	migrationFallbackLinkPerUnit = int64(1e16)
)

func TestVRFv2PlusMigration(t *testing.T) {
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

	linkToken, err := contracts.LoadLinkTokenContract(framework.L, chainClient, common.HexToAddress(c.DeployedContracts.LinkToken))
	require.NoError(t, err, "failed to load LINK token")

	bhs, err := contracts.LoadBlockhashStore(chainClient, c.DeployedContracts.BHS)
	require.NoError(t, err, "failed to load BlockhashStore")

	cl, err := clclient.New(in.NodeSets[0].Out.CLNodes)
	require.NoError(t, err, "failed to create CL clients")

	pollPeriod, pErr := time.ParseDuration(c.VRFJobPollPeriod)
	if pErr != nil || pollPeriod == 0 {
		pollPeriod = time.Second
	}

	// buildNewCoord deploys a fresh upgraded coordinator, registers the proving key,
	// configures it, and creates a VRF job on cl[0] pointing to it.
	buildNewCoord := func(t *testing.T) *contracts.EthereumVRFCoordinatorV2PlusUpgradedVersion {
		t.Helper()

		newCoord, dcErr := contracts.DeployVRFCoordinatorV2PlusUpgradedVersion(chainClient, bhs.Address())
		require.NoError(t, dcErr, "error deploying VRFCoordinatorV2PlusUpgradedVersion")

		provingKey, pkErr := contracts.EncodeOnChainVRFProvingKey(c.VRFKeyData.PubKeyUncompressed)
		require.NoError(t, pkErr, "error encoding VRF proving key")

		gasLaneMaxGas := products.EtherToGwei(big.NewFloat(float64(c.CLNodeMaxGasPriceGWei))).Uint64()
		require.NoError(t, newCoord.RegisterProvingKey(provingKey, gasLaneMaxGas), "error registering proving key on new coordinator")

		require.NoError(t, newCoord.SetConfig(
			c.MinimumConfirmations,
			c.MaxGasLimitCoordinator,
			migrationStalenessSeconds,
			migrationGasAfterPayment,
			big.NewInt(migrationFallbackLinkPerUnit),
			c.FlatFeeNativePPM,
			c.FlatFeeLinkDiscountPPM,
			c.NativePremiumPercentage,
			c.LinkPremiumPercentage,
		), "error setting config on new coordinator")

		require.NoError(t, newCoord.SetLINKAndLINKNativeFeed(
			c.DeployedContracts.LinkToken,
			c.DeployedContracts.MockFeed,
		), "error setting LINK and feed on new coordinator")

		pipelineSpec := &productvrfv2plus.TxPipelineSpec{
			Address:               newCoord.Address(),
			EstimateGasMultiplier: 1.1,
			FromAddress:           c.VRFKeyData.TxKeyAddresses[0],
		}
		observationSource, oErr := pipelineSpec.String()
		require.NoError(t, oErr, "failed to build pipeline spec for new coordinator")

		newJobSpec := &productvrfv2plus.JobSpec{
			Name:                          "vrf-v2-plus-migration-" + uuid.NewString(),
			CoordinatorAddress:            newCoord.Address(),
			BatchCoordinatorAddress:       "",
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
			PollPeriod:                    pollPeriod,
			RequestTimeout:                24 * time.Hour,
		}
		job, jErr := cl[0].MustCreateJob(newJobSpec)
		require.NoError(t, jErr, "error creating VRF job for new coordinator")
		if job != nil && job.Data.ID != "" {
			jobID := job.Data.ID
			t.Cleanup(func() {
				cl[0].MustDeleteJob(jobID) //nolint:errcheck // best-effort cleanup in test teardown
			})
		}

		return newCoord
	}

	// requestAndWaitNewCoord sends a request from the consumer and waits for
	// RandomWordsFulfilled on the upgraded coordinator.
	requestAndWaitNewCoord := func(
		t *testing.T,
		reqCtx context.Context,
		consumer *contracts.EthereumVRFv2PlusLoadTestConsumer,
		newCoord *contracts.EthereumVRFCoordinatorV2PlusUpgradedVersion,
		subID *big.Int,
		isNative bool,
	) {
		t.Helper()
		requestID, rErr := consumer.RequestRandomness(keyHash, subID, c.MinimumConfirmations, defaultCallbackGasLimit, isNative, defaultNumWords, defaultRequestCount)
		require.NoError(t, rErr, "RequestRandomness failed")
		require.NotNil(t, requestID)

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			event, fErr := newCoord.FilterRandomWordsFulfilled(&bind.FilterOpts{Context: reqCtx}, requestID)
			if fErr != nil {
				return false
			}
			return event.Success
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"timed out waiting for RandomWordsFulfilled on new coordinator for requestID %s", requestID)
	}

	// requestAndWaitWrapperNewCoord sends a wrapper request and waits for
	// RandomWordsFulfilled on the upgraded coordinator.
	requestAndWaitWrapperNewCoord := func(
		t *testing.T,
		reqCtx context.Context,
		consumer *contracts.EthereumVRFV2PlusWrapperLoadTestConsumer,
		newCoord *contracts.EthereumVRFCoordinatorV2PlusUpgradedVersion,
		isNative bool,
	) {
		t.Helper()
		var requestID *big.Int
		var reqErr error
		if isNative {
			requestID, reqErr = consumer.RequestRandomWordsNative(c.MinimumConfirmations, defaultCallbackGasLimit, defaultNumWords, defaultRequestCount)
		} else {
			requestID, reqErr = consumer.RequestRandomWords(c.MinimumConfirmations, defaultCallbackGasLimit, defaultNumWords, defaultRequestCount)
		}
		require.NoError(t, reqErr, "wrapper RequestRandomWords failed")
		require.NotNil(t, requestID)

		gomega.NewGomegaWithT(t).Eventually(func() bool {
			event, fErr := newCoord.FilterRandomWordsFulfilled(&bind.FilterOpts{Context: reqCtx}, requestID)
			if fErr != nil {
				return false
			}
			return event.Success
		}, defaultFulfillTimeout, 5*time.Second).Should(gomega.BeTrue(),
			"timed out waiting for wrapper RandomWordsFulfilled on new coordinator for requestID %s", requestID)

		status, sErr := consumer.GetRequestStatus(reqCtx, requestID)
		require.NoError(t, sErr, "error getting wrapper request status")
		require.True(t, status.Fulfilled)
	}

	t.Run("Test migration of Subscription Billing subID", func(t *testing.T) {
		// Deploy 2 consumers and create a funded subscription.
		consumers := make([]*contracts.EthereumVRFv2PlusLoadTestConsumer, 2)
		for i := range consumers {
			consumer, dcErr := contracts.DeployVRFv2PlusLoadTestConsumer(chainClient, coord.Address())
			require.NoError(t, dcErr, "failed to deploy consumer %d", i)
			consumers[i] = consumer
		}

		subID, sErr := createAndFundSub(ctx, chainClient, coord, linkToken, c.SubFundingAmountLink, c.SubFundingAmountNative)
		require.NoError(t, sErr, "failed to create and fund subscription")

		for _, consumer := range consumers {
			require.NoError(t, coord.AddConsumer(subID, consumer.Address()), "failed to add consumer")
		}

		// Verify subID appears in old coordinator's active list.
		activeSubIDsBefore, asErr := coord.GetActiveSubscriptionIDs(ctx, big.NewInt(0), big.NewInt(0))
		require.NoError(t, asErr, "error getting active sub ids before migration")
		require.True(t, containsBigInt(activeSubIDsBefore, subID),
			"subID should be in old coordinator's active list before migration")

		oldSubBefore, gsErr := coord.GetSubscription(ctx, subID)
		require.NoError(t, gsErr, "error getting subscription before migration")

		newCoord := buildNewCoord(t)

		require.NoError(t, coord.RegisterMigratableCoordinator(newCoord.Address()), "error registering migratable coordinator")

		oldCoordLinkBefore, err := coord.GetLinkTotalBalance(ctx)
		require.NoError(t, err, "error getting old coord link balance before migration")
		oldCoordNativeBefore, err := coord.GetNativeTotalBalance(ctx)
		require.NoError(t, err, "error getting old coord native balance before migration")
		newCoordLinkBefore, err := newCoord.GetLinkTotalBalance(ctx)
		require.NoError(t, err, "error getting new coord link balance before migration")
		newCoordNativeBefore, err := newCoord.GetNativeTotalBalance(ctx)
		require.NoError(t, err, "error getting new coord native balance before migration")

		_, mErr := coord.Migrate(subID, newCoord.Address())
		require.NoError(t, mErr, "error migrating subscription")

		oldCoordLinkAfter, err := coord.GetLinkTotalBalance(ctx)
		require.NoError(t, err)
		oldCoordNativeAfter, err := coord.GetNativeTotalBalance(ctx)
		require.NoError(t, err)
		newCoordLinkAfter, err := newCoord.GetLinkTotalBalance(ctx)
		require.NoError(t, err)
		newCoordNativeAfter, err := newCoord.GetNativeTotalBalance(ctx)
		require.NoError(t, err)

		migratedSub, msErr := newCoord.GetSubscription(ctx, subID)
		require.NoError(t, msErr, "error getting migrated subscription from new coordinator")

		// Verify subscription data was transferred correctly.
		require.Equal(t, oldSubBefore.NativeBalance, migratedSub.NativeBalance)
		require.Equal(t, oldSubBefore.Balance, migratedSub.Balance)
		require.Equal(t, oldSubBefore.Owner, migratedSub.Owner)
		require.Equal(t, oldSubBefore.Consumers, migratedSub.Consumers)

		// Old sub should be gone.
		_, gsErr = coord.GetSubscription(ctx, subID)
		require.Error(t, gsErr, "should error when getting deleted sub from old coordinator")

		// Migrated subID should no longer appear in old coordinator's active list.
		activeSubIDsOldAfter, _ := coord.GetActiveSubscriptionIDs(ctx, big.NewInt(0), big.NewInt(0))
		require.False(t, containsBigInt(activeSubIDsOldAfter, subID), "migrated subID should not be in old coordinator's list after migration")

		// New coordinator should have exactly 1 active sub.
		activeSubIDsNew, asErr := newCoord.GetActiveSubscriptionIDs(ctx, big.NewInt(0), big.NewInt(0))
		require.NoError(t, asErr, "error getting active sub ids from new coordinator")
		require.Len(t, activeSubIDsNew, 1, "new coordinator should have exactly 1 active sub")
		require.Equal(t, subID, activeSubIDsNew[0])

		// Verify balance math.
		expectedLinkNew := new(big.Int).Add(oldSubBefore.Balance, newCoordLinkBefore)
		expectedNativeNew := new(big.Int).Add(oldSubBefore.NativeBalance, newCoordNativeBefore)
		expectedLinkOld := new(big.Int).Sub(oldCoordLinkBefore, oldSubBefore.Balance)
		expectedNativeOld := new(big.Int).Sub(oldCoordNativeBefore, oldSubBefore.NativeBalance)
		require.Equal(t, 0, expectedLinkNew.Cmp(newCoordLinkAfter))
		require.Equal(t, 0, expectedNativeNew.Cmp(newCoordNativeAfter))
		require.Equal(t, 0, expectedLinkOld.Cmp(oldCoordLinkAfter))
		require.Equal(t, 0, expectedNativeOld.Cmp(oldCoordNativeAfter))

		// Verify coordinator address was updated in all consumers after migration.
		for _, consumer := range consumers {
			coordAddrInConsumer, gcErr := consumer.GetCoordinator(ctx)
			require.NoError(t, gcErr, "error getting coordinator from consumer")
			require.Equal(t, newCoord.Address(), coordAddrInConsumer.Hex(), "coordinator in consumer should be updated after migration")
		}

		// Verify requests can be fulfilled via the new coordinator.
		requestAndWaitNewCoord(t, ctx, consumers[0], newCoord, subID, false) // LINK billing
		requestAndWaitNewCoord(t, ctx, consumers[1], newCoord, subID, true)  // native billing
	})

	t.Run("Test migration of direct billing using VRFV2PlusWrapper subID", func(t *testing.T) {
		wrapper, wErr := contracts.LoadVRFV2PlusWrapper(chainClient, c.DeployedContracts.Wrapper)
		require.NoError(t, wErr, "failed to load wrapper")

		wrapperConsumer, wcErr := contracts.LoadVRFV2PlusWrapperLoadTestConsumer(chainClient, c.DeployedContracts.WrapperConsumer)
		require.NoError(t, wcErr, "failed to load wrapper consumer")

		wrapperSubID, ok := new(big.Int).SetString(c.DeployedContracts.WrapperSubID, 10)
		require.True(t, ok, "failed to parse wrapper sub ID: %s", c.DeployedContracts.WrapperSubID)

		subID := wrapperSubID
		reconcileConfiguredFunding(ctx, t, chainClient, coord, linkToken, c)

		// After subtest 1 migrated the test sub, the old coordinator should have
		// the wrapper sub present in active subs.
		activeSubIDsBefore, asErr := coord.GetActiveSubscriptionIDs(ctx, big.NewInt(0), big.NewInt(0))
		require.NoError(t, asErr, "error getting active sub ids before wrapper migration")
		require.True(t, containsBigInt(activeSubIDsBefore, subID),
			"old coordinator active subs should include wrapper sub before wrapper migration")

		oldSubBefore, gsErr := coord.GetSubscription(ctx, subID)
		require.NoError(t, gsErr, "error getting wrapper subscription before migration")

		newCoord := buildNewCoord(t)

		require.NoError(t, coord.RegisterMigratableCoordinator(newCoord.Address()), "error registering migratable coordinator for wrapper")

		oldCoordLinkBefore, err := coord.GetLinkTotalBalance(ctx)
		require.NoError(t, err)
		oldCoordNativeBefore, err := coord.GetNativeTotalBalance(ctx)
		require.NoError(t, err)
		newCoordLinkBefore, err := newCoord.GetLinkTotalBalance(ctx)
		require.NoError(t, err)
		newCoordNativeBefore, err := newCoord.GetNativeTotalBalance(ctx)
		require.NoError(t, err)

		_, mErr := coord.Migrate(subID, newCoord.Address())
		require.NoError(t, mErr, "error migrating wrapper subscription")

		oldCoordLinkAfter, err := coord.GetLinkTotalBalance(ctx)
		require.NoError(t, err)
		oldCoordNativeAfter, err := coord.GetNativeTotalBalance(ctx)
		require.NoError(t, err)
		newCoordLinkAfter, err := newCoord.GetLinkTotalBalance(ctx)
		require.NoError(t, err)
		newCoordNativeAfter, err := newCoord.GetNativeTotalBalance(ctx)
		require.NoError(t, err)

		migratedSub, msErr := newCoord.GetSubscription(ctx, subID)
		require.NoError(t, msErr, "error getting migrated wrapper subscription from new coordinator")

		require.Equal(t, oldSubBefore.NativeBalance, migratedSub.NativeBalance)
		require.Equal(t, oldSubBefore.Balance, migratedSub.Balance)
		require.Equal(t, oldSubBefore.Owner, migratedSub.Owner)
		require.Equal(t, oldSubBefore.Consumers, migratedSub.Consumers)

		_, gsErr = coord.GetSubscription(ctx, subID)
		require.Error(t, gsErr, "should error when getting deleted wrapper sub from old coordinator")

		// Old coordinator should have no active subs (wrapper sub was the last one).
		_, activeSubIDsOldAfterErr := coord.GetActiveSubscriptionIDs(ctx, big.NewInt(0), big.NewInt(0))
		require.Error(t, activeSubIDsOldAfterErr, "old coordinator should have no active subs after wrapper sub migration")

		activeSubIDsNew, asErr := newCoord.GetActiveSubscriptionIDs(ctx, big.NewInt(0), big.NewInt(0))
		require.NoError(t, asErr, "error getting active sub ids from new coordinator")
		require.Len(t, activeSubIDsNew, 1, "new coordinator should have exactly 1 active sub")
		require.Equal(t, subID, activeSubIDsNew[0])

		expectedLinkNew := new(big.Int).Add(oldSubBefore.Balance, newCoordLinkBefore)
		expectedNativeNew := new(big.Int).Add(oldSubBefore.NativeBalance, newCoordNativeBefore)
		expectedLinkOld := new(big.Int).Sub(oldCoordLinkBefore, oldSubBefore.Balance)
		expectedNativeOld := new(big.Int).Sub(oldCoordNativeBefore, oldSubBefore.NativeBalance)
		require.Equal(t, 0, expectedLinkNew.Cmp(newCoordLinkAfter))
		require.Equal(t, 0, expectedNativeNew.Cmp(newCoordNativeAfter))
		require.Equal(t, 0, expectedLinkOld.Cmp(oldCoordLinkAfter))
		require.Equal(t, 0, expectedNativeOld.Cmp(oldCoordNativeAfter))

		// Verify the wrapper's coordinator pointer was updated.
		coordInWrapper, gcErr := wrapper.Coordinator(ctx)
		require.NoError(t, gcErr, "error getting coordinator from wrapper")
		require.Equal(t, newCoord.Address(), coordInWrapper.Hex(), "wrapper coordinator should be updated after migration")

		// Verify wrapper requests can be fulfilled via the new coordinator.
		requestAndWaitWrapperNewCoord(t, ctx, wrapperConsumer, newCoord, false) // LINK billing
		requestAndWaitWrapperNewCoord(t, ctx, wrapperConsumer, newCoord, true)  // native billing
	})
}

// containsBigInt returns true if ids contains target.
func containsBigInt(ids []*big.Int, target *big.Int) bool {
	for _, id := range ids {
		if id.Cmp(target) == 0 {
			return true
		}
	}
	return false
}
