package v2_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"
	"github.com/smartcontractkit/chainlink-evm/pkg/chains/legacyevm"

	"github.com/smartcontractkit/chainlink-evm/pkg/assets"
	"github.com/smartcontractkit/chainlink-evm/pkg/config/toml"
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
	gasLimit := uint64(3_000_000)

	consumers := uni.vrfConsumers

	// generate n BHS keys to make sure BHS job rotates sending keys
	bhsKeyAddresses := make([]string, 0, len(consumers))
	keySpecificOverrides := make([]toml.KeySpecific, 0, len(consumers)+1)
	keys := make([]any, 0, len(consumers)+2)
	for range consumers {
		bhsKey := cltest.MustGenerateRandomKey(t)
		bhsKeyAddresses = append(bhsKeyAddresses, bhsKey.Address.String())
		keys = append(keys, bhsKey)
		keySpecificOverrides = append(keySpecificOverrides, toml.KeySpecific{
			Key:          new(bhsKey.EIP55Address),
			GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
		})
		sendEth(t, ownerKey, uni.backend, bhsKey.Address, 10)
	}
	keySpecificOverrides = append(keySpecificOverrides, toml.KeySpecific{
		// Gas lane.
		Key:          new(vrfKey.EIP55Address),
		GasEstimator: toml.KeySpecificGasEstimator{PriceMax: gasLanePriceWei},
	})

	keys = append(keys, ownerKey, vrfKey)

	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		simulatedOverrides(t, gasLanePriceWei, keySpecificOverrides...)(c, s)
		c.EVM[0].MinIncomingConfirmations = new(uint32(2))
		c.Feature.LogPoller = new(true)
		c.EVM[0].FinalityDepth = new(uint32(2))
		c.EVM[0].GasEstimator.LimitDefault = new(gasLimit)
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
		chain, ok := app.GetRelayers().LegacyEVMChains().Slice()[0].(legacyevm.Chain)
		require.True(t, ok)
		require.NoError(t, chain.LogPoller().Ready())
		require.NoError(t, chain.LogPoller().Replay(testutils.Context(t), 1))

		initTxns := 260
		// Wait 260 blocks.
		for range initTxns {
			uni.backend.Commit()
		}
		diff := heartbeatPeriod + 1*time.Second
		t.Logf("Sleeping %.2f seconds before checking blockhash in BHS added by BHS_Heartbeats_Service\n", diff.Seconds())
		time.Sleep(diff)
		// The heartbeat store tx may not reach the mempool before the first
		// Commit under load, so we can't predict which block it mines in.
		// Commit blocks and check current_tip-256 on each attempt until BHS
		// has a blockhash stored at that offset.
		require.Eventually(t, func() bool {
			uni.backend.Commit()
			tip, tipErr := uni.backend.Client().HeaderByNumber(testutils.Context(t), nil)
			if tipErr != nil || tip == nil || tip.Number.Uint64() < 256 {
				return false
			}
			_, err := uni.bhsContract.GetBlockhash(nil, new(big.Int).SetUint64(tip.Number.Uint64()-256))
			return err == nil
		}, testutils.WaitTimeoutCustom(t, 5*time.Minute), time.Second)
	})
}
