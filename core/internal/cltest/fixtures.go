package cltest

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/smartcontractkit/chainlink/core/store/models"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

const (
	key3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea = `{"address":"3cb8e3fd9d27e39a5e9e6852b0e96160061fd4ea","crypto":{"cipher":"aes-128-ctr","ciphertext":"7515678239ccbeeaaaf0b103f0fba46a979bf6b2a52260015f35b9eb5fed5c17","cipherparams":{"iv":"87e5a5db334305e1e4fb8b3538ceea12"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"d89ac837b5dcdce5690af764762fe349d8162bb0086cea2bc3a4289c47853f96"},"mac":"57a7f4ada10d3d89644f541c91f89b5bde73e15e827ee40565e2d1f88bb0ac96"},"id":"c8cb9bc7-0a51-43bd-8348-8a67fd1ec52c","version":3}`
)

// FixtureCreateJobViaWeb creates a job from a fixture using /v2/specs
func FixtureCreateJobViaWeb(t *testing.T, app *TestApplication, path string) models.JobSpec {
	client := app.NewHTTPClient()
	resp, cleanup := client.Post("/v2/specs", bytes.NewBuffer(MustReadFile(t, path)))
	defer cleanup()
	AssertServerResponse(t, resp, 200)

	var job models.JobSpec
	err := ParseJSONAPIResponse(resp, &job)
	require.NoError(t, err)
	return job
}

// FixtureCreateServiceAgreementViaWeb creates a service agreement from a fixture using /v2/service_agreements
func FixtureCreateServiceAgreementViaWeb(
	t *testing.T,
	app *TestApplication,
	path string,
) models.ServiceAgreement {
	client := app.NewHTTPClient()

	agreementWithoutOracle := string(MustReadFile(t, path))
	from := GetAccountAddress(app.ChainlinkApplication.GetStore())
	agreementWithOracle := MustJSONSet(t, agreementWithoutOracle, "oracles", []string{from.Hex()})

	resp, cleanup := client.Post("/v2/service_agreements", bytes.NewBufferString(agreementWithOracle))
	defer cleanup()

	AssertServerResponse(t, resp, 200)
	responseSA := models.ServiceAgreement{}
	err := ParseJSONAPIResponse(resp, &responseSA)
	require.NoError(t, err)

	return FindServiceAgreement(app.Store, responseSA.ID)
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
func LogFromFixture(t *testing.T, path string) models.Log {
	value := gjson.Get(string(MustReadFile(t, path)), "params.result")
	var el models.Log
	require.NoError(t, json.Unmarshal([]byte(value.String()), &el))

	return el
}

// TxReceiptFromFixture create ethtypes.log from file path
func TxReceiptFromFixture(t *testing.T, path string) models.TxReceipt {
	jsonStr := JSONFromFixture(t, path).Get("result").String()

	var receipt models.TxReceipt
	err := json.Unmarshal([]byte(jsonStr), &receipt)
	require.NoError(t, err)

	return receipt
}
