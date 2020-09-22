package pipeline_test

import (
	"context"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
	ormpkg "github.com/smartcontractkit/chainlink/core/store/orm"
)

var jobSpec = []byte(`
contractAddress    = "0x613a38AC1659769640aaE063C651F48E0250454C"
p2pPeerID          = "<libp2p-node-id>"
p2pBootstrapPeers  = [
    {peerID = "<peer id 1>", multiAddr = "<multiaddr1>"},
    {peerID = "<peer id 2>", multiAddr = "<multiaddr2>"},
]
keyBundle          = {encryptedPrivKeyBundle = {asdf = 123}}
monitoringEndpoint = "<ip:port>"
transmitterAddress = "0x613a38AC1659769640aaE063C651F48E0250454C"
observationTimeout = "10s"
blockchainTimeout  = "10s"
contractConfigTrackerPollInterval = "1m"
contractConfigConfirmations = 3
observationSource = """
    // data source 1
    ds1          [type=bridge name=voter_turnout];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\\"hi\\":\\"hello\\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    answer1 [type=median];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer2 [type=bridge name=election_winner];
"""
`)

// `require.Equal` currently has broken handling of `time.Time` values, so we have
// to do equality comparisons of these structs manually.
//
// https://github.com/stretchr/testify/issues/984
// func compareJobSpecs(t *testing.T, expected, actual models.JobSpecV2) {
//     require.Equal(t, expected.ID, actual.ID)
//     require.Equal(t, expected.OffchainreportingOracleSpecID, actual.OffchainreportingOracleSpecID)
//     require.Equal(t, expected.PipelineSpecID, actual.PipelineSpecID)
//     require.Equal(t, expected.OffchainreportingOracleSpec.ID, actual.OffchainreportingOracleSpec.ID)
//     require.Equal(t, expected.OffchainreportingOracleSpec.ContractAddress, actual.OffchainreportingOracleSpec.ContractAddress)
//     require.Equal(t, expected.OffchainreportingOracleSpec.P2PPeerID, actual.OffchainreportingOracleSpec.P2PPeerID)
//     require.Equal(t, expected.OffchainreportingOracleSpec.P2PBootstrapPeers, actual.OffchainreportingOracleSpec.P2PBootstrapPeers)
//     require.Equal(t, expected.OffchainreportingOracleSpec.OffchainreportingKeyBundleID, actual.OffchainreportingOracleSpec.OffchainreportingKeyBundleID)
//     require.Equal(t, expected.OffchainreportingOracleSpec.OffchainreportingKeyBundle.ID, actual.OffchainreportingOracleSpec.OffchainreportingKeyBundle.ID)
//     require.Equal(t, expected.OffchainreportingOracleSpec.OffchainreportingKeyBundle.EncryptedPrivKeyBundle, actual.OffchainreportingOracleSpec.OffchainreportingKeyBundle.EncryptedPrivKeyBundle)
//     require.Equal(t, expected.OffchainreportingOracleSpec.MonitoringEndpoint, actual.OffchainreportingOracleSpec.MonitoringEndpoint)
//     require.Equal(t, expected.OffchainreportingOracleSpec.TransmitterAddress, actual.OffchainreportingOracleSpec.TransmitterAddress)
//     require.Equal(t, expected.OffchainreportingOracleSpec.ObservationTimeout, actual.OffchainreportingOracleSpec.ObservationTimeout)
//     require.Equal(t, expected.OffchainreportingOracleSpec.BlockchainTimeout, actual.OffchainreportingOracleSpec.BlockchainTimeout)
//     require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval, actual.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval)
//     require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigConfirmations, actual.OffchainreportingOracleSpec.ContractConfigConfirmations)
// }

// func TestORM(t *testing.T) {
//     config, cleanup := cltest.NewConfig(t)
//     defer cleanup()

//     db, err := gorm.Open(string(ormpkg.DialectPostgres), config.DatabaseURL())
//     require.NoError(t, err)
//     defer db.Close()

//     orm := job.NewORM(db, config.DatabaseURL())
//     defer orm.Close()

//     var ocrspec offchainreporting.OracleSpec
//     err = toml.Unmarshal(jobSpec, &ocrspec)
//     require.NoError(t, err)

//     spec := models.JobSpecV2{OffchainreportingOracleSpec: &ocrspec.OffchainReportingOracleSpec}
//     pipelineSpec, err := ocrspec.TaskDAG().ToPipelineSpec()
//     require.NoError(t, err)

//     t.Run("it creates job specs", func(t *testing.T) {
//         err := orm.CreateJob(&spec, &pipelineSpec)
//         require.NoError(t, err)

//         var dbSpec models.JobSpecV2
//         err = db.
//             Preload("OffchainreportingOracleSpec").
//             Preload("OffchainreportingOracleSpec.OffchainreportingKeyBundle").
//             Where("id = ?", spec.ID).First(&dbSpec).Error
//         require.NoError(t, err)
//         compareJobSpecs(t, spec, dbSpec)
//     })

//     db2, err := gorm.Open(string(ormpkg.DialectPostgres), config.DatabaseURL())
//     require.NoError(t, err)
//     defer db2.Close()

//     orm2 := job.NewORM(db2, config.DatabaseURL())
//     defer orm2.Close()

//     t.Run("it correctly returns the unclaimed jobs in the DB", func(t *testing.T) {
//         ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel()

//         unclaimed, err := orm.UnclaimedJobs(ctx)
//         require.NoError(t, err)
//         require.Len(t, unclaimed, 1)
//         compareJobSpecs(t, spec, unclaimed[0])

//         ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel2()

//         unclaimed, err = orm2.UnclaimedJobs(ctx2)
//         require.NoError(t, err)
//         require.Len(t, unclaimed, 0)
//     })

//     t.Run("it cannot delete jobs claimed by other nodes", func(t *testing.T) {
//         ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel()

//         err := orm2.DeleteJob(ctx, spec.ID)
//         require.Error(t, err)
//     })

//     t.Run("it deletes its own claimed jobs from the DB", func(t *testing.T) {
//         ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//         defer cancel()

//         err := orm.DeleteJob(ctx, spec.ID)
//         require.NoError(t, err)

//         var dbSpecs []models.JobSpecV2
//         err = db.Find(&dbSpecs).Error
//         require.Len(t, dbSpecs, 0)

//         var oracleSpecs []models.OffchainReportingOracleSpec
//         err = db.Find(&oracleSpecs).Error
//         require.Len(t, oracleSpecs, 0)

//         var pipelineSpecs []pipeline.Spec
//         err = db.Find(&pipelineSpecs).Error
//         require.Len(t, pipelineSpecs, 0)

//         var pipelineTaskSpecs []pipeline.TaskSpec
//         err = db.Find(&pipelineTaskSpecs).Error
//         require.Len(t, pipelineTaskSpecs, 0)
//     })
// }
