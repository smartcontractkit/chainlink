package vrftesthelpers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/keystore/corekeys/vrfkey"
	"github.com/smartcontractkit/chainlink/v2/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockhashstore"
	"github.com/smartcontractkit/chainlink/v2/core/services/blockheaderfeeder"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/vrf/proof"
	"github.com/smartcontractkit/chainlink/v2/core/testdata/testspecs"
)

var (
	WeiPerUnitLink = decimal.RequireFromString("10000000000000000")
)

func GenerateProofResponseFromProof(p vrfkey.Proof, s proof.PreSeedData) (
	proof.MarshaledOnChainResponse, error) {
	return proof.GenerateProofResponseFromProof(p, s)
}

func CreateAndStartBHSJob(
	t *testing.T,
	fromAddresses []string,
	app *cltest.TestApplication,
	bhsAddress, coordinatorV2Address, coordinatorV2PlusAddress string,
	trustedBlockhashStoreAddress string, trustedBlockhashStoreBatchSize int32, lookback int,
	heartbeatPeriod time.Duration, waitBlocks int,
) job.Job {
	jid := uuid.New()
	s := testspecs.GenerateBlockhashStoreSpec(testspecs.BlockhashStoreSpecParams{
		JobID:                          jid.String(),
		Name:                           "blockhash-store",
		CoordinatorV2Address:           coordinatorV2Address,
		CoordinatorV2PlusAddress:       coordinatorV2PlusAddress,
		WaitBlocks:                     waitBlocks,
		LookbackBlocks:                 lookback,
		HeartbeatPeriod:                heartbeatPeriod,
		BlockhashStoreAddress:          bhsAddress,
		TrustedBlockhashStoreAddress:   trustedBlockhashStoreAddress,
		TrustedBlockhashStoreBatchSize: trustedBlockhashStoreBatchSize,
		PollPeriod:                     time.Second,
		RunTimeout:                     10 * time.Second,
		EVMChainID:                     1337,
		FromAddresses:                  fromAddresses,
	})
	jb, err := blockhashstore.ValidatedSpec(s.Toml())
	require.NoError(t, err)

	ctx := testutils.Context(t)
	require.NoError(t, app.JobSpawner().CreateJob(ctx, nil, &jb))
	require.Eventually(t, func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			if jb.Type == job.BlockhashStore {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	return jb
}

func CreateAndStartBlockHeaderFeederJob(
	t *testing.T,
	fromAddresses []string,
	app *cltest.TestApplication,
	bhsAddress, batchBHSAddress, coordinatorV2Address, coordinatorV2PlusAddress string,
) job.Job {
	jid := uuid.New()
	s := testspecs.GenerateBlockHeaderFeederSpec(testspecs.BlockHeaderFeederSpecParams{
		JobID:                      jid.String(),
		Name:                       "block-header-feeder",
		CoordinatorV2Address:       coordinatorV2Address,
		CoordinatorV2PlusAddress:   coordinatorV2PlusAddress,
		WaitBlocks:                 256,
		LookbackBlocks:             1000,
		BlockhashStoreAddress:      bhsAddress,
		BatchBlockhashStoreAddress: batchBHSAddress,
		PollPeriod:                 15 * time.Second,
		RunTimeout:                 15 * time.Second,
		EVMChainID:                 1337,
		FromAddresses:              fromAddresses,
		GetBlockhashesBatchSize:    20,
		StoreBlockhashesBatchSize:  20,
	})
	jb, err := blockheaderfeeder.ValidatedSpec(s.Toml())
	require.NoError(t, err)

	ctx := testutils.Context(t)
	require.NoError(t, app.JobSpawner().CreateJob(ctx, nil, &jb))
	require.Eventually(t, func() bool {
		jbs := app.JobSpawner().ActiveJobs()
		for _, jb := range jbs {
			if jb.Type == job.BlockHeaderFeeder {
				return true
			}
		}
		return false
	}, testutils.WaitTimeout(t), 100*time.Millisecond)

	return jb
}
