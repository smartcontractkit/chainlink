package cltest

import (
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/web"

	"github.com/smartcontractkit/chainlink/core/services/job"

	"github.com/smartcontractkit/chainlink/core/store/models"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

// MustHelloWorldAgreement returns the fixture hello world agreement
func MustHelloWorldAgreement(t *testing.T) string {
	template := MustReadFile(t, "testdata/hello_world_agreement.json")
	return string(template)

}

func FixtureCreateJobSpecV2ViaWeb(t *testing.T, app *TestApplication, path string) job.Job {
	request := web.CreateJobRequest{
		TOML: string(MustReadFile(t, path)),
	}
	output, err := json.Marshal(request)
	require.NoError(t, err)
	return CreateJobViaWeb(t, app, output)
}

// JSONFromFixture create models.JSON from file path
func JSONFromFixture(t *testing.T, path string) models.JSON {
	return JSONFromBytes(t, MustReadFile(t, path))
}

// JSONResultFromFixture create model.JSON with params.result found in the given file path
func JSONResultFromFixture(t *testing.T, path string) models.JSON {
	res := gjson.Get(string(MustReadFile(t, path)), "params.result")
	return JSONFromString(t, res.String())
}

// LogFromFixture create ethtypes.log from file path
func LogFromFixture(t *testing.T, path string) types.Log {
	value := gjson.Get(string(MustReadFile(t, path)), "params.result")
	var el types.Log
	require.NoError(t, json.Unmarshal([]byte(value.String()), &el))

	return el
}

// TxReceiptFromFixture create ethtypes.log from file path
func TxReceiptFromFixture(t *testing.T, path string) *types.Receipt {
	jsonStr := JSONFromFixture(t, path).Get("result").String()

	var receipt types.Receipt
	err := json.Unmarshal([]byte(jsonStr), &receipt)
	require.NoError(t, err)

	return &receipt
}
