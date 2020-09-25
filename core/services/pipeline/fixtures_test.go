package pipeline_test

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

const dotStr = `
    // data source 1
    ds1          [type=bridge name=voter_turnout];
    ds1_parse    [type=jsonparse path="one,two"];
    ds1_multiply [type=multiply times=1.23];

    // data source 2
    ds2          [type=http method=GET url="https://chain.link/voter_turnout/USA-2020" requestData="{\"hi\": \"hello\"}"];
    ds2_parse    [type=jsonparse path="three,four"];
    ds2_multiply [type=multiply times=4.56];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name=election_winner index=1];
`

const ocrJobSpecText = `
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
    ds1_parse    [type=jsonparse path="data,result"];
    ds1_multiply [type=multiply times=100];

    // data source 2
    ds2          [type=httpunrestricted method=POST url="%s" requestData="{\\"hi\\":\\"hello\\"}"];
    ds2_parse    [type=jsonparse path="turnout"];
    ds2_multiply [type=multiply times=100];

    ds1 -> ds1_parse -> ds1_multiply -> answer1;
    ds2 -> ds2_parse -> ds2_multiply -> answer1;

    answer1 [type=median                      index=0];
    answer2 [type=bridge name=election_winner index=1];
"""
`

func makeOCRJobSpec(t *testing.T) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()
	return makeOCRJobSpecWithHTTPURL(t, "https://chain.link/voter_turnout/USA-2020")
}

func makeOCRJobSpecWithHTTPURL(t *testing.T, url string) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()

	jobSpecText := fmt.Sprintf(ocrJobSpecText, cltest.NewAddress().Hex(), url)

	var ocrspec offchainreporting.OracleSpec
	err := toml.Unmarshal([]byte(jobSpecText), &ocrspec)
	require.NoError(t, err)

	dbSpec := models.JobSpecV2{OffchainreportingOracleSpec: &ocrspec.OffchainReportingOracleSpec}
	return &ocrspec, &dbSpec
}
