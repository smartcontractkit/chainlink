package pipeline_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/smartcontractkit/chainlink/core/services"

	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/store/models"
)

const (
	dotStr = `
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
)

func makeMinimalHTTPOracleSpec(t *testing.T, contractAddress, peerID, transmitterAddress, keyBundle, fetchUrl, timeout string) *offchainreporting.OracleSpec {
	t.Helper()
	var os = offchainreporting.OracleSpec{
		OffchainReportingOracleSpec: models.OffchainReportingOracleSpec{
			P2PBootstrapPeers:                      pq.StringArray{},
			ObservationTimeout:                     models.Interval(10 * time.Second),
			BlockchainTimeout:                      models.Interval(20 * time.Second),
			ContractConfigTrackerSubscribeInterval: models.Interval(2 * time.Minute),
			ContractConfigTrackerPollInterval:      models.Interval(1 * time.Minute),
			ContractConfigConfirmations:            uint16(3),
		},
		Pipeline: *pipeline.NewTaskDAG(),
	}
	s := fmt.Sprintf(minimalNonBootstrapTemplate, contractAddress, peerID, transmitterAddress, keyBundle, fetchUrl, timeout)
	_, err := services.ValidatedOracleSpecToml(s)
	require.NoError(t, err)
	err = toml.Unmarshal([]byte(s), &os)
	require.NoError(t, err)
	return &os
}

func makeVoterTurnoutOCRJobSpec(t *testing.T, db *gorm.DB) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()
	return makeVoterTurnoutOCRJobSpecWithHTTPURL(t, db, "https://example.com/foo/bar")
}

func makeVoterTurnoutOCRJobSpecWithHTTPURL(t *testing.T, db *gorm.DB, httpURL string) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()
	peerID := cltest.DefaultP2PPeerID
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	ds := fmt.Sprintf(voterTurnoutDataSourceTemplate, httpURL)
	voterTurnoutJobSpec := fmt.Sprintf(ocrJobSpecTemplate, cltest.NewAddress().Hex(), peerID, ocrKeyID, cltest.DefaultKey, ds)
	return makeOCRJobSpecWithHTTPURL(t, db, voterTurnoutJobSpec)
}

func makeSimpleFetchOCRJobSpecWithHTTPURL(t *testing.T, db *gorm.DB, httpURL string, lax bool) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()
	peerID := cltest.DefaultP2PPeerID
	ocrKeyID := cltest.DefaultOCRKeyBundleID
	ds := fmt.Sprintf(simpleFetchDataSourceTemplate, httpURL, lax)
	simpleFetchJobSpec := fmt.Sprintf(ocrJobSpecTemplate, cltest.NewAddress().Hex(), peerID, ocrKeyID, cltest.DefaultKey, ds)
	return makeOCRJobSpecWithHTTPURL(t, db, simpleFetchJobSpec)
}

func makeOCRJobSpecWithHTTPURL(t *testing.T, db *gorm.DB, jobSpecToml string) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()

	var ocrspec offchainreporting.OracleSpec
	err := toml.Unmarshal([]byte(jobSpecToml), &ocrspec)
	require.NoError(t, err)

	dbSpec := models.JobSpecV2{
		OffchainreportingOracleSpec: &ocrspec.OffchainReportingOracleSpec,
		Type:                        string(offchainreporting.JobType),
		SchemaVersion:               ocrspec.SchemaVersion,
	}
	return &ocrspec, &dbSpec
}

func mustDecimal(t *testing.T, arg string) *decimal.Decimal {
	ret, err := decimal.NewFromString(arg)
	require.NoError(t, err)
	return &ret
}

type adapterRequest struct {
	ID   string                   `json:"id"`
	Data pipeline.HttpRequestData `json:"data"`
	Meta pipeline.HttpRequestData `json:"meta"`
}

type adapterResponseData struct {
	Result *decimal.Decimal `json:"result"`
}

// adapterResponse is the HTTP response as defined by the external adapter:
// https://github.com/smartcontractkit/bnc-adapter
type adapterResponse struct {
	Data         adapterResponseData `json:"data"`
	ErrorMessage null.String         `json:"errorMessage"`
}

func (pr adapterResponse) Result() *decimal.Decimal {
	return pr.Data.Result
}

func fakePriceResponder(t *testing.T, requestData map[string]interface{}, result decimal.Decimal) http.Handler {
	t.Helper()

	body, err := json.Marshal(requestData)
	require.NoError(t, err)
	var expectedRequest adapterRequest
	err = json.Unmarshal(body, &expectedRequest)
	require.NoError(t, err)
	response := adapterResponse{Data: dataWithResult(t, result)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody adapterRequest
		payload, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		err = json.Unmarshal(payload, &reqBody)
		require.NoError(t, err)
		require.Equal(t, expectedRequest.Data, reqBody.Data)
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})
}

func dataWithResult(t *testing.T, result decimal.Decimal) adapterResponseData {
	t.Helper()
	var data adapterResponseData
	body := []byte(fmt.Sprintf(`{"result":%v}`, result))
	require.NoError(t, json.Unmarshal(body, &data))
	return data
}

func mustReadFile(t testing.TB, file string) string {
	t.Helper()

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return string(content)
}
