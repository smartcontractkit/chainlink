package customendpoint_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/pelletier/go-toml"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"github.com/smartcontractkit/libocr/commontypes"
	"github.com/smartcontractkit/libocr/offchainreporting2/reportingplugin/median"
	"github.com/smartcontractkit/libocr/offchainreporting2/types"
	"github.com/smartcontractkit/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/guregu/null.v4"

	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/core/internal/testutils/pgtest"
	"github.com/smartcontractkit/chainlink/core/logger"
	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/smartcontractkit/chainlink/core/services/pipeline"
	"github.com/smartcontractkit/chainlink/core/services/relay/customendpoint"
	"github.com/smartcontractkit/chainlink/core/testdata/testspecs"
	"github.com/smartcontractkit/chainlink/core/utils"
)

var bridgeRequestData = `{"data":{"asset":"ETHUSD"}}`

type adapterRequest struct {
	ID          string            `json:"id"`
	Data        pipeline.MapParam `json:"data"`
	Meta        pipeline.MapParam `json:"meta"`
	ResponseURL string            `json:"responseURL"`
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

func dataWithResult(t *testing.T, result decimal.Decimal) adapterResponseData {
	t.Helper()
	var data adapterResponseData
	body := []byte(fmt.Sprintf(`{"result":%v}`, result))
	require.NoError(t, json.Unmarshal(body, &data))
	return data
}

func addBridgeToDb(t *testing.T, db *sqlx.DB, bridgeName, url string) {
	_, err := db.Exec(`INSERT INTO bridge_types (name, url, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING "bridge_types".*`,
		bridgeName, // from
		url,
		"incoming_token_hash",
		"salt",
		"outgoing_token")
	require.NoError(t, err)
}

func fakePriceResponder(t *testing.T, requestData map[string]interface{}, result decimal.Decimal, inputKey string, expectedInput interface{}) http.Handler {
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

		if inputKey != "" {
			m := utils.MustUnmarshalToMap(string(payload))
			if expectedInput != nil {
				require.Equal(t, expectedInput, m[inputKey])
			} else {
				require.Nil(t, m[inputKey])
			}
		}
	})
}

func makeOCR2JobSpecFromToml(t *testing.T, jobSpecToml string) job.OCR2OracleSpec {
	t.Helper()

	var ocr2spec job.OCR2OracleSpec
	err := toml.Unmarshal([]byte(jobSpecToml), &ocr2spec)
	require.NoError(t, err)

	return ocr2spec
}

