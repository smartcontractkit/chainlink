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

func Test_FetchEncryptedSecrets_Success(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"result": "success",
			"data": {
				"result": "0x616263646566",
				"error": ""
			},
			"statusCode": 200
		}`)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.NoError(t, err, "Unexpected error")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, "abcdef", string(encryptedSecrets), "Unexpected result")
}

func Test_FetchEncryptedSecrets_UserError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"result": "error",
			"data": {
				"result": "",
				"error": "0x616263646566"
			},
			"statusCode": 200
		}`)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.NoError(t, err, "Unexpected error")
	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Equal(t, "abcdef", string(userError), "Unexpected userError")
}

func Test_FetchEncryptedSecrets_UnexpectedResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"invalid": "invalid"
		}`)
	}))

	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, "error fetching encrypted secrets: external adapter response data was empty", err.Error())
}

func Test_FetchEncryptedSecrets_MissingData(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"result": "success",
			"statusCode": 200
		}`)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, "error fetching encrypted secrets: external adapter response data was empty", err.Error())
}

func Test_FetchEncryptedSecrets_InvalidResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
				"result": "success",
				"data": {
					"result": "invalidHexstring",
					"error": ""
				},
				"statusCode": 200
			}`)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, "error fetching encrypted secrets: error decoding result hex string: hex string must have 0x prefix", err.Error())
}

func Test_FetchEncryptedSecrets_InvalidUserError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
				"result": "error",
				"data": {
					"error": "invalidHexstring",
					"result": ""
				},
				"statusCode": 200
			}`)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, "error fetching encrypted secrets: error decoding userError hex string: hex string must have 0x prefix", err.Error())
}

func Test_FetchEncryptedSecrets_UnexpectedResult(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
				"result": "unexpected",
				"data": {
					"result": "0x01",
					"error": ""
				},
				"statusCode": 200
			}`)
	}))
	defer ts.Close()

	adapterUrl, err := url.Parse(ts.URL)
	assert.NoError(t, err, "Unexpected error")

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL:                   *adapterUrl,
		MaxSecretsFetchResponseBytes: 100_000,
	}

	ctx, cancel := context.WithCancel(testutils.Context(t))
	defer cancel()
	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(ctx, encryptedSecretsUrls, requestId, jobName)

	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, "error fetching encrypted secrets: unexpected result in response: 'unexpected'", err.Error())
}
