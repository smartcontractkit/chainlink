package v2_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
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
			Key:          ptr[types.EIP55Address](bhsKey.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})
		sendEth(t, ownerKey, uni.backend, bhsKey.Address, 10)
	}
	keySpecificOverrides = append(keySpecificOverrides, toml.KeySpecific{
		// Gas lane.
		Key:          ptr[types.EIP55Address](vrfKey.EIP55Address),
		GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
	})

	keys = append(keys, ownerKey, vrfKey)

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasLanePriceWei, keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = ptr[uint32](2)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
		c.EVM[0].GasEstimator.LimitDefault = ptr(uint64(gasLimit))
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(time.Second)
	})

	heartbeatPeriod := 5 * time.Second

	t.Run("bhs_feeder_startheartbeats_happy_path", func(t *testing.T) {
		app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, uni.backend, keys...)
		require.NoError(t, app.Start(testutils.Context(t)))

		_ = vrftesthelpers.CreateAndStartBHSJob(
			t, bhsKeyAddresses, app, uni.bhsContractAddress.String(), "",
			uni.rootContractAddress.String(), "", "", 0, 200, heartbeatPeriod, 100)

		// Ensure log poller is ready and has all logs.
		require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Ready())
		require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Replay(testutils.Context(t), 1))

		initTxns := 260
		// Wait 260 blocks.
		for i := 0; i < initTxns; i++ {
			uni.backend.Commit()
		}
		diff := heartbeatPeriod + 1*time.Second
		t.Logf("Sleeping %.2f seconds before checking blockhash in BHS added by BHS_Heartbeats_Service\n", diff.Seconds())
		time.Sleep(diff)
		// storeEarliest in BHS contract stores blocktip - 256 in the Blockhash Store (BHS)
		tipHeader, err := uni.backend.Client().HeaderByNumber(testutils.Context(t), nil)
		require.NoError(t, err)
		// the storeEarliest transaction will end up in a new block, hence the + 1 below.
		blockNumberStored := tipHeader.Number.Uint64() - 256 + 1
		verifyBlockhashStored(t, uni.coordinatorV2UniverseCommon, blockNumberStored)
	})
}