func TestOcr2Provider(t *testing.T) {
	lggr := logger.TestLogger(t)
	cfg := cltest.NewTestGeneralConfig(t)
	db := pgtest.NewSqlxDB(t)

	pipelineORM := pipeline.NewORM(db, lggr, cfg)

	clock := testutils.NewSimulatedClock(time.Now())

	spec := makeOCR2JobSpecFromToml(t, testspecs.OCR2CustomEndpointSpecMinimal)
	var relayConfig customendpoint.RelayConfig
	err := json.Unmarshal(spec.RelayConfig.Bytes(), &relayConfig)
	require.NoError(t, err)
	s1 := httptest.NewServer(fakePriceResponder(t, utils.MustUnmarshalToMap(bridgeRequestData), decimal.NewFromInt(10000000), "result", "0.00000205"))
	defer s1.Close()

	bridgeUrl, err := url.ParseRequestURI(s1.URL)
	require.NoError(t, err)
	addBridgeToDb(t, db, relayConfig.EndpointTarget, bridgeUrl.String())

	// Get the expected digester from TOML file
	digesterFromToml := customendpoint.CreateConfigDigester(
		relayConfig.EndpointName,
		relayConfig.EndpointTarget,
		relayConfig.PayloadType)
	expectedDigestPrefix := digesterFromToml.ConfigDigestPrefix()
	config := types.ContractConfig{}
	expectedDigest, err := digesterFromToml.ConfigDigest(config)
	require.NoError(t, err)

	// Creat Relayer
	backgroundCtx := context.Background()
	relayer := customendpoint.NewRelayer(lggr, cfg, pipelineORM, &clock)
	err = relayer.Start(backgroundCtx)
	require.NoError(t, err)

	// Create new Ocr2Provider
	ocr2Provider, err := relayer.NewOCR2Provider(uuid.UUID{}, customendpoint.OCR2Spec{
		RelayConfig: relayConfig,
		ID:          spec.ID,
		IsBootstrap: false,
	})
	require.NoError(t, err)

	// Test the MedianContract
	medianContract := ocr2Provider.MedianContract()
	digest, epoch, round, latestAnswer, latestTime, err :=
		medianContract.LatestTransmissionDetails(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, expectedDigest, digest)
	assert.Equal(t, uint32(0), epoch)
	assert.Equal(t, uint8(0), round)
	assert.Equal(t, big.NewInt(0), latestAnswer)
	assert.Equal(t, clock.Now(), latestTime)

	// Test the digester
	digester := ocr2Provider.OffchainConfigDigester()
	digest, err = digester.ConfigDigest(config)
	require.NoError(t, err)
	assert.Equal(t, expectedDigestPrefix, digester.ConfigDigestPrefix())
	assert.Equal(t, expectedDigest, digest)

	// Test the config tracker
	tracker := ocr2Provider.ContractConfigTracker()
	changedInBlock, digest, err := tracker.LatestConfigDetails(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), changedInBlock)
	assert.Equal(t, expectedDigest, digest)
	assert.Nil(t, tracker.Notify())
	blockHeight, err := tracker.LatestBlockHeight(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), blockHeight)

	// Test the transmitter
	transmitter := ocr2Provider.ContractTransmitter()
	digest, epoch, err = transmitter.LatestConfigDigestAndEpoch(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, expectedDigest, digest)
	assert.Equal(t, uint32(0), epoch)
	codec := ocr2Provider.ReportCodec()
	observations := []median.ParsedAttributedObservation{
		{
			Value:           big.NewInt(100),
			JuelsPerFeeCoin: big.NewInt(1),
			Observer:        commontypes.OracleID(uint8(1)),
		},
		{
			Value:           big.NewInt(205),
			JuelsPerFeeCoin: big.NewInt(2),
			Observer:        commontypes.OracleID(uint8(2)),
		},
		{
			Value:           big.NewInt(300),
			JuelsPerFeeCoin: big.NewInt(3),
			Observer:        commontypes.OracleID(uint8(3)),
		},
	}
	report, err := codec.BuildReport(observations)
	require.NoError(t, err)
	reportContext := types.ReportContext{
		types.ReportTimestamp{
			ConfigDigest: digest,
			Epoch:        2,
			Round:        1},
		[32]byte{},
	}
	clock.FastForwardBy(time.Second)
	err = transmitter.Transmit(backgroundCtx, reportContext, report, nil)
	require.NoError(t, err)
	customendpoint.WaitForTransmitters(t, transmitter)
	digest, epoch, err = transmitter.LatestConfigDigestAndEpoch(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, expectedDigest, digest)
	assert.Equal(t, uint32(2), epoch)
	digest, epoch, round, latestAnswer, latestTime, err =
		medianContract.LatestTransmissionDetails(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, expectedDigest, digest)
	assert.Equal(t, uint32(2), epoch)
	assert.Equal(t, uint8(1), round)
	assert.Equal(t, observations[1].Value, latestAnswer)
	assert.Equal(t, clock.Now(), latestTime)

	// Transmit again with older epoch should not update latest config
	reportContext.Epoch = 1 // Set epoch to an older value
	oldTime := clock.Now()
	clock.FastForwardBy(time.Second) // Move clock ahead
	oldResult := observations[1].Value
	observations[1].Value = big.NewInt(207) // Get a new median
	err = transmitter.Transmit(backgroundCtx, reportContext, report, nil)
	require.NoError(t, err)
	customendpoint.WaitForTransmitters(t, transmitter)
	digest, epoch, err = transmitter.LatestConfigDigestAndEpoch(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, expectedDigest, digest)
	assert.Equal(t, uint32(2), epoch) // Epoch is unchanged from previous value
	digest, epoch, round, latestAnswer, latestTime, err =
		medianContract.LatestTransmissionDetails(backgroundCtx)
	require.NoError(t, err)
	assert.Equal(t, expectedDigest, digest)
	assert.Equal(t, uint32(2), epoch)
	assert.Equal(t, uint8(1), round)
	assert.Equal(t, oldResult, latestAnswer)
	assert.Equal(t, oldTime, latestTime)

	// Close relaer
	err = relayer.Close()
	require.NoError(t, err)
}
