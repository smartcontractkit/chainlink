package changeset

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/utils/pointer"
)

func TestDistributeStreamJobSpecs(t *testing.T) {
	t.Parallel()

	const donID = 1
	const donName = "don"
	const envName = "envName"

	env := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      false,
		ShouldDeployLinkToken: false,
		NumNodes:              3,
		NodeLabels:            testutil.GetNodeLabels(donID, donName, envName),
		CustomDBSetup: []string{
			// Seed the database with the list of bridges we're using.
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api1', 'http://url', 0, '', '', '', now(), now());`,
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api2', 'http://url', 0, '', '', '', now(), now());`,
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api3', 'http://url', 0, '', '', '', now(), now());`,
			`INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)
				VALUES ('bridge-api4', 'http://url', 0, '', '', '', now(), now());`,
		},
	}).Environment

	// pick the first EVM chain selector
	chainSelector := env.AllChainSelectors()[0]

	// insert a Configurator address for the given DON
	configuratorAddr := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err := env.ExistingAddresses.Save(chainSelector, configuratorAddr, //nolint: staticcheck // I don't care that ExistingAddresses is deprecated. We will fix it later.
		deployment.TypeAndVersion{
			Type:    "Configurator",
			Version: deployment.Version1_0_0,
			Labels:  deployment.NewLabelSet("don-1"),
		})
	require.NoError(t, err)

	renderedSpec := `name = 'ICP/USD-RefPrice | 1000001038'
type = 'stream'
schemaVersion = 1
externalJobID = 'a818bfcb-28cd-4eb2-a4af-6ff9e35b0fb0'
streamID = 1000001038
observationSource = """
// data source 1
ds1_payload [type=bridge name=\"bridge-api1\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"ICP\\\",\\\"to\\\":\\\"USD\\\"}}\"];

ds1_benchmark [type=jsonparse path=\"data,mid\"];
ds1_bid [type=jsonparse path=\"data,bid\"];
ds1_ask [type=jsonparse path=\"data,ask\"];
// data source 2
ds2_payload [type=bridge name=\"bridge-api2\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"ICP\\\",\\\"to\\\":\\\"USD\\\"}}\"];

ds2_benchmark [type=jsonparse path=\"data,mid\"];
ds2_bid [type=jsonparse path=\"data,bid\"];
ds2_ask [type=jsonparse path=\"data,ask\"];
// data source 3
ds3_payload [type=bridge name=\"bridge-api3\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"ICP\\\",\\\"to\\\":\\\"USD\\\"}}\"];

ds3_benchmark [type=jsonparse path=\"data,mid\"];
ds3_bid [type=jsonparse path=\"data,bid\"];
ds3_ask [type=jsonparse path=\"data,ask\"];
// data source 4
ds4_payload [type=bridge name=\"bridge-api4\" timeout=\"50s\" requestData=\"{\\\"data\\\":{\\\"endpoint\\\":\\\"cryptolwba\\\",\\\"from\\\":\\\"ICP\\\",\\\"to\\\":\\\"USD\\\"}}\"];

ds4_benchmark [type=jsonparse path=\"data,mid\"];
ds4_bid [type=jsonparse path=\"data,bid\"];
ds4_ask [type=jsonparse path=\"data,ask\"];
ds1_payload -> ds1_benchmark -> benchmark_price;
ds2_payload -> ds2_benchmark -> benchmark_price;
ds3_payload -> ds3_benchmark -> benchmark_price;
ds4_payload -> ds4_benchmark -> benchmark_price;
benchmark_price [type=median allowedFaults=3 index=0];

ds1_payload -> ds1_bid -> bid_price;
ds2_payload -> ds2_bid -> bid_price;
ds3_payload -> ds3_bid -> bid_price;
ds4_payload -> ds4_bid -> bid_price;
bid_price [type=median allowedFaults=3 index=1];

ds1_payload -> ds1_ask -> ask_price;
ds2_payload -> ds2_ask -> ask_price;
ds3_payload -> ds3_ask -> ask_price;
ds4_payload -> ds4_ask -> ask_price;
ask_price [type=median allowedFaults=3 index=2];
"""
`

	config := CsDistributeStreamJobSpecsConfig{
		Filter: &jd.ListFilter{
			DONID:          donID,
			DONName:        donName,
			EnvLabel:       envName,
			NumOracleNodes: 3,
		},
		Streams: []StreamSpecConfig{
			{
				StreamID:   1000001038,
				Name:       "ICP/USD-RefPrice",
				StreamType: jobs.StreamTypeQuote,
				ReportFields: jobs.QuoteReportFields{
					Bid: jobs.ReportFieldLLO{
						ResultPath: "data,bid",
					},
					Benchmark: jobs.ReportFieldLLO{
						ResultPath: "data,mid",
					},
					Ask: jobs.ReportFieldLLO{
						ResultPath: "data,ask",
					},
				},
				EARequestParams: EARequestParams{
					Endpoint: "cryptolwba",
					From:     "ICP",
					To:       "USD",
				},
				APIs: []string{"api1", "api2", "api3", "api4"},
			},
		},
		NodeNames: []string{"node-0", "node-1", "node-2"},
	}

	tests := []struct {
		name        string
		prepConfFn  func(CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig
		wantErr     *string
		wantSpec    string
		wantNumJobs int
	}{
		{
			name:        "success",
			wantSpec:    renderedSpec,
			wantNumJobs: 3,
		},
		{
			// This test only makes sense when run after "success" because the two use the same ExternalJobID.
			name:        "success proposing updates to existing jobs",
			wantSpec:    renderedSpec,
			wantNumJobs: 3,
		},
		{
			// This happens to also be a job update when run after "success" because the two use the same ExternalJobID.
			name: "success sending jobs to a subset of nodes",
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.NodeNames = []string{"node-0"}
				c.Filter = &jd.ListFilter{
					DONID:          donID,
					DONName:        donName,
					EnvLabel:       envName,
					NumOracleNodes: 1,
				}
				return c
			},
			wantSpec:    renderedSpec,
			wantNumJobs: 1,
		},
		{
			// This happens to also be a job update when run after "success" because the two use the same ExternalJobID.
			name: "success sending jobs to a different subset of nodes",
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.NodeNames = []string{"node-1", "node-2"}
				c.Filter = &jd.ListFilter{
					DONID:          donID,
					DONName:        donName,
					EnvLabel:       envName,
					NumOracleNodes: 2,
				}
				return c
			},
			wantSpec:    renderedSpec,
			wantNumJobs: 2,
		},
		{
			name: "failure when the node name is not found",
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.NodeNames = []string{"non-existing-node"}
				c.Filter = &jd.ListFilter{
					DONID:          donID,
					DONName:        donName,
					EnvLabel:       envName,
					NumOracleNodes: 1,
				}
				return c
			},
			wantErr: pointer.To("failed to get oracle nodes"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			conf := config
			if tc.prepConfFn != nil {
				conf = tc.prepConfFn(conf)
			}
			_, out, err := changeset.ApplyChangesetsV2(t,
				env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(CsDistributeStreamJobSpecs{}, conf),
				},
			)

			if tc.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), *tc.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, tc.wantNumJobs)
			for i := 0; i < tc.wantNumJobs; i++ {
				require.Equal(t,
					testutil.StripLineContaining(tc.wantSpec, []string{"externalJobID"}),
					testutil.StripLineContaining(out[0].Jobs[i].Spec, []string{"externalJobID"}),
				)
			}
		})
	}
}

