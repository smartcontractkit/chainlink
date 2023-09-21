package v2_test

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/v2/core/assets"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/config/toml"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"

	"github.com/stretchr/testify/require"
)

func TestStartHeartbeats(t *testing.T) {
	t.Parallel()
	ownerKey := cltest.MustGenerateRandomKey(t)
	uni := newVRFCoordinatorV2Universe(t, ownerKey, 2)

	vrfKey := cltest.MustGenerateRandomKey(t)
	sendEth(t, ownerKey, uni.backend, vrfKey.Address, 10)
	gasLanePriceWei := assets.GWei(1)
	gasLimit := 3_000_000

	consumers := uni.vrfConsumers

	// generate n BHS keys to make sure BHS job rotates sending keys
	var bhsKeyAddresses []string
	var keySpecificOverrides []toml.KeySpecific
	var keys []interface{}
	for i := 0; i < len(consumers); i++ {
		bhsKey := cltest.MustGenerateRandomKey(t)
		bhsKeyAddresses = append(bhsKeyAddresses, bhsKey.Address.String())
		keys = append(keys, bhsKey)
		keySpecificOverrides = append(keySpecificOverrides, toml.KeySpecific{
			Key:          ptr(bhsKey.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})
		sendEth(t, ownerKey, uni.backend, bhsKey.Address, 10)
	}
	keySpecificOverrides = append(keySpecificOverrides, toml.KeySpecific{
		// Gas lane.
		Key:          ptr(vrfKey.EIP55Address),
		GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
	})

	keys = append(keys, ownerKey, vrfKey)

	config, _ := heavyweight.FullTestDBV2(t, "vrfv2_needs_blockhash_store", func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasLanePriceWei, keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint32(gasLimit))
		c.EVM[0].LogPollInterval = models.MustNewDuration(time.Second)
	})

	heartbeatPeriod := 5 * time.Second
	// lggr, logs := logger.TestLoggerObserved(t, zapcore.DebugLevel)

	t.Run("bhs_feeder_startheartbeats_happy_path", func(tt *testing.T) {
		// consumerContracts := uni.consumerContracts
		// consumerContractAddresses := uni.consumerContractAddresses
		coordinator := uni.rootContract
		coordinatorAddress := uni.rootContractAddress
		batchCoordinatorAddress := uni.batchCoordinatorContractAddress
		vrfOwnerAddress := ptr(uni.vrfOwnerAddress)
		vrfVersion := vrfcommon.V2

		app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, keys...)
		require.NoError(t, app.Start(testutils.Context(t)))

		// Create VRF job.
		vrfJobs := createVRFJobs(
			t,
			[][]ethkey.KeyV2{{vrfKey}},
			app,
			coordinator,
			coordinatorAddress,
			batchCoordinatorAddress,
			uni.coordinatorV2UniverseCommon,
			vrfOwnerAddress,
			vrfVersion,
			false,
			gasLanePriceWei)
		_ = vrfJobs[0].VRFSpec.PublicKey.MustHash()

		var (
			v2CoordinatorAddress     string
			v2PlusCoordinatorAddress string
		)

		if vrfVersion == vrfcommon.V2 {
			v2CoordinatorAddress = coordinatorAddress.String()
		} else if vrfVersion == vrfcommon.V2Plus {
			v2PlusCoordinatorAddress = coordinatorAddress.String()
		}

		_ = vrftesthelpers.CreateAndStartBHSJob(
			t, bhsKeyAddresses, app, uni.bhsContractAddress.String(), "",
			v2CoordinatorAddress, v2PlusCoordinatorAddress, "", 0, 200, heartbeatPeriod)

		// Ensure log poller is ready and has all logs.
		require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Ready())
		require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Replay(testutils.Context(t), 1))

		initTxns := 260
		// Wait 260 blocks.
		for i := 0; i < initTxns+1; i++ {
			uni.backend.Commit()
		}
		diff := heartbeatPeriod + 2*time.Second
		t.Logf("Sleeping %.2f seconds before checking blockhash in BHS added by BHS_Heartbeats_Service\n", diff.Seconds())
		time.Sleep(diff)
		verifyBlockhashStored(t, uni.coordinatorV2UniverseCommon, uint64(initTxns+20-256))
	})
}

// Send eth from prefunded account.
// Amount is number of ETH not wei.
// func sendWei(t *testing.T, key ethkey.KeyV2, ec *backends.SimulatedBackend, to common.Address, wei *big.Int) {
// 	nonce, err := ec.PendingNonceAt(testutils.Context(t), key.Address)
// 	require.NoError(t, err)
// 	tx := gethtypes.NewTx(&gethtypes.DynamicFeeTx{
// 		ChainID:   big.NewInt(1337),
// 		Nonce:     nonce,
// 		GasTipCap: big.NewInt(1),
// 		GasFeeCap: assets.GWei(10).ToInt(), // block base fee in sim
// 		Gas:       uint64(21_000),
// 		To:        &to,
// 		Value:     wei,
// 		Data:      nil,
// 	})
// 	signedTx, err := gethtypes.SignTx(tx, gethtypes.NewLondonSigner(big.NewInt(1337)), key.ToEcdsaPrivKey())
// 	require.NoError(t, err)
// 	err = ec.SendTransaction(testutils.Context(t), signedTx)
// 	require.NoError(t, err)
// 	ec.Commit()
// }
