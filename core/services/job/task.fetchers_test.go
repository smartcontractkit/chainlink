package job_test

// import (
// 	"encoding/json"
// 	"testing"

// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/require"

// 	"github.com/smartcontractkit/chainlink/core/services/job"
// )

// func TestFetchers_UnmarshalJSON(t *testing.T) {
// 	json := `{
//         "type": "median",
//         "fetchers": [
//             {
//                 "type": "http",
//                 "url": "http://chain.link",
//                 "method": "GET",
//                 "requestData": {
//                     "one": "asdf",
//                     "two": "xyzzy"
//                 },
//                 "transformPipeline": [
//                     { "type": "jsonparse", "path": ["one", "two"] },
//                     { "type": "multiply", "times": 1.23 }
//                 ]
//             },
//             {
//                 "type": "bridge",
//                 "name": "t00f4r",
//                 "requestData": {
//                     "one": "asdf",
//                     "two": "xyzzy"
//                 },
//                 "transformPipeline": [
//                     { "type": "jsonparse", "path": ["one", "two"] },
//                     { "type": "multiply", "times": 1.23 }
//                 ]
//             }
//         ]
//     }`

// 	expected := job.MedianTask{
// 		Fetchers: []job.Fetcher{
// 			job.HTTPTask{
// 				URL:    "http://chain.link",
// 				Method: "GET",
// 				RequestData: map[string]interface{}{
// 					"one": "asdf",
// 					"two": "xyzzy",
// 				},
// 				Transformers: job.Transformers{
// 					job.JSONParseTask{Path: []string{"one", "two"}},
// 					job.MultiplyTask{Times: decimal.NewFromFloat(1.23)},
// 				},
// 			},
// 			job.BridgeTask{
// 				BridgeName: "t00f4r",
// 				RequestData: map[string]interface{}{
// 					"one": "asdf",
// 					"two": "xyzzy",
// 				},
// 				Transformers: job.Transformers{
// 					job.JSONParseTask{Path: []string{"one", "two"}},
// 					job.MultiplyTask{Times: decimal.NewFromFloat(1.23)},
// 				},
// 			},
// 		},
// 	}

// 	fetcher, err := job.UnmarshalFetcherJSON([]byte(json))
// 	require.NoError(t, err)
// 	require.Equal(t, expected, fetcher)
// }

// func TestFetchers_MarshalJSON(t *testing.T) {
// 	f := job.MedianTask{
// 		Fetchers: []job.Fetcher{
// 			job.HTTPTask{
// 				URL:    "http://chain.link",
// 				Method: "GET",
// 				RequestData: map[string]interface{}{
// 					"one": "asdf",
// 					"two": "xyzzy",
// 				},
// 				Transformers: job.Transformers{
// 					job.JSONParseTask{Path: []string{"one", "two"}},
// 					job.MultiplyTask{Times: decimal.NewFromFloat(1.23)},
// 				},
// 			},
// 			job.BridgeTask{
// 				BridgeName: "t00f4r",
// 				RequestData: map[string]interface{}{
// 					"one": "asdf",
// 					"two": "xyzzy",
// 				},
// 				Transformers: job.Transformers{
// 					job.JSONParseTask{Path: []string{"one", "two"}},
// 					job.MultiplyTask{Times: decimal.NewFromFloat(1.23)},
// 				},
// 			},
// 		},
// 	}

// 	bs, err := json.MarshalIndent(f, "", "    ")
// 	require.NoError(t, err)

// 	expected := `{
//     "type": "median",
//     "fetchers": [
//         {
//             "type": "http",
//             "url": "http://chain.link",
//             "method": "GET",
//             "requestData": {
//                 "one": "asdf",
//                 "two": "xyzzy"
//             },
//             "transformPipeline": [
//                 {
//                     "type": "jsonparse",
//                     "path": [
//                         "one",
//                         "two"
//                     ]
//                 },
//                 {
//                     "type": "multiply",
//                     "times": "1.23"
//                 }
//             ]
//         },
//         {
//             "type": "bridge",
//             "name": "t00f4r",
//             "requestData": {
//                 "one": "asdf",
//                 "two": "xyzzy"
//             },
//             "transformPipeline": [
//                 {
//                     "type": "jsonparse",
//                     "path": [
//                         "one",
//                         "two"
//                     ]
//                 },
//                 {
//                     "type": "multiply",
//                     "times": "1.23"
//                 }
//             ]
//         }
//     ]
// }`

// 	require.Equal(t, expected, string(bs))
// }
