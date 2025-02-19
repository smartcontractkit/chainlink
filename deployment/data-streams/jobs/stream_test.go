package jobs

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestStreamJobSpec_Median_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name       string
		spec       StreamJobSpec
		obs        MedianObservationSource
		wantSubstr []string
	}{
		{
			name: "multiple datasources with valid paths",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "BTC/USD-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				StreamID: "1000",
			},
			obs: MedianObservationSource{
				BaseObservationSource: BaseObservationSource{
					Datasources: []Datasource{
						{
							BridgeName: "bridge1",
							ReqData:    `{"data":{"endpoint":"test1"}}`,
						},
						{
							BridgeName: "bridge2",
							ReqData:    `{"data":{"endpoint":"test2"}}`,
						},
					},
					AllowedFaults: 2,
					Benchmark: ReportFieldLLO{
						ResultPath: "data,median",
					},
				},
			},
			wantSubstr: []string{
				"bridge-bridge1",
				"bridge-bridge2",
				"allowedFaults=2",
				`data,median`,
			},
		},
		{
			name: "empty datasource list",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "Empty-Median-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				StreamID: "2000",
			},
			obs: MedianObservationSource{
				BaseObservationSource: BaseObservationSource{
					Datasources:   []Datasource{},
					AllowedFaults: 1,
					Benchmark: ReportFieldLLO{
						ResultPath: "data,empty",
					},
				},
			},
			wantSubstr: []string{
				"allowedFaults=1",
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
				require.Contains(t, result, substr,
					"result %q does not contain expected substring %q", result, substr)
			}
		})
	}
}

func TestStreamJobSpec_Quote_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name       string
		spec       StreamJobSpec
		obs        QuoteObservationSource
		wantSubstr []string
	}{
		{
			name: "multiple datasources with valid paths",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "BTC/USD-Quote",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				StreamID: "3000",
			},
			obs: QuoteObservationSource{
				BaseObservationSource: BaseObservationSource{
					Datasources: []Datasource{
						{
							BridgeName: "bridge1",
							ReqData:    `{"data":{"endpoint":"quote1"}}`,
						},
						{
							BridgeName: "bridge2",
							ReqData:    `{"data":{"endpoint":"quote2"}}`,
						},
					},
					AllowedFaults: 3,
					Benchmark: ReportFieldLLO{
						ResultPath: "data,benchmark",
					},
				},
				Bid: ReportFieldLLO{
					ResultPath: "data,bid",
				},
				Ask: ReportFieldLLO{
					ResultPath: "data,ask",
				},
			},
			wantSubstr: []string{
				"bridge-bridge1",
				"bridge-bridge2",
				"allowedFaults=3",
				`data,benchmark`,
				`data,bid`,
				`data,ask`,
			},
		},
		{
			name: "empty datasource list",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "Empty-Quote-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.New(),
				},
				StreamID: "4000",
			},
			obs: QuoteObservationSource{
				BaseObservationSource: BaseObservationSource{
					Datasources:   []Datasource{},
					AllowedFaults: 1,
					Benchmark: ReportFieldLLO{
						ResultPath: "data,empty",
					},
				},
				Bid: ReportFieldLLO{
					ResultPath: "data,emptyBid",
				},
				Ask: ReportFieldLLO{
					ResultPath: "data,emptyAsk",
				},
			},
			wantSubstr: []string{
				"allowedFaults=1",
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
				require.Contains(t, result, substr,
					"result %q does not contain expected substring %q", result, substr)
			}
		})
	}
}

type errorPipeline struct{}

func (e errorPipeline) Render() (string, error) {
	return "", errors.New("forced error")
}

func TestStreamJobSpec_SetObservationSource_Error(t *testing.T) {
	spec := StreamJobSpec{
		Base: Base{
			Name:          "Error-Test",
			Type:          "stream",
			SchemaVersion: 1,
			ExternalJobID: uuid.New(),
		},
		StreamID: "5000",
	}

	err := spec.SetObservationSource(errorPipeline{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "forced error")
}
