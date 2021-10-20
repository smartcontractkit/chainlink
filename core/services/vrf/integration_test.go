package vrf_test

import (
	"context"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/internal/cltest/heavyweight"

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
	config, _, cleanupDB := heavyweight.FullTestORM(t, "vrf_v2", true)
	defer cleanupDB()
	key := cltest.MustGenerateRandomKey(t)
	cu := newVRFCoordinatorUniverse(t, key)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
	defer cleanup()
	require.NoError(t, app.StartAndConnect())

	_, err := app.KeyStore.VRF().Unlock(cltest.Password)
	require.NoError(t, err)
	vrfkey, err := app.KeyStore.VRF().CreateKey()
	require.NoError(t, err)
	// Let's use a real onchain job ID to ensure it'll work with
	// existing contract state.
	jid := uuid.FromStringOrNil("96a8a26f-d426-4784-8d8f-fb387d4d8345")
	expectedOnChainJobID, err := hex.DecodeString("3936613861323666643432363437383438643866666233383764346438333435")
	require.NoError(t, err)
	incomingConfs := 2
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:              jid.String(),
		Name:               "vrf-primary",
		CoordinatorAddress: cu.rootContractAddress.String(),
		Confirmations:      incomingConfs,
		PublicKey:          vrfkey.String()}).Toml()
	jb, err := vrf.ValidatedVRFSpec(s)
	require.NoError(t, err)
	assert.Equal(t, expectedOnChainJobID, jb.ExternalIDEncodeStringToTopic().Bytes())
	jb, err = app.JobORM().CreateJob(context.Background(), &jb, jb.Pipeline)
	require.NoError(t, err)

	p, err := vrfkey.Point()
	require.NoError(t, err)
	_, err = cu.rootContract.RegisterProvingKey(
		cu.neil, big.NewInt(7), cu.neil.From, pair(secp256k1.Coordinates(p)), jb.ExternalIDEncodeStringToTopic())
	require.NoError(t, err)
	cu.backend.Commit()
	_, err = cu.consumerContract.TestRequestRandomness(cu.carol,
		vrfkey.MustHash(), big.NewInt(100))
	require.NoError(t, err)
	cu.backend.Commit()
	t.Log("Sent test request")
	// Mine the required number of blocks
	// So our request gets confirmed.
	for i := 0; i < incomingConfs; i++ {
		cu.backend.Commit()
	}
	var runs []pipeline.Run
	gomega.NewGomegaWithT(t).Eventually(func() bool {
		runs, err = app.PipelineORM().GetAllRuns()
		require.NoError(t, err)
		// It possible that we send the test request
		// before the job spawner has started the vrf services, which is fine
		// the lb will backfill the logs. However we need to
		// keep blocks coming in for the lb to send the backfilled logs.
		cu.backend.Commit()
		return len(runs) == 1
	}, 5*time.Second, 1*time.Second).Should(gomega.BeTrue())
	assert.Equal(t, pipeline.RunErrors([]null.String{{}}), runs[0].Errors)
	assert.Equal(t, 4, len(runs[0].PipelineTaskRuns))
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
