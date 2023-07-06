package functions_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

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
	userResult, userError, domains, err := ea.RunComputation(testutils.Context(t), "requestID1234", "TestJob", "SubOwner", 1, "", &functions.RequestData{})

	if expectedError != nil {
		assert.Equal(t, expectedError.Error(), err.Error(), "Unexpected error")
	} else {
		assert.Nil(t, err)
	}
	assert.Equal(t, expectedUserResult, string(userResult), "Unexpected user result")
	assert.Equal(t, expectedUserError, string(userError), "Unexpected user error")
	assert.Equal(t, expectedDomains, domains, "Unexpected domains")
}

func TestFetchEncryptedSecrets_Success(t *testing.T) {
	runFetcherTest(t, `{
			"result": "success",
			"data": {
				"result": "0x616263646566",
				"error": ""
			},
			"statusCode": 200
		}`, "abcdef", "", nil)
}

func TestFetchEncryptedSecrets_UserError(t *testing.T) {
	runFetcherTest(t, `{
			"result": "error",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 200
		}`, "", "abcdef", nil)
}

func TestFetchEncryptedSecrets_UnexpectedResponse(t *testing.T) {
	runFetcherTest(t, `{
			"invalid": "invalid",
			"statusCode": 200
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter response data was empty"))
}

func TestFetchEncryptedSecrets_FailedStatusCode(t *testing.T) {
	runFetcherTest(t, `{
			"result": "success",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 400
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter invalid StatusCode 400"))
}

func TestFetchEncryptedSecrets_MissingData(t *testing.T) {
	runFetcherTest(t, `{
			"result": "success",
			"statusCode": 200
		}`, "", "", fmt.Errorf("error fetching encrypted secrets: external adapter response data was empty"))
}

func TestFetchEncryptedSecrets_InvalidResponse(t *testing.T) {
	runFetcherTest(t, `{
				"result": "success",
				"data": {
					"result": "invalidHexstring",
					"error": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: error decoding result hex string: hex string must have 0x prefix"))
}

func TestFetchEncryptedSecrets_InvalidUserError(t *testing.T) {
	runFetcherTest(t, `{
				"result": "error",
				"data": {
					"error": "invalidHexstring",
					"result": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: error decoding userError hex string: hex string must have 0x prefix"))
}

func TestFetchEncryptedSecrets_UnexpectedResult(t *testing.T) {
	runFetcherTest(t, `{
				"result": "unexpected",
				"data": {
					"result": "0x01",
					"error": ""
				},
				"statusCode": 200
			}`, "", "", fmt.Errorf("error fetching encrypted secrets: unexpected result in response: 'unexpected'"))
}

func TestRunComputation_Success(t *testing.T) {
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

func TestRunComputation_MissingData(t *testing.T) {
	runRequestTest(t, `{
				"result": "success",
				"statusCode": 200
			}`, "", "", nil, fmt.Errorf("error running computation: external adapter response data was empty"))
}

func TestRunComputation_CorrectAdapterRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)
		expectedData := `{"source":"abcd","language":7,"codeLocation":42,"secrets":"qrvM","secretsLocation":88,"args":["arg1","arg2"]}`
		expectedBody := fmt.Sprintf(`{"endpoint":"lambda","requestId":"requestID1234","jobName":"TestJob","subscriptionOwner":"SubOwner","subscriptionId":1,"nodeProvidedSecrets":"secRETS","data":%s}`, expectedData)
		assert.Equal(t, expectedBody, string(body))

		fmt.Fprintln(w, "}}invalidJSON")
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	ea := functions.NewExternalAdapterClient(*adapterUrl, 100_000)
	reqData := &functions.RequestData{
		Source:          "abcd",
		Language:        7,
		CodeLocation:    42,
		Secrets:         []byte{0xaa, 0xbb, 0xcc}, // "qrvM" base64 encoded
		SecretsLocation: 88,
		Args:            []string{"arg1", "arg2"},
	}
	_, _, _, err = ea.RunComputation(testutils.Context(t), "requestID1234", "TestJob", "SubOwner", 1, "secRETS", reqData)
	assert.Error(t, err)
}

func TestRunComputation_HTTP500(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	ea := functions.NewExternalAdapterClient(*adapterUrl, 100_000)
	_, _, _, err = ea.RunComputation(testutils.Context(t), "requestID1234", "TestJob", "SubOwner", 1, "secRETS", &functions.RequestData{})
	assert.Error(t, err)
}

func TestRunComputation_ContextRespected(t *testing.T) {
	done := make(chan struct{})
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-done
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err)

	ea := functions.NewExternalAdapterClient(*adapterUrl, 100_000)
	ctx, cancel := context.WithTimeout(testutils.Context(t), 10*time.Millisecond)
	defer cancel()
	_, _, _, err = ea.RunComputation(ctx, "requestID1234", "TestJob", "SubOwner", 1, "secRETS", &functions.RequestData{})
	assert.Error(t, err)
	close(done)
}
