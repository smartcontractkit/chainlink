package jobs

import (
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestLLOSpec_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name       string
		spec       LLOSpec
		obs        ObservationSource
		wantSubstr []string
	}{
		{
			name: "multiple datasources with valid paths",
			spec: LLOSpec{
				Base: Base{
					Name:          "ETH/USD-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				StreamID: "1000000001",
			},
			obs: ObservationSource{
				Datasources: []Datasource{
					{
						BridgeName: "coinmetrics",
						ReqData:    `{"data":{"endpoint":"cryptolwba","from":"ETH","to":"USD"}}`,
					},
					{
						BridgeName: "ncfx",
						ReqData:    `{"data":{"endpoint":"cryptolwba","from":"ETH","to":"USD"}}`,
					},
				},
				AllowedFaults: 2,
				Benchmark: ReportField{
					ResultPath: "data,mid",
				},
				Bid: ReportField{
					ResultPath: "data,bid",
				},
				Ask: ReportField{
					ResultPath: "data,ask",
				},
			},
			wantSubstr: []string{
				`bridge-coinmetrics`,
				`bridge-ncfx`,
				`allowedFaults=2`,
				`path=\"data,mid\"`,
				`path=\"data,bid\"`,
				`path=\"data,ask\"`,
			},
		},
		{
			name: "empty datasource list",
			spec: LLOSpec{
				Base: Base{
					Name:          "Empty-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				StreamID: "2000000002",
			},
			obs: ObservationSource{
				Datasources:   []Datasource{},
				AllowedFaults: 1,
				Benchmark:     ReportField{ResultPath: "data,benchmark"},
				Bid:           ReportField{ResultPath: "data,bid"},
				Ask:           ReportField{ResultPath: "data,ask"},
			},
			wantSubstr: []string{
				`allowedFaults=1`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.SetObservationSource(tc.obs)
			require.NoError(t, err)
			tomlBytes, err := tc.spec.MarshalTOML()
			require.NoError(t, err)
			result := string(tomlBytes)
			for _, substr := range tc.wantSubstr {
				require.True(t, strings.Contains(result, substr),
					"result %q does not contain expected substring %q", result, substr)
			}
		})
	}
}

func TestLLOSpec_Error(t *testing.T) {
	originalTmpl := observationTmpl
	defer func() {
		observationTmpl = originalTmpl
	}()

	faultyTmpl := template.Must(template.New("faulty").
		Funcs(template.FuncMap{
			"error": func(msg string) (string, error) {
				return "", fmt.Errorf("%s", msg)
			},
		}).
		Parse(`{{ error "forced error" }}`))
	observationTmpl = faultyTmpl

	spec := LLOSpec{}
	obs := ObservationSource{}
	_, err := spec.buildObservationSource(obs)
	require.Error(t, err, "expected error from faulty template execution")
}
