package services_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/smartcontractkit/chainlink/internal/cltest"
	"github.com/smartcontractkit/chainlink/services"
	"github.com/smartcontractkit/chainlink/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestStore_FormatLogJSON(t *testing.T) {
	t.Parallel()

	var clData models.JSON
	clDataFixture := `{"url":"https://etherprice.com/api","path":["recent","usd"],"address":"0x3cCad4715152693fE3BC4460591e3D3Fbd071b42","dataPrefix":"0x0000000000000000000000000000000000000000000000000000000000000001","functionSelector":"76005c26"}`
	assert.Nil(t, json.Unmarshal([]byte(clDataFixture), &clData))

	hwLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")
	exampleLog := cltest.LogFromFixture("../internal/fixtures/eth/subscription_logs.json")
	tests := []struct {
		name        string
		el          types.Log
		initr       models.Initiator
		wantErrored bool
		wantData    models.JSON
	}{
		{"example ethLog", exampleLog, models.Initiator{Type: "ethlog"}, false,
			jsonFromFixture("../internal/fixtures/eth/subscription_logs.json")},
		{"hello world ethLog", hwLog, models.Initiator{Type: "ethlog"}, false,
			jsonFromFixture("../internal/fixtures/eth/subscription_logs_hello_world.json")},
		{"hello world runLog", hwLog, models.Initiator{Type: "runlog"}, false,
			clData},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := services.FormatLogJSON(test.initr, test.el)
			assert.JSONEq(t, strings.ToLower(test.wantData.String()), strings.ToLower(output.String()))
			assert.Equal(t, test.wantErrored, (err != nil))
		})
	}
}

func jsonFromFixture(path string) models.JSON {
	res := gjson.Get(string(cltest.LoadJSON(path)), "params.result")
	out := cltest.JSONFromString(res.String())
	return out
}

// If updating this test, be sure to update the truffle suite's "expected event signature" test.
func TestServices_RunLogTopic_ExpectedEventSignature(t *testing.T) {
	t.Parallel()

	expected := "0x06f4bf36b4e011a5c499cef1113c2d166800ce4013f6c2509cab1a0e92b83fb2"
	assert.Equal(t, expected, services.RunLogTopic.Hex())
}
