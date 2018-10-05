package web_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/smartcontractkit/chainlink/store/models"

	"github.com/ethereum/go-ethereum/common"
	"github.com/smartcontractkit/chainlink/adapters"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/stretchr/testify/assert"
)

func TestAssignmentsController_Create_V1_Format(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()

	j := cltest.FixtureCreateJobWithAssignmentViaWeb(t, app, "../internal/fixtures/web/v1_format_job.json")

	adapter1, err := adapters.For(j.Tasks[0], app.Store)
	assert.NoError(t, err)
	httpGet := adapter1.BaseAdapter.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.URL.String(), "https://bitstamp.net/api/ticker/")

	adapter2, err := adapters.For(j.Tasks[1], app.Store)
	assert.NoError(t, err)
	jsonParse := adapter2.BaseAdapter.(*adapters.JSONParse)
	assert.Equal(t, []string(jsonParse.Path), []string{"last"})

	adapter3, err := adapters.For(j.Tasks[2], app.Store)
	assert.NoError(t, err)
	assert.Equal(t, "*adapters.EthBytes32", reflect.TypeOf(adapter3.BaseAdapter).String())

	adapter4, err := adapters.For(j.Tasks[3], app.Store)
	assert.NoError(t, err)
	ethTx := adapter4.BaseAdapter.(*adapters.EthTx)
	assert.Equal(t, ethTx.Address, common.HexToAddress("0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f"))
}

func TestAssignmentsController_Show_V1_Format(t *testing.T) {
	t.Parallel()
	app, cleanup := cltest.NewApplication()
	defer cleanup()
	client := app.NewHTTPClient()

	j := cltest.FixtureCreateJobWithAssignmentViaWeb(t, app, "../internal/fixtures/web/v1_format_job_with_schedule.json")
	a1, err := models.ConvertToAssignment(j)
	assert.NoError(t, err)

	resp, cleanup := client.Get("/v1/assignments/" + j.ID)
	defer cleanup()
	assert.Equal(t, 200, resp.StatusCode, "Response should be successful")

	var respAssignment models.AssignmentSpec
	json.Unmarshal(cltest.ParseResponseBody(resp), &respAssignment)

	for i, v := range a1.Assignment.Subtasks {
		assert.Equal(t, v.Type, respAssignment.Assignment.Subtasks[i].Type)
		assert.JSONEq(t, v.Params.String(), respAssignment.Assignment.Subtasks[i].Params.String())
	}

	for i, v := range a1.Schedule.RunAt {
		assert.Equal(t, respAssignment.Schedule.RunAt[i], v)
	}

	assert.Equal(t, a1.Schedule.DayOfMonth, respAssignment.Schedule.DayOfMonth)
	assert.Equal(t, a1.Schedule.DayOfWeek, respAssignment.Schedule.DayOfWeek)
	assert.Equal(t, a1.Schedule.EndAt, respAssignment.Schedule.EndAt)
	assert.Equal(t, a1.Schedule.Hour, respAssignment.Schedule.Hour)
	assert.Equal(t, a1.Schedule.Minute, respAssignment.Schedule.Minute)
	assert.Equal(t, a1.Schedule.MonthOfYear, respAssignment.Schedule.MonthOfYear)

}