func TestValidatePreconditions(t *testing.T) {
	t.Parallel()

	e := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{}).Environment

	config := CsDistributeStreamJobSpecsConfig{
		Filter: &jd.ListFilter{
			DONID:          1,
			DONName:        "don",
			EnvLabel:       "env",
			NumOracleNodes: 1,
		},
		Streams: []StreamSpecConfig{
			{
				StreamID:   1000001038,
				Name:       "ICP/USD-RefPrice",
				StreamType: jobs.StreamTypeQuote,
				ReportFields: jobs.QuoteReportFields{
					Bid: jobs.ReportFieldLLO{
						ResultPath: "data,bid",
					},
					Benchmark: jobs.ReportFieldLLO{
						ResultPath: "data,mid",
					},
					Ask: jobs.ReportFieldLLO{
						ResultPath: "data,ask",
					},
				},
				EARequestParams: EARequestParams{
					Endpoint: "cryptolwba",
					From:     "ICP",
					To:       "USD",
				},
				APIs: []string{"api1", "api2", "api3", "api4"},
			},
		},
		NodeNames: []string{"testnode-0"},
	}

	tests := []struct {
		name       string
		config     CsDistributeStreamJobSpecsConfig
		prepConfFn func(CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig
		wantErr    *string
	}{
		{
			name:    "valid config",
			config:  config,
			wantErr: nil,
		},
		{
			name:   "no filter",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Filter = nil
				return c
			},
			wantErr: pointer.To("filter is required"),
		},
		{
			name:   "no DONID",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Filter = &jd.ListFilter{
					DONID:   0,
					DONName: "don",
				}
				return c
			},
			wantErr: pointer.To("DONID and DONName are required"),
		},
		{
			name:   "no DONName",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Filter = &jd.ListFilter{
					DONID:   111,
					DONName: "",
				}
				return c
			},
			wantErr: pointer.To("DONID and DONName are required"),
		},
		{
			name:   "empty streams",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{}
				return c
			},
			wantErr: pointer.To("streams are required"),
		},
		{
			name:   "no streamID",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{{}}
				return c
			},
			wantErr: pointer.To("streamID is required for each stream"),
		},
		{
			name:   "no name",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{
					{
						StreamID: 1000001038,
					},
				}
				return c
			},
			wantErr: pointer.To("name is required for each stream"),
		},
		{
			name:   "invalid stream type",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{
					{
						StreamID: 1000001038,
						Name:     "ICP/USD-RefPrice",
					},
				}
				return c
			},
			wantErr: pointer.To("stream type is not valid"),
		},
		{
			name:   "no report fields",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{
					{
						StreamID:   1000001038,
						Name:       "ICP/USD-RefPrice",
						StreamType: jobs.StreamTypeQuote,
					},
				}
				return c
			},
			wantErr: pointer.To("report fields are required for each stream"),
		},
		{
			name:   "no EARequestParams",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{
					{
						StreamID:     1000001038,
						Name:         "ICP/USD-RefPrice",
						StreamType:   jobs.StreamTypeQuote,
						ReportFields: jobs.QuoteReportFields{},
					},
				}
				return c
			},
			wantErr: pointer.To("endpoint is required for each EARequestParam on each stream"),
		},
		{
			name:   "no endpoint",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{
					{
						StreamID:        1000001038,
						Name:            "ICP/USD-RefPrice",
						StreamType:      jobs.StreamTypeQuote,
						ReportFields:    jobs.QuoteReportFields{},
						EARequestParams: EARequestParams{},
					},
				}
				return c
			},
			wantErr: pointer.To("endpoint is required for each EARequestParam on each stream"),
		},
		{
			name:   "no APIs",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Streams = []StreamSpecConfig{
					{
						StreamID:     1000001038,
						Name:         "ICP/USD-RefPrice",
						StreamType:   jobs.StreamTypeQuote,
						ReportFields: jobs.QuoteReportFields{},
						EARequestParams: EARequestParams{
							Endpoint: "cryptolwba",
						},
					},
				}
				return c
			},
			wantErr: pointer.To("at least one API is required for each stream"),
		},
		{
			name:   "no node names",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.NodeNames = []string{}
				return c
			},
			wantErr: pointer.To("at least one node name is required"),
		},
		{
			name:   "filter size does not match node names",
			config: config,
			prepConfFn: func(c CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig {
				c.Filter.NumOracleNodes = 2
				c.NodeNames = []string{"node-0"}
				return c
			},
			wantErr: pointer.To("number of node names (1) does not match filter size (2)"),
		},
	}

	cs := CsDistributeStreamJobSpecs{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepConfFn != nil {
				tt.config = tt.prepConfFn(tt.config)
			}
			err := cs.VerifyPreconditions(e, tt.config)
			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), *tt.wantErr)
				return
			}
			require.NoError(t, err)
		})
	}
}
