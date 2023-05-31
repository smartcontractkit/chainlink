package functions

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
)

func runTest(t *testing.T, adapterJSONResponse, expectedSecrets, expectedUserError string, expectedError error) {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, adapterJSONResponse)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterClient{
		AdapterURL:       *adapterUrl,
		MaxResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	if expectedError != nil {
		assert.Equal(t, expectedError.Error(), err.Error(), "Unexpected error")
	} else {
		assert.Nil(t, err)
	}
	assert.Equal(t, expectedUserError, string(userError), "Unexpected userError")
	assert.Equal(t, expectedSecrets, string(encryptedSecrets), "Unexpected secrets")
}

func Test_FetchEncryptedSecrets_Success(t *testing.T) {
	runTest(t, `{
			"result": "success",
			"data": {
				"result": "0x616263646566",
				"error": ""
			},
			"statusCode": 200
		}`, "abcdef", "", nil)
}

func Test_FetchEncryptedSecrets_UserError(t *testing.T) {
	runTest(t, `{
			"result": "error",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 200
		}`, "", "abcdef", nil)
}

func Test_FetchEncryptedSecrets_UnexpectedResponse(t *testing.T) {
	runTest(t, `{
			"invalid": "invalid",
			"statusCode": 200
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter response data was empty"))
}

func Test_FetchEncryptedSecrets_FailedStatusCode(t *testing.T) {
	runTest(t, `{
			"result": "success",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 400
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter responded with error code 400"))
}

func Test_FetchEncryptedSecrets_MissingData(t *testing.T) {
	runTest(t, `{
			"result": "success",
			"statusCode": 200
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter response data was empty"))
}

func Test_FetchEncryptedSecrets_InvalidResponse(t *testing.T) {
	runTest(t, `{
				"result": "success",
				"data": {
					"result": "invalidHexstring",
					"error": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: error decoding result hex string: hex string must have 0x prefix"))
}

func Test_FetchEncryptedSecrets_InvalidUserError(t *testing.T) {
	runTest(t, `{
				"result": "error",
				"data": {
					"error": "invalidHexstring",
					"result": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: error decoding userError hex string: hex string must have 0x prefix"))
}

func Test_FetchEncryptedSecrets_UnexpectedResult(t *testing.T) {
	runTest(t, `{
				"result": "unexpected",
				"data": {
					"result": "0x01",
					"error": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: unexpected result in response: 'unexpected'"))
}
