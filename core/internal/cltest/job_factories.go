package cltest

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink-common/pkg/sqlutil"
	"github.com/smartcontractkit/chainlink-evm/pkg/types"
	"github.com/smartcontractkit/chainlink/v2/core/bridges"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils/configtest"
	"github.com/smartcontractkit/chainlink/v2/core/logger"
	"github.com/smartcontractkit/chainlink/v2/core/services/job"
	"github.com/smartcontractkit/chainlink/v2/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/v2/core/services/pipeline"
)

const (
	minimalOCRNonBootstrapTemplate = `
			type               = "offchainreporting"
			schemaVersion      = 1
			contractAddress    = "%s"
			evmChainID		   = "0"
			p2pPeerID          = "%s"
			p2pv2Bootstrappers = ["12D3KooWHfYFQ8hGttAYbMCevQVESEQhzJAqFZokMVtom8bNxwGq@127.0.0.1:5001"]
			isBootstrapPeer    = false
			transmitterAddress = "%s"
			keyBundleID = "%s"
			observationTimeout = "10s"
			observationSource = """
	ds1          [type=http method=GET url="http://data.com"];
	ds1_parse    [type=jsonparse path="USD" lax=true];
	ds1 -> ds1_parse;
	"""
	`
)

func MinimalOCRNonBootstrapSpec(contractAddress, transmitterAddress types.EIP55Address, peerID p2pkey.PeerID, keyBundleID string) string {
	return fmt.Sprintf(minimalOCRNonBootstrapTemplate, contractAddress, peerID, transmitterAddress.Hex(), keyBundleID)
}

func MustInsertWebhookSpec(t *testing.T, ds sqlutil.DataSource) (job.Job, job.WebhookSpec) {
	ctx := testutils.Context(t)
	jobORM, pipelineORM := getORMs(t, ds)
	webhookSpec := job.WebhookSpec{}
	require.NoError(t, jobORM.InsertWebhookSpec(ctx, &webhookSpec))

	pSpec := pipeline.Pipeline{}
	pipelineSpecID, err := pipelineORM.CreateSpec(ctx, pSpec, 0)
	require.NoError(t, err)

	createdJob := job.Job{WebhookSpecID: &webhookSpec.ID, WebhookSpec: &webhookSpec, SchemaVersion: 1, Type: "webhook",
		ExternalJobID: uuid.New(), PipelineSpecID: pipelineSpecID}
	require.NoError(t, jobORM.InsertJob(ctx, &createdJob))

	return createdJob, webhookSpec
}

func getORMs(t *testing.T, ds sqlutil.DataSource) (jobORM job.ORM, pipelineORM pipeline.ORM) {
	config := configtest.NewTestGeneralConfig(t)
	keyStore := NewKeyStore(t, ds)
	lggr := logger.TestLogger(t)
	pipelineORM = pipeline.NewORM(ds, lggr, config.JobPipeline().MaxSuccessfulRuns())
	bridgeORM := bridges.NewORM(ds)
	jobORM = job.NewORM(ds, pipelineORM, bridgeORM, keyStore, lggr)
	t.Cleanup(func() { jobORM.Close() })
	return
}
