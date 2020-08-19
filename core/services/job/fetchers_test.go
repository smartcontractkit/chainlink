package job_test

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/core/services/job"
)

func TestFetchers_JSONMarshalUnmarshal(t *testing.T) {
	json := `{
        "type": "median",
        "fetchers": [
            {
                "type": "http",
                "url": "http://chain.link",
                "method": "GET",
                "requestData": {
                    "one": "asdf",
                    "two": "xyzzy"
                },
                "transformPipeline": [
                    { "type": "jsonparse", "path": ["one", "two"] },
                    { "type": "multiply", "times": 1.23 }
                ]
            },
            {
                "type": "bridge",
                "name": "t00f4r",
                "requestData": {
                    "one": "asdf",
                    "two": "xyzzy"
                },
                "transformPipeline": [
                    { "type": "jsonparse", "path": ["one", "two"] },
                    { "type": "multiply", "times": 1.23 }
                ]
            }
        ]
    }`

	expected := job.MedianFetcher{
		Fetchers: []job.Fetcher{
			job.HttpFetcher{
				URL:    "http://chain.link",
				Method: "GET",
				RequestData: map[string]interface{}{
					"one": "asdf",
					"two": "xyzzy",
				},
				Transformers: job.Transformers{
					job.JSONParseTransformer{Path: []string{"one", "two"}},
					job.MultiplyTransformer{Times: decimal.NewFromFloat(1.23)},
				},
			},
			job.BridgeFetcher{
				BridgeName: "t00f4r",
				RequestData: map[string]interface{}{
					"one": "asdf",
					"two": "xyzzy",
				},
				Transformers: job.Transformers{
					job.JSONParseTransformer{Path: []string{"one", "two"}},
					job.MultiplyTransformer{Times: decimal.NewFromFloat(1.23)},
				},
			},
		},
	}

	fetcher, err := job.UnmarshalFetcherJSON([]byte(json))
	require.NoError(t, err)
	require.Equal(t, expected, fetcher)
}
