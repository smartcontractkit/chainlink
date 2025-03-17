package v1_test

import (
	"encoding/hex"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	commonconfig "github.com/smartcontractkit/chainlink-common/pkg/config"

	ubig "github.com/smartcontractkit/chainlink-integrations/evm/utils/big"
	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrfcommon"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/vrftesthelpers"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/v2/core/utils/testutils/heavyweight"
)

func TestIntegration_VRF_JPV2(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		eip1559 bool
	}{
		{"legacy", false},
		{"eip1559", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			if tt.name == "eip1559" {
				t.Skip("fails after geth upgrade https://github.com/smartcontractkit/chainlink/pull/11809")
			}
			ctx := testutils.Context(t)
			config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.EIP1559DynamicFees = &test.eip1559
				c.EVM[0].ChainID = (*ubig.Big)(testutils.SimulatedChainID)
			})
			key1 := cltest.MustGenerateRandomKey(t)
			key2 := cltest.MustGenerateRandomKey(t)
			cu := vrftesthelpers.NewVRFCoordinatorUniverse(t, key1, key2)
			incomingConfs := 2
			app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, cu.Backend, key1, key2)
			require.NoError(t, app.Start(ctx))

			jb, vrfKey := createVRFJobRegisterKey(t, cu, app, incomingConfs)
			require.NoError(t, app.JobSpawner().CreateJob(ctx, nil, &jb))

			_, err := cu.ConsumerContract.TestRequestRandomness(cu.Carol,
				vrfKey.PublicKey.MustHash(), big.NewInt(100))
			require.NoError(t, err)

			_, err = cu.ConsumerContract.TestRequestRandomness(cu.Carol,
				vrfKey.PublicKey.MustHash(), big.NewInt(100))
			require.NoError(t, err)
			cu.Backend.Commit()
			t.Log("Sent 2 test requests")
			// Mine the required number of blocks
			// So our request gets confirmed.
			for i := 0; i < incomingConfs; i++ {
				cu.Backend.Commit()
			}
			var runs []pipeline.Run
			require.Eventually(t, func() bool {
				runs, err = app.PipelineORM().GetAllRuns(ctx)
				require.NoError(t, err)
				// It possible that we send the test request
				// before the Job spawner has started the vrf services, which is fine
				// the lb will backfill the logs. However we need to
				// keep blocks coming in for the lb to send the backfilled logs.
				cu.Backend.Commit()
				return len(runs) == 2 && runs[0].State == pipeline.RunStatusCompleted && runs[1].State == pipeline.RunStatusCompleted
			}, testutils.WaitTimeout(t), 1*time.Second)
			assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].FatalErrors)
			assert.Equal(t, 4, len(runs[0].PipelineTaskRuns))
			assert.Equal(t, 4, len(runs[1].PipelineTaskRuns))
			assert.NotNil(t, 0, runs[0].Outputs.Val)
			assert.NotNil(t, 0, runs[1].Outputs.Val)

			// stop jobs as to not cause a race condition in geth simulated backend
			// between job creating new tx and fulfillment logs polling below
			require.NoError(t, app.JobSpawner().DeleteJob(ctx, nil, jb.ID))

			// Ensure the eth transaction gets confirmed on chain.
			require.Eventually(t, func() bool {
				orm := txmgr.NewTxStore(app.GetDB(), app.GetLogger())
				uc, err2 := orm.CountUnconfirmedTransactions(ctx, key1.Address, testutils.SimulatedChainID)
				require.NoError(t, err2)
				return uc == 0
			}, testutils.WaitTimeout(t), 100*time.Millisecond)

			// Assert the request was fulfilled on-chain.
			var rf []*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled
			require.Eventually(t, func() bool {
				rfIterator, err2 := cu.RootContract.FilterRandomnessRequestFulfilled(nil)
				require.NoError(t, err2, "failed to subscribe to RandomnessRequest logs")
				rf = nil
				for rfIterator.Next() {
					rf = append(rf, rfIterator.Event)
				}
				return len(rf) == 2
			}, testutils.WaitTimeout(t), 500*time.Millisecond)

			// Check that each sending address sent one transaction
			n1, err := cu.Backend.Client().PendingNonceAt(ctx, key1.Address)
			require.NoError(t, err)
			require.EqualValues(t, 1, n1)

			n2, err := cu.Backend.Client().PendingNonceAt(ctx, key2.Address)
			require.NoError(t, err)
			require.EqualValues(t, 1, n2)
		})
	}
}

