package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink/core/services/job"
	"github.com/stretchr/testify/assert"

	webpresenters "github.com/smartcontractkit/chainlink/core/web/presenters"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/core/adapters"
	"github.com/smartcontractkit/chainlink/core/assets"
	"github.com/smartcontractkit/chainlink/core/internal/cltest"
	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/smartcontractkit/chainlink/core/store/presenters"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"gopkg.in/guregu/null.v4"
)

func TestMigrateController_Migrate(t *testing.T) {
	app, cleanup := cltest.NewApplicationWithKey(t)
	t.Cleanup(cleanup)
	require.NoError(t, app.Start())

	client := app.NewHTTPClient()
	// Create a v1 FM job.
	// TODO one with spec errors
	jsr := models.JobSpecRequest{
		Name: "a v1 fm job",
		Initiators: []models.InitiatorRequest{
			{
				Type: models.InitiatorFluxMonitor,
				InitiatorParams: models.InitiatorParams{
					Address:           common.HexToAddress("0x5A0b54D5dc17e0AadC383d2db43B0a0D3E029c4c"),
					RequestData:       models.JSON{Result: gjson.Parse(`{"data":{"coin":"ETH","market":"USD"}}`)},
					Feeds:             models.JSON{Result: gjson.Parse(`["https://lambda.staging.devnet.tools/bnc/call"]`)},
					Threshold:         0.5,
					AbsoluteThreshold: 0.01,
					IdleTimer: models.IdleTimerConfig{
						Duration: models.MustMakeDuration(2 * time.Minute),
					},
					PollTimer: models.PollTimerConfig{
						Period: models.MustMakeDuration(time.Minute),
					},
					Precision: 2,
				},
			},
		},
		Tasks: []models.TaskSpecRequest{
			{
				Type:   adapters.TaskTypeMultiply,
				Params: models.MustParseJSON([]byte(`{"times":"10"}`)),
			},
			{
				Type: adapters.TaskTypeEthUint256,
			},
			{
				Type: adapters.TaskTypeEthTx,
			},
		},
		MinPayment: assets.NewLink(100),
		StartAt:    null.TimeFrom(time.Now()),
		EndAt:      null.TimeFrom(time.Now().Add(time.Second)),
	}

	// Create a v1 job
	b, err := json.Marshal(&jsr)
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/specs", bytes.NewReader(b))
	t.Cleanup(cleanup)
	var jobV1 presenters.JobSpec
	cltest.ParseJSONAPIResponse(t, resp, &jobV1)

	// Migrate it
	resp, cleanup = client.Post(fmt.Sprintf("/v2/migrate/%s", jobV1.ID.String()), nil)
	t.Cleanup(cleanup)
	var createdJobV2 webpresenters.JobResource
	cltest.ParseJSONAPIResponse(t, resp, &createdJobV2)

	// v2 job migrated should be identical to v1.
	assert.Equal(t, uint32(1), createdJobV2.SchemaVersion)
	assert.Equal(t, job.FluxMonitor.String(), createdJobV2.Type.String())
	assert.Equal(t, createdJobV2.Name, jobV1.Name)
	require.NotNil(t, createdJobV2.FluxMonitorSpec)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.CreatedAt, jobV1.CreatedAt)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.MinPayment, jobV1.MinPayment)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.AbsoluteThreshold, jobV1.Initiators[0].AbsoluteThreshold)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.Precision, jobV1.Initiators[0].Precision)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.Threshold, jobV1.Initiators[0].Threshold)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.IdleTimerDisabled, jobV1.Initiators[0].IdleTimer.Disabled)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.IdleTimerPeriod, jobV1.Initiators[0].IdleTimer.Duration.String())
	assert.Equal(t, createdJobV2.FluxMonitorSpec.PollTimerDisabled, jobV1.Initiators[0].PollTimer.Disabled)
	assert.Equal(t, createdJobV2.FluxMonitorSpec.PollTimerPeriod, jobV1.Initiators[0].PollTimer.Period.String())

	// v1 FM job should be archived
	resp, cleanup = client.Get(fmt.Sprintf("/v2/specs/%s", jobV1.ID.String()), nil)
	t.Cleanup(cleanup)
	errs := cltest.ParseJSONAPIErrors(t, resp.Body)
	require.NotNil(t, errs)
	require.Len(t, errs.Errors, 1)
	require.Equal(t, "JobSpec not found", errs.Errors[0].Detail)

	// v2 job read should be identical to created.
	resp, cleanup = client.Get(fmt.Sprintf("/v2/jobs/%s", createdJobV2.ID), nil)
	t.Cleanup(cleanup)
	var migratedJobV2 webpresenters.JobResource
	cltest.ParseJSONAPIResponse(t, resp, &migratedJobV2)

}
