package matrix

import (
	"context"
	"fmt"
	"strconv"
)

// CRES smoke & regression topology structure.
type TopologyConfig struct {
	Topology string `json:"topology,omitempty"`
	Configs  string `json:"configs"`
}

// CRES smoke options.
type CRESmokeOptions struct {
	Dir        string
	RunID      string
	RunAttempt string
	SpotFlag   string
}

// CRES smoke matrix entry.
type CRESmokeEntry struct {
	TestName string `json:"test_name"`
	Topology string `json:"topology"`
	Configs  string `json:"configs"`
	TestID   int    `json:"test_id"`
	RunsOn   string `json:"runs_on"`
}

// CRERegression options.
type CRERegressionOptions struct {
	Dir        string
	RunID      string
	RunAttempt string
	SpotFlag   string
}

// CRERegression matrix entry.
type CRERegressionEntry struct {
	TestName string `json:"test_name"`
	TestID   int    `json:"test_id"`
	Configs  string `json:"configs"`
	RunsOn   string `json:"runs_on"`
}

var defaultCRESmokePerTestTopologies = map[string][]TopologyConfig{
	"Test_CRE_V2_Sharding": {
		{Topology: "workflow-gateway-sharded", Configs: "configs/workflow-gateway-sharded-don.toml"},
	},
	"Test_CRE_V2_ShardingWithHttpTrigger": {
		{Topology: "workflow-gateway-sharded", Configs: "configs/workflow-gateway-sharded-don.toml"},
	},
	"Test_CRE_V2_ShardManualAssignment": {
		{Topology: "workflow-gateway-sharded-manual", Configs: "configs/workflow-gateway-sharded-manual.toml"},
	},
	"Test_CRE_V2_ShardRingOCROverrides": {
		{Topology: "workflow-gateway-sharded-ringocr-overrides", Configs: "configs/workflow-gateway-sharded-ringocr-overrides.toml"},
	},
	"Test_CRE_V2_Module_Cache": {
		{Topology: "workflow-gateway-cache-test", Configs: "configs/workflow-gateway-don-cache-test.toml"},
	},
	"Test_CRE_V2_HTTP_Action_Multi_Gateway": {
		{Topology: "workflow-gateway-capabilities-multi-gateway", Configs: "configs/workflow-gateway-capabilities-multi-gateway-don.toml"},
	},
	"Test_CRE_V2_ConfidentialWorkflows_Relay": {
		{Topology: "workflow-gateway-capabilities-confidential-workflows", Configs: "configs/workflow-gateway-capabilities-don-confidential-workflows.toml"},
	},
}

var defaultCRERegressionPerTestConfigs = map[string]string{
	"Test_CRE_V2_Stellar_Regression": "configs/workflow-gateway-don-stellar.toml",
}

// BuildCRESmokeMatrix discovers tests in dir and constructs CRE smoke test matrix.
func BuildCRESmokeMatrix(ctx context.Context, opts CRESmokeOptions) ([]CRESmokeEntry, error) {
	testNames, err := DiscoverGoTestNames(opts.Dir, DiscoverOptions{
		IgnoredPatterns: []string{"^Test_Upgrade", "^TestMain$"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed discovering tests in %s: %w", opts.Dir, err)
	}

	spotFlag := opts.SpotFlag
	if spotFlag == "" {
		spotFlag = "spot=co"
	}
	runAttempt := opts.RunAttempt
	if runAttempt == "" {
		runAttempt = "1"
	}

	var entries []CRESmokeEntry
	testID := 0

	for _, name := range testNames {
		topologies, hasOverride := defaultCRESmokePerTestTopologies[name]
		if !hasOverride {
			topologies = []TopologyConfig{
				{Topology: "workflow-gateway-capabilities", Configs: "configs/workflow-gateway-capabilities-don.toml"},
			}
		}

		for _, top := range topologies {
			runsOn := fmt.Sprintf("runs-on=%s-%d-%s/cpu=16/ram=64/family=m7i+m8i/%s/image=ubuntu24-full-x64/extras=s3-cache+tmpfs",
				opts.RunID, testID, runAttempt, spotFlag)
			entries = append(entries, CRESmokeEntry{
				TestName: name,
				Topology: top.Topology,
				Configs:  top.Configs,
				TestID:   testID,
				RunsOn:   runsOn,
			})
			testID++
		}
	}

	return entries, nil
}

// BuildCRERegressionMatrix discovers tests in dir and constructs CRE regression test matrix.
func BuildCRERegressionMatrix(ctx context.Context, opts CRERegressionOptions) ([]CRERegressionEntry, error) {
	testNames, err := DiscoverGoTestNames(opts.Dir, DiscoverOptions{
		IgnoredPatterns: []string{"^TestMain$"},
	})
	if err != nil {
		return nil, fmt.Errorf("failed discovering tests in %s: %w", opts.Dir, err)
	}

	spotFlag := opts.SpotFlag
	if spotFlag == "" {
		spotFlag = "spot=co"
	}
	runAttempt := opts.RunAttempt
	if runAttempt == "" {
		runAttempt = "1"
	}

	var entries []CRERegressionEntry
	for i, name := range testNames {
		configs, ok := defaultCRERegressionPerTestConfigs[name]
		if !ok {
			configs = "configs/workflow-gateway-capabilities-don.toml"
		}

		runsOn := fmt.Sprintf("runs-on=%s-%s-%s/cpu=16/ram=64/family=m7i+m8i/%s/image=ubuntu24-full-x64/extras=s3-cache+tmpfs",
			opts.RunID, strconv.Itoa(i), runAttempt, spotFlag)

		entries = append(entries, CRERegressionEntry{
			TestName: name,
			TestID:   i,
			Configs:  configs,
			RunsOn:   runsOn,
		})
	}

	return entries, nil
}
