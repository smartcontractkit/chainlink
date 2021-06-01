package vrf_test

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/signatures/secp256k1"
	"github.com/smartcontractkit/chainlink/core/services/vrf"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
	"time"
)

func TestIntegration_VRFV2(t *testing.T) {
	config, _, cleanup := cltest.BootstrapThrowawayORM(t, "vrfv2", true)
	defer cleanup()
	key := cltest.MustGenerateRandomKey(t)
	cu := newVRFCoordinatorUniverse(t, key)
	app, cleanup := cltest.NewApplicationWithConfigAndKeyOnSimulatedBlockchain(t, config, cu.backend, key)
	defer cleanup()
	app.Start()

	vrfkey, err := app.Store.VRFKeyStore.CreateKey(cltest.Password)
	require.NoError(t, err)
	unlocked, err := app.Store.VRFKeyStore.Unlock(cltest.Password)
	require.NoError(t, err)
	t.Log("unlocked", unlocked)
	jid := uuid.NewV4()
	s := testspecs.GenerateVRFSpec(testspecs.VRFSpecParams{
		JobID:              jid.String(),
		Name:               "vrf-primary",
		CoordinatorAddress: cu.rootContractAddress.String(),
		Confirmations:      2,
		PublicKey:          unlocked[0].String()}).Toml()
	jb, _ := vrf.ValidateVRFSpec(s)
	require.NoError(t, app.JobORM.CreateJob(context.Background(), &jb, jb.Pipeline))

	p, err := vrfkey.Point()
	require.NoError(t, err)
	_, err = cu.rootContract.RegisterProvingKey(
		cu.neil, big.NewInt(7), cu.neil.From, pair(secp256k1.Coordinates(p)), jb.ExternalIDToTopicHash())
	require.NoError(t, err)
	_, err = cu.consumerContract.TestRequestRandomness(cu.carol,
		vrfkey.MustHash(), big.NewInt(100), big.NewInt(1))
	require.NoError(t, err, "problem during initial VRF randomness request")
	cu.backend.Commit()
	// We should mine blocks until we see a run
	for i := 0; i < 3; i++ {
		cu.backend.Commit()
	}
	time.Sleep(5 * time.Second)
}
