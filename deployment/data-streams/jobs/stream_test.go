package jobs

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const medianSpecTOMLMultiple = `name = 'BTC/USD-Test'
type = 'stream'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
streamID = '1000'
observationSource = """
// data source 1
ds1_payload [type=bridge name=\"bridge-bridge1\" timeout=\"50s\" requestData={\"data\":{\"endpoint\":\"test1\"}}];

ds1_benchmark [type=jsonparse path=\"data,median\"];
// data source 2
ds2_payload [type=bridge name=\"bridge-bridge2\" timeout=\"50s\" requestData={\"data\":{\"endpoint\":\"test2\"}}];

ds2_benchmark [type=jsonparse path=\"data,median\"];
ds1_payload -> ds1_benchmark -> benchmark_price;
ds2_payload -> ds2_benchmark -> benchmark_price;
benchmark_price [type=median allowedFaults=2 index=0];
"""
`

const medianSpecTOMLEmpty = `name = 'Empty-Median-Test'
type = 'stream'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
streamID = '2000'
observationSource = """
benchmark_price [type=median allowedFaults=1 index=0];
"""
`

const quoteSpecTOMLMultiple = `name = 'BTC/USD-Quote'
type = 'stream'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
streamID = '3000'
observationSource = """
// data source 1
ds1_payload [type=bridge name=\"bridge-bridge1\" timeout=\"50s\" requestData={\"data\":{\"endpoint\":\"quote1\"}}];

ds1_benchmark [type=jsonparse path=\"data,benchmark\"];
ds1_bid [type=jsonparse path=\"data,bid\"];
ds1_ask [type=jsonparse path=\"data,ask\"];
// data source 2
ds2_payload [type=bridge name=\"bridge-bridge2\" timeout=\"50s\" requestData={\"data\":{\"endpoint\":\"quote2\"}}];

ds2_benchmark [type=jsonparse path=\"data,benchmark\"];
ds2_bid [type=jsonparse path=\"data,bid\"];
ds2_ask [type=jsonparse path=\"data,ask\"];
ds1_payload -> ds1_benchmark -> benchmark_price;
ds2_payload -> ds2_benchmark -> benchmark_price;
benchmark_price [type=median allowedFaults=3 index=0];

ds1_payload -> ds1_bid -> bid_price;
ds2_payload -> ds2_bid -> bid_price;
bid_price [type=median allowedFaults=3 index=1];

ds1_payload -> ds1_ask -> ask_price;
ds2_payload -> ds2_ask -> ask_price;
ask_price [type=median allowedFaults=3 index=2];
"""
`

const quoteSpecTOMLEmpty = `name = 'Empty-Quote-Test'
type = 'stream'
schemaVersion = 1
externalJobID = '00000000-0000-0000-0000-000000000000'
streamID = '4000'
observationSource = """
benchmark_price [type=median allowedFaults=1 index=0];

bid_price [type=median allowedFaults=1 index=1];

ask_price [type=median allowedFaults=1 index=2];
"""
`

func TestStreamJobSpec_Median_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name string
		spec StreamJobSpec
		obs  MedianObservationSource
		want string
	}{
		{
			name: "multiple datasources with valid paths",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "BTC/USD-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
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
			want: medianSpecTOMLMultiple,
		},
		{
			name: "empty datasource list",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "Empty-Median-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
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
			want: medianSpecTOMLEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.SetObservationSource(tc.obs)
			require.NoError(t, err)
			tomlBytes, err := tc.spec.MarshalTOML()
			require.NoError(t, err)
			got := string(tomlBytes)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestStreamJobSpec_Quote_MarshalTOML(t *testing.T) {
	testCases := []struct {
		name string
		spec StreamJobSpec
		obs  QuoteObservationSource
		want string
	}{
		{
			name: "multiple datasources with valid paths",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "BTC/USD-Quote",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
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
			want: quoteSpecTOMLMultiple,
		},
		{
			name: "empty datasource list",
			spec: StreamJobSpec{
				Base: Base{
					Name:          "Empty-Quote-Test",
					Type:          "stream",
					SchemaVersion: 1,
					ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
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
			want: quoteSpecTOMLEmpty,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.spec.SetObservationSource(tc.obs)
			require.NoError(t, err)

			tomlBytes, err := tc.spec.MarshalTOML()
			require.NoError(t, err)
			got := string(tomlBytes)
			require.Equal(t, tc.want, got)
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
			ExternalJobID: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
		},
		StreamID: "5000",
	}
	err := spec.SetObservationSource(errorPipeline{})
	require.Error(t, err)
	require.Contains(t, err.Error(), "forced error")
}
