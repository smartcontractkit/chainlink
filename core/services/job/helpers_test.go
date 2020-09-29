package job_test

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var ocrJobSpecText = `
contractAddress    = "%s"
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
`

func makeOCRJobSpec(t *testing.T) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()

	jobSpecText := fmt.Sprintf(ocrJobSpecText, cltest.NewAddress().Hex())

	var ocrspec offchainreporting.OracleSpec
	err := toml.Unmarshal([]byte(jobSpecText), &ocrspec)
	require.NoError(t, err)

	dbSpec := models.JobSpecV2{OffchainreportingOracleSpec: &ocrspec.OffchainReportingOracleSpec}
	return &ocrspec, &dbSpec
}

// `require.Equal` currently has broken handling of `time.Time` values, so we have
// to do equality comparisons of these structs manually.
//
// https://github.com/stretchr/testify/issues/984
func compareOCRJobSpecs(t *testing.T, expected, actual models.JobSpecV2) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.OffchainreportingOracleSpecID, actual.OffchainreportingOracleSpecID)
	require.Equal(t, expected.PipelineSpecID, actual.PipelineSpecID)
	require.Equal(t, expected.OffchainreportingOracleSpec.ID, actual.OffchainreportingOracleSpec.ID)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractAddress, actual.OffchainreportingOracleSpec.ContractAddress)
	require.Equal(t, expected.OffchainreportingOracleSpec.P2PPeerID, actual.OffchainreportingOracleSpec.P2PPeerID)
	require.Equal(t, expected.OffchainreportingOracleSpec.P2PBootstrapPeers, actual.OffchainreportingOracleSpec.P2PBootstrapPeers)
	require.Equal(t, expected.OffchainreportingOracleSpec.OffchainreportingKeyBundleID, actual.OffchainreportingOracleSpec.OffchainreportingKeyBundleID)
	require.Equal(t, expected.OffchainreportingOracleSpec.OffchainreportingKeyBundle.ID, actual.OffchainreportingOracleSpec.OffchainreportingKeyBundle.ID)
	require.Equal(t, expected.OffchainreportingOracleSpec.OffchainreportingKeyBundle.EncryptedPrivKeyBundle, actual.OffchainreportingOracleSpec.OffchainreportingKeyBundle.EncryptedPrivKeyBundle)
	require.Equal(t, expected.OffchainreportingOracleSpec.MonitoringEndpoint, actual.OffchainreportingOracleSpec.MonitoringEndpoint)
	require.Equal(t, expected.OffchainreportingOracleSpec.TransmitterAddress, actual.OffchainreportingOracleSpec.TransmitterAddress)
	require.Equal(t, expected.OffchainreportingOracleSpec.ObservationTimeout, actual.OffchainreportingOracleSpec.ObservationTimeout)
	require.Equal(t, expected.OffchainreportingOracleSpec.BlockchainTimeout, actual.OffchainreportingOracleSpec.BlockchainTimeout)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval, actual.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigConfirmations, actual.OffchainreportingOracleSpec.ContractConfigConfirmations)
}
