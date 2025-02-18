package clients

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEntityToolClient_GetOverrides(t *testing.T) {
	type testCase struct {
		name              string
		asset             string
		quote             string
		product           *string
		responseStatus    int
		responseBody      string
		expectedOverrides GetOverridesResponse
		expectError       bool
	}

	testCases := []testCase{
		{
			name:           "success default product",
			asset:          "BTC",
			quote:          "USD",
			product:        nil,
			responseStatus: http.StatusOK,
			responseBody: `{
                "externalAdapterRequestParams": {
                    "overrides": {
                        "api1": "value1",
                        "api2": "value2"
                    }
                },
                "apis": []
            }`,
			expectedOverrides: GetOverridesResponse{
				"api1": "value1",
				"api2": "value2",
			},
			expectError: false,
		},
		{
			name:           "non-200 response",
			asset:          "BTC",
			quote:          "USD",
			product:        nil,
			responseStatus: http.StatusInternalServerError,
			responseBody:   "internal error",
			expectError:    true,
		},
		{
			name:           "no overrides present",
			asset:          "BTC",
			quote:          "USD",
			product:        nil,
			responseStatus: http.StatusOK,
			responseBody: `{
                "externalAdapterRequestParams": {
                    "overrides": null
                },
                "apis": []
            }`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				if q.Get("base") != tc.asset || q.Get("quote") != tc.quote {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				fmt.Fprint(w, tc.responseBody)
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			client := NewEntityToolClient(server.URL, server.Client())

			var req *GetOverridesRequest
			if tc.product != nil {
				req = NewGetOverridesRequest(tc.asset, tc.quote, *tc.product)
			} else {
				req = NewGetOverridesRequest(tc.asset, tc.quote)
			}

			ctx := context.Background()
			resp, err := client.GetOverrides(ctx, req)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(*resp) != len(tc.expectedOverrides) {
				t.Fatalf("expected %d override entries, got %d", len(tc.expectedOverrides), len(*resp))
			}
			for key, expectedVal := range tc.expectedOverrides {
				if val, ok := (*resp)[key]; !ok || val != expectedVal {
					t.Errorf("expected %s to be %s, got %v", key, expectedVal, val)
				}
			}
		})
	}
}

func TestEntityToolClient_GetAssetEAs(t *testing.T) {
	type testCase struct {
		name           string
		asset          string
		quote          string
		product        *string
		responseStatus int
		responseBody   string
		expectedAPIs   []string
		expectError    bool
	}

	testCases := []testCase{
		{
			name:           "success default product",
			asset:          "BTC",
			quote:          "USD",
			product:        nil,
			responseStatus: http.StatusOK,
			responseBody: `{
                "externalAdapterRequestParams": { "overrides": {} },
                "apis": ["tiingo", "ncfx"]
            }`,
			expectedAPIs: []string{"tiingo", "ncfx"},
			expectError:  false,
		},
		{
			name:           "non-200 response",
			asset:          "BTC",
			quote:          "USD",
			product:        nil,
			responseStatus: http.StatusInternalServerError,
			responseBody:   "internal error",
			expectError:    true,
		},
		{
			name:           "empty APIs list",
			asset:          "BTC",
			quote:          "USD",
			product:        nil,
			responseStatus: http.StatusOK,
			responseBody: `{
                "externalAdapterRequestParams": { "overrides": {"dummy": "dummy"} },
                "apis": []
            }`,
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				q := r.URL.Query()
				if q.Get("base") != tc.asset || q.Get("quote") != tc.quote {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				fmt.Fprint(w, tc.responseBody)
			})
			server := httptest.NewServer(handler)
			defer server.Close()

			client := NewEntityToolClient(server.URL, server.Client())

			var req *GetAssetEAsRequest
			if tc.product != nil {
				req = NewGetAssetEAsRequest(tc.asset, tc.quote, *tc.product)
			} else {
				req = NewGetAssetEAsRequest(tc.asset, tc.quote)
			}

			ctx := context.Background()
			resp, err := client.GetAssetEAs(ctx, req)
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(*resp) != len(tc.expectedAPIs) {
				t.Fatalf("expected %d APIs, got %d", len(tc.expectedAPIs), len(*resp))
			}
			for i, expected := range tc.expectedAPIs {
				if (*resp)[i] != expected {
					t.Errorf("expected API at index %d to be %s, got %s", i, expected, (*resp)[i])
				}
			}
		})
	}
}
