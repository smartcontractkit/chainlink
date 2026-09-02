package matrix

import (
	"context"
	"fmt"
)

// CREMixedEnvOptions contains parameters for generating CRE mixed-env test matrix.
type CREMixedEnvOptions struct {
	RunID      string
	RunAttempt string
	SpotFlag   string
}

// CREMixedEnvEntry is a single test entry in the CRE mixed-env matrix.
type CREMixedEnvEntry struct {
	TestName string `json:"test_name"`
	Configs  string `json:"configs"`
	TestID   int    `json:"test_id"`
	RunsOn   string `json:"runs_on"`
}

var defaultCREMixedEnvTests = []string{
	"Test_CRE_V2_Suite_Bucket_A",
	"Test_CRE_V2_Suite_Bucket_B",
	"Test_CRE_V2_EVM_Read_HeavyCalls",
	"Test_CRE_V2_EVM_Read_StateQueries",
	"Test_CRE_V2_EVM_Read_TxArtifacts",
	"Test_CRE_V2_ConfidentialWorkflows_Relay",
}

var defaultCREMixedEnvConfigs = map[string]string{
	"Test_CRE_V2_ConfidentialWorkflows_Relay": "configs/mixed-env-confidential-workflows.toml",
}

// BuildCREMixedEnvMatrix generates the matrix for CRE mixed-env tests.
func BuildCREMixedEnvMatrix(ctx context.Context, opts CREMixedEnvOptions) ([]CREMixedEnvEntry, error) {
	spotFlag := opts.SpotFlag
	if spotFlag == "" {
		spotFlag = "spot=co"
	}
	runAttempt := opts.RunAttempt
	if runAttempt == "" {
		runAttempt = "1"
	}

	entries := make([]CREMixedEnvEntry, len(defaultCREMixedEnvTests))
	for i, name := range defaultCREMixedEnvTests {
		configs, ok := defaultCREMixedEnvConfigs[name]
		if !ok {
			configs = "configs/mixed-env-don.toml"
		}

		runsOn := fmt.Sprintf("runs-on=%s-%d-%s/cpu=16/ram=64/family=m7i+m8i/%s/image=ubuntu24-full-x64/extras=s3-cache+tmpfs",
			opts.RunID, i, runAttempt, spotFlag)

		entries[i] = CREMixedEnvEntry{
			TestName: name,
			Configs:  configs,
			TestID:   i,
			RunsOn:   runsOn,
		}
	}

	return entries, nil
}
