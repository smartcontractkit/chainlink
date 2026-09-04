package matrix

import (
	"context"
	"fmt"
)

// CCIPSystemOptions contains parameters for generating CCIP system test matrix.
type CCIPSystemOptions struct {
	RunID      string
	RunAttempt string
	SpotFlag   string
}

// CCIPSystemEntry is a single test entry in the CCIP system matrix.
type CCIPSystemEntry struct {
	TestName            string `json:"test_name"`
	Timeout             string `json:"timeout,omitempty"`
	JobTimeout          int    `json:"job_timeout,omitempty"`
	SelectedNetwork     string `json:"selected_network,omitempty"`
	RMNRageProxyVersion string `json:"rmn_rageproxy_version,omitempty"`
	RMNAFN2ProxyVersion string `json:"rmn_afn2proxy_version,omitempty"`
	TestID              int    `json:"test_id"`
	RunsOn              string `json:"runs_on"`
}

// BuildCCIPSystemMatrix generates the matrix for CCIP system tests.
func BuildCCIPSystemMatrix(ctx context.Context, opts CCIPSystemOptions) ([]CCIPSystemEntry, error) {
	spotFlag := opts.SpotFlag
	if spotFlag == "" {
		spotFlag = "spot=co"
	}
	runAttempt := opts.RunAttempt
	if runAttempt == "" {
		runAttempt = "1"
	}

	definitions := []struct {
		TestName            string
		Timeout             string
		JobTimeout          int
		SelectedNetwork     string
		RMNRageProxyVersion string
		RMNAFN2ProxyVersion string
	}{
		{
			TestName:        "Test_CCIPGasPriceUpdatesWriteFrequency",
			Timeout:         "15m",
			SelectedNetwork: "SIMULATED_1,SIMULATED_2",
		},
		{
			TestName:            "TestRMN_GlobalCurseTwoMessagesOnTwoLanes",
			Timeout:             "15m",
			SelectedNetwork:     "SIMULATED_1,SIMULATED_2",
			RMNRageProxyVersion: "master-amd6416f5d86",
			RMNAFN2ProxyVersion: "master-amd64-10b42b2",
		},
		{
			TestName:        "TestDeleteCCIPJobs-TestRevokeJobs",
			Timeout:         "15m",
			JobTimeout:      20,
			SelectedNetwork: "SIMULATED_1,SIMULATED_2",
		},
	}

	entries := make([]CCIPSystemEntry, len(definitions))
	for i, def := range definitions {
		runsOn := fmt.Sprintf("runs-on=%s-%d-%s/cpu=8/ram=64/family=r6i+r7i+r8i/%s/image=ubuntu24-full-x64/extras=s3-cache+tmpfs",
			opts.RunID, i, runAttempt, spotFlag)

		entries[i] = CCIPSystemEntry{
			TestName:            def.TestName,
			Timeout:             def.Timeout,
			JobTimeout:          def.JobTimeout,
			SelectedNetwork:     def.SelectedNetwork,
			RMNRageProxyVersion: def.RMNRageProxyVersion,
			RMNAFN2ProxyVersion: def.RMNAFN2ProxyVersion,
			TestID:              i,
			RunsOn:              runsOn,
		}
	}

	return entries, nil
}
