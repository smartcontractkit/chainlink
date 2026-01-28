package automation

import (
	"testing"
	"time"

	"github.com/smartcontractkit/chainlink-testing-framework/framework/leak"
	de "github.com/smartcontractkit/chainlink/devenv"
	"github.com/smartcontractkit/chainlink/devenv/contracts"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	testCases := []testcase{
		{
			registryVersion:          contracts.RegistryVersion_2_1,
			name:                     "registry_2_1_logtrigger",
			upkeepCount:              5,
			upkeepExecutionTimeout:   "2h",
			expectedUpkeepExecutions: 100, //TODO increase maybe to 1000 once we have figured out which funds are missing and make the counter stop around 95-99 executions
			testKeyFundingEth:        15,
			upkeepFundingLink:        15,
		},
	}

	for _, tc := range testCases {
		start := time.Now()

		basicAutomationTest(t, tc)

		outputFile := "../../env-out.toml"
		in, err := de.LoadOutput[de.Cfg](outputFile)
		require.NoError(t, err)

		l, err := leak.NewCLNodesLeakDetector(leak.NewResourceLeakChecker())
		require.NoError(t, err)
		errs := l.Check(&leak.CLNodesCheck{
			// since the test is stable we assert absolute values
			// no more than 30% CPU and 200Mb (last 5m)
			ComparisonMode:  leak.ComparisonModeAbsolute,
			NumNodes:        in.NodeSets[0].Nodes,
			Start:           start,
			End:             time.Now(),
			WarmUpDuration:  30 * time.Minute,
			CPUThreshold:    30.0,
			MemoryThreshold: 200.0,
		})
		require.NoError(t, errs)
	}
}
