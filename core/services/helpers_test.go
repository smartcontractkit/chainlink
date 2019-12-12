package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustReadFile(t testing.TB, file string) string {
	t.Helper()

	content, err := ioutil.ReadFile(file)
	require.NoError(t, err)
	return string(content)
}

type fixedFetcher struct {
	price float64
}

func newFixedPricedFetcher(price float64) *fixedFetcher {
	return &fixedFetcher{price: price}
}

func (ps *fixedFetcher) Fetch() (float64, error) {
	return ps.price, nil
}

type erroringFetcher struct{}

func newErroringPricedFetcher() *erroringFetcher {
	return &erroringFetcher{}
}

func (*erroringFetcher) Fetch() (float64, error) {
	return 0, errors.New("failed to fetch; I always error")
}

func fakePriceResponder(t *testing.T, requestData string, result float64) http.Handler {
	t.Helper()

	response := adapterResponse{Data: dataWithResult(t, result)}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		payload, err := ioutil.ReadAll(r.Body)
		require.NoError(t, err)
		defer r.Body.Close()
		assert.Equal(t, requestData, string(payload))
		w.Header().Set("Content-Type", "application/json")
		require.NoError(t, json.NewEncoder(w).Encode(response))
	})
}

func dataWithResult(t *testing.T, result float64) adapterResponseData {
	t.Helper()
	var data adapterResponseData
	body := []byte(fmt.Sprintf(`{"result":%v}`, result))
	require.NoError(t, json.Unmarshal(body, &data))
	return data
}
