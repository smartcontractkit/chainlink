package functions_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/smartcontractkit/chainlink/v2/core/internal/testutils"
	"github.com/smartcontractkit/chainlink/v2/core/services/functions"
)

func runFetcherTest(t *testing.T, adapterJSONResponse, expectedSecrets, expectedUserError string, expectedError error) {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, adapterJSONResponse)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	ea := functions.NewExternalAdapterClient(*adapterUrl, 100_000)
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(testutils.Context(t), []byte("urls to secrets"), "requestID1234", "TestJob")

	if expectedError != nil {
		assert.Equal(t, expectedError.Error(), err.Error(), "Unexpected error")
	} else {
		assert.Nil(t, err)
	}
	assert.Equal(t, expectedUserError, string(userError), "Unexpected userError")
	assert.Equal(t, expectedSecrets, string(encryptedSecrets), "Unexpected secrets")
}

func runRequestTest(t *testing.T, adapterJSONResponse, expectedUserResult, expectedUserError string, expectedDomains []string, expectedError error) {
	t.Helper()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, adapterJSONResponse)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	ea := functions.NewExternalAdapterClient(*adapterUrl, 100_000)
	userResult, userError, domains, err := ea.RunComputation(testutils.Context(t), "requestID1234", "TestJob", "SubOwner", 1, "", []byte("{}"))

	if expectedError != nil {
		assert.Equal(t, expectedError.Error(), err.Error(), "Unexpected error")
	} else {
		assert.Nil(t, err)
	}
	assert.Equal(t, expectedUserResult, string(userResult), "Unexpected user result")
	assert.Equal(t, expectedUserError, string(userError), "Unexpected user error")
	assert.Equal(t, expectedDomains, domains, "Unexpected domains")
}

func Test_FetchEncryptedSecrets_Success(t *testing.T) {
	runFetcherTest(t, `{
			"result": "success",
			"data": {
				"result": "0x616263646566",
				"error": ""
			},
			"statusCode": 200
		}`, "abcdef", "", nil)
}

func Test_FetchEncryptedSecrets_UserError(t *testing.T) {
	runFetcherTest(t, `{
			"result": "error",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 200
		}`, "", "abcdef", nil)
}

func Test_FetchEncryptedSecrets_UnexpectedResponse(t *testing.T) {
	runFetcherTest(t, `{
			"invalid": "invalid",
			"statusCode": 200
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter response data was empty"))
}

func Test_FetchEncryptedSecrets_FailedStatusCode(t *testing.T) {
	runFetcherTest(t, `{
			"result": "success",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 400
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter responded with error code 400"))
}

func Test_FetchEncryptedSecrets_MissingData(t *testing.T) {
	runFetcherTest(t, `{
			"result": "success",
			"statusCode": 200
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter response data was empty"))
}

func Test_FetchEncryptedSecrets_InvalidResponse(t *testing.T) {
	runFetcherTest(t, `{
				"result": "success",
				"data": {
					"result": "invalidHexstring",
					"error": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: error decoding result hex string: hex string must have 0x prefix"))
}

func Test_FetchEncryptedSecrets_InvalidUserError(t *testing.T) {
	runFetcherTest(t, `{
				"result": "error",
				"data": {
					"error": "invalidHexstring",
					"result": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: error decoding userError hex string: hex string must have 0x prefix"))
}

func Test_FetchEncryptedSecrets_UnexpectedResult(t *testing.T) {
	runFetcherTest(t, `{
				"result": "unexpected",
				"data": {
					"result": "0x01",
					"error": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: unexpected result in response: 'unexpected'"))
}

func Test_RunComputation_Success(t *testing.T) {
	runRequestTest(t, `{
	    	"result": "success",
				"data": {
					"result": "0x616263646566",
					"error": "",
					"domains": ["domain1", "domain2"]
				},
				"statusCode": 200
			}`, "abcdef", "", []string{"domain1", "domain2"}, nil)
}

func Test_RunComputation_MissingData(t *testing.T) {
	runRequestTest(t, `{
				"result": "success",
				"statusCode": 200
			}`, "", "", nil, fmt.Errorf("error running computation: external adapter response data was empty"))
}

func Test_RunComputation_CorrectAdapterRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		expectedData := `{"source":"","language":0,"codeLocation":0,"secrets":"","secretsLocation":0,"args":null}`
		expectedBody := fmt.Sprintf(`{"endpoint":"lambda","requestId":"requestID1234","jobName":"TestJob","subscriptionOwner":"SubOwner","subscriptionId":"1","nodeProvidedSecrets":"secRETS","data":%s}`, expectedData)
		assert.Equal(t, expectedBody, string(body))

		fmt.Fprintln(w, "}}invalidJSON")
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	ea := functions.NewExternalAdapterClient(*adapterUrl, 100_000)
	_, _, _, err = ea.RunComputation(testutils.Context(t), "requestID1234", "TestJob", "SubOwner", 1, "secRETS", []byte("{}"))
	assert.Error(t, err)
}
