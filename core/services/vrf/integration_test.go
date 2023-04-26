package vrf_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/v2/core/chains/evm/txmgr"
	"github.com/smartcontractkit/chainlink/v2/core/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/chainlink"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/v2/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf"
	"github.com/smartcontractkit/chainlink/v2/core/store/models"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
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
			config, _ := heavyweight.FullTestDBV2(t, fmt.Sprintf("vrf_jpv2_%v", test.eip1559), func(c *chainlink.Config, s *chainlink.Secrets) {
				c.EVM[0].GasEstimator.EIP1559DynamicFees = &test.eip1559
			})
			key1 := cltest.MustGenerateRandomKey(t)
			key2 := cltest.MustGenerateRandomKey(t)
			cu := newVRFCoordinatorUniverse(t, key1, key2)
			incomingConfs := 2
			app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, cu.backend, key1, key2)
			require.NoError(t, app.Start(testutils.Context(t)))

			jb, vrfKey := createVRFJobRegisterKey(t, cu, app, incomingConfs)
			require.NoError(t, app.JobSpawner().CreateJob(&jb))

			_, err := cu.consumerContract.TestRequestRandomness(cu.carol,
				vrfKey.PublicKey.MustHash(), big.NewInt(100))
			require.NoError(t, err)

			_, err = cu.consumerContract.TestRequestRandomness(cu.carol,
				vrfKey.PublicKey.MustHash(), big.NewInt(100))
			require.NoError(t, err)
			cu.backend.Commit()
			t.Log("Sent 2 test requests")
			// Mine the required number of blocks
			// So our request gets confirmed.
			for i := 0; i < incomingConfs; i++ {
				cu.backend.Commit()
			}
			var runs []pipeline.Run
			gomega.NewWithT(t).Eventually(func() bool {
				runs, err = app.PipelineORM().GetAllRuns()
				require.NoError(t, err)
				// It possible that we send the test request
				// before the job spawner has started the vrf services, which is fine
				// the lb will backfill the logs. However we need to
				// keep blocks coming in for the lb to send the backfilled logs.
				cu.backend.Commit()
				return len(runs) == 2 && runs[0].State == pipeline.RunStatusCompleted && runs[1].State == pipeline.RunStatusCompleted
			}, testutils.WaitTimeout(t), 1*time.Second).Should(gomega.BeTrue())
			assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].FatalErrors)
			assert.Equal(t, 4, len(runs[0].PipelineTaskRuns))
			assert.Equal(t, 4, len(runs[1].PipelineTaskRuns))
			assert.NotNil(t, 0, runs[0].Outputs.Val)
			assert.NotNil(t, 0, runs[1].Outputs.Val)

			// Ensure the eth transaction gets confirmed on chain.
			gomega.NewWithT(t).Eventually(func() bool {
				orm := txmgr.NewTxStore(app.GetSqlxDB(), app.GetLogger(), app.GetConfig())
				uc, err2 := orm.CountUnconfirmedTransactions(key1.Address, testutils.SimulatedChainID)
				require.NoError(t, err2)
				return uc == 0
			}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

			// Assert the request was fulfilled on-chain.
			var rf []*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled
			gomega.NewWithT(t).Eventually(func() bool {
				rfIterator, err := cu.rootContract.FilterRandomnessRequestFulfilled(nil)
				require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
				rf = nil
				for rfIterator.Next() {
					rf = append(rf, rfIterator.Event)
				}
				return len(rf) == 2
			}, testutils.WaitTimeout(t), 500*time.Millisecond).Should(gomega.BeTrue())

			// Check that each sending address sent one transaction
			n1, err := cu.backend.PendingNonceAt(testutils.Context(t), key1.Address)
			require.NoError(t, err)
			require.EqualValues(t, 1, n1)

			n2, err := cu.backend.PendingNonceAt(testutils.Context(t), key2.Address)
			require.NoError(t, err)
			require.EqualValues(t, 1, n2)
		})
	}
}

