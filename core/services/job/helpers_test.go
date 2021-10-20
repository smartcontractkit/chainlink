package job_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

const (
	ocrJobSpecTemplate = `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "%s"
p2pPeerID          = "%s"
p2pBootstrapPeers  = [
    "/dns4/chain.link/tcp/1234/p2p/16Uiu2HAm58SP7UL8zsnpeuwHfytLocaqgnyaYKP8wu7qRdrixLju",
]
isBootstrapPeer    = false
keyBundleID        = "%s"
monitoringEndpoint = "chain.link:4321"
transmitterAddress = "%s"
observationTimeout = "10s"
blockchainTimeout  = "20s"
contractConfigTrackerSubscribeInterval = "2m"
contractConfigTrackerPollInterval = "1m"
contractConfigConfirmations = 3
observationSource = """
	%s
"""
`
	voterTurnoutDataSourceTemplate = `
// data source 1
ds1          [type=bridge name=voter_turnout];
ds1_parse    [type=jsonparse path="data,result"];
ds1_multiply [type=multiply times=100];

// data source 2
ds2          [type=http method=POST url="%s" requestData="{\\"hi\\": \\"hello\\"}"];
ds2_parse    [type=jsonparse path="turnout"];
ds2_multiply [type=multiply times=100];

ds1 -> ds1_parse -> ds1_multiply -> answer1;
ds2 -> ds2_parse -> ds2_multiply -> answer1;

answer1 [type=median                      index=0];
answer2 [type=bridge name=election_winner index=1];
`

	simpleFetchDataSourceTemplate = `
// data source 1
ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true"];
ds1_parse    [type=jsonparse path="USD" lax=%t];
ds1_multiply [type=multiply times=100];
ds1 -> ds1_parse -> ds1_multiply;
`
	minimalNonBootstrapTemplate = `
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
ds1          [type=http method=GET url="%s" allowunrestrictednetworkaccess="true" %s];
ds1_parse    [type=jsonparse path="USD" lax=true];
ds1 -> ds1_parse;
"""
`
	minimalBootstrapTemplate = `
		type               = "offchainreporting"
		schemaVersion      = 1
		contractAddress    = "%s"
		p2pPeerID          = "%s"
		p2pBootstrapPeers  = []
		isBootstrapPeer    = true
`
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
monitoringEndpoint = "chain.link:4321"
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
)

func makeOCRJobSpec(t *testing.T, transmitterAddress common.Address) *job.Job {
	t.Helper()

	peerID := cltest.DefaultP2PPeerID
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	jobSpecText := fmt.Sprintf(ocrJobSpecText, cltest.NewAddress().Hex(), peerID.String(), ocrKeyID, transmitterAddress.Hex())

	dbSpec := job.Job{
		ExternalJobID: uuid.NewV4(),
	}
	err := toml.Unmarshal([]byte(jobSpecText), &dbSpec)
	require.NoError(t, err)
	var ocrspec job.OffchainReportingOracleSpec
	err = toml.Unmarshal([]byte(jobSpecText), &ocrspec)
	require.NoError(t, err)
	dbSpec.OffchainreportingOracleSpec = &ocrspec

	return &dbSpec
}

// `require.Equal` currently has broken handling of `time.Time` values, so we have
// to do equality comparisons of these structs manually.
//
// https://github.com/stretchr/testify/issues/984
func compareOCRJobSpecs(t *testing.T, expected, actual job.Job) {
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

func makeMinimalHTTPOracleSpec(t *testing.T, db *gorm.DB, cfg config.GeneralConfig, contractAddress, peerID, transmitterAddress, keyBundle, fetchUrl, timeout string) *job.Job {
	var ocrSpec = job.OffchainReportingOracleSpec{
		P2PBootstrapPeers:                      pq.StringArray{},
		ObservationTimeout:                     models.Interval(10 * time.Second),
		BlockchainTimeout:                      models.Interval(20 * time.Second),
		ContractConfigTrackerSubscribeInterval: models.Interval(2 * time.Minute),
		ContractConfigTrackerPollInterval:      models.Interval(1 * time.Minute),
		ContractConfigConfirmations:            uint16(3),
	}
	var os = job.Job{
		Name:          null.NewString("a job", true),
		Type:          job.OffchainReporting,
		SchemaVersion: 1,
		ExternalJobID: uuid.NewV4(),
	}
	s := fmt.Sprintf(minimalNonBootstrapTemplate, contractAddress, peerID, transmitterAddress, keyBundle, fetchUrl, timeout)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: cltest.NewEthClientMockWithDefaultChain(t), GeneralConfig: cfg})
	_, err := offchainreporting.ValidatedOracleSpecToml(cc, s)
	require.NoError(t, err)
	err = toml.Unmarshal([]byte(s), &os)
	require.NoError(t, err)
	err = toml.Unmarshal([]byte(s), &ocrSpec)
	require.NoError(t, err)
	os.OffchainreportingOracleSpec = &ocrSpec
	return &os
}

func makeVoterTurnoutOCRJobSpec(t *testing.T, db *gorm.DB, transmitterAddress common.Address) *job.Job {
	t.Helper()
	return MakeVoterTurnoutOCRJobSpecWithHTTPURL(t, db, transmitterAddress, "https://example.com/foo/bar")
}

func MakeVoterTurnoutOCRJobSpecWithHTTPURL(t *testing.T, db *gorm.DB, transmitterAddress common.Address, httpURL string) *job.Job {
	t.Helper()
	peerID := cltest.DefaultP2PPeerID
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	ds := fmt.Sprintf(voterTurnoutDataSourceTemplate, httpURL)
	voterTurnoutJobSpec := fmt.Sprintf(ocrJobSpecTemplate, cltest.NewAddress().Hex(), peerID, ocrKeyID, transmitterAddress.Hex(), ds)
	return makeOCRJobSpecFromToml(t, db, voterTurnoutJobSpec)
}

func makeSimpleFetchOCRJobSpecWithHTTPURL(t *testing.T, db *gorm.DB, transmitterAddress common.Address, httpURL string, lax bool) *job.Job {
	t.Helper()
	peerID := cltest.DefaultP2PPeerID
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	ds := fmt.Sprintf(simpleFetchDataSourceTemplate, httpURL, lax)
	simpleFetchJobSpec := fmt.Sprintf(ocrJobSpecTemplate, cltest.NewAddress().Hex(), peerID, ocrKeyID, transmitterAddress.Hex(), ds)
	return makeOCRJobSpecFromToml(t, db, simpleFetchJobSpec)
}

func makeOCRJobSpecFromToml(t *testing.T, db *gorm.DB, jobSpecToml string) *job.Job {
	t.Helper()

	id := uuid.NewV4()
	var jb = job.Job{
		Name:          null.StringFrom(id.String()),
		ExternalJobID: id,
	}
	err := toml.Unmarshal([]byte(jobSpecToml), &jb)
	require.NoError(t, err)
	var ocrspec job.OffchainReportingOracleSpec
	err = toml.Unmarshal([]byte(jobSpecToml), &ocrspec)
	require.NoError(t, err)
	jb.OffchainreportingOracleSpec = &ocrspec

	return &jb
}
