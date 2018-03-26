package web_test

import (
	"reflect"
	"testing"

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
	assert.Nil(t, err)
	httpGet := adapter1.(*adapters.HTTPGet)
	assert.Equal(t, httpGet.URL.String(), "https://bitstamp.net/api/ticker/")

	adapter2, err := adapters.For(j.Tasks[1], app.Store)
	assert.Nil(t, err)
	jsonParse := adapter2.(*adapters.JSONParse)
	assert.Equal(t, jsonParse.Path, []string{"last"})

	adapter3, err := adapters.For(j.Tasks[2], app.Store)
	assert.Nil(t, err)
	assert.Equal(t, "*adapters.EthBytes32", reflect.TypeOf(adapter3).String())

	adapter4, err := adapters.For(j.Tasks[3], app.Store)
	assert.Nil(t, err)
	ethTx := adapter4.(*adapters.EthTx)
	assert.Equal(t, ethTx.Address, common.HexToAddress("0x9CA9d2D5E04012C9Ed24C0e513C9bfAa4A2dD77f"))
}
