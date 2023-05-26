package functions

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
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

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL: ts.URL,
	}

	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(encryptedSecretsUrls, requestId, jobName)

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

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL: ts.URL,
	}

	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(encryptedSecretsUrls, requestId, jobName)

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

	encryptedSecretsUrls := []byte("test")
	requestId := "1234"
	jobName := "TestJob"

	ea := ExternalAdapterInterface{
		AdapterURL: ts.URL,
	}

	encryptedSecrets, userError, err := ea.FetchEncryptedSecrets(encryptedSecretsUrls, requestId, jobName)

	assert.Nil(t, encryptedSecrets, "Unexpected encryptedSecrets")
	assert.Nil(t, userError, "Unexpected userError")
	assert.Equal(t, err.Error(), "unexpected response ")
}
