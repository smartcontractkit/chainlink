package changeset

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/smartcontractkit/chainlink/deployment"
	"github.com/smartcontractkit/chainlink/deployment/common/changeset"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/changeset/testutil"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jd"
	"github.com/smartcontractkit/chainlink/deployment/data-streams/jobs"
)

func TestDistributeStreamJobSpecs(t *testing.T) {
	t.Parallel()

	const donID = 1
	const donName = "don"
	const env = "env"

	memEnv := testutil.NewMemoryEnvV2(t, testutil.MemoryEnvConfig{
		ShouldDeployMCMS:      false,
		ShouldDeployLinkToken: false,
		NumNodes:              1,
		NodeLabels:            testutil.GetNodeLabels(donID, donName, env),
		CustomDBSetup: []string{
			// Setup the database with the list of bridges we're using.
			"INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)" +
				" VALUES ('bridge-api1', 'http://url', 0, '', '', '', now(), now());",
			"INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)" +
				" VALUES ('bridge-api2', 'http://url', 0, '', '', '', now(), now());",
			"INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)" +
				" VALUES ('bridge-api3', 'http://url', 0, '', '', '', now(), now());",
			"INSERT INTO bridge_types (name, url, confirmations, incoming_token_hash, salt, outgoing_token, created_at, updated_at)" +
				" VALUES ('bridge-api4', 'http://url', 0, '', '', '', now(), now());",
		},
	})

	// pick the first EVM chain selector
	chainSelector := memEnv.Environment.AllChainSelectors()[0]

	// insert a Configurator address for the given DON
	configuratorAddr := "0x4170ed0880ac9a755fd29b2688956bd959f923f4"
	err := memEnv.Environment.ExistingAddresses.Save(chainSelector, configuratorAddr,
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
		ChainSelectorEVM: chainSelector,
		Filter: &jd.ListFilter{
			DONID:    donID,
			DONName:  donName,
			EnvLabel: "env",
			Size:     1,
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
	}

	tests := []struct {
		name       string
		env        deployment.Environment
		config     CsDistributeStreamJobSpecsConfig
		prepConfFn func(CsDistributeStreamJobSpecsConfig) CsDistributeStreamJobSpecsConfig
		wantErr    *string
		wantSpec   string
	}{
		{
			name:     "success",
			env:      memEnv.Environment,
			config:   config,
			wantSpec: renderedSpec,
		},
		// TODO: Cover all failure cases.
	}

	cs := CsDistributeStreamJobSpecs{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf := tt.config
			if tt.prepConfFn != nil {
				conf = tt.prepConfFn(tt.config)
			}
			_, out, err := changeset.ApplyChangesetsV2(t,
				tt.env,
				[]changeset.ConfiguredChangeSet{
					changeset.Configure(cs, conf),
				},
			)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), *tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.Len(t, out, 1)
			require.Len(t, out[0].Jobs, 1)
			require.Equal(t, stripExternalJobID(tt.wantSpec), stripExternalJobID(out[0].Jobs[0].Spec))
		})
	}
}

// Remove the externalJobID line from the spec. This is needed because the externalJobID is generated randomly
// and we want to exclude it from the comparison.
func stripExternalJobID(spec string) string {
	idx := strings.Index(spec, "externalJobID = ")
	strLen := len(fmt.Sprintf(`externalJobID = "%s"`, uuid.New()))
	return spec[:idx] + spec[idx+strLen:]
}
