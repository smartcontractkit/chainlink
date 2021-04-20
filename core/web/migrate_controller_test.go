package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

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
	// TODO:
	// - an archived one
	// - one with spec errors.
	// - one with bridge
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
				Type:   adapters.TaskTypeJSONParse,
				Params: models.MustParseJSON([]byte(`{"path":"data"}`)),
			},
			{
				Type:   adapters.TaskTypeMultiply,
				Params: models.MustParseJSON([]byte(`{"times":"10"}`)),
			},
			{
				Type: adapters.TaskTypeEthTx,
			},
		},
		MinPayment: assets.NewLink(100),
		StartAt:    null.TimeFrom(time.Now()),
		EndAt:      null.TimeFrom(time.Now().Add(time.Second)),
	}

	b, err := json.Marshal(&jsr)
	require.NoError(t, err)
	resp, cleanup := client.Post("/v2/specs", bytes.NewReader(b))
	t.Cleanup(cleanup)
	var js presenters.JobSpec
	cltest.ParseJSONAPIResponse(t, resp, &js)

	// Migrate it
	t.Log(js.ID.String())
	resp, cleanup = client.Post(fmt.Sprintf("/v2/migrate/%s", js.ID.String()), nil)
	t.Cleanup(cleanup)
	var out webpresenters.JobResource
	cltest.ParseJSONAPIResponse(t, resp, &out)
	t.Log(out)

	// Should be no v1 FM job left

	// v2 FM job should be identical
}
