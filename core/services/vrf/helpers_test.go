package vrf_test

import (
	"testing"
	"time"

	"github.com/onsi/gomega"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/vrfkey"
	"github.com/smartcontractkit/chainlink/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
)

func GenerateProofResponseFromProof(p vrfkey.Proof, s proof.PreSeedData) (
	proof.MarshaledOnChainResponse, error) {
	return proof.GenerateProofResponseFromProof(p, s)
}

func createAndStartBHSJob(
	t *testing.T,
	fromAddresses []string,
	app *cltest.TestApplication,
	bhsAddress, coordinatorV1Address, coordinatorV2Address string,
) job.Job {
	jid := uuid.NewV4()
	s := testspecs.GenerateBlockhashStoreSpec(testspecs.BlockhashStoreSpecParams{
		JobID:                 jid.String(),
		Name:                  "blockhash-store",
		CoordinatorV1Address:  coordinatorV1Address,
		CoordinatorV2Address:  coordinatorV2Address,
		WaitBlocks:            100,
		LookbackBlocks:        200,
		BlockhashStoreAddress: bhsAddress,
		PollPeriod:            time.Second,
		RunTimeout:            100 * time.Millisecond,
		EVMChainID:            1337,
		FromAddresses:         fromAddresses,
	})
	jb, err := blockhashstore.ValidatedSpec(s.Toml())
	require.NoError(t, err)

	require.NoError(t, app.JobSpawner().CreateJob(&jb))
	gomega.NewWithT(t).Eventually(func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			if jb.Type == job.BlockhashStore {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond).Should(gomega.BeTrue())

	return jb
}
