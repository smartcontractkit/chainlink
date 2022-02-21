package job_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/lib/pq"
	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/config"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/evmtest"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/ocr"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/sqlx"
)

const (
	ocrJobSpecTemplate = `
type               = "offchainreporting"
schemaVersion      = 1
contractAddress    = "%s"
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
ds1          [type=bridge name="%s"];
ds1_parse    [type=jsonparse path="data,result"];
ds1_multiply [type=multiply times=100];

// data source 2
ds2          [type=http method=POST url="%s" requestData="{\\"hi\\": \\"hello\\"}"];
ds2_parse    [type=jsonparse path="turnout"];
ds2_multiply [type=multiply times=100];

ds1 -> ds1_parse -> ds1_multiply -> answer1;
ds2 -> ds2_parse -> ds2_multiply -> answer1;

answer1 [type=median                      index=0];
answer2 [type=bridge name="%s" index=1];
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
    ds1          [type=bridge name="%s"];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\\"hi\\": \\"hello\\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name="%s" index=1];
"""
`
)

func makeOCRJobSpec(t *testing.T, transmitterAddress common.Address, b1, b2 string) *job.Job {
	t.Helper()

	peerID := cltest.DefaultP2PPeerID
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	jobSpecText := fmt.Sprintf(ocrJobSpecText, testutils.NewAddress().Hex(), peerID, ocrKeyID, transmitterAddress.Hex(), b1, b2)

	dbSpec := job.Job{
		ExternalJobID: uuid.NewV4(),
	}
	err := toml.Unmarshal([]byte(jobSpecText), &dbSpec)
	require.NoError(t, err)
	var ocrspec job.OCROracleSpec
	err = toml.Unmarshal([]byte(jobSpecText), &ocrspec)
	require.NoError(t, err)
	dbSpec.OCROracleSpec = &ocrspec

	return &dbSpec
}

// `require.Equal` currently has broken handling of `time.Time` values, so we have
// to do equality comparisons of these structs manually.
//
// https://github.com/stretchr/testify/issues/984
func compareOCRJobSpecs(t *testing.T, expected, actual job.Job) {
	require.NotNil(t, expected.OCROracleSpec)
	require.Equal(t, expected.OCROracleSpec.ContractAddress, actual.OCROracleSpec.ContractAddress)
	require.Equal(t, expected.OCROracleSpec.P2PBootstrapPeers, actual.OCROracleSpec.P2PBootstrapPeers)
	require.Equal(t, expected.OCROracleSpec.IsBootstrapPeer, actual.OCROracleSpec.IsBootstrapPeer)
	require.Equal(t, expected.OCROracleSpec.EncryptedOCRKeyBundleID, actual.OCROracleSpec.EncryptedOCRKeyBundleID)
	require.Equal(t, expected.OCROracleSpec.TransmitterAddress, actual.OCROracleSpec.TransmitterAddress)
	require.Equal(t, expected.OCROracleSpec.ObservationTimeout, actual.OCROracleSpec.ObservationTimeout)
	require.Equal(t, expected.OCROracleSpec.BlockchainTimeout, actual.OCROracleSpec.BlockchainTimeout)
	require.Equal(t, expected.OCROracleSpec.ContractConfigTrackerSubscribeInterval, actual.OCROracleSpec.ContractConfigTrackerSubscribeInterval)
	require.Equal(t, expected.OCROracleSpec.ContractConfigTrackerPollInterval, actual.OCROracleSpec.ContractConfigTrackerPollInterval)
	require.Equal(t, expected.OCROracleSpec.ContractConfigConfirmations, actual.OCROracleSpec.ContractConfigConfirmations)
}

func makeMinimalHTTPOracleSpec(t *testing.T, db *sqlx.DB, cfg config.GeneralConfig, contractAddress, transmitterAddress, keyBundle, fetchUrl, timeout string) *job.Job {
	var ocrSpec = job.OCROracleSpec{
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
	s := fmt.Sprintf(minimalNonBootstrapTemplate, contractAddress, transmitterAddress, keyBundle, fetchUrl, timeout)
	cc := evmtest.NewChainSet(t, evmtest.TestChainOpts{DB: db, Client: cltest.NewEthClientMockWithDefaultChain(t), GeneralConfig: cfg})
	_, err := ocr.ValidatedOracleSpecToml(cc, s)
	require.NoError(t, err)
	err = toml.Unmarshal([]byte(s), &os)
	require.NoError(t, err)
	err = toml.Unmarshal([]byte(s), &ocrSpec)
	require.NoError(t, err)
	os.OCROracleSpec = &ocrSpec
	return &os
}

func makeVoterTurnoutOCRJobSpec(t *testing.T, transmitterAddress common.Address, b1, b2 string) *job.Job {
	t.Helper()
	return MakeVoterTurnoutOCRJobSpecWithHTTPURL(t, transmitterAddress, "https://example.com/foo/bar", b1, b2)
}

func MakeVoterTurnoutOCRJobSpecWithHTTPURL(t *testing.T, transmitterAddress common.Address, httpURL, b1, b2 string) *job.Job {
	t.Helper()
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	ds := fmt.Sprintf(voterTurnoutDataSourceTemplate, b1, httpURL, b2)
	voterTurnoutJobSpec := fmt.Sprintf(ocrJobSpecTemplate, testutils.NewAddress().Hex(), ocrKeyID, transmitterAddress.Hex(), ds)
	return makeOCRJobSpecFromToml(t, voterTurnoutJobSpec)
}

func makeSimpleFetchOCRJobSpecWithHTTPURL(t *testing.T, transmitterAddress common.Address, httpURL string, lax bool) *job.Job {
	t.Helper()
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	ds := fmt.Sprintf(simpleFetchDataSourceTemplate, httpURL, lax)
	simpleFetchJobSpec := fmt.Sprintf(ocrJobSpecTemplate, testutils.NewAddress().Hex(), ocrKeyID, transmitterAddress.Hex(), ds)
	return makeOCRJobSpecFromToml(t, simpleFetchJobSpec)
}

func makeOCRJobSpecFromToml(t *testing.T, jobSpecToml string) *job.Job {
	t.Helper()

	id := uuid.NewV4()
	var jb = job.Job{
		Name:          null.StringFrom(id.String()),
		ExternalJobID: id,
	}
	err := toml.Unmarshal([]byte(jobSpecToml), &jb)
	require.NoError(t, err)
	var ocrspec job.OCROracleSpec
	err = toml.Unmarshal([]byte(jobSpecToml), &ocrspec)
	require.NoError(t, err)
	jb.OCROracleSpec = &ocrspec

	return &jb
}
