package vrf_test

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"
	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/pg"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
)

func TestIntegration_VRF_JPV2(t *testing.T) {
	tests := []struct {
		name    string
		eip1559 bool
	}{
		{"legacy mode", false},
		{"eip1559 mode", true},
	}

	for _, tt := range tests {
		test := tt
		t.Run(test.name, func(t *testing.T) {
			config, _ := heavyweight.FullTestDB(t, fmt.Sprintf("vrf_jpv2_%v", test.eip1559), true, true)
			config.Overrides.GlobalEvmEIP1559DynamicFees = null.BoolFrom(test.eip1559)
			key := cltest.MustGenerateRandomKey(t)
			cu := newVRFCoordinatorUniverse(t, key)
			app := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
			require.NoError(t, app.Start())

			vrfkey, err := app.KeyStore.VRF().Create()
			require.NoError(t, err)
			// Let's use a real onchain job ID to ensure it'll work with
			// existing contract state.
			jid := uuid.FromStringOrNil("96a8a26f-d426-4784-8d8f-fb387d4d8345")
			expectedOnChainJobID, err := hex.DecodeString("3936613861323666643432363437383438643866666233383764346438333435")
			require.NoError(t, err)
			incomingConfs := 2
			s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
				JobID:                    jid.String(),
				Name:                     "vrf-primary",
				CoordinatorAddress:       cu.rootContractAddress.String(),
				MinIncomingConfirmations: incomingConfs,
				PublicKey:                vrfkey.PublicKey.String()}).Toml()
			jb, err := vrf.ValidatedVRFSpec(s)
			require.NoError(t, err)
			assert.Equal(t, expectedOnChainJobID, jb.ExternalIDEncodeStringToTopic().Bytes())
			err = app.JobSpawner().CreateJob(&jb)
			require.NoError(t, err)

			p, err := vrfkey.PublicKey.Point()
			require.NoError(t, err)
			_, err = cu.rootContract.RegisterProvingKey(
				cu.neil, big.NewInt(7), cu.neil.From, pair(secp256k1.Coordinates(p)), jb.ExternalIDEncodeStringToTopic())
			require.NoError(t, err)
			cu.backend.Commit()
			_, err = cu.consumerContract.TestRequestRandomness(cu.carol,
				vrfkey.PublicKey.MustHash(), big.NewInt(100))
			require.NoError(t, err)
			cu.backend.Commit()
			t.Log("Sent test request")
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
				return len(runs) == 1 && runs[0].State == pipeline.RunStatusCompleted
			}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())
			assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].FatalErrors)
			assert.Equal(t, 4, len(runs[0].PipelineTaskRuns))
			assert.NotNil(t, 0, runs[0].Outputs.Val)

			// Ensure the eth transaction gets confirmed on chain.
			gomega.NewWithT(t).Eventually(func() bool {
				q := pg.NewQ(app.GetSqlxDB(), app.GetLogger(), app.GetConfig())
				uc, err2 := bulletprooftxmanager.CountUnconfirmedTransactions(q, key.Address.Address(), cltest.FixtureChainID)
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
		})
	}
}