func TestIntegration_VRF_WithBHS(t *testing.T) {
	t.Parallel()
	ctx := testutils.Context(t)
	config, _ := heavyweight.FullTestDBV2(t, func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
		c.EVM[0].BlockBackfillDepth = ptr[uint32](500)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
		c.EVM[0].LogPollInterval = commonconfig.MustNewDuration(time.Second)
		c.EVM[0].ChainID = (*ubig.Big)(testutils.SimulatedChainID)
	})
	key := cltest.MustGenerateRandomKey(t)
	cu := vrftesthelpers.NewVRFCoordinatorUniverse(t, key)
	incomingConfs := 2
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, cu.Backend, key)
	require.NoError(t, app.Start(ctx))

	// Create VRF Job but do not start it yet
	jb, vrfKey := createVRFJobRegisterKey(t, cu, app, incomingConfs)

	sendingKeys := []string{key.Address.String()}

	// Create BHS Job and start it
	bhsJob := vrftesthelpers.CreateAndStartBHSJob(t, sendingKeys, app, cu.BHSContractAddress.String(),
		cu.RootContractAddress.String(), "", "", "", 0, 200, 0, 100)

	// Ensure log poller is ready and has all logs.
	require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Ready())
	require.NoError(t, app.GetRelayers().LegacyEVMChains().Slice()[0].LogPoller().Replay(ctx, 1))

	// Create a VRF request
	_, err := cu.ConsumerContract.TestRequestRandomness(cu.Carol,
		vrfKey.PublicKey.MustHash(), big.NewInt(100))
	require.NoError(t, err)

	cu.Backend.Commit()
	h, err := cu.Backend.Client().HeaderByNumber(testutils.Context(t), nil)
	require.NoError(t, err)
	requestBlock := h.Number

	// Wait 101 blocks.
	for i := 0; i < 100; i++ {
		cu.Backend.Commit()
	}

	// Wait for the blockhash to be stored
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		cu.Backend.Commit()
		_, err2 := cu.BHSContract.GetBlockhash(&bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: nil,
			Context:     nil,
		}, requestBlock)
		if err2 == nil {
			return true
		} else if strings.Contains(err2.Error(), "execution reverted") {
			return false
		}
		t.Fatal(err2)
		return false
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Wait another 160 blocks so that the request is outside the 256 block window
	for i := 0; i < 160; i++ {
		cu.Backend.Commit()
	}

	// Start the VRF Job and wait until it's processed
	require.NoError(t, app.JobSpawner().CreateJob(ctx, nil, &jb))

	var runs []pipeline.Run
	require.Eventually(t, func() bool {
		runs, err = app.PipelineORM().GetAllRuns(ctx)
		require.NoError(t, err)
		cu.Backend.Commit()
		return len(runs) == 1 && runs[0].State == pipeline.RunStatusCompleted
	}, 10*time.Second, 1*time.Second)
	assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].FatalErrors)
	assert.Equal(t, 4, len(runs[0].PipelineTaskRuns))
	assert.NotNil(t, 0, runs[0].Outputs.Val)

	// stop jobs as to not cause a race condition in geth simulated backend
	// between job creating new tx and fulfillment logs polling below
	require.NoError(t, app.JobSpawner().DeleteJob(ctx, nil, jb.ID))
	require.NoError(t, app.JobSpawner().DeleteJob(ctx, nil, bhsJob.ID))

	// Ensure the eth transaction gets confirmed on chain.
	require.Eventually(t, func() bool {
		orm := txmgr.NewTxStore(app.GetDB(), app.GetLogger())
		uc, err2 := orm.CountUnconfirmedTransactions(ctx, key.Address, testutils.SimulatedChainID)
		require.NoError(t, err2)
		return uc == 0
	}, 5*time.Second, 100*time.Millisecond)

	// Assert the request was fulfilled on-chain.
	require.Eventually(t, func() bool {
		rfIterator, err := cu.RootContract.FilterRandomnessRequestFulfilled(nil)
		require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
		var rf []*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, 5*time.Second, 500*time.Millisecond)
}

func createVRFJobRegisterKey(t *testing.T, u vrftesthelpers.CoordinatorUniverse, app *cltest.TestApplication, incomingConfs int) (job.Job, vrfkey.KeyV2) {
	ctx := testutils.Context(t)
	vrfKey, err := app.KeyStore.VRF().Create(ctx)
	require.NoError(t, err)

	jid := uuid.MustParse("96a8a26f-d426-4784-8d8f-fb387d4d8345")
	expectedOnChainJobID, err := hex.DecodeString("3936613861323666643432363437383438643866666233383764346438333435")
	require.NoError(t, err)
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		CoordinatorAddress:       u.RootContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		PublicKey:                vrfKey.PublicKey.String(),
		EVMChainID:               testutils.SimulatedChainID.String(),
	}).Toml()
	jb, err := vrfcommon.ValidatedVRFSpec(s)
	require.NoError(t, err)
	assert.Equal(t, expectedOnChainJobID, jb.ExternalIDEncodeStringToTopic().Bytes())

	p, err := vrfKey.PublicKey.Point()
	require.NoError(t, err)
	_, err = u.RootContract.RegisterProvingKey(
		u.Neil, big.NewInt(7), u.Neil.From, pair(secp256k1.Coordinates(p)), jb.ExternalIDEncodeStringToTopic())
	require.NoError(t, err)
	u.Backend.Commit()
	return jb, vrfKey
}

func ptr[T any](t T) *T { return &t }

func pair(x, y *big.Int) [2]*big.Int { return [2]*big.Int{x, y} }
