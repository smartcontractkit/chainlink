package vrf_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/gethwrappers/generated/solidity_vrf_coordinator_interface"

	"github.com/onsi/gomega"
	"github.com/smartcontractkit/chainlink/core/services/bulletprooftxmanager"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/stretchr/testify/assert"
	"gopkg.in/guregu/null.v4"

	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/stretchr/testify/require"
)

func TestIntegration_VRFV2(t *testing.T) {
	config, _, cleanupDB := cltest.BootstrapThrowawayORM(t, "vrf_v2", true)
	defer cleanupDB()
	key := cltest.MustGenerateRandomKey(t)
	cu := newVRFCoordinatorUniverse(t, key)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())

	vrfkey, err := app.Store.VRFKeyStore.CreateKey(cltest.Password)
	require.NoError(t, err)
	unlocked, err := app.Store.VRFKeyStore.Unlock(cltest.Password)
	require.NoError(t, err)
	jid := uuid.NewV4()
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:              jid.String(),
		Name:               "vrf-primary",
		CoordinatorAddress: cu.rootContractAddress.String(),
		Confirmations:      incomingConfs,
		PublicKey:          unlocked[0].String()}).Toml()
	jb, err := vrf.ValidateVRFSpec(s)
	require.NoError(t, err)
	require.NoError(t, app.GetJobORM().CreateJob(context.Background(), &jb, jb.Pipeline))

	p, err := vrfkey.Point()
	require.NoError(t, err)
	_, err = cu.rootContract.RegisterProvingKey(
		cu.neil, big.NewInt(7), cu.neil.From, pair(secp256k1.Coordinates(p)), jb.ExternalIDToTopicHash())
	require.NoError(t, err)
	cu.backend.Commit()
	_, err = cu.consumerContract.TestRequestRandomness(cu.carol,
		vrfkey.MustHash(), big.NewInt(100), big.NewInt(1))
	require.NoError(t, err)
	cu.backend.Commit()
	// Mine the required number of blocks
	// So our request gets confirmed.
	for i := 0; i < incomingConfs; i++ {
		cu.backend.Commit()
	}
	var runs []pipeline.Run
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM.GetAllRuns()
		require.NoError(t, err)
		return len(runs) == 1
	}, 5*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())
	assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].Errors)
	assert.Equal(t, 0, len(runs[0].PipelineTaskRuns))
	assert.NotNil(t, 0, runs[0].Outputs.Val)

	// Ensure the eth transaction gets confirmed on chain.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		uc, err2 := bulletprooftxmanager.CountUnconfirmedTransactions(app.Store.DB, key.Address.Address())
		require.NoError(t, err2)
		return uc == 0
	}, 5*time.Second, 100*time.Millisecond).Should(gomega.BeTrue())

	// Assert the request was fulfilled on-chain.
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		rfIterator, err := cu.rootContract.FilterRandomnessRequestFulfilled(nil)
		require.NoError(t, err, "failed to subscribe to RandomnessRequest logs")
		var rf []*solidity_vrf_coordinator_interface.VRFCoordinatorRandomnessRequestFulfilled
		for rfIterator.Next() {
			rf = append(rf, rfIterator.Event)
		}
		return len(rf) == 1
	}, 5*time.Second, 500*time.Millisecond).Should(gomega.BeTrue())
}
