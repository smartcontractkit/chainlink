package job_test

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

var ocrJobSpecText = `
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

func makeOCRJobSpec(t *testing.T, db *gorm.DB) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()

	// Insert keys into the store
	keystore := offchainreporting.NewKeyStore(db)
	t.Fatal("FIXME: Don't do this, use the fixture keys instead")
	p2pkey, _, err := keystore.GenerateEncryptedP2PKey("password")
	require.NoError(t, err)
	ocrkey, _, err := keystore.GenerateEncryptedOCRKeyBundle("password")
	require.NoError(t, err)
	peerID, err := p2pkey.GetPeerID()
	require.NoError(t, err)

	jobSpecText := fmt.Sprintf(ocrJobSpecText, cltest.NewAddress().Hex(), peer.ID(peerID), ocrkey.ID, cltest.DefaultKey)

	var ocrspec offchainreporting.OracleSpec
	err = toml.Unmarshal([]byte(jobSpecText), &ocrspec)
	require.NoError(t, err)

	dbSpec := models.JobSpecV2{OffchainreportingOracleSpec: &ocrspec.OffchainReportingOracleSpec}
	return &ocrspec, &dbSpec
}

// `require.Equal` currently has broken handling of `time.Time` values, so we have
// to do equality comparisons of these structs manually.
//
// https://github.com/stretchr/testify/issues/984
func compareOCRJobSpecs(t *testing.T, expected, actual models.JobSpecV2) {
	t.Helper()
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractAddress, actual.OffchainreportingOracleSpec.ContractAddress)
	require.Equal(t, expected.OffchainreportingOracleSpec.P2PPeerID, actual.OffchainreportingOracleSpec.P2PPeerID)
	require.Equal(t, expected.OffchainreportingOracleSpec.P2PBootstrapPeers, actual.OffchainreportingOracleSpec.P2PBootstrapPeers)
	require.Equal(t, expected.OffchainreportingOracleSpec.IsBootstrapPeer, actual.OffchainreportingOracleSpec.IsBootstrapPeer)
	require.Equal(t, expected.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID, actual.OffchainreportingOracleSpec.EncryptedOCRKeyBundleID)
	require.Equal(t, expected.OffchainreportingOracleSpec.MonitoringEndpoint, actual.OffchainreportingOracleSpec.MonitoringEndpoint)
	require.Equal(t, expected.OffchainreportingOracleSpec.TransmitterAddress, actual.OffchainreportingOracleSpec.TransmitterAddress)
	require.Equal(t, expected.OffchainreportingOracleSpec.ObservationTimeout, actual.OffchainreportingOracleSpec.ObservationTimeout)
	require.Equal(t, expected.OffchainreportingOracleSpec.BlockchainTimeout, actual.OffchainreportingOracleSpec.BlockchainTimeout)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval, actual.OffchainreportingOracleSpec.ContractConfigTrackerSubscribeInterval)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval, actual.OffchainreportingOracleSpec.ContractConfigTrackerPollInterval)
	require.Equal(t, expected.OffchainreportingOracleSpec.ContractConfigConfirmations, actual.OffchainreportingOracleSpec.ContractConfigConfirmations)
}
