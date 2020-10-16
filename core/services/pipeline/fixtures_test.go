package pipeline_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	peer "github.com/libp2p/go-libp2p-core/peer"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/services/offchainreporting"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
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
"""
`

func makeOCRJobSpec(t *testing.T, db *gorm.DB) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()
	return makeOCRJobSpecWithHTTPURL(t, db, "https://chain.link/voter_turnout/USA-2020")
}

func makeOCRJobSpecWithHTTPURL(t *testing.T, db *gorm.DB, httpUrl string) (*offchainreporting.OracleSpec, *models.JobSpecV2) {
	t.Helper()

	// Insert keys into the store
	t.Fatal("FIXME: use fixture keys instead")
	keystore := offchainreporting.NewKeyStore(db)
	p2pkey, _, err := keystore.GenerateEncryptedP2PKey("password")
	require.NoError(t, err)
	ocrkey, _, err := keystore.GenerateEncryptedOCRKeyBundle("password")
	require.NoError(t, err)
	peerID, err := p2pkey.GetPeerID()
	require.NoError(t, err)
	err = db.Create(&models.Key{
		Address: cltest.DefaultKey,
		JSON:    cltest.JSONFromString(t, "{}"),
	}).Error
	require.NoError(t, err)

	jobSpecText := fmt.Sprintf(ocrJobSpecText, cltest.NewAddress().Hex(), peer.ID(peerID), ocrkey.ID, cltest.DefaultKey, httpUrl)

	var ocrspec offchainreporting.OracleSpec
	err = toml.Unmarshal([]byte(jobSpecText), &ocrspec)
	require.NoError(t, err)

	dbSpec := models.JobSpecV2{OffchainreportingOracleSpec: &ocrspec.OffchainReportingOracleSpec}
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
