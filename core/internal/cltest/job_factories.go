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
	ocrJobSpecText = `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "%s"
p2pPeerID          = "%s"
p2pBootstrapPeers  = [
    "/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "%s"
transmitterAddress = "%s"
observationTimeout = "10s"
blockchainTimeout  = "20s"
contractConfigTrackerSubscribeInterval = "2m"
contractConfigTrackerPollInterval = "1m"
contractConfigConfirmations = 3
observationSource = """
    // data source 1
    ds1          [type=bridge name=voter_turnout];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\\"hi\\": \\"hello\\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name=election_winner index=1];
"""
`
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

// `require.Equal` currently has broken handling of `time.Time` values, so we have
// to do equality comparisons of these structs manually.
//
// https://github.com/stretchr/testify/issues/984
func CompareOCRJobSpecs(t *testing.T, expected, actual job.Job) {
	t.Helper()
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractAddress, actual.OffchainreportingOracleSpec.ContractAddress)
	require.Equal(t, expected.OffchainreportingOracleSpec.P2PPeerID, actual.OffchainreportingOracleSpec.P2PPeerID)
	require.Equal(t, expected.OffchainreportingOracleSpec.P2PBootstrapPeers, actual.OffchainreportingOracleSpec.P2PBootstrapPeers)
	require.Equal(t, expected.OffchainreportingOracleSpec.IsBootstrapPeer, actual.OffchainreportingOracleSpec.IsBootstrapPeer)
	require.Equal(t, expected.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID, actual.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
	require.Equal(t, expected.OffchainreportingOracleSpec.TransmitterAddress, actual.OffchainreportingOracleSpec.TransmitterAddress)
	require.Equal(t, expected.OffchainreportingOracleSpec.ObservationTimeout, actual.OffchainreportingOracleSpec.ObservationTimeout)
	require.Equal(t, expected.OffchainreportingOracleSpec.BlockchainTimeout, actual.OffchainreportingOracleSpec.BlockchainTimeout)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, actual.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval, actual.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigConfirmations, actual.OffchainreportingOracleSpec.ContractConfigConfirmations)
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
	keyStore := NewKeyStore(t, db)
	pipelineORM = pipeline.NewORM(db, logger.TestLogger(t), config)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, GeneralConfig: config})
	jobORM = job.NewORM(db, cc, pipelineORM, keyStore, logger.TestLogger(t), config)
	t.Cleanup(func() { jobORM.Close() })
	return
}
