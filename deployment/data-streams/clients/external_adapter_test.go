package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type queryTestCase struct {
	name                string
	asset               string
	quote               string
	endpoint            string
	responseStatus      int
	responseBody        string
	expectedPass        bool
	expectedStatusCode  int
	expectedData        map[string]interface{}
	expectedErrorSubstr string
}

func TestExternalAdapterClient_Query(t *testing.T) {
	testCases := []queryTestCase{
		{
			name:               "success",
			asset:              "BTC",
			quote:              "USD",
			endpoint:           "price",
			responseStatus:     http.StatusOK,
			responseBody:       `{"result": "ok", "value": 123}`,
			expectedPass:       true,
			expectedStatusCode: http.StatusOK,
			expectedData: map[string]interface{}{
				"result": "ok",
				"value":  float64(123),
			},
		},
		{
			name:                "failure - non-200 response",
			asset:               "BTC",
			quote:               "USD",
			endpoint:            "price",
			responseStatus:      http.StatusInternalServerError,
			responseBody:        `internal error`,
			expectedPass:        false,
			expectedStatusCode:  http.StatusInternalServerError,
			expectedErrorSubstr: "500",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var reqData map[string]map[string]interface{}
				if err := json.NewDecoder(r.Body).Decode(&reqData); err != nil {
					http.Error(w, "bad request", http.StatusBadRequest)
					return
				}
				data := reqData["data"]
				if data["from"] != tc.asset || data["to"] != tc.quote || data["endpoint"] != tc.endpoint {
					http.Error(w, "invalid payload", http.StatusBadRequest)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.responseStatus)
				if _, err := w.Write([]byte(tc.responseBody)); err != nil {
					return
				}
			}))
			defer server.Close()

			client := NewExternalAdapterClient(server.Client())
			req := ExternalAdapterRequest{
				AdapterURL: server.URL,
				Asset:      tc.asset,
				Quote:      tc.quote,
				Endpoint:   tc.endpoint,
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			resp, err := client.Query(ctx, req)
			if err != nil {
				t.Fatalf("unexpected error from Query: %v", err)
			}

			if resp.Pass != tc.expectedPass {
				t.Errorf("expected Pass to be %v, got %v", tc.expectedPass, resp.Pass)
			}
			if resp.StatusCode != tc.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tc.expectedStatusCode, resp.StatusCode)
			}

			if !tc.expectedPass {
				if resp.ErrorMessage == "" {
					t.Error("expected an error message, but got empty")
				}
				if !contains(resp.ErrorMessage, tc.expectedErrorSubstr) {
					t.Errorf("expected error message to contain %q, got %q", tc.expectedErrorSubstr, resp.ErrorMessage)
				}
				return
			}

			gotData, err := json.Marshal(resp.Data)
			if err != nil {
				t.Fatalf("failed to marshal response data: %v", err)
			}
			expectedData, err := json.Marshal(tc.expectedData)
			if err != nil {
				t.Fatalf("failed to marshal expected data: %v", err)
			}
			if !bytes.Equal(gotData, expectedData) {
				t.Errorf("expected data %s, got %s", expectedData, gotData)
			}
		})
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (substr == "" || (len(substr) > 0 && (str == substr || (len(str) > len(substr) && (str[len(str)-len(substr):] == substr || str[:len(substr)] == substr)))))
}
