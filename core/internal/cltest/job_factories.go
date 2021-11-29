package cltest

import (
	"fmt"
	"testing"

	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/ethkey"
	"github.com/smartcontractkit/chainlink/core/services/keystore/keys/p2pkey"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/sqlx"

	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

const (
	minimalOCRNonBootstrapTemplate = `
			type               = "offchainreporting"
			schemaVersion      = 1
			contractAddress    = "%s"
			p2pPeerID          = "%s"
			p2pBootstrapPeers  = ["/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju"]
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

func MinimalOCRNonBootstrapSpec(contractAddress, transmitterAddress ethkey.EIP55Address, peerID p2pkey.PeerID, keyBundleID string) string {
	return fmt.Sprintf(minimalOCRNonBootstrapTemplate, contractAddress, peerID, transmitterAddress.Hex(), keyBundleID)
}

func MustInsertWebhookSpec(t *testing.T, db *sqlx.DB) (job.Job, job.WebhookSpec) {
	jobORM, pipelineORM := getORMs(t, db)
	webhookSpec := job.WebhookSpec{}
	require.NoError(t, jobORM.InsertWebhookSpec(&webhookSpec))

	pSpec := pipeline.Pipeline{}
	pipelineSpecID, err := pipelineORM.CreateSpec(pSpec, 0)
	require.NoError(t, err)

	job := job.Job{WebhookSpecID: &webhookSpec.ID, WebhookSpec: &webhookSpec, SchemaVersion: 1, Type: "webhook", ExternalJobID: uuid.NewV4(), PipelineSpecID: pipelineSpecID}
	require.NoError(t, jobORM.InsertJob(&job))

	return job, webhookSpec
}

func getORMs(t *testing.T, db *sqlx.DB) (jobORM job.ORM, pipelineORM pipeline.ORM) {
	config := NewTestGeneralConfig(t)
	keyStore := NewKeyStore(t, db, config)
	pipelineORM = pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM = job.NewORM(db, cc, pipelineORM, keyStore, logger.TestLogger(t), config)
	t.Cleanup(func() { jobORM.Close() })
	return
}