func TestIntegration_VRF_WithBHS(t *testing.T) {
	t.Parallel()
	config, _ := heavyweight.FullTestDBV2(t, "vrf_with_bhs", func(c *chainlink.Config, s *chainlink.Secrets) {
		c.EVM[0].GasEstimator.EIP1559DynamicFees = ptr(true)
		c.EVM[0].BlockBackfillDepth = ptr[uint32](500)
		c.Feature.LogPoller = ptr(true)
		c.EVM[0].FinalityDepth = ptr[uint32](2)
		c.EVM[0].LogPollInterval = models.MustNewDuration(time.Second)
	})
	key := cltest.MustGenerateRandomKey(t)
	cu := newVRFCoordinatorUniverse(t, key)
	incomingConfs := 2
	app := cltest.NewApplicationWithConfigV2AndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
	require.NoError(t, app.Start(testutils.Context(t)))

	// Create VRF job but do not start it yet
	jb, vrfKey := createVRFJobRegisterKey(t, cu, app, incomingConfs)

	sendingKeys := []string{key.Address.String()}

	// Create BHS job and start it
	_ = createAndStartBHSJob(t, sendingKeys, app, cu.bhsContractAddress.String(),
		cu.rootContractAddress.String(), "")

	// Ensure log poller is ready and has all logs.
	require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Ready())
	require.NoError(t, app.Chains.EVM.Chains()[0].LogPoller().Replay(testutils.Context(t), 1))

	// Create a VRF request
	_, err := cu.consumerContract.TestRequestRandomness(cu.carol,
		vrfKey.PublicKey.MustHash(), big.NewInt(100))
	require.NoError(t, err)

	cu.backend.Commit()
	requestBlock := cu.backend.Blockchain().CurrentHeader().Number

	// Wait 101 blocks.
	for i := 0; i < 100; i++ {
		cu.backend.Commit()
	}

	// Wait for the blockhash to be stored
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		cu.backend.Commit()
		_, err := cu.bhsContract.GetBlockhash(&bind.CallOpts{
			Pending:     false,
			From:        common.Address{},
			BlockNumber: nil,
			Context:     nil,
		}, requestBlock)
		if err == nil {
			return true
		} else if strings.Contains(err.Error(), "execution reverted") {
			return false
		} else {
			t.Fatal(err)
			return false
		}
	}, testutils.WaitTimeout(t), time.Second).Should(gomega.BeTrue())

	// Wait another 160 blocks so that the request is outside the 256 block window
	for i := 0; i < 160; i++ {
		cu.backend.Commit()
	}

	// Start the VRF job and wait until it's processed
	require.NoError(t, app.JobSpawner().CreateJob(&jb))

	var runs []pipeline.Run
	gomega.NewWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		cu.backend.Commit()
		return len(runs) == 1 && runs[0].State == pipeline.RunStatusCompleted
	}, 10*time.Second, 1*time.Second).Should(gomega.BeTrue())
	assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].FatalErrors)
	assert.Equal(t, 4, len(runs[0].PipelineTaskRuns))
	assert.NotNil(t, 0, runs[0].Outputs.Val)

	// Ensure the eth transaction gets confirmed on chain.
	gomega.NewWithT(t).Eventually(func() bool {
		orm := txmgr.NewTxStore(app.GetSqlxDB(), app.GetLogger(), app.GetConfig())
		uc, err2 := orm.CountUnconfirmedTransactions(key.Address, testutils.SimulatedChainID)
		require.NoError(t, err2)
		return uc == 0
	}, 5*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())

	// Assert the request was fulfilled on-chain.
	gomega.NewWithT(t).Eventually(func() bool {
		rfIterator, err := cu.rootContract.FilterRandomnessRequestFulfilled(nil)
		require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
		var rf []*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, 5*time.Second, 500*time.Millisecond).Should(gomega.BeTrue())
}

func createVRFJobRegisterKey(t *testing.T, u coordinatorUniverse, app *cltest.TestApplication, incomingConfs int) (job.Job, vrfkey.KeyV2) {
	vrfKey, err := app.KeyStore.VRF().Create()
	require.NoError(t, err)

	jid := uuid.FromStringOrNil("96a8a26f-d426-4784-8d8f-fb387d4d8345")
	expectedOnChainJobID, err := hex.DecodeString("3936613861323666643432363437383438643866666233383764346438333435")
	require.NoError(t, err)
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:                    jid.String(),
		Name:                     "vrf-primary",
		CoordinatorAddress:       u.rootContractAddress.String(),
		MinIncomingConfirmations: incomingConfs,
		PublicKey:                vrfKey.PublicKey.String()}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	assert.Equal(t, expectedOnChainJobID, jb.ExternalIDEncodeStringToTopic().Bytes())

	p, err := vrfKey.PublicKey.Point()
	require.NoError(t, err)
	_, err = u.rootContract.RegisterProvingKey(
		u.neil, big.NewInt(7), u.neil.From, pair(secp256k1.Coordinates(p)), jb.ExternalIDEncodeStringToTopic())
	require.NoError(t, err)
	u.backend.Commit()
	return jb, vrfKey
}
